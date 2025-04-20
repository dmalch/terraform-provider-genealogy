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

	profileModel.Gender = types.StringPointerValue(profile.Gender)
	profileModel.About = types.StringPointerValue(profile.AboutMe)
	profileModel.Public = types.BoolValue(profile.Public)
	profileModel.Alive = types.BoolValue(profile.IsAlive)

	names, diags := NameValueFrom(ctx, profile.Names)
	d.Append(diags...)
	profileModel.Names = names

	unions, diags := types.SetValueFrom(ctx, types.StringType, profile.Unions)
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

	profileModel.CauseOfDeath = types.StringPointerValue(profile.CauseOfDeath)
	profileModel.Deleted = types.BoolValue(profile.Deleted)
	profileModel.MergedInto = types.StringValue(profile.MergedInto)

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
			FistName:      types.StringPointerValue(name.FirstName),
			MiddleName:    types.StringPointerValue(name.MiddleName),
			LastName:      types.StringPointerValue(name.LastName),
			BirthLastName: types.StringPointerValue(name.MaidenName),
			DisplayName:   types.StringPointerValue(name.DisplayName),
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

	var convertedNames map[string]geni.NameElement
	if len(resourceModel.Names.Elements()) > 0 {
		convertedNames, diags = NameElementsFrom(ctx, resourceModel.Names)
		d.Append(diags...)
	}

	profileRequest := &geni.ProfileRequest{
		Names:        convertedNames,
		Gender:       resourceModel.Gender.ValueStringPointer(),
		Birth:        birth,
		Baptism:      baptism,
		Death:        death,
		Burial:       burial,
		CauseOfDeath: resourceModel.CauseOfDeath.ValueStringPointer(),
		AboutMe:      resourceModel.About.ValueStringPointer(),
		Public:       resourceModel.Public.ValueBool(),
		IsAlive:      resourceModel.Alive.ValueBool(),
	}

	return profileRequest, d
}

func NameElementsFrom(ctx context.Context, names types.Map) (map[string]geni.NameElement, diag.Diagnostics) {
	nameModels, diags := NameModelsFrom(ctx, names)

	var profileNames = make(map[string]geni.NameElement, len(nameModels))

	for locale, nameModel := range nameModels {
		profileNames[locale] = geni.NameElement{
			FirstName:   nameModel.FistName.ValueStringPointer(),
			MiddleName:  nameModel.MiddleName.ValueStringPointer(),
			LastName:    nameModel.LastName.ValueStringPointer(),
			MaidenName:  nameModel.BirthLastName.ValueStringPointer(),
			DisplayName: nameModel.DisplayName.ValueStringPointer(),
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

func UpdateComputedFields(ctx context.Context, profile *geni.ProfileResponse, profileModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	profileModel.ID = types.StringValue(profile.Id)

	unions, diags := types.SetValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	birth, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Birth, profile.Birth)
	d.Append(diags...)
	profileModel.Birth = birth

	baptism, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Baptism, profile.Baptism)
	d.Append(diags...)
	profileModel.Baptism = baptism

	death, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Death, profile.Death)
	d.Append(diags...)
	profileModel.Death = death

	burial, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Burial, profile.Burial)
	d.Append(diags...)
	profileModel.Burial = burial

	profileModel.Deleted = types.BoolValue(profile.Deleted)
	profileModel.MergedInto = types.StringValue(profile.MergedInto)
	profileModel.CreatedAt = types.StringValue(profile.CreatedAt)

	return d
}

func convertToSlice(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	if len(set.Elements()) == 0 {
		return nil, diag.Diagnostics{}
	}

	var slice = make([]string, len(set.Elements()))
	diags := set.ElementsAs(ctx, &slice, false)

	return slice, diags
}
