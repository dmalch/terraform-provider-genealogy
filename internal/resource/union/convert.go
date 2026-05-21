package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	geniunion "github.com/dmalch/go-geni/union"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ValueFrom(ctx context.Context, union *geniunion.Union, unionModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	if union.ID != "" {
		unionModel.ID = types.StringValue(union.ID)
	}

	tagged := make(map[string]struct{}, len(union.FosterChildren)+len(union.AdoptedChildren))
	for _, id := range union.FosterChildren {
		tagged[id] = struct{}{}
	}
	for _, id := range union.AdoptedChildren {
		tagged[id] = struct{}{}
	}

	// Assign every collection unconditionally so Read reflects collections
	// that have drained to empty on Geni. Skipping the assignment for an
	// empty API list would leave the stale prior-state value in place and
	// produce a permanent phantom diff on every plan (issue #106).
	biological := make([]string, 0, len(union.Children))
	for _, id := range union.Children {
		if _, isTagged := tagged[id]; !isTagged {
			biological = append(biological, id)
		}
	}
	unionModel.Children = setOrNull(ctx, biological, &d)
	unionModel.FosterChildren = setOrNull(ctx, union.FosterChildren, &d)
	unionModel.AdoptedChildren = setOrNull(ctx, union.AdoptedChildren, &d)
	unionModel.Partners = setOrNull(ctx, union.Partners, &d)

	marriage, diags := event.ValueFrom(ctx, union.Marriage)
	d.Append(diags...)
	unionModel.Marriage = marriage

	divorce, diags := event.ValueFrom(ctx, union.Divorce)
	d.Append(diags...)
	unionModel.Divorce = divorce

	return d
}

func UpdateComputedFields(ctx context.Context, union *geniunion.Union, unionModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	unionModel.ID = types.StringValue(union.ID)

	if union.Marriage != nil {
		marriage, diags := event.UpdateComputedFieldsInEvent(ctx, unionModel.Marriage, union.Marriage)
		d.Append(diags...)
		unionModel.Marriage = marriage
	}

	if union.Divorce != nil {
		divorce, diags := event.UpdateComputedFieldsInEvent(ctx, unionModel.Divorce, union.Divorce)
		d.Append(diags...)
		unionModel.Divorce = divorce
	}

	return d
}

func RequestFrom(ctx context.Context, plan ResourceModel) (*geniunion.Request, diag.Diagnostics) {
	var d diag.Diagnostics

	marriage, diags := event.ElementFrom(ctx, plan.Marriage)
	d.Append(diags...)

	divorce, diags := event.ElementFrom(ctx, plan.Divorce)
	d.Append(diags...)

	unionRequest := geniunion.Request{
		Marriage: marriage,
		Divorce:  divorce,
	}

	return &unionRequest, d
}

type childModifierChange struct {
	id   string
	from string
	to   string
}

// childrenWithChangedModifier returns any child ids that exist in both
// plan and state but moved between the biological / foster / adopted
// buckets. Geni's public API has no endpoint to change a child's
// relationship modifier on an existing edge, so these changes cannot be
// applied and the caller should surface a warning.
func childrenWithChangedModifier(ctx context.Context, plan, state ResourceModel) []childModifierChange {
	planLabel := labelsFor(ctx, plan)
	stateLabel := labelsFor(ctx, state)

	var out []childModifierChange
	for id, to := range planLabel {
		if from, ok := stateLabel[id]; ok && from != to {
			out = append(out, childModifierChange{id: id, from: from, to: to})
		}
	}
	return out
}

func labelsFor(ctx context.Context, m ResourceModel) map[string]string {
	labels := make(map[string]string)
	add := func(set types.Set, label string) {
		ids, _ := convertToSlice(ctx, set)
		for _, id := range ids {
			labels[id] = label
		}
	}
	add(m.Children, "biological")
	add(m.FosterChildren, "foster")
	add(m.AdoptedChildren, "adopted")
	return labels
}

// modifierFor returns the relationship_modifier value to send when adding
// childId to a union. "adopt" and "foster" correspond to the matching
// subset memberships; biological children return "".
func modifierFor(childId string, fosterChildren, adoptedChildren map[string]struct{}) string {
	if _, ok := fosterChildren[childId]; ok {
		return "foster"
	}
	if _, ok := adoptedChildren[childId]; ok {
		return "adopt"
	}
	return ""
}

func hashMapFrom(slice []string) map[string]struct{} {
	hashMap := make(map[string]struct{}, len(slice))
	for _, elem := range slice {
		hashMap[elem] = struct{}{}
	}
	return hashMap
}

// setOrNull builds a string Set from ids, or returns a typed null Set when ids
// is empty so that Read reflects collections that drained to empty on Geni.
func setOrNull(ctx context.Context, ids []string, d *diag.Diagnostics) types.Set {
	if len(ids) == 0 {
		return types.SetNull(types.StringType)
	}
	set, diags := types.SetValueFrom(ctx, types.StringType, ids)
	d.Append(diags...)
	return set
}

func convertToSlice(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	slice := make([]string, 0, len(set.Elements()))
	diags := set.ElementsAs(ctx, &slice, false)
	return slice, diags
}
