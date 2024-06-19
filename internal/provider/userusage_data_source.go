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
	_ datasource.DataSource              = &UserUsageDataSource{}
	_ datasource.DataSourceWithConfigure = &UserUsageDataSource{}
)

func NewUserUsageDataSource() datasource.DataSource {
	return &UserUsageDataSource{}
}

// UserUsageDataSource defines the data source implementation.
type UserUsageDataSource struct {
	client *loriot.APIClient
}

// UserUsageDataSourceModel describes the data source data model.
type UserUsageDataSourceModel struct {
	Apps              types.Float64 `tfsdk:"apps"`
	SignedDevices     types.Float64 `tfsdk:"signed_devices"`
	DevicesLimit      types.Float64 `tfsdk:"devices_limit"`
	DevicesUsed       types.Float64 `tfsdk:"devices_used"`
	GatewaysUsed      types.Float64 `tfsdk:"gateways_used"`
	GatewaysLimit     types.Float64 `tfsdk:"gateways_limit"`
	MCastDevicesLimit types.Float64 `tfsdk:"mcast_devices_limit"`
	MCastDevicesUsed  types.Float64 `tfsdk:"mcast_devices_used"`
}

func (d *UserUsageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_userusage"
}

func (d *UserUsageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User Usage data source for the current user",

		Attributes: map[string]schema.Attribute{
			"apps": schema.Float64Attribute{
				MarkdownDescription: "Number of apps in use",
				Computed:            true,
			},
			"signed_devices": schema.Float64Attribute{
				MarkdownDescription: "Number of signed devices",
				Computed:            true,
			},
			"devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Maximum number of devices",
				Computed:            true,
			},
			"devices_used": schema.Float64Attribute{
				MarkdownDescription: "Number of devices in use",
				Computed:            true,
			},
			"gateways_used": schema.Float64Attribute{
				MarkdownDescription: "Gateways in use",
				Computed:            true,
			},
			"gateways_limit": schema.Float64Attribute{
				MarkdownDescription: "Limit of gateways in account",
				Computed:            true,
			},
			"mcast_devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Limit of multicast devices in account",
				Computed:            true,
			},
			"mcast_devices_used": schema.Float64Attribute{
				MarkdownDescription: "Multicast devices in use",
				Computed:            true,
			},
		},
	}
}

func (d *UserUsageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserUsageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserUsageDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	userusage, _, err := d.client.UserApi.V1NwkUserUsageGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read User, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	data.Apps = types.Float64Value(userusage.Apps)
	data.SignedDevices = types.Float64Value(userusage.Devices)
	data.DevicesLimit = types.Float64Value(userusage.Devlimit)
	data.DevicesUsed = types.Float64Value(userusage.Devuse)
	data.GatewaysUsed = types.Float64Value(userusage.Gateways)
	data.GatewaysLimit = types.Float64Value(userusage.Gwlimit)
	data.MCastDevicesLimit = types.Float64Value(userusage.Mcastdevices)
	data.MCastDevicesUsed = types.Float64Value(userusage.Mcastdevuse)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
