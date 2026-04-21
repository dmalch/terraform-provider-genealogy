package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ValueFrom(ctx context.Context, union *geni.UnionResponse, unionModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	if union.Id != "" {
		unionModel.ID = types.StringValue(union.Id)
	}

	tagged := make(map[string]struct{}, len(union.FosterChildren)+len(union.AdoptedChildren))
	for _, id := range union.FosterChildren {
		tagged[id] = struct{}{}
	}
	for _, id := range union.AdoptedChildren {
		tagged[id] = struct{}{}
	}

	if len(union.Children) > 0 {
		biological := make([]string, 0, len(union.Children))
		for _, id := range union.Children {
			if _, isTagged := tagged[id]; !isTagged {
				biological = append(biological, id)
			}
		}
		if len(biological) > 0 {
			children, diags := types.SetValueFrom(ctx, types.StringType, biological)
			d.Append(diags...)
			unionModel.Children = children
		}
	}

	if len(union.FosterChildren) > 0 {
		foster, diags := types.SetValueFrom(ctx, types.StringType, union.FosterChildren)
		d.Append(diags...)
		unionModel.FosterChildren = foster
	}

	if len(union.AdoptedChildren) > 0 {
		adopted, diags := types.SetValueFrom(ctx, types.StringType, union.AdoptedChildren)
		d.Append(diags...)
		unionModel.AdoptedChildren = adopted
	}

	if len(union.Partners) > 0 {
		partners, diags := types.SetValueFrom(ctx, types.StringType, union.Partners)
		d.Append(diags...)
		unionModel.Partners = partners
	}

	marriage, diags := event.ValueFrom(ctx, union.Marriage)
	d.Append(diags...)
	unionModel.Marriage = marriage

	divorce, diags := event.ValueFrom(ctx, union.Divorce)
	d.Append(diags...)
	unionModel.Divorce = divorce

	return d
}

func UpdateComputedFields(ctx context.Context, union *geni.UnionResponse, unionModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	unionModel.ID = types.StringValue(union.Id)

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

func RequestFrom(ctx context.Context, plan ResourceModel) (*geni.UnionRequest, diag.Diagnostics) {
	var d diag.Diagnostics

	marriage, diags := event.ElementFrom(ctx, plan.Marriage)
	d.Append(diags...)

	divorce, diags := event.ElementFrom(ctx, plan.Divorce)
	d.Append(diags...)

	unionRequest := geni.UnionRequest{
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

func convertToSlice(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	slice := make([]string, 0, len(set.Elements()))
	diags := set.ElementsAs(ctx, &slice, false)
	return slice, diags
}
