package scylla

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	stdpath "path"
	"strconv"
	"time"
)

var (
	defaultTimeout              = 60 * time.Second
	retriesAllowed              = 3
	maxResponseBodyLength int64 = 1 << 20
)

// APIError represents an error that occurred while calling the API.
type APIError struct {
	URL        string
	Code       string
	Message    string
	StatusCode int
}

func makeError(text string, errCodes map[string]string, r *http.Response) *APIError {
	var err APIError
	if _, e := strconv.Atoi(text); e == nil {
		err.Code = text

		switch text := errCodes[text]; {
		case err.Message == "" && text == "":
			err.Message = "Request has failed. For more details consult the error code"
		case err.Message == "":
			err.Message = text
		case text != "":
			err.Message = err.Message + " (" + text + ")"
		}
	} else {
		err.Message = text
	}
	if err.URL == "" {
		err.URL = r.Request.URL.String()
	}
	if err.StatusCode == 0 {
		err.StatusCode = r.StatusCode
	}
	return &err
}

func (err *APIError) Error() string {
	return fmt.Sprintf("Error %q: %s (http status %d, url %q)", err.Code, err.Message, err.StatusCode, err.URL)
}

// Client represents a client to call the Scylla Cloud API
type Client struct {
	Meta *Cloudmeta
	// headers holds headers that will be set for all http requests.
	Headers http.Header
	// AccountID holds the account ID used in requests to the API.
	AccountID int64
	// API endpoint
	Endpoint *url.URL
	// HTTPClient is the underlying HTTP client used to run the requests.
	// It may be overloaded but a default one is provided in ``NewClient`` by default.
	HTTPClient *http.Client
}

// NewClient represents a new client to call the API
func (c *Client) Auth(ctx context.Context, token string) error {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: defaultTimeout}
	}

	if c.Headers == nil {
		c.Headers = make(http.Header)
	}

	c.Headers.Set("Authorization", "Bearer "+token)

	if err := c.findAndSaveAccountID(); err != nil {
		return err
	}

	if c.Meta == nil {
		var err error
		if c.Meta, err = BuildCloudmeta(ctx, c); err != nil {
			return fmt.Errorf("error building metadata: %w", err)
		}
	}

	return nil
}

func (c *Client) newHttpRequest(method, path string, reqBody interface{}, query ...string) (*http.Request, error) {
	var body []byte
	var err error

	if reqBody != nil {
		body, err = json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		fmt.Printf("[DEBUG] %s %s body:\n%s\n", method, path, body)
	}

	url := *c.Endpoint
	url.Path = stdpath.Join("/", url.Path, path)

	req, err := http.NewRequest(method, url.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	if body != nil {
		req.Header = req.Header.Clone()
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	if len(query) != 0 {
		if len(query)%2 != 0 {
			return nil, errors.New("odd number of query arguments")
		}

		for i := 0; i < len(query); i += 2 {
			q := req.URL.Query()
			q.Set(query[i], query[i+1])
			req.URL.RawQuery = q.Encode()
		}
	}

	return req, nil
}

func (c *Client) doHttpRequest(req *http.Request) (resp *http.Response, temporaryErr bool, err error) {
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		if oe, ok := err.(*net.OpError); ok {
			temporaryErr = oe.Temporary()
		}

		return
	}

	temporaryErr = resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusGatewayTimeout
	return
}

func (c *Client) doHttpRequestWithRetries(req *http.Request, retries int, retryBackoffDuration time.Duration) (*http.Response, error) {
	resp, temporaryErr, err := c.doHttpRequest(req)
	if temporaryErr && retries > 0 {
		if err == nil {
			_ = resp.Body.Close() // We want to retry anyway.
		}
		return c.doHttpRequestWithRetries(req, retries-1, retryBackoffDuration*2)
	}

	return resp, err
}

func (c *Client) callAPI(ctx context.Context, method, path string, reqBody, resType interface{}, query ...string) error {
	req, err := c.newHttpRequest(method, path, reqBody, query...)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.doHttpRequestWithRetries(req, retriesAllowed, time.Second)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "[DEBUG] (%d) %s %s:\n\n", resp.StatusCode, req.Method, req.URL)
	defer func() {
		fmt.Printf("%s\n\n", &buf)
	}()

	d := json.NewDecoder(io.TeeReader(io.LimitReader(resp.Body, maxResponseBodyLength), &buf))
	d.UseNumber()

	var data = struct {
		Error string      `json:"error"`
		Data  interface{} `json:"data"`
	}{"", resType}

	if err := d.Decode(&data); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if data.Error == "" {
			data.Error = http.StatusText(resp.StatusCode)
		}
	}

	if data.Error != "" {
		return makeError(data.Error, c.Meta.ErrCodes, resp)
	}

	return nil
}

func (c *Client) get(path string, resultType interface{}, query ...string) error {
	return c.callAPI(context.TODO(), http.MethodGet, path, nil, resultType, query...)
}

func (c *Client) post(path string, requestBody, resultType interface{}) error {
	return c.callAPI(context.TODO(), http.MethodPost, path, requestBody, resultType)
}

func (c *Client) delete(path string) error {
	return c.callAPI(context.TODO(), http.MethodDelete, path, nil, nil)
}

func (c *Client) findAndSaveAccountID() error {
	var result struct {
		AccountID int64 `json:"accountId"`
	}

	if err := c.get("/account/default", &result); err != nil {
		return err
	}

	c.AccountID = result.AccountID

	return nil
}
