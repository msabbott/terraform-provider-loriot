package provider

import (
	"context"
	"fmt"

	"bitbucket.org/msabbott/loriot-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &AppTokenDataSource{}
	_ datasource.DataSourceWithConfigure = &AppTokenDataSource{}
)

func NewAppTokenDataSource() datasource.DataSource {
	return &AppTokenDataSource{}
}

// AppTokenDataSource defines the data source implementation.
type AppTokenDataSource struct {
	client *loriot.APIClient
}

// AppTokenDataSourceModel describes the data source data model.
type AppTokenDataSourceModel struct {
	AppId types.String `tfsdk:"app_id"`
	Token types.String `tfsdk:"token"`
}

func (d *AppTokenDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apptoken"
}

func (d *AppTokenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for the token for an application",

		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				MarkdownDescription: "Application ID in hexadecimal format",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Application token",
				Computed:            true,
			},
		},
	}
}

func (d *AppTokenDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*loriot.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *loriot.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *AppTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppTokenDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Fetching App Token with App ID: %s", data.AppId.ValueString()))

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	appToken, _, err := d.client.LoRaApplicationApi.V1NwkAppAPPIDTokenGet(ctx, data.AppId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read App Token, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Token = types.StringValue(appToken[0])

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
