package iotcentral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	iotcentral "github.com/kenspur/azure-iot-central-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &adGroupUserResource{}
	_ resource.ResourceWithConfigure   = &adGroupUserResource{}
	_ resource.ResourceWithImportState = &adGroupUserResource{}
)

// NewADGroupUserResource is a helper function to simplify the provider implementation.
func NewADGroupUserResource() resource.Resource {
	return &adGroupUserResource{}
}

// adGroupUserResource is the resource implementation.
type adGroupUserResource struct {
	client *iotcentral.Client
}

// adGroupUserResourceModel maps ad group user schema data.
type adGroupUserResourceModel struct {
	ID       types.String                  `tfsdk:"id"`
	ObjectID types.String                  `tfsdk:"object_id"`
	TenantID types.String                  `tfsdk:"tenant_id"`
	Roles    []roleAssignmentResourceModel `tfsdk:"roles"`
}

// Metadata returns the resource type name.
func (r *adGroupUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ad_group_user"
}

// Schema defines the schema for the resource.
func (r *adGroupUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique ID of the user.",
				Computed:    true,
			},
			"object_id": schema.StringAttribute{
				Description: "The AAD object ID of the AD Group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The AAD tenant ID of the AD Group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetNestedAttribute{
				Required:    true,
				Description: "List of role assignments that specify the permissions to access the application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							Description: "ID of the role for this role assignment.",
							Required:    true,
						},
						"organization": schema.StringAttribute{
							Description: "ID of the organization for this role assignment.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *adGroupUserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*iotcentral.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *adGroupUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan adGroupUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var adGroupUserRequest = iotcentral.ADGroupUserRequest{
		ObjectID: plan.ObjectID.ValueString(),
		TenantID: plan.TenantID.ValueString(),
	}

	for _, role := range plan.Roles {
		var roleToAdd = iotcentral.RoleAssignment{
			Role: role.Role.ValueString(),
		}

		if !role.Organization.IsNull() {
			roleToAdd.Organization = role.Organization.ValueString()
		}

		adGroupUserRequest.Roles = append(adGroupUserRequest.Roles, roleToAdd)
	}

	// Create new ad group user
	adGroupUser, err := r.client.CreateADGroupUser(adGroupUserRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ad group user",
			"Could not create ad group user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(adGroupUser.ID)
	plan.ObjectID = types.StringValue(adGroupUser.ObjectID)
	plan.TenantID = types.StringValue(adGroupUser.TenantID)

	plan.Roles = []roleAssignmentResourceModel{}
	for _, role := range adGroupUser.Roles {
		var roleToAdd = roleAssignmentResourceModel{
			Role: types.StringValue(role.Role),
		}

		if role.Organization != "" {
			roleToAdd.Organization = types.StringValue(role.Organization)
		}

		plan.Roles = append(plan.Roles, roleToAdd)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *adGroupUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state adGroupUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed ad group user value from IotCentral
	adGroupUser, err := r.client.GetADGroupUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading IotCentral User",
			"Could not read IotCentral ad group user ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state.ID = types.StringValue(adGroupUser.ID)
	state.ObjectID = types.StringValue(adGroupUser.ObjectID)
	state.TenantID = types.StringValue(adGroupUser.TenantID)

	state.Roles = []roleAssignmentResourceModel{}
	for _, role := range adGroupUser.Roles {
		var roleToAdd = roleAssignmentResourceModel{
			Role: types.StringValue(role.Role),
		}

		if role.Organization != "" {
			roleToAdd.Organization = types.StringValue(role.Organization)
		}

		state.Roles = append(state.Roles, roleToAdd)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *adGroupUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan adGroupUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var adGroupUserRequest = iotcentral.ADGroupUserRequest{
		ObjectID: plan.ObjectID.ValueString(),
		TenantID: plan.TenantID.ValueString(),
	}

	for _, role := range plan.Roles {
		var roleToAdd = iotcentral.RoleAssignment{
			Role: role.Role.ValueString(),
		}

		if !role.Organization.IsNull() {
			roleToAdd.Organization = role.Organization.ValueString()
		}

		adGroupUserRequest.Roles = append(adGroupUserRequest.Roles, roleToAdd)
	}

	var state adGroupUserResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing ad group user
	adGroupUser, err := r.client.UpdateADGroupUser(state.ID.ValueString(), adGroupUserRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IotCentraL User",
			"Could not update ad group user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(adGroupUser.ID)
	plan.ObjectID = types.StringValue(adGroupUser.ObjectID)
	plan.TenantID = types.StringValue(adGroupUser.TenantID)

	plan.Roles = []roleAssignmentResourceModel{}
	for _, role := range adGroupUser.Roles {
		var roleToAdd = roleAssignmentResourceModel{
			Role: types.StringValue(role.Role),
		}

		if role.Organization != "" {
			roleToAdd.Organization = types.StringValue(role.Organization)
		}

		plan.Roles = append(plan.Roles, roleToAdd)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *adGroupUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state adGroupUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting IotCentral User",
			"Could not delete ad group user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *adGroupUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
