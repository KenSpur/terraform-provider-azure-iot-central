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
	_ datasource.DataSource              = &devicesDataSource{}
	_ datasource.DataSourceWithConfigure = &devicesDataSource{}
)

// NewDevicesDataSource is a helper function to simplify the provider implementation.
func NewDevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

// devicesDataSource is the data source implementation.
type devicesDataSource struct {
	client *iotcentral.Client
}

// devicesDataSourceModel maps the data source schema data.
type devicesDataSourceModel struct {
	Devices []deviceModel `tfsdk:"devices"`
}

// deviceModel maps device schema data.
type deviceModel struct {
	ID            types.String   `tfsdk:"id"`
	Etag          types.String   `tfsdk:"etag"`
	DisplayName   types.String   `tfsdk:"display_name"`
	Template      types.String   `tfsdk:"template"`
	Simulated     types.Bool     `tfsdk:"simulated"`
	Provisioned   types.Bool     `tfsdk:"provisioned"`
	Enabled       types.Bool     `tfsdk:"enabled"`
	Organizations []types.String `tfsdk:"organizations"`
}

// Metadata returns the data source type name.
func (d *devicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

// Schema defines the schema for the data source.
func (d *devicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"etag": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"template": schema.StringAttribute{
							Computed: true,
						},
						"simulated": schema.BoolAttribute{
							Computed: true,
						},
						"provisioned": schema.BoolAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"organizations": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *devicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*iotcentral.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state devicesDataSourceModel

	devices, err := d.client.GetDevices()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read IotCentral Devices",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, device := range devices {
		deviceState := deviceModel{
			ID:          types.StringValue(device.ID),
			Etag:        types.StringValue(device.Etag),
			DisplayName: types.StringValue(device.DisplayName),
			Simulated:   types.BoolValue(device.Simulated),
			Provisioned: types.BoolValue(device.Provisioned),
			Enabled:     types.BoolValue(device.Enabled),
		}

		for _, organization := range device.Organizations {
			deviceState.Organizations = append(deviceState.Organizations, types.StringValue(organization))
		}

		state.Devices = append(state.Devices, deviceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
