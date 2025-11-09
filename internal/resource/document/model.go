package document

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	ContentType types.String `tfsdk:"content_type"`
	Text        types.String `tfsdk:"text"`
	File        types.String `tfsdk:"file"`
	FileName    types.String `tfsdk:"file_name"`
	SourceUrl   types.String `tfsdk:"source_url"`
	Date        types.Object `tfsdk:"date"`
	Location    types.Object `tfsdk:"location"`
	Profiles    types.Set    `tfsdk:"profiles"`
	Projects    types.Set    `tfsdk:"projects"`
	Labels      types.Set    `tfsdk:"labels"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

type ResourceIdentityModel struct {
	ID types.String `tfsdk:"id"`
}
