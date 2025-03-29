package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ValueFrom(ctx context.Context, profile *geni.ProfileResponse, profileModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	if profile.Id != "" {
		profileModel.ID = types.StringValue(profile.Id)
	}

	if profile.FirstName != "" {
		profileModel.FirstName = types.StringValue(profile.FirstName)
	}

	if profile.LastName != "" {
		profileModel.LastName = types.StringValue(profile.LastName)
	}

	if profile.MiddleName != "" {
		profileModel.MiddleName = types.StringValue(profile.MiddleName)
	}

	if profile.MaidenName != "" {
		profileModel.MaidenName = types.StringValue(profile.MaidenName)
	}

	if profile.Gender != "" {
		profileModel.Gender = types.StringValue(profile.Gender)
	}

	names, diags := NameValueFrom(ctx, profile.Names)
	d.Append(diags...)
	profileModel.Names = names

	unions, diags := types.ListValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	birth, diags := event.ValueFrom(ctx, profile.Birth)
	d.Append(diags...)
	profileModel.Birth = birth

	baptism, diags := event.ValueFrom(ctx, profile.Baptism)
	d.Append(diags...)
	profileModel.Baptism = baptism

	death, diags := event.ValueFrom(ctx, profile.Death)
	d.Append(diags...)
	profileModel.Death = death

	burial, diags := event.ValueFrom(ctx, profile.Burial)
	d.Append(diags...)
	profileModel.Burial = burial

	if profile.CreatedAt != "" {
		profileModel.CreatedAt = types.StringValue(profile.CreatedAt)
	}

	return d
}

func NameValueFrom(ctx context.Context, profileNames map[string]geni.NameElement) (basetypes.MapValue, diag.Diagnostics) {
	if len(profileNames) == 0 {
		return basetypes.NewMapNull(types.ObjectType{AttrTypes: NameAttributeTypes()}), diag.Diagnostics{}
	}

	nameModels := make(map[string]NameModel, len(profileNames))

	for locale, name := range profileNames {
		nameModels[locale] = NameModel{
			FistName:   types.StringPointerValue(name.FirstName),
			MiddleName: types.StringPointerValue(name.MiddleName),
			LastName:   types.StringPointerValue(name.LastName),
			MaidenName: types.StringPointerValue(name.MaidenName),
		}
	}

	return types.MapValueFrom(ctx, types.ObjectType{AttrTypes: NameAttributeTypes()}, nameModels)
}

func RequestFrom(ctx context.Context, resourceModel ResourceModel) (*geni.ProfileRequest, diag.Diagnostics) {
	var d diag.Diagnostics

	birth, diags := event.ElementFrom(ctx, resourceModel.Birth)
	d.Append(diags...)

	baptism, diags := event.ElementFrom(ctx, resourceModel.Baptism)
	d.Append(diags...)

	death, diags := event.ElementFrom(ctx, resourceModel.Death)
	d.Append(diags...)

	burial, diags := event.ElementFrom(ctx, resourceModel.Burial)
	d.Append(diags...)

	convertedNames, diags := NameElementsFrom(ctx, resourceModel.Names)
	d.Append(diags...)

	profileRequest := &geni.ProfileRequest{
		FirstName:  resourceModel.FirstName.ValueString(),
		LastName:   resourceModel.LastName.ValueString(),
		MiddleName: resourceModel.MiddleName.ValueString(),
		MaidenName: resourceModel.MaidenName.ValueString(),
		Names:      convertedNames,
		Gender:     resourceModel.Gender.ValueString(),
		Birth:      birth,
		Baptism:    baptism,
		Death:      death,
		Burial:     burial,
	}

	return profileRequest, d
}

func NameElementsFrom(ctx context.Context, names types.Map) (map[string]geni.NameElement, diag.Diagnostics) {
	nameModels, diags := NameModelsFrom(ctx, names)

	var profileNames = make(map[string]geni.NameElement, len(nameModels))

	for locale, nameModel := range nameModels {
		profileNames[locale] = geni.NameElement{
			FirstName:  nameModel.FistName.ValueStringPointer(),
			MiddleName: nameModel.MiddleName.ValueStringPointer(),
			LastName:   nameModel.LastName.ValueStringPointer(),
			MaidenName: nameModel.MaidenName.ValueStringPointer(),
		}
	}

	return profileNames, diags
}

func NameModelsFrom(ctx context.Context, names types.Map) (map[string]NameModel, diag.Diagnostics) {
	if len(names.Elements()) == 0 {
		return nil, diag.Diagnostics{}
	}

	var nameModels = make(map[string]NameModel, len(names.Elements()))
	diags := names.ElementsAs(ctx, &nameModels, false)

	return nameModels, diags
}
