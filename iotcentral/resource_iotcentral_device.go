package iotcentral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	iotcentral "github.com/kenspur/azure-iot-central-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

// NewDeviceResource is a helper function to simplify the provider implementation.
func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceResource is the resource implementation.
type deviceResource struct {
	client *iotcentral.Client
}

// deviceResourceModel maps device schema data.
type deviceResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Etag          types.String `tfsdk:"etag"`
	DisplayName   types.String `tfsdk:"display_name"`
	Template      types.String `tfsdk:"template"`
	Simulated     types.Bool   `tfsdk:"simulated"`
	Provisioned   types.Bool   `tfsdk:"provisioned"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Organizations types.List   `tfsdk:"organizations"`
}

// Metadata returns the resource type name.
func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the resource.
func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique ID of the device.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"etag": schema.StringAttribute{
				Description: "ETag used to prevent conflict in device updates.",
				Optional:    true,
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Display name of the device.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"template": schema.StringAttribute{
				Description: "The device template definition for the device.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"simulated": schema.BoolAttribute{
				Description: "Whether the device is simulated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"provisioned": schema.BoolAttribute{
				Description: "Whether resources have been allocated for the device.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "List of organization IDs that the device is a part of, only one organization is supported today, multiple organizations will be supported soon.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"organizations": schema.ListAttribute{
				Description: "List of organization IDs that the device is a part of, only one organization is supported today, multiple organizations will be supported soon.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*iotcentral.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var deviceID = plan.ID.ValueString()
	var deviceRequest = iotcentral.DeviceRequest{
		DisplayName: plan.DisplayName.ValueString(),
		Template:    plan.Template.ValueString(),
		Simulated:   plan.Simulated.ValueBool(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	for _, organization := range plan.Organizations.Elements() {
		deviceRequest.Organizations = append(deviceRequest.Organizations, organization.(types.String).ValueString())
	}

	// Create new device
	device, err := r.client.CreateDevice(deviceID, deviceRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device",
			"Could not create device, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(device.ID)
	plan.Etag = types.StringValue(device.Etag)
	plan.DisplayName = types.StringValue(device.DisplayName)
	plan.Template = types.StringValue(device.Template)
	plan.Simulated = types.BoolValue(device.Simulated)
	plan.Provisioned = types.BoolValue(device.Provisioned)
	plan.Enabled = types.BoolValue(device.Enabled)

	plan.Organizations, diags = types.ListValueFrom(ctx, types.StringType, device.Organizations)
	resp.Diagnostics.Append(diags...)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed device value from IotCentral
	device, err := r.client.GetDevice(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading IotCentral Device",
			"Could not read IotCentral device ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(device.ID)
	state.Etag = types.StringValue(device.Etag)
	state.DisplayName = types.StringValue(device.DisplayName)
	state.Template = types.StringValue(device.Template)
	state.Simulated = types.BoolValue(device.Simulated)
	state.Provisioned = types.BoolValue(device.Provisioned)
	state.Enabled = types.BoolValue(device.Enabled)

	state.Organizations, diags = types.ListValueFrom(ctx, types.StringType, device.Organizations)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var deviceID = plan.ID.ValueString()
	var deviceRequest = iotcentral.DeviceRequest{
		DisplayName: plan.DisplayName.ValueString(),
		Template:    plan.Template.ValueString(),
		Simulated:   plan.Simulated.ValueBool(),
		Enabled:     plan.Enabled.ValueBool(),
	}

	for _, organization := range plan.Organizations.Elements() {
		deviceRequest.Organizations = append(deviceRequest.Organizations, organization.(types.String).ValueString())
	}

	// Update existing device
	device, err := r.client.UpdateDevice(deviceID, deviceRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IotCentraL Device",
			"Could not update device, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(device.ID)
	plan.Etag = types.StringValue(device.Etag)
	plan.DisplayName = types.StringValue(device.DisplayName)
	plan.Template = types.StringValue(device.Template)
	plan.Simulated = types.BoolValue(device.Simulated)
	plan.Provisioned = types.BoolValue(device.Provisioned)
	plan.Enabled = types.BoolValue(device.Enabled)

	plan.Organizations, diags = types.ListValueFrom(ctx, types.StringType, device.Organizations)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteDevice(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting IotCentral Device",
			"Could not delete device, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
