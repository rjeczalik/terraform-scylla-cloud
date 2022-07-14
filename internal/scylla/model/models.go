package model

type UserAccount struct {
	UserId            int64  `json:"UserID"`
	AccountId         int64  `json:"AccountID"`
	Name              string `json:"Name"`
	OwnerUserId       int64  `json:"OwnerUserID"`
	AccountStatus     string `json:"AccountStatus"`
	Role              string `json:"Role"`
	UserAccountStatus string `json:"UserAccountStatus"`
}

type CloudProvider struct {
	Id            int64  `json:"ID"`
	Name          string `json:"Name"`
	RootAccountId string `json:"RootAccountID"`
}

type CloudProviderRegion struct {
	Id                          int64  `json:"ID"`
	CloudProviderId             int64  `json:"CloudProviderID"`
	Name                        string `json:"Name"`
	FullName                    string `json:"FullName"`
	ExternalId                  string `json:"ExternalID"`
	MultiRegionExternalId       string `json:"MultiRegionExternalID"`
	DcName                      string `json:"DCName"`
	BackupStorageGbCost         string `json:"BackupStorageGBCost"`
	TrafficSameRegionInGbCost   string `json:"TrafficSameRegionInGBCost"`
	TrafficSameRegionOutGbCost  string `json:"TrafficSameRegionOutGBCost"`
	TrafficCrossRegionOutGbCost string `json:"TrafficCrossRegionOutGBCost"`
	TrafficInternetOutGbCost    string `json:"TrafficInternetOutGBCost"`
	Continent                   string `json:"Continent"`
}

type DataCenter struct {
	Id                               int64  `json:"ID"`
	ClusterId                        int64  `json:"ClusterID"`
	CloudProviderId                  int64  `json:"CloudProviderID"`
	CloudProviderRegionId            int64  `json:"CloudProviderRegionID"`
	ReplicationFactor                int64  `json:"ReplicationFactor"`
	Ipv4Cidr                         string `json:"IPv4CIDR"`
	AccountCloudProviderCredentialId int64  `json:"AccountCloudProviderCredentialID"`
	Status                           string `json:"Status"`
	Name                             string `json:"Name"`
	ManagementNetwork                string `json:"ManagementNetwork"`
	InstanceTypeId                   int64  `json:"InstanceTypeID"`
}

type DataCenterWithClientConnections struct {
	DataCenter
	ClientConnection []string `json:"ClientConnection"`
}

type FreeTier struct {
	ExpirationDate    string `json:"ExpirationDate"`
	ExpirationSeconds int64  `json:"ExpirationSeconds"`
	CreationTime      string `json:"CreationTime"`
}

type Cluster struct {
	Id                        int64                             `json:"ID"`
	Name                      string                            `json:"Name"`
	ClusterNameOnConfigFile   string                            `json:"ClusterNameOnConfigFile"`
	Status                    string                            `json:"Status"`
	CloudProviderId           int64                             `json:"CloudProviderID"`
	ReplicationFactor         int64                             `json:"ReplicationFactor"`
	BroadcastType             string                            `json:"BroadcastType"`
	ScyllaVersionId           int64                             `json:"ScyllaVersionID"`
	ScyllaVersion             string                            `json:"ScyllaVersion"`
	Dc                        []DataCenterWithClientConnections `json:"DC"`
	GrafanaUrl                string                            `json:"GrafanaURL"`
	GrafanaRootUrl            string                            `json:"GrafanaRootURL"`
	BackofficeGrafanaUrl      string                            `json:"BackofficeGrafanaURL"`
	BackofficePrometheusUrl   string                            `json:"BackofficePrometheusURL"`
	BackofficeAlertManagerUrl string                            `json:"BackofficeAlertManagerURL"`
	FreeTier                  FreeTier                          `json:"FreeTier"`
	EncryptionMode            string                            `json:"EncryptionMode"`
	UserApiInterface          string                            `json:"UserAPIInterface"`
	PricingModel              int64                             `json:"PricingModel"`
	MaxAllowedCidrRange       int64                             `json:"MaxAllowedCidrRange"`
	CreatedAt                 string                            `json:"CreatedAt"`
	Dns                       bool                              `json:"DNS"`
	PromProxyEnabled          bool                              `json:"PromProxyEnabled"`
}

type AllowlistRule struct {
	Id            int64  `json:"ID"`
	ClusterId     int64  `json:"ClusterID"`
	SourceAddress string `json:"SourceAddress"`
}

type Node struct {
	Id                          int64  `json:"ID"`
	ClusterId                   int64  `json:"ClusterID"`
	CloudProviderId             int64  `json:"CloudProviderID"`
	CloudProviderInstanceTypeId int64  `json:"CloudProviderInstanceTypeID"`
	CloudProviderRegionId       int64  `json:"CloudProviderRegionID"`
	PublicIP                    string `json:"PublicIP"`
	PrivateIP                   string `json:"PrivateIP"`
	ClusterJoinDate             string `json:"ClusterJoinDate"`
	ServiceId                   int64  `json:"ServiceID"`
	ServiceVersionId            int64  `json:"ServiceVersionID"`
	ServiceVersion              string `json:"ServiceVersion"`
	BillingStartDate            string `json:"BillingStartDate"`
	Hostname                    string `json:"Hostname"`
	ClusterHostId               string `json:"ClusterHostID"`
	Status                      string `json:"Status"`
	NodeState                   string `json:"NodeState"`
	ClusterDcId                 int64  `json:"ClusterDCID"`
	ServerActionId              int64  `json:"ServerActionID"`
	Distribution                string `json:"Distribution"`
	Dns                         string `json:"DNS"`
}
