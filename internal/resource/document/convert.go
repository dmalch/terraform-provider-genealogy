package document

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/resource/event"
)

func ValueFrom(ctx context.Context, response *geni.DocumentResponse, model *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	model.ID = types.StringValue(response.Id)
	model.Title = types.StringValue(response.Title)
	model.Description = types.StringPointerValue(response.Description)
	model.ContentType = types.StringPointerValue(response.ContentType)

	dateObjectValue, diags := event.DateValueFrom(ctx, response.Date)
	d.Append(diags...)
	model.Date = dateObjectValue

	locationObjectValue, diags := event.LocationValueFrom(ctx, response.Location)
	d.Append(diags...)
	model.Location = locationObjectValue

	tags, diags := types.SetValueFrom(ctx, types.StringType, response.Tags)
	d.Append(diags...)
	model.Profiles = tags

	labels, diags := types.SetValueFrom(ctx, types.StringType, filterOutDuplicateLabelsFrom(response.Labels))
	d.Append(diags...)
	model.Labels = labels

	model.CreatedAt = types.StringValue(response.CreatedAt)

	return d
}

func RequestFrom(ctx context.Context, resourceModel ResourceModel) (*geni.DocumentRequest, diag.Diagnostics) {
	var d diag.Diagnostics

	dateModel, diags := event.DateObjectValueFrom(ctx, resourceModel.Date)
	d.Append(diags...)

	locationModel, diags := event.LocationObjectValueFrom(ctx, resourceModel.Location)
	d.Append(diags...)

	labelModels, diags := LabelModelsFrom(ctx, resourceModel.Labels)
	d.Append(diags...)

	var labels *string

	if len(labelModels) != 0 {
		join := strings.Join(labelModels, ",")
		labels = &join
	}

	documentRequest := &geni.DocumentRequest{
		Title:       resourceModel.Title.ValueString(),
		Description: resourceModel.Description.ValueStringPointer(),
		ContentType: resourceModel.ContentType.ValueStringPointer(),
		Text:        resourceModel.Text.ValueStringPointer(),
		Date:        event.DateElementFrom(dateModel),
		Location:    event.LocationElementFrom(locationModel),
		Labels:      labels,
	}

	return documentRequest, d
}

func LabelModelsFrom(ctx context.Context, labels types.Set) ([]string, diag.Diagnostics) {
	if len(labels.Elements()) == 0 {
		return nil, diag.Diagnostics{}
	}

	var labelModels = make([]string, len(labels.Elements()))
	diags := labels.ElementsAs(ctx, &labelModels, false)

	return labelModels, diags
}

func UpdateComputedFields(ctx context.Context, response *geni.DocumentResponse, resourceModel *ResourceModel) diag.Diagnostics {
	d := diag.Diagnostics{}

	resourceModel.ID = types.StringValue(response.Id)
	resourceModel.ContentType = types.StringPointerValue(response.ContentType)

	tags, diags := types.SetValueFrom(ctx, types.StringType, response.Tags)
	d.Append(diags...)
	resourceModel.Profiles = tags

	// Filter out duplicate labels
	labels, diags := types.SetValueFrom(ctx, types.StringType, filterOutDuplicateLabelsFrom(response.Labels))
	d.Append(diags...)
	resourceModel.Labels = labels

	resourceModel.CreatedAt = types.StringValue(response.CreatedAt)
	return d
}

func filterOutDuplicateLabelsFrom(res []string) []string {
	uniqueLabels := make([]string, 0, len(res))
	seen := make(map[string]struct{})
	for _, label := range res {
		if _, ok := seen[label]; !ok {
			seen[label] = struct{}{}
			uniqueLabels = append(uniqueLabels, label)
		}
	}
	return uniqueLabels
}
