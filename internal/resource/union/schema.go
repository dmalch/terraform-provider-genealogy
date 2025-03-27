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

// Schema defines the schema for the resource.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(unionIdFormat, "must be in the format union-1 or union-g1")},
				Description:   "The unique identifier for the union. This is a string that starts with 'union-' followed by a number.",
			},
			fieldPartners: schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of partner IDs.",
			},
			fieldChildren: schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of children IDs.",
			},
			"marriage": event.Schema(event.SchemaOptions{
				NameComputed:        true,
				DescriptionComputed: true,
				Description:         "Marriage event information.",
			}),
			"divorce": event.Schema(event.SchemaOptions{
				NameComputed:        true,
				DescriptionComputed: true,
				Description:         "Divorce event information.",
			}),
		},
	}
}
