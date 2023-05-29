package iotcentral

import "github.com/hashicorp/terraform-plugin-framework/types"

// roleAssignmentResourceModel maps user schema data.
type roleAssignmentResourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Role         types.String `tfsdk:"role"`
}
