package project

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

// Read reads the resource.
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	documentResponse, err := d.client.GetProject(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, geni.ErrResourceNotFound) {
			resp.Diagnostics.AddWarning("Project not found", "The project was not found in the Geni API.")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	diags := ValueFrom(ctx, documentResponse, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
