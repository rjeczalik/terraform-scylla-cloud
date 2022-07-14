package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.DataSourceType = allowlistDataSourceType{}
var _ tfsdk.DataSource = allowlistDataSource{}

type allowlistDataSourceType struct{}

var allowlistAttrs = markAttrsAsComputed(map[string]tfsdk.Attribute{
	"id": {
		MarkdownDescription: "ID of the rule",
		Type:                types.Int64Type,
	},
	"cluster_id": {
		MarkdownDescription: "ID of the cluster",
		Type:                types.Int64Type,
	},
	"source_address": {
		MarkdownDescription: "Source address of allowed traffic",
		Type:                types.StringType,
	},
})

var allowlistAttrsTypes = extractAttrsTypes(allowlistAttrs)

func (t allowlistDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Cluster firewall rules data source",

		Attributes: map[string]tfsdk.Attribute{
			"cluster_id": {
				MarkdownDescription: "ID of the cluster",
				Required:            true,
				Type:                types.Int64Type,
			},
			"all": {
				MarkdownDescription: "List of all firewall rules",
				Computed:            true,
				Type: types.ListType{
					ElemType: types.ObjectType{AttrTypes: allowlistAttrsTypes},
				},
			},
		},
	}, nil
}

func (t allowlistDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return allowlistDataSource{provider: provider}, diags
}

type allowlistDataSourceData struct {
	ClusterId types.Int64 `tfsdk:"cluster_id"`
	All       types.List  `tfsdk:"all"`
}

type allowlistDataSource struct {
	provider provider
}

func (d allowlistDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data allowlistDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := d.provider.client.ListAllowlistRules(data.ClusterId.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list cloud provider regions, got error: %s", err))
		return
	}

	wrappedRules := make([]attr.Value, 0, len(rules))
	for _, rule := range rules {
		wrappedRules = append(wrappedRules, types.Object{
			Attrs: map[string]attr.Value{
				"id":             types.Int64{Value: rule.Id},
				"cluster_id":     types.Int64{Value: rule.ClusterId},
				"source_address": types.String{Value: rule.SourceAddress},
			},
			AttrTypes: allowlistAttrsTypes,
		})
	}

	data.All = types.List{Elems: wrappedRules, ElemType: types.ObjectType{AttrTypes: allowlistAttrsTypes}}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
