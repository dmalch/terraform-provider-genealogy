package profile

import (
	"context"
	"strings"

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

	aboutMeMap := make(map[string]string)
	for locale, localeDetails := range profile.DetailStrings {
		if localeDetails.AboutMe != nil {
			aboutMeMap[locale] = *localeDetails.AboutMe
		}
	}
	// Fallback to AboutMe if DetailStrings is empty
	if len(aboutMeMap) == 0 && profile.AboutMe != nil && *profile.AboutMe != "" {
		aboutMeMap["en-US"] = *profile.AboutMe
	}

	detailStrings, diags := types.MapValueFrom(ctx, types.StringType, aboutMeMap)
	d.Append(diags...)
	profileModel.About = detailStrings

	profileModel.Public = types.BoolValue(profile.Public)
	profileModel.Alive = types.BoolValue(profile.IsAlive)

	currentResidence, diags := event.LocationValueFrom(ctx, profile.CurrentResidence)
	d.Append(diags...)
	profileModel.CurrentResidence = currentResidence

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

	var d diag.Diagnostics
	for locale, name := range profileNames {
		nicknames := types.SetNull(types.StringType)
		if name.Nicknames != nil {
			nicknamesList := strings.Split(*name.Nicknames, ",")
			var diags diag.Diagnostics
			nicknames, diags = types.SetValueFrom(ctx, types.StringType, nicknamesList)
			d.Append(diags...)
		}

		nameModels[locale] = NameModel{
			FistName:      types.StringPointerValue(name.FirstName),
			MiddleName:    types.StringPointerValue(name.MiddleName),
			LastName:      types.StringPointerValue(name.LastName),
			BirthLastName: types.StringPointerValue(name.MaidenName),
			DisplayName:   types.StringPointerValue(name.DisplayName),
			Nicknames:     nicknames,
		}
	}

	nameMap, diags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: NameAttributeTypes()}, nameModels)
	d.Append(diags...)
	return nameMap, d
}

func RequestFrom(ctx context.Context, resourceModel ResourceModel) (*geni.ProfileRequest, diag.Diagnostics) {
	var d diag.Diagnostics

	birth, diags := event.ElementFrom(ctx, resourceModel.Birth)
	d.Append(diags...)
	if birth == nil {
		birth = &geni.EventElement{}
	}

	baptism, diags := event.ElementFrom(ctx, resourceModel.Baptism)
	d.Append(diags...)
	if baptism == nil {
		baptism = &geni.EventElement{}
	}

	death, diags := event.ElementFrom(ctx, resourceModel.Death)
	d.Append(diags...)
	if death == nil {
		death = &geni.EventElement{}
	}

	burial, diags := event.ElementFrom(ctx, resourceModel.Burial)
	d.Append(diags...)
	if burial == nil {
		burial = &geni.EventElement{}
	}

	currentResidence, diags := event.LocationObjectValueFrom(ctx, resourceModel.CurrentResidence)
	d.Append(diags...)

	var convertedNames map[string]geni.NameElement
	if len(resourceModel.Names.Elements()) > 0 {
		convertedNames, diags = NameElementsFrom(ctx, resourceModel.Names)
		d.Append(diags...)
	}

	convertedDetails := make(map[string]geni.DetailsString)
	if len(resourceModel.About.Elements()) > 0 {
		var aboutMeMap map[string]*string
		diags := resourceModel.About.ElementsAs(ctx, &aboutMeMap, false)
		d.Append(diags...)

		for locale, detailsString := range aboutMeMap {
			convertedDetails[locale] = geni.DetailsString{
				AboutMe: detailsString,
			}
		}
	}

	profileRequest := &geni.ProfileRequest{
		Names:            convertedNames,
		Gender:           resourceModel.Gender.ValueStringPointer(),
		Birth:            birth,
		Baptism:          baptism,
		Death:            death,
		Burial:           burial,
		CauseOfDeath:     resourceModel.CauseOfDeath.ValueStringPointer(),
		CurrentResidence: event.LocationElementFrom(currentResidence),
		DetailStrings:    convertedDetails,
		Public:           resourceModel.Public.ValueBool(),
		IsAlive:          resourceModel.Alive.ValueBool(),
	}

	return profileRequest, d
}

func NameElementsFrom(ctx context.Context, names types.Map) (map[string]geni.NameElement, diag.Diagnostics) {
	nameModels, diags := NameModelsFrom(ctx, names)

	var profileNames = make(map[string]geni.NameElement, len(nameModels))

	for locale, nameModel := range nameModels {
		nicknamesSlice, d := convertToSlice(ctx, nameModel.Nicknames)
		diags = append(diags, d...)
		nicknamesCsv := strings.Join(nicknamesSlice, ",")

		profileNames[locale] = geni.NameElement{
			FirstName:   nameModel.FistName.ValueStringPointer(),
			MiddleName:  nameModel.MiddleName.ValueStringPointer(),
			LastName:    nameModel.LastName.ValueStringPointer(),
			MaidenName:  nameModel.BirthLastName.ValueStringPointer(),
			DisplayName: nameModel.DisplayName.ValueStringPointer(),
			Nicknames:   &nicknamesCsv,
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

	currentResidence, diags := event.UpdateComputedFieldsInLocationObject(ctx, profileModel.CurrentResidence, profile.CurrentResidence)
	d.Append(diags...)
	profileModel.CurrentResidence = currentResidence

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
