package document

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/geniplanmodifier"
)

var (
	documentIdFormat       = regexp.MustCompile(`^document-\d+$`)
	createdAtFormat        = regexp.MustCompile(`^\d+$`)
	documentMimeTypeFormat = regexp.MustCompile(`^(text/(plain|html)|(application/pdf)|image/(jpg|png|tif))$`)
	urlFormat              = regexp.MustCompile(`^https?://`)
)

// Schema defines the schema for the resource.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Optional:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Validators: []validator.String{
					stringvalidator.RegexMatches(documentIdFormat, "must be in the format document-1"),
				},
				Description: "The unique identifier for the document. This is a string that starts with 'document-' followed by a number.",
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
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.RegexMatches(documentMimeTypeFormat, "must be a supported mime type")},
				Description: "The document's original content type.",
			},
			"text": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{stringvalidator.ExactlyOneOf(
					path.MatchRelative().AtParent().AtName("file"),
					path.MatchRelative().AtParent().AtName("source_url"),
				)},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Description:   "The document's text content.",
			},
			"file": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("text"),
						path.MatchRelative().AtParent().AtName("source_url"),
					),
					stringvalidator.AlsoRequires(
						path.MatchRelative().AtParent().AtName("file_name"),
						path.MatchRelative().AtParent().AtName("content_type"),
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Description:   "The document's file content. This is a base64 encoded string.",
			},
			"file_name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("text"),
						path.MatchRelative().AtParent().AtName("source_url"),
					),
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("file")),
				},
				Description: "The document's filename. Required if the file is set.",
			},
			"source_url": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("text"),
						path.MatchRelative().AtParent().AtName("file"),
					),
					stringvalidator.RegexMatches(urlFormat, "must be a valid URL"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
				Description:   "The document's source URL. This is the URL where the document was found.",
			},
			"date":     event.DateSchema("Document's date."),
			"location": event.LocationSchema("Document's location."),
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Description: "The list of profiles associated with the document.",
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
			"labels": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				PlanModifiers: []planmodifier.Set{setplanmodifier.RequiresReplaceIf(
					geniplanmodifier.ValuesAreRemovedFromState,
					"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
					"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.",
				)},
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

func (r *Resource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "The unique identifier for the document.",
			},
		},
	}
}
