// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure bambooProvider satisfies various provider interfaces.
var _ provider.Provider = &bambooProvider{}

// bambooProvider defines the provider implementation.
type bambooProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Url      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *bambooProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bamboo"
	resp.Version = p.version
}

func (p *bambooProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *bambooProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Bamboo API url",
			"The provider cannot create the Bamboo client as there is an unknown configuration value for the Bamboo url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_URL environment variable.",
		)
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Bamboo username",
			"The provider cannot create the Bamboo client as there is an unknown configuration value for the Bamboo username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_USERNAME environment variable.",
		)
	}
	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Bamboo API password",
			"The provider cannot create the Bamboo client as there is an unknown configuration value for the Bamboo password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("BAMBOO_URL")
	username := os.Getenv("BAMBOO_USERNAME")
	password := os.Getenv("BAMBOO_PASSWORD")

	if !data.Url.IsNull() {
		url = data.Url.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Bamboo API url",
			"The provider cannot create the Bamboo client as there is a missing or empty value for the Bamboo url. "+
				"Set the host value in the configuration or use the BAMBOO_URL environment variable.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Bamboo API username",
			"The provider cannot create the Bamboo client as there is a missing or empty value for the Bamboo username. "+
				"Set the host value in the configuration or use the BAMBOO_USERNAME environment variable.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Bamboo API password",
			"The provider cannot create the Bamboo client as there is a missing or empty value for the Bamboo password. "+
				"Set the host value in the configuration or use the BAMBOO_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	url = url + "rest/api/latest/"

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	httpRequest, err := http.NewRequest("GET", url, nil)
	httpRequest.SetBasicAuth(username, password)
	httpRequest.Header.Add("Accept", "application/json")

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Bamboo API client",
			"An unexpected error ocurred when creating the bamboo client. "+
				"Bamboo client error: "+err.Error(),
		)
		return
	}

	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)

	if httpResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unable to create Bamboo API client",
			"Bamboo returned status code "+httpResponse.Status+"when validating the api. Response:\n"+string(body),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *bambooProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *bambooProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCoffeesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bambooProvider{
			version: version,
		}
	}
}
