package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func ValueFrom(_ context.Context, response *geni.ProjectResponse, model *Model) diag.Diagnostics {
	var d diag.Diagnostics

	model.ID = types.StringValue(response.Id)
	model.Name = types.StringValue(response.Name)
	model.Description = types.StringPointerValue(response.Description)
	model.UpdatedAt = types.StringValue(response.UpdatedAt)
	model.CreatedAt = types.StringValue(response.CreatedAt)

	return d
}
