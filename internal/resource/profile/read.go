package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

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

	newState := state
	diags := ValueFrom(ctx, profileResponse, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If names in the new state are empty, and the names in the state contain one
	// element for en-US, then use the state names.
	if len(newState.Names.Elements()) == 0 {
		// Get names from the current state
		names, diags := NameModelsFrom(ctx, state.Names)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if _, ok := names["en-US"]; ok && len(names) == 1 {
			newState.Names = state.Names
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
