package profile

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID        types.String `tfsdk:"id"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Gender    types.String `tfsdk:"gender"`
	Birth     types.Object `tfsdk:"birth"`
	Baptism   types.Object `tfsdk:"baptism"`
	Death     types.Object `tfsdk:"death"`
	Burial    types.Object `tfsdk:"burial"`
	Unions    types.List   `tfsdk:"unions"`
	CreatedAt types.String `tfsdk:"created_at"`
}
