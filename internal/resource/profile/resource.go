package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/dmalch/terraform-provider-geni/internal/config"
	"github.com/dmalch/terraform-provider-geni/internal/geni"
	"github.com/dmalch/terraform-provider-geni/internal/resource/event"
)

type Resource struct {
	resource.ResourceWithConfigure
	accessToken types.String
}

func NewProfileResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "geni_profile"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.GeniProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.GeniProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.accessToken = cfg.AccessToken
}

// Create creates the resource
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	birth, diags := EventElementFrom(ctx, plan.Birth)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileRequest := &geni.ProfileRequest{
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
		Gender:    plan.Gender.ValueString(),
		Birth:     birth,
	}

	profileResponse, err := geni.CreateProfile(r.accessToken.ValueString(), profileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating profile", err.Error())
		return
	}

	diags = updateComputedFields(ctx, plan, profileResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func updateComputedFields(ctx context.Context, profileModel ResourceModel, profile *geni.ProfileResponse) diag.Diagnostics {
	var d diag.Diagnostics

	profileModel.ID = types.StringValue(profile.Id)

	unions, diags := types.ListValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	profileModel.CreatedAt = types.StringValue(profile.CreatedAt)

	return d
}

func EventElementFrom(ctx context.Context, eventObject types.Object) (*geni.EventElement, diag.Diagnostics) {
	var d diag.Diagnostics

	if !eventObject.IsNull() && !eventObject.IsUnknown() {
		var eventModel event.Model

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
		var locationModel event.LocationModel

		d.Append(locationObject.As(ctx, &locationModel, basetypes.ObjectAsOptions{})...)

		return &geni.LocationElement{
			City:           locationModel.City.ValueString(),
			Country:        locationModel.Country.ValueString(),
			County:         locationModel.County.ValueString(),
			Latitude:       locationModel.Latitude.ValueBigFloat(),
			Longitude:      locationModel.Longitude.ValueBigFloat(),
			PlaceName:      locationModel.PlaceName.ValueString(),
			State:          locationModel.State.ValueString(),
			StreetAddress1: locationModel.StreetAddress1.ValueString(),
			StreetAddress2: locationModel.StreetAddress2.ValueString(),
			StreetAddress3: locationModel.StreetAddress3.ValueString(),
		}, d
	}

	return nil, d
}

func DateElementFrom(ctx context.Context, dateObject types.Object) (*geni.DateElement, diag.Diagnostics) {
	var d diag.Diagnostics

	if !dateObject.IsNull() && !dateObject.IsUnknown() {
		var dateModel event.DateModel

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

// Read reads the resource
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := geni.GetProfile(r.accessToken.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}

	diags := ValueFrom(ctx, profileResponse, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

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

	unions, diags := types.ListValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	eventObjectValue, diags := EventValueFrom(ctx, profile.Birth)
	d.Append(diags...)
	profileModel.Birth = eventObjectValue

	if profile.CreatedAt != "" {
		profileModel.CreatedAt = types.StringValue(profile.CreatedAt)
	}

	return d
}

func EventValueFrom(ctx context.Context, birth *geni.EventElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if birth != nil {
		var d diag.Diagnostics
		dateObjectValue, diags := DateValueFrom(ctx, birth.Date)
		d.Append(diags...)

		locationObjectValue, diags := LocationValueFrom(ctx, birth.Location)
		d.Append(diags...)

		eventModel := event.Model{
			Description: types.StringValue(birth.Description),
			Name:        types.StringValue(birth.Name),
			Date:        dateObjectValue,
			Location:    locationObjectValue,
		}

		eventObjectValue, diags := types.ObjectValueFrom(ctx, eventModel.AttributeTypes(), eventModel)
		d.Append(diags...)

		return eventObjectValue, d
	}

	return types.ObjectNull(event.EventModelAttributeTypes()), nil
}

func DateValueFrom(ctx context.Context, dateElement *geni.DateElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if dateElement != nil {
		dateModel := event.DateModel{
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

	return types.ObjectNull(event.DateModelAttributeTypes()), nil
}

func LocationValueFrom(ctx context.Context, location *geni.LocationElement) (basetypes.ObjectValue, diag.Diagnostics) {
	if location != nil {
		locationModel := event.LocationModel{
			City:           types.StringValue(location.City),
			Country:        types.StringValue(location.Country),
			County:         types.StringValue(location.County),
			Latitude:       types.NumberValue(location.Latitude),
			Longitude:      types.NumberValue(location.Longitude),
			PlaceName:      types.StringValue(location.PlaceName),
			State:          types.StringValue(location.State),
			StreetAddress1: types.StringValue(location.StreetAddress1),
			StreetAddress2: types.StringValue(location.StreetAddress2),
			StreetAddress3: types.StringValue(location.StreetAddress3),
		}

		return types.ObjectValueFrom(ctx, locationModel.AttributeTypes(), locationModel)
	}

	return types.ObjectNull(event.LocationModelAttributeTypes()), nil
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := geni.UpdateProfile(r.accessToken.ValueString(), plan.ID.ValueString(), &geni.ProfileRequest{
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
		Gender:    plan.Gender.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	plan.CreatedAt = types.StringValue(response.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := geni.DeleteProfile(r.accessToken.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting profile", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
