package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/config"
	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

type Resource struct {
	resource.ResourceWithConfigure
	client *geni.Client
}

func NewProfileResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name.
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "geni_profile"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.ClientData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.ClientData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

// Create creates the resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileRequest, diags := RequestFrom(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := r.client.CreateProfile(ctx, profileRequest)
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

	profileModel.FirstName = types.StringValue(profile.FirstName)
	profileModel.LastName = types.StringValue(profile.LastName)
	profileModel.MiddleName = types.StringValue(profile.MiddleName)
	profileModel.MaidenName = types.StringValue(profile.MaidenName)

	unions, diags := types.ListValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	if profile.Birth != nil {
		birth, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Birth, profile.Birth)
		d.Append(diags...)
		profileModel.Birth = birth
	}

	if profile.Baptism != nil {
		baptism, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Baptism, profile.Baptism)
		d.Append(diags...)
		profileModel.Baptism = baptism
	}

	if profile.Death != nil {
		death, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Death, profile.Death)
		d.Append(diags...)
		profileModel.Death = death
	}

	if profile.Burial != nil {
		burial, diags := event.UpdateComputedFieldsInEvent(ctx, profileModel.Burial, profile.Burial)
		d.Append(diags...)
		profileModel.Burial = burial
	}

	profileModel.CreatedAt = types.StringValue(profile.CreatedAt)

	return d
}

// Read reads the resource.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := r.client.GetProfile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}

	plan := state
	diags := ValueFrom(ctx, profileResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If names in the plan are empty, and the names in the state contain one element
	// for en-US, then use the state names.
	if len(plan.Names.Elements()) == 0 {
		// Get names from the current state
		names, diags := NameModelsFrom(ctx, state.Names)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if _, ok := names["en-US"]; ok && len(names) == 1 {
			plan.Names = state.Names
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileRequest, diags := RequestFrom(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profileResponse, err := r.client.UpdateProfile(ctx, plan.ID.ValueString(), profileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	diags = updateComputedFields(ctx, &plan, profileResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProfile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting profile", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
