package profile

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

var (
	profileIdFormat = regexp.MustCompile(`^profile-(g)?\d+$`)
	createdAtFormat = regexp.MustCompile(`^\d+$`)
)

// Schema defines the schema for the resource
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(profileIdFormat, "must be in the format profile-1 or profile-g1")},
			},
			"first_name": schema.StringAttribute{
				Optional: true,
			},
			"last_name": schema.StringAttribute{
				Optional: true,
			},
			"gender": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{stringvalidator.OneOf("female", "male")},
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
				Validators:    []validator.String{stringvalidator.RegexMatches(createdAtFormat, "must be a Unix epoch time in seconds")},
			},
		},
	}
}
