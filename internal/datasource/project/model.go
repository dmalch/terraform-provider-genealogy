package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Model struct {
	ID          types.String          `tfsdk:"id"`
	Name        basetypes.StringValue `tfsdk:"name"`
	Description basetypes.StringValue `tfsdk:"description"`
	UpdatedAt   basetypes.StringValue `tfsdk:"updated_at"`
	CreatedAt   basetypes.StringValue `tfsdk:"created_at"`
}
