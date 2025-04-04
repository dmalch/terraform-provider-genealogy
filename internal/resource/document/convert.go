package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ValueFrom(ctx context.Context, response *geni.DocumentResponse, model *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	model.ID = types.StringValue(response.Id)
	model.Title = types.StringValue(response.Title)
	model.Description = types.StringValue(response.Description)
	model.ContentType = types.StringValue(response.ContentType)

	dateObjectValue, diags := event.DateValueFrom(ctx, response.Date)
	d.Append(diags...)
	model.Date = dateObjectValue

	locationObjectValue, diags := event.LocationValueFrom(ctx, response.Location)
	d.Append(diags...)
	model.Location = locationObjectValue

	tags, diags := types.SetValueFrom(ctx, types.StringType, response.Tags)
	d.Append(diags...)
	model.Profiles = tags

	labels, diags := types.SetValueFrom(ctx, types.StringType, response.Labels)
	d.Append(diags...)
	model.Labels = labels

	model.CreatedAt = types.StringValue(response.CreatedAt)

	return d
}
