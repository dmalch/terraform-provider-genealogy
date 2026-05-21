package photo

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var photoIdFormat = regexp.MustCompile(`^photo-\d+$`)

// Schema defines the schema for the resource.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(photoIdFormat, "must be in the format photo-1")},
				Description:   "The unique identifier for the photo. A string that starts with 'photo-' followed by a number.",
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The photo's title.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The photo's description.",
			},
			"date": schema.StringAttribute{
				Optional:    true,
				Description: "The photo's date, as a free-form string (Geni does not impose a fixed format).",
			},
			"file": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:   "The photo's image content, as a base64-encoded string. Changing it replaces the photo.",
			},
			"file_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:   "The photo's file name. Changing it replaces the photo.",
			},
			"album": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "The id of the album that holds the photo. Set only at creation; changing it replaces the photo.",
			},
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "The set of profile IDs tagged in the photo.",
			},
			"guid": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "The photo's legacy global identifier.",
			},
			"content_type": schema.StringAttribute{
				Computed:    true,
				Description: "The photo's content type, as detected by Geni from the upload.",
			},
			"attribution": schema.StringAttribute{
				Computed:    true,
				Description: "The photo's attribution string.",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The Geni API URL for the photo.",
			},
			"sizes": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: `Image URLs keyed by Geni size name (e.g. "small", "medium", "large").`,
			},
			"created_at": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Description:   "When the photo was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "When the photo was last updated.",
			},
		},
	}
}

func (r *Resource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The unique identifier for the photo.",
			},
		},
	}
}
