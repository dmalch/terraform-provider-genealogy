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

	if !eventObject.IsNull() && !eventObject.IsUnknown() {
		var eventModel Model

		diags := eventObject.As(ctx, &eventModel, basetypes.ObjectAsOptions{})
		d.Append(diags...)

		date, diags := DateElementFrom(ctx, eventModel.Date)
		d.Append(diags...)

		location, diags := LocationElementFrom(ctx, eventModel.Location)
		d.Append(diags...)

		return &geni.EventElement{
			Name:        eventModel.Name.ValueString(),
			Description: eventModel.Description.ValueString(),
			Date:        date,
			Location:    location,
		}, d
	}

	return nil, d
}

func LocationElementFrom(ctx context.Context, locationObject types.Object) (*geni.LocationElement, diag.Diagnostics) {
	var d diag.Diagnostics

	if !locationObject.IsNull() && !locationObject.IsUnknown() {
		var locationModel LocationModel

		d.Append(locationObject.As(ctx, &locationModel, basetypes.ObjectAsOptions{})...)

		return &geni.LocationElement{
			City:           locationModel.City.ValueStringPointer(),
			Country:        locationModel.Country.ValueStringPointer(),
			County:         locationModel.County.ValueStringPointer(),
			Latitude:       locationModel.Latitude.ValueBigFloat(),
			Longitude:      locationModel.Longitude.ValueBigFloat(),
			PlaceName:      locationModel.PlaceName.ValueStringPointer(),
			State:          locationModel.State.ValueStringPointer(),
			StreetAddress1: locationModel.StreetAddress1.ValueStringPointer(),
			StreetAddress2: locationModel.StreetAddress2.ValueStringPointer(),
			StreetAddress3: locationModel.StreetAddress3.ValueStringPointer(),
		}, d
	}

	return nil, d
}

func DateElementFrom(ctx context.Context, dateObject types.Object) (*geni.DateElement, diag.Diagnostics) {
	var d diag.Diagnostics

	if !dateObject.IsNull() && !dateObject.IsUnknown() {
		var dateModel DateModel

		d.Append(dateObject.As(ctx, &dateModel, basetypes.ObjectAsOptions{})...)

		return &geni.DateElement{
			Range:    dateModel.Range.ValueString(),
			Circa:    dateModel.Circa.ValueBool(),
			Day:      int(dateModel.Day.ValueInt32()),
			Month:    int(dateModel.Month.ValueInt32()),
			Year:     int(dateModel.Year.ValueInt32()),
			EndCirca: dateModel.EndCirca.ValueBool(),
			EndDay:   int(dateModel.EndDay.ValueInt32()),
			EndMonth: int(dateModel.EndMonth.ValueInt32()),
			EndYear:  int(dateModel.EndYear.ValueInt32()),
		}, d
	}

	return nil, d
}

func ValueFrom(ctx context.Context, eventElement *geni.EventElement) (basetypes.ObjectValue, diag.Diagnostics) {
	var d diag.Diagnostics

	if eventElement != nil {
		dateObjectValue, diags := DateValueFrom(ctx, eventElement.Date)
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
			Range:    types.StringValue(dateElement.Range),
			Circa:    types.BoolValue(dateElement.Circa),
			Day:      types.Int32Value(int32(dateElement.Day)),
			Month:    types.Int32Value(int32(dateElement.Month)),
			Year:     types.Int32Value(int32(dateElement.Year)),
			EndCirca: types.BoolValue(dateElement.EndCirca),
			EndDay:   types.Int32Value(int32(dateElement.EndDay)),
			EndMonth: types.Int32Value(int32(dateElement.EndMonth)),
			EndYear:  types.Int32Value(int32(dateElement.EndYear)),
		}

		return types.ObjectValueFrom(ctx, dateModel.AttributeTypes(), dateModel)
	}

	return types.ObjectNull(DateModelAttributeTypes()), diag.Diagnostics{}
}

func LocationValueFrom(ctx context.Context, location *geni.LocationElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if location != nil {
		locationModel := LocationModel{
			City:           types.StringPointerValue(location.City),
			Country:        types.StringPointerValue(location.Country),
			County:         types.StringPointerValue(location.County),
			Latitude:       types.NumberValue(location.Latitude),
			Longitude:      types.NumberValue(location.Longitude),
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

func updateComputedFieldsInEventObject(_ context.Context, eventObject *Model, eventElement *geni.EventElement) diag.Diagnostics {
	var d diag.Diagnostics

	if eventObject.Name.IsNull() || eventObject.Name.IsUnknown() {
		eventObject.Name = types.StringValue(eventElement.Name)
	}
	if eventObject.Description.IsNull() || eventObject.Description.IsUnknown() {
		eventObject.Description = types.StringValue(eventElement.Description)
	}

	return d
}
