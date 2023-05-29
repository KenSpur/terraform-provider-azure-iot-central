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
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *iotcentral.Client
}

// userResourceModel maps user schema data.
type userResourceModel struct {
	ID    types.String        `tfsdk:"id"`
	Email types.String        `tfsdk:"email"`
	Roles []roleResourceModel `tfsdk:"roles"`
}

type roleResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Role         types.String `tfsdk:"role"`
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique ID of the user.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "Email address of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetNestedAttribute{
				Required: true,
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
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*iotcentral.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var userRequest = iotcentral.UserRequest{
		Email: plan.Email.ValueString(),
	}

	for _, role := range plan.Roles {
		var roleToAdd = iotcentral.RoleAssignment{
			Role: role.Role.ValueString(),
		}

		if !role.Organization.IsNull() {
			roleToAdd.Organization = role.Organization.ValueString()
		}

		userRequest.Roles = append(userRequest.Roles, roleToAdd)
	}

	// Create new user
	user, err := r.client.CreateUser(userRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(user.ID)
	plan.Email = types.StringValue(user.Email)

	plan.Roles = []roleResourceModel{}
	for _, role := range user.Roles {
		var roleToAdd = roleResourceModel{
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
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed user value from IotCentral
	user, err := r.client.GetUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading IotCentral User",
			"Could not read IotCentral user ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state.ID = types.StringValue(user.ID)
	state.Email = types.StringValue(user.Email)

	state.Roles = []roleResourceModel{}
	for _, role := range user.Roles {
		var roleToAdd = roleResourceModel{
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
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var userRequest = iotcentral.UserRequest{
		Email: plan.Email.ValueString(),
	}

	for _, role := range plan.Roles {
		var roleToAdd = iotcentral.RoleAssignment{
			Role: role.Role.ValueString(),
		}

		if !role.Organization.IsNull() {
			roleToAdd.Organization = role.Organization.ValueString()
		}

		userRequest.Roles = append(userRequest.Roles, roleToAdd)
	}

	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing user
	user, err := r.client.UpdateUser(state.ID.ValueString(), userRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IotCentraL User",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(user.ID)
	plan.Email = types.StringValue(user.Email)

	plan.Roles = []roleResourceModel{}
	for _, role := range user.Roles {
		var roleToAdd = roleResourceModel{
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
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userResourceModel
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
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
