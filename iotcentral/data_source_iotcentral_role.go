package iotcentral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	iotcentral "github.com/kenspur/azure-iot-central-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

// NewRoleDataSource is a helper function to simplify the provider implementation.
func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

// roleDataSource is the data source implementation.
type roleDataSource struct {
	client *iotcentral.Client
}

// roleDataSourceModel maps the data source schema data.
type roleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.
func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the data source.
func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique ID of the role.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Display name of the role.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*iotcentral.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state roleDataSourceModel

	// Read the config
	var cfg roleDataSourceModel
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := d.client.GetRoleByName(cfg.DisplayName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read IotCentral Role",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(role.ID)
	state.DisplayName = types.StringValue(role.DisplayName)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
