package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-geni/internal/resource/event"
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
			"unions": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"birth":   event.Schema(event.SchemaOptions{NameComputed: true, DescriptionComputed: true}),
			"baptism": event.Schema(event.SchemaOptions{NameComputed: true, DescriptionComputed: true}),
			"death":   event.Schema(event.SchemaOptions{NameComputed: true, DescriptionComputed: true}),
			"burial":  event.Schema(event.SchemaOptions{NameComputed: true, DescriptionComputed: true}),
			"created_at": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}
