package iotcentral

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

// Metadata returns the provider type name.
func (p *iotcentralProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "iotcentral"
}

// Schema defines the provider-level schema for configuration data.
func (p *iotcentralProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure prepares a iotcentral API client for data sources and resources.
func (p *iotcentralProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *iotcentralProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *iotcentralProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
