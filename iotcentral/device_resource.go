package iotcentral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	iotcentral "github.com/kenspur/azure-iot-central-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deviceResource{}
	_ resource.ResourceWithConfigure = &deviceResource{}
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
	ID          types.String `tfsdk:"id"`
	Etag        types.String `tfsdk:"etag"`
	DisplayName types.String `tfsdk:"display_name"`
	Template    types.String `tfsdk:"template"`
	Simulated   types.Bool   `tfsdk:"simulated"`
	Provisioned types.Bool   `tfsdk:"provisioned"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	//Organizations []types.String `tfsdk:"organizations"`
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
				Required: true,
			},
			"etag": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				Required: true,
			},
			"template": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"simulated": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"provisioned": schema.BoolAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			// "organizations": schema.SetAttribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	ElementType: types.StringType,
			// },
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

	// for _, organization := range plan.Organizations {
	// 	deviceRequest.Organizations = append(deviceRequest.Organizations, organization.ValueString())
	// }

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
	// for _, organization := range device.Organizations {
	// 	plan.Organizations = append(plan.Organizations, types.StringValue(organization))
	// }

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

	// Get refreshed order value from HashiCups
	order, err := r.client.GetDevice(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading IotCentral Device",
			"Could not read IotCentral device ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(order.ID)
	state.Etag = types.StringValue(order.Etag)
	state.DisplayName = types.StringValue(order.DisplayName)
	state.Template = types.StringValue(order.Template)
	state.Simulated = types.BoolValue(order.Simulated)
	state.Provisioned = types.BoolValue(order.Provisioned)
	state.Enabled = types.BoolValue(order.Enabled)
	// state.Organizations = []types.String{}
	// for _, organization := range order.Organizations {
	// 	state.Organizations = append(state.Organizations, types.StringValue(organization))
	// }

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
