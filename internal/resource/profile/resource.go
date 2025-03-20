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

	birth, diags := event.ElementFrom(ctx, plan.Birth)
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

	diags = updateComputedFields(ctx, &plan, profileResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func updateComputedFields(ctx context.Context, profileModel *ResourceModel, profile *geni.ProfileResponse) diag.Diagnostics {
	var d diag.Diagnostics

	profileModel.ID = types.StringValue(profile.Id)

	unions, diags := types.ListValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	eventObject := profileModel.Birth
	if !eventObject.IsNull() && !eventObject.IsUnknown() && profile.Birth != nil {
		var eventModel event.Model

		diags := eventObject.As(ctx, &eventModel, basetypes.ObjectAsOptions{})
		d.Append(diags...)

		diags = updateComputedFieldsInEvent(ctx, &eventModel, profile.Birth)
		d.Append(diags...)

		eventObject, diags = types.ObjectValueFrom(ctx, eventModel.AttributeTypes(), eventModel)
		d.Append(diags...)

		profileModel.Birth = eventObject
	}

	profileModel.CreatedAt = types.StringValue(profile.CreatedAt)

	return d
}

func updateComputedFieldsInEvent(_ context.Context, eventObjectValue *event.Model, eventElement *geni.EventElement) diag.Diagnostics {
	var d diag.Diagnostics

	if eventObjectValue.Name.IsNull() || eventObjectValue.Name.IsUnknown() {
		eventObjectValue.Name = types.StringValue(eventElement.Name)
	}
	if eventObjectValue.Description.IsNull() || eventObjectValue.Description.IsUnknown() {
		eventObjectValue.Description = types.StringValue(eventElement.Description)
	}

	return d
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
