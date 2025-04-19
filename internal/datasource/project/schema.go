package project

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	projectIdFormat = regexp.MustCompile(`^project-\d+$`)
)

func (d *DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(projectIdFormat, "must be in the format project-1"),
				},
				Description: "The unique identifier for the project. This is a string that starts with 'project-' followed by a number.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the project.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the project.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the project was last updated.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the project was created.",
			},
		},
	}
}
