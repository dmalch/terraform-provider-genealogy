package profile

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID         types.String `tfsdk:"id"`
	FirstName  types.String `tfsdk:"first_name"`
	LastName   types.String `tfsdk:"last_name"`
	MiddleName types.String `tfsdk:"middle_name"`
	MaidenName types.String `tfsdk:"maiden_name"`
	Names      types.Map    `tfsdk:"names"`
	Gender     types.String `tfsdk:"gender"`
	Birth      types.Object `tfsdk:"birth"`
	Baptism    types.Object `tfsdk:"baptism"`
	Death      types.Object `tfsdk:"death"`
	Burial     types.Object `tfsdk:"burial"`
	Unions     types.List   `tfsdk:"unions"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

type NameModel struct {
	FistName   types.String `tfsdk:"first_name"`
	MiddleName types.String `tfsdk:"middle_name"`
	LastName   types.String `tfsdk:"last_name"`
	MaidenName types.String `tfsdk:"maiden_name"`
}

func (m NameModel) AttributeTypes() map[string]attr.Type {
	return NameAttributeTypes()
}

func NameAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_name":  types.StringType,
		"middle_name": types.StringType,
		"last_name":   types.StringType,
		"maiden_name": types.StringType,
	}
}
