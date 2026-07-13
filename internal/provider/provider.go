// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure UnraidProvider satisfies various provider interfaces.
var _ provider.Provider = &UnraidProvider{}
var _ provider.ProviderWithFunctions = &UnraidProvider{}
var _ provider.ProviderWithEphemeralResources = &UnraidProvider{}
var _ provider.ProviderWithActions = &UnraidProvider{}

// UnraidProvider defines the provider implementation.
type UnraidProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// UnraidProviderModel describes the provider data model.
type UnraidProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiToken types.String `tfsdk:"api_token"`
}

func (p *UnraidProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "unraid"
	resp.Version = p.version
}

func (p *UnraidProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint of the Unraid server. Can also be set via the `UNRAID_ENDPOINT` environment variable.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token for authenticating with the Unraid server. Can also be set via the `UNRAID_API_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *UnraidProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data UnraidProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fall back to environment variables if not set in configuration.
	endpoint := os.Getenv("UNRAID_ENDPOINT")
	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	apiToken := os.Getenv("UNRAID_API_TOKEN")
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing Unraid Endpoint",
			"The provider requires an endpoint to be set via the `endpoint` attribute or the `UNRAID_ENDPOINT` environment variable.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing Unraid API Token",
			"The provider requires an API token to be set via the `api_token` attribute or the `UNRAID_API_TOKEN` environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	_ = endpoint
	_ = apiToken

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *UnraidProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *UnraidProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *UnraidProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *UnraidProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *UnraidProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UnraidProvider{
			version: version,
		}
	}
}
