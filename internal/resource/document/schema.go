package document

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
	documentIdFormat = regexp.MustCompile(`^profile-\d+$`)
	createdAtFormat  = regexp.MustCompile(`^\d+$`)
	documentMimeType = regexp.MustCompile(`^((text/plain)|(application/pdf))$`)
)

// Schema defines the schema for the resource.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(documentIdFormat, "must be in the format document-1")},
				Description:   "The unique identifier for the document. This is a string that starts with 'document-' followed by a number.",
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "The document's title.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The document's description.",
			},
			"content_type": schema.StringAttribute{
				Required:    true,
				Validators:  []validator.String{stringvalidator.RegexMatches(documentMimeType, "must be a valid mime type")},
				Description: "The document's original content type.",
			},
			"date":     event.DateSchema("Document's date."),
			"location": event.LocationSchema("Document's location."),
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "The list of profiles associated with the document.",
			},
			"labels": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "The list of labels associated with the document.",
			},
			"created_at": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators:    []validator.String{stringvalidator.RegexMatches(createdAtFormat, "must be a Unix epoch time in seconds")},
				Description:   "The Unix epoch time in seconds when the document was created.",
			},
		},
	}
}
