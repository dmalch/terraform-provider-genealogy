package event

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func Schema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"date": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					// Range (before, after, or between)
					"range": schema.StringAttribute{
						Optional: true},
					// Indicates whether the date is an approximation
					"circa": schema.BoolAttribute{
						Optional: true},
					// Date's day
					"day": schema.Int32Attribute{
						Optional: true},
					// Date's month
					"month": schema.Int32Attribute{
						Optional: true},
					// Date's year
					"year": schema.Int32Attribute{
						Optional: true},
					// Indicates whether the end date is an approximation
					"end_circa": schema.BoolAttribute{
						Optional: true},
					// Date's end day (only valid if range is between)
					"end_day": schema.Int32Attribute{
						Optional: true},
					// Date's end month (only valid if range is between)
					"end_month": schema.Int32Attribute{
						Optional: true},
					// Date's end year (only valid if range is between)
					"end_year": schema.Int32Attribute{
						Optional: true},
				},
			},
			"location": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"city": schema.StringAttribute{
						Optional: true,
					},
					"country": schema.StringAttribute{
						Optional: true,
					},
					"county": schema.StringAttribute{
						Optional: true,
					},
					"latitude": schema.NumberAttribute{
						Optional: true,
					},
					"longitude": schema.NumberAttribute{
						Optional: true,
					},
					"place_name": schema.StringAttribute{
						Optional: true,
					},
					"state": schema.StringAttribute{
						Optional: true,
					},
					"street_address1": schema.StringAttribute{
						Optional: true,
					},
					"street_address2": schema.StringAttribute{
						Optional: true,
					},
					"street_address3": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}
