// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"bitbucket.org/msabbott/loriot-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure LoriotProvider satisfies various provider interfaces.
var _ provider.Provider = &LoriotProvider{}
var _ provider.ProviderWithFunctions = &LoriotProvider{}

// LoriotProvider defines the provider implementation.
type LoriotProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LoriotProviderModel describes the provider data model.
type LoriotProviderModel struct {
	Host   types.String `tfsdk:"host"`
	APIKey types.String `tfsdk:"key"`
}

func (p *LoriotProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "loriot"
	resp.Version = p.version
}

func (p *LoriotProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Hostname of the Loriot instance",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "API Key used to authenticate with the instance",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *LoriotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LoriotProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }
	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Loriot instance Host",
			"The provider cannot create the Loriot API client as there is an unknown configuration value for the Loriot instance. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LORIOT_HOST environment variable.",
		)
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown Loriot API Key",
			"The provider cannot create the Loriot API client as there is an unknown configuration value for the Loriot API key "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LORIOT_API_KEY environment variable.",
		)
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("LORIOT_HOST")
	key := os.Getenv("LORIOT_API_KEY")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.APIKey.IsNull() {
		key = data.APIKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Loriot instance Host",
			"The provider cannot create the Loriot API client as there is a missing or empty value for the Loriot instance host. "+
				"Set the host value in the configuration or use the LORIOT_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Missing Loriot API key",
			"The provider cannot create the Loriot API client as there is a missing or empty value for the Loriot API key. "+
				"Set the apikey value in the configuration or use the LORIOT_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	cfg := loriot.NewConfiguration()

	cfg.BasePath = host
	cfg.AddDefaultHeader("Authorization", "Bearer "+key)

	client := loriot.NewAPIClient(cfg)

	// Example client configuration for data sources and resources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LoriotProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
		NewAppResource,
	}
}

func (p *LoriotProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		NewUserDataSource,
		NewUserUsageDataSource,
		NewAppDataSource,
		NewAppTokenDataSource,
	}
}

func (p *LoriotProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LoriotProvider{
			version: version,
		}
	}
}
