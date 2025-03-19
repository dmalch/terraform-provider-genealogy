package union

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Children types.Set    `tfsdk:"children"`
	Partners types.Set    `tfsdk:"partners"`
}
