package event

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var rangePath = path.MatchRelative().AtParent().AtName("range")

type SchemaOptions struct {
	NameComputed        bool
	DescriptionComputed bool
	Description         string
}

func Schema(opts ...SchemaOptions) schema.SingleNestedAttribute {
	// Get first opt or default
	opt := SchemaOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    opt.NameComputed,
				Description: "Event's name.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    opt.DescriptionComputed,
				Description: "Event's description.",
			},
			"date":     DateSchema(""),
			"location": LocationSchema(""),
		},
		Description: opt.Description,
	}
}

func LocationSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"city": schema.StringAttribute{
				Optional:    true,
				Description: "City name.",
			},
			"country": schema.StringAttribute{
				Optional:    true,
				Description: "Country name.",
			},
			"county": schema.StringAttribute{
				Optional:    true,
				Description: "County name.",
			},
			"latitude": schema.NumberAttribute{
				Optional:    true,
				Description: "Latitude coordinate.",
			},
			"longitude": schema.NumberAttribute{
				Optional:    true,
				Description: "Longitude coordinate.",
			},
			"place_name": schema.StringAttribute{
				Optional:    true,
				Description: "Place name.",
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Description: "State name.",
			},
			"street_address1": schema.StringAttribute{
				Optional:    true,
				Description: "First line of the street address.",
			},
			"street_address2": schema.StringAttribute{
				Optional:    true,
				Description: "Second line of the street address.",
			},
			"street_address3": schema.StringAttribute{
				Optional:    true,
				Description: "Third line of the street address.",
			},
		},
		Description: description,
	}
}

func DateSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"range": schema.StringAttribute{
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("before", "after", "between")},
				Description: "Range (before, after, or between).",
			},
			"circa": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the date is an approximation.",
			},
			"day": schema.Int32Attribute{
				Optional:    true,
				Description: "Day of the month.",
			},
			"month": schema.Int32Attribute{
				Optional:    true,
				Description: "Month of the year.",
			},
			"year": schema.Int32Attribute{
				Optional:    true,
				Description: "Date's year.",
			},
			"end_circa": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the end date is an approximation.",
			},
			"end_day": schema.Int32Attribute{
				Optional:    true,
				Validators:  []validator.Int32{int32validator.AlsoRequires(rangePath)},
				Description: "Date's end day (only valid if range is between).",
			},
			"end_month": schema.Int32Attribute{
				Optional:    true,
				Validators:  []validator.Int32{int32validator.AlsoRequires(rangePath)},
				Description: "Date's end month (only valid if range is between).",
			},
			"end_year": schema.Int32Attribute{
				Optional:    true,
				Validators:  []validator.Int32{int32validator.AlsoRequires(rangePath)},
				Description: "Date's end year (only valid if range is between).",
			},
		},
		Description: description,
	}
}
