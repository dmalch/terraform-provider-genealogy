package union

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

const fieldPartners = "partners"
const fieldChildren = "children"

var unionIdFormat = regexp.MustCompile(`^union-(g)?\d+$`)

// Schema defines the schema for the resource
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(unionIdFormat, "must be in the format union-1 or union-g1")},
			},
			fieldPartners: schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			fieldChildren: schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"marriage": event.Schema(event.SchemaOptions{NameComputed: true, DescriptionComputed: true}),
			"divorce":  event.Schema(event.SchemaOptions{NameComputed: true, DescriptionComputed: true}),
		},
	}
}
