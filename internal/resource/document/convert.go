package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func ValueFrom(ctx context.Context, response *geni.DocumentResponse, model *ResourceModel) diag.Diagnostics {
	return diag.Diagnostics{}
}
