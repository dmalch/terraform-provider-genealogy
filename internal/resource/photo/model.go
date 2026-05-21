package photo

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Date        types.String `tfsdk:"date"`
	File        types.String `tfsdk:"file"`
	FileName    types.String `tfsdk:"file_name"`
	Album       types.String `tfsdk:"album"`
	Profiles    types.Set    `tfsdk:"profiles"`
	Guid        types.String `tfsdk:"guid"`
	ContentType types.String `tfsdk:"content_type"`
	Attribution types.String `tfsdk:"attribution"`
	URL         types.String `tfsdk:"url"`
	Sizes       types.Map    `tfsdk:"sizes"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

type ResourceIdentityModel struct {
	ID types.String `tfsdk:"id"`
}
