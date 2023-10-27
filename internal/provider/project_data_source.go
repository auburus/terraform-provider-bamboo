package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

type projectDataSource struct {
	bambooClient *BambooClient
}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSourceModel struct {
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var projectConfig projectDataSourceModel

	diags := req.Config.Get(ctx, &projectConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.bambooClient.ReadProject(projectConfig.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Bamboo Project",
			err.Error(),
		)

		return
	}

	projectState := projectDataSourceModel{
		Key:         types.StringValue(project.Key),
		Name:        types.StringValue(project.Name),
		Description: types.StringValue(project.Description),
	}

	diags = resp.State.Set(ctx, &projectState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	bambooClient, ok := req.ProviderData.(*BambooClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *BambooClient, got: %T. please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.bambooClient = bambooClient

}
