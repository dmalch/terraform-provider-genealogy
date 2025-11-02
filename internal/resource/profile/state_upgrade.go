package profile

import (
	"context"

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

func (r *Resource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
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
								"nicknames": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "The nicknames of the person.",
								},
							},
						},
						Description: "Nested maps of locales to name fields to values.",
					},
					"unions": schema.SetAttribute{
						ElementType:   types.StringType,
						Computed:      true,
						Optional:      true,
						PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
						Description:   "List of union IDs.",
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
					"current_residence": event.LocationSchema("Event's location."),
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
					"created_at": schema.StringAttribute{
						Computed:      true,
						Optional:      true,
						PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						Validators:    []validator.String{stringvalidator.RegexMatches(createdAtFormat, "must be a Unix epoch time in seconds")},
						Description:   "The Unix epoch time in seconds when the profile was created.",
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData ResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
				if resp.Diagnostics.HasError() {
					return
				}

				upgradedStateData := ResourceModel{
					ID:               priorStateData.ID,
					Names:            priorStateData.Names,
					Gender:           priorStateData.Gender,
					Birth:            priorStateData.Birth,
					Baptism:          priorStateData.Baptism,
					Death:            priorStateData.Death,
					Burial:           priorStateData.Burial,
					CauseOfDeath:     priorStateData.CauseOfDeath,
					Unions:           priorStateData.Unions,
					Projects:         priorStateData.Projects,
					CurrentResidence: priorStateData.CurrentResidence,
					Public:           priorStateData.Public,
					Alive:            priorStateData.Alive,
					Deleted:          priorStateData.Deleted,
					MergedInto:       priorStateData.MergedInto,
					CreatedAt:        priorStateData.CreatedAt,
				}

				if priorStateData.About.ValueString() != "" {
					upgradedStateData.About, resp.Diagnostics = types.MapValueFrom(ctx, types.StringType, map[string]string{
						"en-US": priorStateData.About.ValueString(),
					})
					if resp.Diagnostics.HasError() {
						return
					}
				} else {
					upgradedStateData.About = types.MapNull(types.StringType)
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}
