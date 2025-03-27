package profile

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

var (
	profileIdFormat = regexp.MustCompile(`^profile-\d+$`)
	createdAtFormat = regexp.MustCompile(`^\d+$`)
)

// Schema defines the schema for the resource.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(profileIdFormat, "must be in the format profile-1")},
				Description:   "The unique identifier for the profile. This is a string that starts with 'profile-' followed by a number.",
			},
			"first_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					// This validator ensures that the first_name field is not set if the names field is set.
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("names").AtAnyMapKey()),
				},
				Description: "The first name of the person.",
			},
			"last_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					// This validator ensures that the last_name field is not set if the names field is set.
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("names").AtAnyMapKey()),
				},
				Description: "The last name of the person.",
			},
			"gender": schema.StringAttribute{
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("female", "male")},
				Description: "Profile's gender.",
			},
			"names": schema.MapNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_name": schema.StringAttribute{
							Optional:    true,
							Description: "The first name of the person.",
						},
						"middle_name": schema.StringAttribute{
							Optional:    true,
							Description: "The middle name of the person.",
						},
						"last_name": schema.StringAttribute{
							Optional:    true,
							Description: "The last name of the person.",
						},
					},
				},
				Description: "Nested maps of locales to name fields to values.",
			},
			"unions": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "List of union IDs.",
			},
			"birth": event.Schema(event.SchemaOptions{
				NameComputed:        true,
				DescriptionComputed: true,
				Description:         "Birth event information.",
			}),
			"baptism": event.Schema(event.SchemaOptions{
				NameComputed:        true,
				DescriptionComputed: true,
				Description:         "Baptism event information.",
			}),
			"death": event.Schema(event.SchemaOptions{
				NameComputed:        true,
				DescriptionComputed: true,
				Description:         "Death event information.",
			}),
			"burial": event.Schema(event.SchemaOptions{
				NameComputed:        true,
				DescriptionComputed: true,
				Description:         "Burial event information.",
			}),
			"created_at": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(createdAtFormat, "must be a Unix epoch time in seconds")},
				Description:   "The Unix epoch time in seconds when the profile was created.",
			},
		},
	}
}
