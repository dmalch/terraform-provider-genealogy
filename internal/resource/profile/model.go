package profile

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Names            types.Map    `tfsdk:"names"`
	Gender           types.String `tfsdk:"gender"`
	Birth            types.Object `tfsdk:"birth"`
	Baptism          types.Object `tfsdk:"baptism"`
	Death            types.Object `tfsdk:"death"`
	Burial           types.Object `tfsdk:"burial"`
	CauseOfDeath     types.String `tfsdk:"cause_of_death"`
	Unions           types.Set    `tfsdk:"unions"`
	Projects         types.Set    `tfsdk:"projects"`
	CurrentResidence types.Object `tfsdk:"current_residence"`
	About            types.Map    `tfsdk:"about"`
	Public           types.Bool   `tfsdk:"public"`
	Alive            types.Bool   `tfsdk:"alive"`
	Deleted          types.Bool   `tfsdk:"deleted"`
	MergedInto       types.String `tfsdk:"merged_into"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

type ResourceModelV0 struct {
	ID               types.String `tfsdk:"id"`
	Names            types.Map    `tfsdk:"names"`
	Gender           types.String `tfsdk:"gender"`
	Birth            types.Object `tfsdk:"birth"`
	Baptism          types.Object `tfsdk:"baptism"`
	Death            types.Object `tfsdk:"death"`
	Burial           types.Object `tfsdk:"burial"`
	CauseOfDeath     types.String `tfsdk:"cause_of_death"`
	Unions           types.Set    `tfsdk:"unions"`
	Projects         types.Set    `tfsdk:"projects"`
	CurrentResidence types.Object `tfsdk:"current_residence"`
	About            types.String `tfsdk:"about"`
	Public           types.Bool   `tfsdk:"public"`
	Alive            types.Bool   `tfsdk:"alive"`
	Deleted          types.Bool   `tfsdk:"deleted"`
	MergedInto       types.String `tfsdk:"merged_into"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

type NameModel struct {
	FistName      types.String `tfsdk:"first_name"`
	MiddleName    types.String `tfsdk:"middle_name"`
	LastName      types.String `tfsdk:"last_name"`
	BirthLastName types.String `tfsdk:"birth_last_name"`
	DisplayName   types.String `tfsdk:"display_name"`
	Nicknames     types.Set    `tfsdk:"nicknames"`
}

func (m NameModel) AttributeTypes() map[string]attr.Type {
	return NameAttributeTypes()
}

func NameAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"first_name":      types.StringType,
		"middle_name":     types.StringType,
		"last_name":       types.StringType,
		"birth_last_name": types.StringType,
		"display_name":    types.StringType,
		"nicknames": types.SetType{
			ElemType: types.StringType,
		},
	}
}
