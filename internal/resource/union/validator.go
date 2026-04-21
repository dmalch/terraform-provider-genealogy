package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(data.Partners.Elements()) > 2 {
		resp.Diagnostics.AddAttributeError(path.Root(fieldPartners),
			"Too Many Partners",
			"Only 2 partners are allowed in the union. Please remove any extra partners.",
		)
	}

	totalProfiles := len(data.Partners.Elements()) +
		len(data.Children.Elements()) +
		len(data.FosterChildren.Elements()) +
		len(data.AdoptedChildren.Elements())
	if totalProfiles < 2 {
		resp.Diagnostics.AddAttributeError(path.Root(fieldPartners),
			"Insufficient Attribute Configuration",
			"At least two profiles must be configured in either partners or children. "+
				"Please ensure that the resource has the required profiles to function correctly.",
		)
	}

	childSets := []struct {
		field string
		root  path.Path
		set   types.Set
	}{
		{fieldChildren, path.Root(fieldChildren), data.Children},
		{fieldFosterChildren, path.Root(fieldFosterChildren), data.FosterChildren},
		{fieldAdoptedChildren, path.Root(fieldAdoptedChildren), data.AdoptedChildren},
	}
	for i := 0; i < len(childSets); i++ {
		left := setIDs(childSets[i].set)
		for j := i + 1; j < len(childSets); j++ {
			right := setIDs(childSets[j].set)
			if id := firstCommon(left, right); id != "" {
				resp.Diagnostics.AddAttributeError(childSets[j].root,
					"Overlapping Child Sets",
					"Profile "+id+" appears in both "+childSets[i].field+" and "+childSets[j].field+". "+
						"Each child must belong to exactly one relationship category.",
				)
			}
		}
	}
}

func setIDs(s types.Set) map[string]struct{} {
	out := make(map[string]struct{}, len(s.Elements()))
	for _, e := range s.Elements() {
		if str, ok := e.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
			out[str.ValueString()] = struct{}{}
		}
	}
	return out
}

func firstCommon(a, b map[string]struct{}) string {
	for id := range a {
		if _, ok := b[id]; ok {
			return id
		}
	}
	return ""
}
