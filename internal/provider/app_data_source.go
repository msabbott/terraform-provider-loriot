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
	_ datasource.DataSource              = &AppDataSource{}
	_ datasource.DataSourceWithConfigure = &AppDataSource{}
)

func NewAppDataSource() datasource.DataSource {
	return &AppDataSource{}
}

// AppDataSource defines the data source implementation.
type AppDataSource struct {
	client *loriot.APIClient
}

// AppDataSourceModel describes the data source data model.
type AppDataSourceModel struct {
	ID                types.String                        `tfsdk:"id"`
	AppId             types.String                        `tfsdk:"app_id"`
	DecimalId         types.Float64                       `tfsdk:"decimal_id"`
	Name              types.String                        `tfsdk:"name"`
	OwnerId           types.Float64                       `tfsdk:"owner_id"`
	OrganizationId    types.Float64                       `tfsdk:"organization_id"`
	Visibility        types.String                        `tfsdk:"visibility"`
	CreatedDate       types.String                        `tfsdk:"created_date"`
	DevicesUsed       types.Float64                       `tfsdk:"devices_used"`
	DevicesLimit      types.Float64                       `tfsdk:"devices_limit"`
	MCastDevicesUsed  types.Float64                       `tfsdk:"mcast_devices_used"`
	MCastDevicesLimit types.Float64                       `tfsdk:"mcast_devices_limit"`
	ConfigDeviceBase  *AppConfigDeviceBaseDataSourceModel `tfsdk:"config_device_base"`
}

type AppConfigDeviceBaseDataSourceModel struct {
	DeviceClass        types.String `tfsdk:"device_class"`
	RxW                types.Int64  `tfsdk:"rxw"`
	DutyCycle          types.Int64  `tfsdk:"duty_cycle"`
	Address            types.Bool   `tfsdk:"address"`
	AddressMin         types.Int64  `tfsdk:"address_min"`
	AddressMax         types.Int64  `tfsdk:"address_max"`
	AddressFix         types.Int64  `tfsdk:"address_fix"`
	SequenceRelax      types.Bool   `tfsdk:"sequence_relax"`
	SequenceDoNotReset types.Bool   `tfsdk:"sequence_do_not_reset"`
	AddressCountLimit  types.Int64  `tfsdk:"address_count_limit"`
}

func (d *AppDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (d *AppDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source for a configured application",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Synonym for app_id",
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "Application ID in hexadecimal format",
				Required:            true,
			},
			"decimal_id": schema.Float64Attribute{
				MarkdownDescription: "Application ID in decimal format",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name",
				Optional:            true,
			},
			"owner_id": schema.Float64Attribute{
				MarkdownDescription: "User ID of the application owner",
				Computed:            true,
			},
			"organization_id": schema.Float64Attribute{
				MarkdownDescription: "Identifier of the organization the application belongs to",
				Computed:            true,
			},
			"visibility": schema.StringAttribute{
				MarkdownDescription: "Visibility of the application",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "Creation date",
				Optional:            true,
			},
			"devices_used": schema.Float64Attribute{
				MarkdownDescription: "Number of devices registered with the application",
				Computed:            true,
			},
			"devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Limit of devices which can be registered",
				Computed:            true,
			},
			"mcast_devices_used": schema.Float64Attribute{
				MarkdownDescription: "Number of multicast devices registered with the application",
				Computed:            true,
			},
			"mcast_devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Limit of multicast devices which can be registered",
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"config_device_base": schema.SingleNestedBlock{
				MarkdownDescription: "Configuration of Devices",
				Attributes: map[string]schema.Attribute{
					"device_class": schema.StringAttribute{
						MarkdownDescription: "Device class",
						Computed:            true,
					},
					"rxw": schema.Int64Attribute{
						MarkdownDescription: "rwx",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *AppDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Fetching App with ID: %s", data.AppId.ValueString()))

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	app, _, err := d.client.LoRaApplicationApi.V1NwkAppAPPIDGet(ctx, data.AppId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read App, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue(app.AppHexId)
	data.AppId = types.StringValue(app.AppHexId)
	data.DecimalId = types.Float64Value(app.Id)
	data.Name = types.StringValue(app.Name)
	data.OwnerId = types.Float64Value(app.Ownerid)
	data.OrganizationId = types.Float64Value(app.OrganizationId)
	data.Visibility = types.StringValue(app.Visibility)
	data.CreatedDate = types.StringValue(app.Created)
	data.DevicesUsed = types.Float64Value(app.Devices)
	data.DevicesLimit = types.Float64Value(app.DeviceLimit)
	data.MCastDevicesUsed = types.Float64Value(app.Mcastdevices)
	data.MCastDevicesLimit = types.Float64Value(app.Mcastdevlimit)

	if app.CfgDevBase != nil {
		var configDevBase AppConfigDeviceBaseDataSourceModel
		configDevBase.DeviceClass = types.StringPointerValue(app.CfgDevBase.Devclass)
		configDevBase.RxW = types.Int64Value(int64(*app.CfgDevBase.Rxw))
		configDevBase.DutyCycle = types.Int64Value(int64(*app.CfgDevBase.Dutycycle))
		configDevBase.Address = types.BoolValue(*app.CfgDevBase.Adr)
		configDevBase.AddressMin = types.Int64Value(int64(*app.CfgDevBase.AdrMin))
		configDevBase.AddressMax = types.Int64Value(int64(*app.CfgDevBase.AdrMax))
		configDevBase.AddressFix = types.Int64Value(int64(*app.CfgDevBase.AdrFix))
		configDevBase.SequenceRelax = types.BoolValue(*app.CfgDevBase.Seqrelax)
		configDevBase.SequenceDoNotReset = types.BoolValue(*app.CfgDevBase.Seqdnreset)
		configDevBase.AddressCountLimit = types.Int64Value(int64(*app.CfgDevBase.AdrCntLimit))
		data.ConfigDeviceBase = &configDevBase
	} else {
		data.ConfigDeviceBase = nil
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
