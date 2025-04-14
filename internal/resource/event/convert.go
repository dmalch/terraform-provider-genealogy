package event

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func ElementFrom(ctx context.Context, eventObject types.Object) (*geni.EventElement, diag.Diagnostics) {
	var d diag.Diagnostics

	eventModel, diags := ObjectValueFrom(ctx, eventObject)
	d.Append(diags...)

	if eventModel == nil {
		return nil, d
	}

	dateModel, diags := DateRangeObjectValueFrom(ctx, eventModel.Date)
	d.Append(diags...)

	locationModel, diags := LocationObjectValueFrom(ctx, eventModel.Location)
	d.Append(diags...)

	return &geni.EventElement{
		Name:        eventModel.Name.ValueString(),
		Description: eventModel.Description.ValueString(),
		Date:        DateRangeElementFrom(dateModel),
		Location:    LocationElementFrom(locationModel),
	}, d
}

func ObjectValueFrom(ctx context.Context, eventObject types.Object) (*Model, diag.Diagnostics) {
	if eventObject.IsNull() || eventObject.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var eventModel Model

	diags := eventObject.As(ctx, &eventModel, basetypes.ObjectAsOptions{})

	return &eventModel, diags
}

func LocationElementFrom(locationModel *LocationModel) *geni.LocationElement {
	if locationModel == nil {
		return nil
	}

	return &geni.LocationElement{
		City:           locationModel.City.ValueStringPointer(),
		Country:        locationModel.Country.ValueStringPointer(),
		County:         locationModel.County.ValueStringPointer(),
		Latitude:       locationModel.Latitude.ValueFloat64Pointer(),
		Longitude:      locationModel.Longitude.ValueFloat64Pointer(),
		PlaceName:      locationModel.PlaceName.ValueStringPointer(),
		State:          locationModel.State.ValueStringPointer(),
		StreetAddress1: locationModel.StreetAddress1.ValueStringPointer(),
		StreetAddress2: locationModel.StreetAddress2.ValueStringPointer(),
		StreetAddress3: locationModel.StreetAddress3.ValueStringPointer(),
	}
}

func LocationObjectValueFrom(ctx context.Context, locationObject types.Object) (*LocationModel, diag.Diagnostics) {
	if locationObject.IsNull() || locationObject.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var locationModel LocationModel

	diags := locationObject.As(ctx, &locationModel, basetypes.ObjectAsOptions{})

	return &locationModel, diags
}

func DateElementFrom(model *DateModel) *geni.DateElement {
	if model == nil {
		return nil
	}

	return &geni.DateElement{
		Circa: model.Circa.ValueBoolPointer(),
		Day:   model.Day.ValueInt32Pointer(),
		Month: model.Month.ValueInt32Pointer(),
		Year:  model.Year.ValueInt32Pointer(),
	}
}

func DateRangeElementFrom(model *DateRangeModel) *geni.DateElement {
	if model == nil {
		return nil
	}

	return &geni.DateElement{
		Range:    model.Range.ValueStringPointer(),
		Circa:    model.Circa.ValueBoolPointer(),
		Day:      model.Day.ValueInt32Pointer(),
		Month:    model.Month.ValueInt32Pointer(),
		Year:     model.Year.ValueInt32Pointer(),
		EndCirca: model.EndCirca.ValueBoolPointer(),
		EndDay:   model.EndDay.ValueInt32Pointer(),
		EndMonth: model.EndMonth.ValueInt32Pointer(),
		EndYear:  model.EndYear.ValueInt32Pointer(),
	}
}

func DateObjectValueFrom(ctx context.Context, dateObject types.Object) (*DateModel, diag.Diagnostics) {
	if dateObject.IsNull() || dateObject.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var dateModel DateModel

	diags := dateObject.As(ctx, &dateModel, basetypes.ObjectAsOptions{})

	return &dateModel, diags
}

func DateRangeObjectValueFrom(ctx context.Context, dateObject types.Object) (*DateRangeModel, diag.Diagnostics) {
	if dateObject.IsNull() || dateObject.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var dateModel DateRangeModel

	diags := dateObject.As(ctx, &dateModel, basetypes.ObjectAsOptions{})

	return &dateModel, diags
}

func ValueFrom(ctx context.Context, eventElement *geni.EventElement) (basetypes.ObjectValue, diag.Diagnostics) {
	var d diag.Diagnostics

	if eventElement != nil {
		dateObjectValue, diags := DateRangeValueFrom(ctx, eventElement.Date)
		d.Append(diags...)

		locationObjectValue, diags := LocationValueFrom(ctx, eventElement.Location)
		d.Append(diags...)

		eventModel := Model{
			Description: types.StringValue(eventElement.Description),
			Name:        types.StringValue(eventElement.Name),
			Date:        dateObjectValue,
			Location:    locationObjectValue,
		}

		eventObjectValue, diags := types.ObjectValueFrom(ctx, eventModel.AttributeTypes(), eventModel)
		d.Append(diags...)

		return eventObjectValue, d
	}

	return types.ObjectNull(EventModelAttributeTypes()), d
}

func DateValueFrom(ctx context.Context, dateElement *geni.DateElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if dateElement != nil {
		dateModel := DateModel{
			Circa: types.BoolPointerValue(dateElement.Circa),
			Day:   types.Int32PointerValue(dateElement.Day),
			Month: types.Int32PointerValue(dateElement.Month),
			Year:  types.Int32PointerValue(dateElement.Year),
		}

		return types.ObjectValueFrom(ctx, dateModel.AttributeTypes(), dateModel)
	}

	return types.ObjectNull(DateModelAttributeTypes()), diag.Diagnostics{}
}

func DateRangeValueFrom(ctx context.Context, dateElement *geni.DateElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if dateElement != nil {
		dateModel := DateRangeModel{
			DateModel: DateModel{
				Circa: types.BoolPointerValue(dateElement.Circa),
				Day:   types.Int32PointerValue(dateElement.Day),
				Month: types.Int32PointerValue(dateElement.Month),
				Year:  types.Int32PointerValue(dateElement.Year),
			},
			Range:    types.StringPointerValue(dateElement.Range),
			EndCirca: types.BoolPointerValue(dateElement.EndCirca),
			EndDay:   types.Int32PointerValue(dateElement.EndDay),
			EndMonth: types.Int32PointerValue(dateElement.EndMonth),
			EndYear:  types.Int32PointerValue(dateElement.EndYear),
		}

		return types.ObjectValueFrom(ctx, dateModel.AttributeTypes(), dateModel)
	}

	return types.ObjectNull(DateRangeModelAttributeTypes()), diag.Diagnostics{}
}

func LocationValueFrom(ctx context.Context, location *geni.LocationElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if location != nil {
		locationModel := LocationModel{
			City:           types.StringPointerValue(location.City),
			Country:        types.StringPointerValue(location.Country),
			County:         types.StringPointerValue(location.County),
			Latitude:       types.Float64PointerValue(location.Latitude),
			Longitude:      types.Float64PointerValue(location.Longitude),
			PlaceName:      types.StringPointerValue(location.PlaceName),
			State:          types.StringPointerValue(location.State),
			StreetAddress1: types.StringPointerValue(location.StreetAddress1),
			StreetAddress2: types.StringPointerValue(location.StreetAddress2),
			StreetAddress3: types.StringPointerValue(location.StreetAddress3),
		}

		return types.ObjectValueFrom(ctx, locationModel.AttributeTypes(), locationModel)
	}

	return types.ObjectNull(LocationModelAttributeTypes()), diag.Diagnostics{}
}

func UpdateComputedFieldsInEvent(ctx context.Context, eventObject types.Object, eventElement *geni.EventElement) (types.Object, diag.Diagnostics) {
	var d diag.Diagnostics

	if !eventObject.IsNull() && !eventObject.IsUnknown() {
		var eventModel Model

		diags := eventObject.As(ctx, &eventModel, basetypes.ObjectAsOptions{})
		d.Append(diags...)

		diags = updateComputedFieldsInEventObject(ctx, &eventModel, eventElement)
		d.Append(diags...)

		eventObject, diags = types.ObjectValueFrom(ctx, eventModel.AttributeTypes(), eventModel)
		d.Append(diags...)
	}

	return eventObject, d
}

func updateComputedFieldsInEventObject(ctx context.Context, eventObject *Model, eventElement *geni.EventElement) diag.Diagnostics {
	var d diag.Diagnostics

	if eventObject.Name.IsNull() || eventObject.Name.IsUnknown() {
		eventObject.Name = types.StringValue(eventElement.Name)
	}

	if eventObject.Description.IsNull() || eventObject.Description.IsUnknown() {
		eventObject.Description = types.StringValue(eventElement.Description)
	}

	if eventObject.Location.IsNull() || eventObject.Location.IsUnknown() {
		location, diags := UpdateComputedFieldsInLocationObject(ctx, eventObject.Location, eventElement.Location)
		d.Append(diags...)
		eventObject.Location = location
	}

	return d
}

func UpdateComputedFieldsInLocationObject(ctx context.Context, locationObject types.Object, locationElement *geni.LocationElement) (types.Object, diag.Diagnostics) {
	var d diag.Diagnostics

	var locationModel LocationModel
	diags := locationObject.As(ctx, &locationModel, basetypes.ObjectAsOptions{})
	d.Append(diags...)

	if locationModel.Latitude.IsUnknown() {
		if locationElement == nil || locationElement.Latitude != nil && *locationElement.Latitude == 0.0 {
			locationModel.Latitude = types.Float64Null()
		} else {
			locationModel.Latitude = types.Float64PointerValue(locationElement.Latitude)
		}
	}

	if locationModel.Longitude.IsUnknown() {
		if locationElement == nil || locationElement.Longitude != nil && *locationElement.Longitude == 0.0 {
			locationModel.Longitude = types.Float64Null()
		} else {
			locationModel.Longitude = types.Float64PointerValue(locationElement.Longitude)
		}
	}

	locationObject, diags = types.ObjectValueFrom(ctx, locationModel.AttributeTypes(), locationModel)
	d.Append(diags...)

	return locationObject, d
}
