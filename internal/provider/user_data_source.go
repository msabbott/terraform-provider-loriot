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
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	client *loriot.APIClient
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	UserId            types.Float64 `tfsdk:"id"`
	Email             types.String  `tfsdk:"email"`
	Alerts            types.Bool    `tfsdk:"alerts"`
	DevicesLimit      types.Float64 `tfsdk:"devices_limit"`
	FirstName         types.String  `tfsdk:"first_name"`
	GatewaysLimit     types.Float64 `tfsdk:"gateways_limit"`
	HasCard           types.Bool    `tfsdk:"has_credit_card"`
	LastName          types.String  `tfsdk:"last_name"`
	Level             types.Float64 `tfsdk:"level"`
	MCastDevicesLimit types.Float64 `tfsdk:"mcast_devices_limit"`
	OrganizationRole  types.String  `tfsdk:"organization_role"`
	OrganizationUUID  types.String  `tfsdk:"organization_uuid"`
	OutputLimit       types.Float64 `tfsdk:"output_limit"`
	Tier              types.Float64 `tfsdk:"tier"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source for the current user",

		Attributes: map[string]schema.Attribute{
			"id": schema.Float64Attribute{
				MarkdownDescription: "Internal user ID in Loriot Network Server",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email address",
				Computed:            true,
			},
			"alerts": schema.BoolAttribute{
				MarkdownDescription: "Notifcations alerts configuration",
				Computed:            true,
			},
			"devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Devices limit for this user",
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name or Forename",
				Computed:            true,
			},
			"gateways_limit": schema.Float64Attribute{
				MarkdownDescription: "Gateways limit for this user",
				Computed:            true,
			},
			"has_credit_card": schema.BoolAttribute{
				MarkdownDescription: "User has a credit card in their account",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name or surname",
				Computed:            true,
			},
			"level": schema.Float64Attribute{
				MarkdownDescription: "Level of the user for admin rights (1 to 100)",
				Computed:            true,
			},
			"mcast_devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Multicast devices limit by user",
				Computed:            true,
			},
			"organization_role": schema.StringAttribute{
				MarkdownDescription: "Organization role of the user",
				Computed:            true,
			},
			"organization_uuid": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the organization",
				Computed:            true,
			},
			"output_limit": schema.Float64Attribute{
				MarkdownDescription: "Maximum number of outputs allowed",
				Computed:            true,
			},
			"tier": schema.Float64Attribute{
				MarkdownDescription: "Tier of the user",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	user, _, err := d.client.UserApi.V1NwkUserGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read User, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	//data.UserID = types.NumberValue(user.Userid)
	data.UserId = types.Float64Value(user.Userid)
	data.Email = types.StringValue(user.Email)
	data.Alerts = types.BoolValue(user.Alerts)
	data.DevicesLimit = types.Float64Value(user.Devlimit)
	data.FirstName = types.StringValue(user.FirstName)
	data.GatewaysLimit = types.Float64Value(user.Gwlimit)
	data.HasCard = types.BoolValue(user.Hascard)
	data.LastName = types.StringValue(user.LastName)
	data.Level = types.Float64Value(user.Level)
	data.MCastDevicesLimit = types.Float64Value(user.Mcastdevlimit)
	data.OrganizationRole = types.StringValue(user.OrganizationRole)
	data.OrganizationUUID = types.StringValue(user.OrganizationUuid)
	data.OutputLimit = types.Float64Value(user.OutputLimit)
	data.Tier = types.Float64Value(user.Tier)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
