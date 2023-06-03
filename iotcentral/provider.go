package iotcentral

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	iotcentral "github.com/kenspur/azure-iot-central-client-go"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &iotcentralProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &iotcentralProvider{}
}

// iotcentralProvider is the provider implementation.
type iotcentralProvider struct{}

// iotcentralProviderModel maps provider schema data to a Go type.
type iotcentralProviderModel struct {
	Host types.String `tfsdk:"host"`
}

// Metadata returns the provider type name.
func (p *iotcentralProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "iotcentral"
}

// Schema defines the provider-level schema for configuration data.
func (p *iotcentralProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "IoT Central Application URL. May also be provided via IOTCENTRAL_HOST environment variable.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a iotcentral API client for data sources and resources.
func (p *iotcentralProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring IotCentral client")

	// Retrieve provider data from configuration
	var config iotcentralProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown IotCentral API Host",
			"The provider cannot create the IotCentral API client as there is an unknown configuration value for the IotCentral API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the IOTCENTRAL_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("IOTCENTRAL_HOST")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing IotCentral API Host",
			"The provider cannot create the IotCentral API client as there is a missing or empty value for the IotCentral API host. "+
				"Set the host value in the configuration or use the IOTCENTRAL_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "iotcentral_host", host)

	// Create default Azure credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create default Azure credential",
			"An unexpected error occurred when creating default Azure credential. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Azure identity Client Error: "+err.Error(),
		)
		return
	}

	// Prepare the token request options
	policy := policy.TokenRequestOptions{
		Scopes: []string{"https://apps.azureiotcentral.com/.default"},
	}

	// Get the access token
	token, err := cred.GetToken(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Azure access token",
			"An unexpected error occurred when getting the Azure access token. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Azure identity Client Error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Creating IotCentral client")

	// Create a new IotCentral client using the configuration values
	client, err := iotcentral.NewClient(&host, &token.Token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create IotCentral API Client",
			"An unexpected error occurred when creating the IotCentral API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"IotCentral Client Error: "+err.Error(),
		)
		return
	}

	// Make the IotCentral client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured IotCentral client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *iotcentralProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRoleDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *iotcentralProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewUserResource,
		NewADGroupUserResource,
		NewServicePrincipalUserResource,
	}
}
