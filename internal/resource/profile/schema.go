package profile

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/geniplanmodifier"
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
						"birth_last_name": schema.StringAttribute{
							Optional:    true,
							Description: "The birth last name of the person.",
						},
						"display_name": schema.StringAttribute{
							Optional:    true,
							Description: "The display name of the person.",
						},
					},
				},
				Description: "Nested maps of locales to name fields to values.",
			},
			"unions": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "List of union IDs.",
			},
			"projects": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{setplanmodifier.RequiresReplaceIf(
					geniplanmodifier.ValuesAreRemovedFromState,
					"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
					"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
				)},
				Description: "List of project IDs.",
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
			"cause_of_death": schema.StringAttribute{
				Optional:    true,
				Description: "Profile's death cause",
			},
			"about": schema.StringAttribute{
				Optional:    true,
				Description: "Profile's about me section.",
			},
			"public": schema.BoolAttribute{
				Required:    true,
				Description: "Profile's public visibility.",
			},
			"alive": schema.BoolAttribute{
				Required:    true,
				Description: "Profile's alive status.",
			},
			"deleted": schema.BoolAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
				Description:   "Profile's deleted status.",
			},
			"merged_into": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "The ID of the profile this profile was merged into.",
			},
			"auto_update_when_merged": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to automatically update the profile when it is merged with another profile.",
			},
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
