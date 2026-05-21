package profile

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	geniprofile "github.com/dmalch/go-geni/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
	"github.com/dmalch/terraform-provider-genealogy/internal/tfset"
)

func ValueFrom(ctx context.Context, profile *geniprofile.Profile, profileModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	if profile.ID != "" {
		profileModel.ID = types.StringValue(profile.ID)
	}

	if profile.Guid != "" {
		profileModel.Guid = types.StringValue(profile.Guid)
	} else {
		profileModel.Guid = types.StringNull()
	}

	profileModel.Title = optionalString(profile.Title)
	profileModel.Occupation = optionalString(profile.Occupation)
	profileModel.Suffix = optionalString(profile.Suffix)

	profileModel.Gender = types.StringPointerValue(profile.Gender)

	detailStrings, diags := detailStringsValueFrom(ctx, profile)
	d.Append(diags...)

	if len(detailStrings.Elements()) > 0 {
		profileModel.About = detailStrings
	}

	profileModel.Public = types.BoolValue(profile.Public)
	profileModel.Alive = types.BoolValue(profile.IsAlive)

	currentResidence, diags := event.LocationValueFrom(ctx, profile.CurrentResidence)
	d.Append(diags...)
	profileModel.CurrentResidence = currentResidence

	names, diags := NameValueFrom(ctx, namesWithFlatFallback(profile))
	d.Append(diags...)
	profileModel.Names = names

	unions, diags := types.SetValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	projects, diags := types.SetValueFrom(ctx, types.StringType, profile.Projects)
	d.Append(diags...)
	profileModel.Projects = projects

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

func detailStringsValueFrom(ctx context.Context, profile *geniprofile.Profile) (basetypes.MapValue, diag.Diagnostics) {
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

	return types.MapValueFrom(ctx, types.StringType, aboutMeMap)
}

// namesWithFlatFallback returns the API's localized name map when present, or
// synthesizes an en-US entry from the response's top-level flat name fields
// (first_name/last_name/...) when the map is empty. The Geni API returns flat
// fields for profiles created without a locale tag — without this fallback,
// importing such a profile produces a null `names` map and the next plan
// shows a spurious in-place update that recreates the locale entry.
func namesWithFlatFallback(profile *geniprofile.Profile) map[string]geniprofile.NameElement {
	if len(profile.Names) > 0 {
		return profile.Names
	}
	if profile.FirstName == nil && profile.LastName == nil && profile.MiddleName == nil &&
		profile.MaidenName == nil && profile.DisplayName == nil && len(profile.Nicknames) == 0 {
		return profile.Names
	}
	var nicknames *string
	if len(profile.Nicknames) > 0 {
		joined := strings.Join(profile.Nicknames, ",")
		nicknames = &joined
	}
	return map[string]geniprofile.NameElement{
		"en-US": {
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			MiddleName:  profile.MiddleName,
			MaidenName:  profile.MaidenName,
			DisplayName: profile.DisplayName,
			Nicknames:   nicknames,
		},
	}
}

func NameValueFrom(ctx context.Context, profileNames map[string]geniprofile.NameElement) (basetypes.MapValue, diag.Diagnostics) {
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
			FirstName:     types.StringPointerValue(name.FirstName),
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

func RequestFrom(ctx context.Context, resourceModel ResourceModel) (*geniprofile.Request, diag.Diagnostics) {
	var d diag.Diagnostics

	birth, diags := event.ElementFrom(ctx, resourceModel.Birth)
	d.Append(diags...)
	if birth == nil {
		birth = &geniprofile.EventElement{}
	}

	baptism, diags := event.ElementFrom(ctx, resourceModel.Baptism)
	d.Append(diags...)
	if baptism == nil {
		baptism = &geniprofile.EventElement{}
	}

	death, diags := event.ElementFrom(ctx, resourceModel.Death)
	d.Append(diags...)
	if death == nil {
		death = &geniprofile.EventElement{}
	}

	burial, diags := event.ElementFrom(ctx, resourceModel.Burial)
	d.Append(diags...)
	if burial == nil {
		burial = &geniprofile.EventElement{}
	}

	currentResidence, diags := event.LocationObjectValueFrom(ctx, resourceModel.CurrentResidence)
	d.Append(diags...)

	var convertedNames map[string]geniprofile.NameElement
	if len(resourceModel.Names.Elements()) > 0 {
		convertedNames, diags = NameElementsFrom(ctx, resourceModel.Names)
		d.Append(diags...)
	}

	convertedDetails := make(map[string]geniprofile.DetailsString)
	if len(resourceModel.About.Elements()) > 0 {
		var aboutMeMap map[string]*string
		diags := resourceModel.About.ElementsAs(ctx, &aboutMeMap, false)
		d.Append(diags...)

		for locale, detailsString := range aboutMeMap {
			convertedDetails[locale] = geniprofile.DetailsString{
				AboutMe: detailsString,
			}
		}
	}

	profileRequest := &geniprofile.Request{
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
		Title:            resourceModel.Title.ValueString(),
		Occupation:       resourceModel.Occupation.ValueString(),
		Suffix:           resourceModel.Suffix.ValueString(),
	}

	return profileRequest, d
}

func NameElementsFrom(ctx context.Context, names types.Map) (map[string]geniprofile.NameElement, diag.Diagnostics) {
	nameModels, diags := NameModelsFrom(ctx, names)

	var profileNames = make(map[string]geniprofile.NameElement, len(nameModels))

	for locale, nameModel := range nameModels {
		nicknamesSlice, d := tfset.Strings(ctx, nameModel.Nicknames)
		diags = append(diags, d...)
		nicknamesCsv := strings.Join(nicknamesSlice, ",")

		profileNames[locale] = geniprofile.NameElement{
			FirstName:   nameModel.FirstName.ValueStringPointer(),
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

func UpdateComputedFields(ctx context.Context, profile *geniprofile.Profile, profileModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	profileModel.ID = types.StringValue(profile.ID)

	if profile.Guid != "" {
		profileModel.Guid = types.StringValue(profile.Guid)
	} else {
		profileModel.Guid = types.StringNull()
	}

	unions, diags := types.SetValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	// Projects is intentionally not updated here: Create/Update pass the
	// pre-link API response (the AddProfileToProject calls fire afterwards),
	// so profile.Projects is stale and would overwrite plan.Projects with
	// null. The Read path's ValueFrom is responsible for populating
	// projects from the API; the plan is authoritative on Create/Update.

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

	detailStrings, diags := detailStringsValueFrom(ctx, profile)
	d.Append(diags...)

	if len(detailStrings.Elements()) > 0 {
		profileModel.About = detailStrings
	}

	profileModel.Deleted = types.BoolValue(profile.Deleted)
	profileModel.MergedInto = types.StringValue(profile.MergedInto)
	profileModel.CreatedAt = types.StringValue(profile.CreatedAt)

	return d
}

// optionalString maps a Geni API string field (where "" means "not set") to a
// Terraform-framework optional string: empty becomes null, non-empty becomes
// the value. Used for plain string fields like title/occupation/suffix that
// the API returns omitempty.
func optionalString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
