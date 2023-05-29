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
	_ resource.Resource                = &servicePrincipalUserResource{}
	_ resource.ResourceWithConfigure   = &servicePrincipalUserResource{}
	_ resource.ResourceWithImportState = &servicePrincipalUserResource{}
)

// NewServicePrincipalUserResource is a helper function to simplify the provider implementation.
func NewServicePrincipalUserResource() resource.Resource {
	return &servicePrincipalUserResource{}
}

// servicePrincipalUserResource is the resource implementation.
type servicePrincipalUserResource struct {
	client *iotcentral.Client
}

// servicePrincipalUserResourceModel maps service principal user schema data.
type servicePrincipalUserResourceModel struct {
	ID       types.String                  `tfsdk:"id"`
	ObjectID types.String                  `tfsdk:"object_id"`
	TenantID types.String                  `tfsdk:"tenant_id"`
	Roles    []roleAssignmentResourceModel `tfsdk:"roles"`
}

// Metadata returns the resource type name.
func (r *servicePrincipalUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_principal_user"
}

// Schema defines the schema for the resource.
func (r *servicePrincipalUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique ID of the user.",
				Computed:    true,
			},
			"object_id": schema.StringAttribute{
				Description: "The AAD object ID of the service principal.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The AAD tenant ID of the service principal.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetNestedAttribute{
				Description: "List of role assignments that specify the permissions to access the application.",
				Required:    true,
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
func (r *servicePrincipalUserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*iotcentral.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *servicePrincipalUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan servicePrincipalUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var servicePrincipalUserRequest = iotcentral.ServicePrincipalUserRequest{
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

		servicePrincipalUserRequest.Roles = append(servicePrincipalUserRequest.Roles, roleToAdd)
	}

	// Create new service principal user
	servicePrincipalUser, err := r.client.CreateServicePrincipalUser(servicePrincipalUserRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service principal user",
			"Could not create service principal user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(servicePrincipalUser.ID)
	plan.ObjectID = types.StringValue(servicePrincipalUser.ObjectID)
	plan.TenantID = types.StringValue(servicePrincipalUser.TenantID)

	plan.Roles = []roleAssignmentResourceModel{}
	for _, role := range servicePrincipalUser.Roles {
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
func (r *servicePrincipalUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state servicePrincipalUserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed service principal user value from IotCentral
	servicePrincipalUser, err := r.client.GetServicePrincipalUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading IotCentral User",
			"Could not read IotCentral service principal user ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state.ID = types.StringValue(servicePrincipalUser.ID)
	state.ObjectID = types.StringValue(servicePrincipalUser.ObjectID)
	state.TenantID = types.StringValue(servicePrincipalUser.TenantID)

	state.Roles = []roleAssignmentResourceModel{}
	for _, role := range servicePrincipalUser.Roles {
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
func (r *servicePrincipalUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan servicePrincipalUserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var servicePrincipalUserRequest = iotcentral.ServicePrincipalUserRequest{
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

		servicePrincipalUserRequest.Roles = append(servicePrincipalUserRequest.Roles, roleToAdd)
	}

	var state servicePrincipalUserResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing service principal user
	servicePrincipalUser, err := r.client.UpdateServicePrincipalUser(state.ID.ValueString(), servicePrincipalUserRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating IotCentraL User",
			"Could not update service principal user, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(servicePrincipalUser.ID)
	plan.ObjectID = types.StringValue(servicePrincipalUser.ObjectID)
	plan.TenantID = types.StringValue(servicePrincipalUser.TenantID)

	plan.Roles = []roleAssignmentResourceModel{}
	for _, role := range servicePrincipalUser.Roles {
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
func (r *servicePrincipalUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicePrincipalUserResourceModel
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
			"Could not delete service principal user, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *servicePrincipalUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
