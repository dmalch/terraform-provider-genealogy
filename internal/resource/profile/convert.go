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
	nameModels := make(map[string]NameModel, len(profileNames))

	for locale, name := range profileNames {
		nameModels[locale] = NameModel{
			FistName:   types.StringPointerValue(name.FirstName),
			MiddleName: types.StringPointerValue(name.MiddleName),
			LastName:   types.StringPointerValue(name.LastName),
		}
	}

	return types.MapValueFrom(ctx, types.ObjectType{AttrTypes: NameAttributeTypes()}, nameModels)
}

func RequestFrom(ctx context.Context, plan ResourceModel) (*geni.ProfileRequest, diag.Diagnostics) {
	var d diag.Diagnostics

	birth, diags := event.ElementFrom(ctx, plan.Birth)
	d.Append(diags...)

	baptism, diags := event.ElementFrom(ctx, plan.Baptism)
	d.Append(diags...)

	death, diags := event.ElementFrom(ctx, plan.Death)
	d.Append(diags...)

	burial, diags := event.ElementFrom(ctx, plan.Burial)
	d.Append(diags...)

	convertedNames, diags := NameElementFrom(ctx, plan)
	d.Append(diags...)

	profileRequest := &geni.ProfileRequest{
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
		Names:     convertedNames,
		Gender:    plan.Gender.ValueString(),
		Birth:     birth,
		Baptism:   baptism,
		Death:     death,
		Burial:    burial,
	}

	return profileRequest, d
}

func NameElementFrom(ctx context.Context, plan ResourceModel) (map[string]geni.NameElement, diag.Diagnostics) {
	var nameModels = make(map[string]NameModel)
	diags := plan.Names.ElementsAs(ctx, &nameModels, false)

	var profileNames = make(map[string]geni.NameElement, len(nameModels))

	for locale, nameModel := range nameModels {
		profileNames[locale] = geni.NameElement{
			FirstName:  nameModel.FistName.ValueStringPointer(),
			MiddleName: nameModel.MiddleName.ValueStringPointer(),
			LastName:   nameModel.LastName.ValueStringPointer(),
		}
	}

	return profileNames, diags
}
