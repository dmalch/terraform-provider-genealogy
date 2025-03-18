package union

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If attribute_one is not configured, return without warning.
	if len(data.Children.Elements())+len(data.Partners.Elements()) < 2 {
		resp.Diagnostics.AddAttributeError(path.Root(fieldPartners),
			"Insufficient Attribute Configuration",
			"At least two profiles must be configured in either partners or children. "+
				"Please ensure that the resource has the required profiles to function correctly.",
		)
	}
}
