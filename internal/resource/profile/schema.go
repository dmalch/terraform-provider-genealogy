package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Schema defines the schema for the resource
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"first_name": schema.StringAttribute{
				Optional: true,
			},
			"last_name": schema.StringAttribute{
				Optional: true,
			},
			"gender": schema.StringAttribute{
				Optional: true,
			},
			"father_id": schema.StringAttribute{
				Optional: true,
			},
			"mother_id": schema.StringAttribute{
				Optional: true,
			},
			"individual_id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"events": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:      true,
							Optional:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"date": schema.StringAttribute{
							Optional: true,
						},
						"additional_content": schema.StringAttribute{
							Optional: true,
						},
						"formatted_place": schema.StringAttribute{
							Optional: true,
						},
						"title": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"cause_of_death": schema.StringAttribute{
							Optional: true,
						},
						"spouse_id": schema.StringAttribute{
							Optional: true,
						},
						"content": schema.StringAttribute{
							Optional: true,
						},
						"notes": notesSchema(),
						"media": mediaSchema(),
					},
				},
			},
			"notes": notesSchema(),
		},
	}
}

func notesSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed:      true,
					Optional:      true,
					PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				},
				"text": schema.StringAttribute{
					Required: true,
				},
			},
		},
	}
}

func mediaSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Computed:      true,
					Optional:      true,
					PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				},
				"link": schema.StringAttribute{
					Computed:      true,
					Optional:      true,
					PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				},
			},
		},
	}
}
