package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure StringTemplateProvider satisfies various provider interfaces.
var _ provider.Provider = &StringTemplateProvider{}
var _ provider.ProviderWithFunctions = &StringTemplateProvider{}

// StringTemplateProvider defines the provider implementation.
type StringTemplateProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// StringTemplateProviderModel describes the provider data model.
type StringTemplateProviderModel struct {
}

func (p *StringTemplateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "string-template"
	resp.Version = p.version
}

func (p *StringTemplateProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *StringTemplateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data StringTemplateProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *StringTemplateProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *StringTemplateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *StringTemplateProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewTemplateFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &StringTemplateProvider{
			version: version,
		}
	}
}
