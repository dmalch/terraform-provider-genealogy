package event

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func Schema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{},
			"name":        schema.StringAttribute{},
		},
		Blocks: map[string]schema.Block{
			"date": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"circa":     schema.BoolAttribute{},
					"day":       schema.NumberAttribute{},
					"end_circa": schema.BoolAttribute{},
					"end_day":   schema.NumberAttribute{},
					"end_month": schema.NumberAttribute{},
					"end_year":  schema.NumberAttribute{},
					"month":     schema.NumberAttribute{},
					"range":     schema.StringAttribute{},
					"year":      schema.NumberAttribute{},
				},
			},
			"location": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"city":            schema.StringAttribute{},
					"country":         schema.StringAttribute{},
					"county":          schema.StringAttribute{},
					"latitude":        schema.NumberAttribute{},
					"longitude":       schema.NumberAttribute{},
					"place_name":      schema.StringAttribute{},
					"state":           schema.StringAttribute{},
					"street_address1": schema.StringAttribute{},
					"street_address2": schema.StringAttribute{},
					"street_address3": schema.StringAttribute{},
				},
			},
		},
	}
}
