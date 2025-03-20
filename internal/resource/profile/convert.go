package profile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-geni/internal/geni"
	"github.com/dmalch/terraform-provider-geni/internal/resource/event"
)

func ValueFrom(ctx context.Context, profile *geni.ProfileResponse, profileModel *ResourceModel) diag.Diagnostics {
	var d diag.Diagnostics

	if profile.Id != "" {
		profileModel.ID = types.StringValue(profile.Id)
	}

	if profile.FirstName != "" {
		profileModel.FirstName = types.StringValue(profile.FirstName)
	}

	if profile.LastName != "" {
		profileModel.LastName = types.StringValue(profile.LastName)
	}

	if profile.Gender != "" {
		profileModel.Gender = types.StringValue(profile.Gender)
	}

	unions, diags := types.ListValueFrom(ctx, types.StringType, profile.Unions)
	d.Append(diags...)
	profileModel.Unions = unions

	eventObjectValue, diags := event.ValueFrom(ctx, profile.Birth)
	d.Append(diags...)
	profileModel.Birth = eventObjectValue

	if profile.CreatedAt != "" {
		profileModel.CreatedAt = types.StringValue(profile.CreatedAt)
	}

	return d
}

func ElementFrom(ctx context.Context, plan ResourceModel) (*geni.ProfileRequest, diag.Diagnostics) {
	var d diag.Diagnostics

	birth, diags := event.ElementFrom(ctx, plan.Birth)
	d.Append(diags...)

	profileRequest := &geni.ProfileRequest{
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
		Gender:    plan.Gender.ValueString(),
		Birth:     birth,
	}

	return profileRequest, d
}
