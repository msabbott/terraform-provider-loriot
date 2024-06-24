// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"math"

	"bitbucket.org/msabbott/loriot-go-client"
	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AppResource{}
var _ resource.ResourceWithImportState = &AppResource{}

func NewAppResource() resource.Resource {
	return &AppResource{}
}

// AppResource defines the resource implementation.
type AppResource struct {
	client *loriot.APIClient
}

// AppResourceModel describes the resource data model.
type AppResourceModel struct {
	AppId          types.String  `tfsdk:"app_id"`
	DecimalId      types.Float64 `tfsdk:"decimal_id"`
	Name           types.String  `tfsdk:"name"`
	OwnerId        types.Float64 `tfsdk:"owner_id"`
	OrganizationId types.Float64 `tfsdk:"organization_id"`
	//visibility
	CreatedDate       types.String  `tfsdk:"created_date"`
	DevicesUsed       types.Float64 `tfsdk:"devices_used"`
	DevicesLimit      types.Float64 `tfsdk:"devices_limit"`
	MCastDevicesUsed  types.Float64 `tfsdk:"mcast_devices_used"`
	MCastDevicesLimit types.Float64 `tfsdk:"mcast_devices_limit"`
}

func (r *AppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *AppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "App resource",

		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				MarkdownDescription: "Application ID in hexadecimal format",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"decimal_id": schema.Float64Attribute{
				MarkdownDescription: "Application ID in decimal format",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name",
				Required:            false,
				Optional:            true,
				Computed:            false,
			},
			"owner_id": schema.Float64Attribute{
				MarkdownDescription: "User ID of the application owner",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"organization_id": schema.Float64Attribute{
				MarkdownDescription: "Identifier of the organization the application belongs to",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "Creation date",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"devices_used": schema.Float64Attribute{
				MarkdownDescription: "Number of devices registered with the application",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Limit of devices which can be registered",
				Required:            true,
				Optional:            false,
				Computed:            false,
			},
			"mcast_devices_used": schema.Float64Attribute{
				MarkdownDescription: "Number of multicast devices registered with the application",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"mcast_devices_limit": schema.Float64Attribute{
				MarkdownDescription: "Limit of multicate devices which can be registered",
				Required:            true,
				Optional:            false,
				Computed:            false,
			},
		},
	}
}

func (r *AppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*loriot.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *loriot.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AppResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := loriot.NwkAppsBody{
		Title: data.Name.ValueString(),
		// Value must be a multiple of 10. In this case, round up.
		Capacity:      math.Round(data.DevicesLimit.ValueFloat64()/10) * 10,
		Visibility:    "private",
		Mcastdevlimit: data.MCastDevicesLimit.ValueFloat64(),
	}

	opts := loriot.LoRaApplicationApi1NwkAppsPostOpts{
		Body: optional.NewInterface(body),
	}

	app, _, err := r.client.LoRaApplicationApi.V1NwkAppsPost(ctx, &opts)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create App, got error: %s", err))
		return
	}

	data.AppId = types.StringValue(app.AppHexId)
	data.OrganizationId = types.Float64Value(app.OrganizationId)
	data.OwnerId = types.Float64Value(app.Ownerid)
	data.DecimalId = types.Float64Value(app.Id)
	data.CreatedDate = types.StringValue(app.Created)
	data.DevicesLimit = types.Float64Value(app.DeviceLimit)
	data.DevicesUsed = types.Float64Value(app.Devices)
	data.MCastDevicesLimit = types.Float64Value(app.Mcastdevlimit)
	data.MCastDevicesUsed = types.Float64Value(app.Mcastdevices)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AppResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Fetching App with ID %s", data.AppId.ValueString()))

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.

	app, _, err := r.client.LoRaApplicationApi.V1NwkAppAPPIDGet(ctx, data.AppId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read App, got error: %s", err))
		return
	}

	data.AppId = types.StringValue(app.AppHexId)
	data.DecimalId = types.Float64Value(app.Id)
	data.Name = types.StringValue(app.Name)
	data.OwnerId = types.Float64Value(app.Ownerid)
	data.OrganizationId = types.Float64Value(app.OrganizationId)
	//visibility
	data.CreatedDate = types.StringValue(app.Created)
	data.DevicesUsed = types.Float64Value(app.Devices)
	data.DevicesLimit = types.Float64Value(app.DeviceLimit)
	data.MCastDevicesUsed = types.Float64Value(app.Mcastdevices)
	data.MCastDevicesLimit = types.Float64Value(app.Mcastdevlimit)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data AppResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Changes to the title are updated through the title API
	if !data.Name.Equal(state.Name) {

		titleBody := loriot.AppidTitleBody{
			Title: data.Name.ValueString(),
		}

		_, _, nameErr := r.client.LoRaApplicationApi.V1NwkAppAPPIDTitlePost(ctx, titleBody, state.AppId.ValueString())

		if nameErr != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update name of App, got error: %s", nameErr))
			return
		}
	}

	if !data.DevicesLimit.Equal(state.DevicesLimit) {

		capacityBody := loriot.AppidCapacityBody{
			Inc: 0,
			Dec: 0,
		}

		capacityBody.Inc = data.DevicesLimit.ValueFloat64() - state.DecimalId.ValueFloat64()
		capacityBody.Dec = state.DevicesLimit.ValueFloat64() - data.DecimalId.ValueFloat64()

		if capacityBody.Inc < 0 {
			capacityBody.Inc = 0
		}

		if capacityBody.Dec < 0 {
			capacityBody.Dec = 0
		}

		_, capacityErr := r.client.LoRaApplicationApi.V1NwkAppAPPIDCapacityPost(ctx, capacityBody, state.AppId.ValueString())

		if capacityErr != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update capacity of App, got error: %s", capacityErr))
			return
		}
	}

	// Re-read the application to ensure the most up-to-date version is returned
	app, _, err := r.client.LoRaApplicationApi.V1NwkAppAPPIDGet(ctx, state.AppId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read App, got error: %s", err))
		return
	}

	data.AppId = types.StringValue(app.AppHexId)
	data.DecimalId = types.Float64Value(app.Id)
	data.Name = types.StringValue(app.Name)
	data.OwnerId = types.Float64Value(app.Ownerid)
	data.OrganizationId = types.Float64Value(app.OrganizationId)
	//visibility
	data.CreatedDate = types.StringValue(app.Created)
	data.DevicesUsed = types.Float64Value(app.Devices)
	data.DevicesLimit = types.Float64Value(app.DeviceLimit)
	data.MCastDevicesUsed = types.Float64Value(app.Mcastdevices)
	data.MCastDevicesLimit = types.Float64Value(app.Mcastdevlimit)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AppResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.LoRaApplicationApi.V1NwkAppAPPIDDelete(ctx, data.AppId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete App, got error: %s", err))
		return
	}
}

func (r *AppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_id"), req, resp)
}
