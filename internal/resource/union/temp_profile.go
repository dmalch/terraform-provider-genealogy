package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	geniprofile "github.com/dmalch/go-geni/profile"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

// addAndMerge creates a temporary Geni profile via create and merges it into the
// existing profile realID. Geni has no API to link two existing profiles into a
// union, so every union edge is built create-then-merge.
//
// If the merge fails the temporary profile is an orphan — live on Geni but
// untracked by Terraform — so addAndMerge best-effort deletes it. The temp
// profile is returned even on merge failure (alongside the error) so the caller
// can still recover the union id from tmp.Unions for partial-state persistence.
func (r *Resource) addAndMerge(
	ctx context.Context,
	realID string,
	create func(context.Context) (*geniprofile.Profile, error),
) (*geniprofile.Profile, error) {
	tmp, err := create(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := r.client.Profile().Merge(ctx, realID, tmp.ID); err != nil {
		if delErr := r.client.Profile().Delete(ctx, tmp.ID); delErr != nil {
			tflog.Error(ctx, "failed to delete orphan temporary profile after a failed merge",
				map[string]interface{}{"profile_id": tmp.ID, "error": delErr})
		}
		return tmp, err
	}

	return tmp, nil
}

// unionIDFrom returns the union id once it is known: the current value if it is
// already set, otherwise the union the temporary profile was created in.
func unionIDFrom(current types.String, tmp *geniprofile.Profile) types.String {
	if !current.IsUnknown() && !current.IsNull() {
		return current
	}
	if tmp != nil && len(tmp.Unions) > 0 {
		return types.StringValue(tmp.Unions[0])
	}
	return current
}

// persistPartialUnion writes a union that was created but not fully configured
// into state, so a failed Create does not strand it untracked — which would
// otherwise make the next apply create yet another union. Marriage and Divorce
// are nulled because their computed fields are still unresolved on the failure
// path; the next Read repopulates them from the API.
func persistPartialUnion(ctx context.Context, resp *resource.CreateResponse, plan ResourceModel) {
	if plan.ID.IsUnknown() || plan.ID.IsNull() {
		return // no union was created — nothing to track
	}

	plan.Marriage = types.ObjectNull(event.EventModelAttributeTypes())
	plan.Divorce = types.ObjectNull(event.EventModelAttributeTypes())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	resp.Diagnostics.Append(resp.Identity.Set(ctx, ResourceIdentityModel{ID: plan.ID})...)
}
