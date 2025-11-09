package union

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Children types.Set    `tfsdk:"children"`
	Partners types.Set    `tfsdk:"partners"`
	Marriage types.Object `tfsdk:"marriage"`
	Divorce  types.Object `tfsdk:"divorce"`
}

type ResourceIdentityModel struct {
	ID types.String `tfsdk:"id"`
}
