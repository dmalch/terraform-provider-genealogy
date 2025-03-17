package union

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-geni/internal/config"
	"github.com/dmalch/terraform-provider-geni/internal/geni"
)

type Resource struct {
	resource.ResourceWithConfigure
	accessToken types.String
}

func NewUnionResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "geni_union"
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

type ResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Children types.List   `tfsdk:"children"`
	Partners types.List   `tfsdk:"partners"`
}

// Create creates the resource
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If there are two partners, we can create a union by calling the profile/add-partner API
	if len(plan.Partners.Elements()) == 2 {
		// It is impossible to create a union using the API, so we need to create a
		// temporary partner profile and then merge it with the second partner.

		partnerIds := make([]types.String, 0, len(plan.Partners.Elements()))
		diag := plan.Partners.ElementsAs(ctx, &partnerIds, false)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		tmpProfile, err := geni.AddPartner(r.accessToken.ValueString(), partnerIds[0].ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error adding partner", err.Error())
			return
		}

		// Merge the temporary profile with the second partner
		if err := geni.MergeProfiles(r.accessToken.ValueString(), partnerIds[1].ValueString(), tmpProfile.Id); err != nil {
			resp.Diagnostics.AddError("Error merging profiles", err.Error())
			return
		}

		plan.ID = types.StringValue(tmpProfile.Unions[0])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read reads the resource
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	union, err := geni.GetUnion(r.accessToken.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading union", err.Error())
		return
	}

	if union.Id != "" {
		state.ID = types.StringValue(union.Id)
	}
	if len(union.Children) > 0 {
		listValue, diag := types.ListValueFrom(ctx, types.StringType, union.Children)
		state.Children = listValue
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if len(union.Partners) > 0 {
		listValue, diag := types.ListValueFrom(ctx, types.StringType, union.Partners)
		state.Partners = listValue
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
