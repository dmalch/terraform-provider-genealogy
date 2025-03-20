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

	if len(union.Children) > 0 {
		children, diags := types.SetValueFrom(ctx, types.StringType, union.Children)
		d.Append(diags...)
		unionModel.Children = children
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
