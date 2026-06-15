package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// currentUnionMembers fetches unionId's live membership so callers can add only
// the edges that are genuinely missing. Geni auto-merges unions with identical
// partner sets server-side; if the resource's union was absorbed (no partners
// and no children survive on it), this remaps to the surviving union discovered
// from the partner set — the same signal Read uses (findExistingUnionForPartners).
//
// The returned resolvedID is what add-calls should target. Callers MUST keep
// their state `id` untouched: rewriting it to resolvedID would violate
// Terraform's plan/apply consistency contract on `.id` (it is pinned by
// UseStateForUnknown). The persisted remap is left to the next Read, exactly as
// it happens today (#138).
func (r *Resource) currentUnionMembers(ctx context.Context, unionID string, statePartners types.Set) (resolvedID string, partners, children map[string]struct{}, diags diag.Diagnostics) {
	resolvedID = unionID

	live, err := r.batchClient.GetUnion(ctx, unionID)
	if err != nil {
		diags.AddError("Error reading union", err.Error())
		return resolvedID, nil, nil, diags
	}

	// The union was auto-merged away — find the surviving canonical union from
	// the partner set and target that for the add-calls instead.
	if len(live.Partners) == 0 && len(live.Children) == 0 {
		canonicalID, d := r.findExistingUnionForPartners(ctx, statePartners)
		diags.Append(d...)
		if diags.HasError() {
			return resolvedID, nil, nil, diags
		}

		if canonicalID != "" && canonicalID != unionID {
			diags.AddWarning("Found existing union",
				"The union in the state has no partners and children. Adding members to the existing union with ID "+canonicalID+" instead.")
			resolvedID = canonicalID
			live, err = r.client.Union().Get(ctx, canonicalID)
			if err != nil {
				diags.AddError("Error reading union", err.Error())
				return resolvedID, nil, nil, diags
			}
		}
	}

	partners = make(map[string]struct{}, len(live.Partners))
	for _, p := range live.Partners {
		partners[p] = struct{}{}
	}

	// Geni reports every child (biological, foster, adopted) in Children, with
	// FosterChildren/AdoptedChildren as labelled subsets, so Children alone is
	// the complete add-target set.
	children = make(map[string]struct{}, len(live.Children))
	for _, c := range live.Children {
		children[c] = struct{}{}
	}

	return resolvedID, partners, children, diags
}

// missingEdges returns the planned ids that are not already present in current,
// preserving the order of planned.
func missingEdges(planned []string, current map[string]struct{}) []string {
	var missing []string
	for _, id := range planned {
		if _, ok := current[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing
}
