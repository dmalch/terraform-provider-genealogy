package profile

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	profileIdFormat = regexp.MustCompile(`^profile-\d+$`)
	guidFormat      = regexp.MustCompile(`^[a-f0-9]+$`)
)

func (d *DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	exactlyOne := []validator.String{
		stringvalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("guid"),
		),
	}

	resp.Schema = schema.Schema{
		Description: "Look up a single Geni profile by `id` or `guid`. Exactly one must be set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Validators:  append(exactlyOne, stringvalidator.RegexMatches(profileIdFormat, "must be in the format profile-1")),
				Description: "The unique identifier for the profile. This is a string that starts with 'profile-' followed by a number.",
			},
			"guid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Validators:  append(exactlyOne, stringvalidator.RegexMatches(guidFormat, "must be a lowercase hexadecimal GUID")),
				Description: "The globally unique identifier (GUID) for the profile, as assigned by Geni.",
			},
			"gender": schema.StringAttribute{
				Computed:    true,
				Description: "Profile's gender.",
			},
			"title": schema.StringAttribute{
				Computed:    true,
				Description: "Profile's name title (e.g. \"Dr.\", \"Sir\").",
			},
			"suffix": schema.StringAttribute{
				Computed:    true,
				Description: "Profile's name suffix (e.g. \"Jr.\", \"III\").",
			},
			"occupation": schema.StringAttribute{
				Computed:    true,
				Description: "Profile's occupation.",
			},
			"names": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_name":      schema.StringAttribute{Computed: true, Description: "The first name of the person."},
						"middle_name":     schema.StringAttribute{Computed: true, Description: "The middle name of the person."},
						"last_name":       schema.StringAttribute{Computed: true, Description: "The last name of the person."},
						"birth_last_name": schema.StringAttribute{Computed: true, Description: "The birth last name of the person."},
						"display_name":    schema.StringAttribute{Computed: true, Description: "The display name of the person."},
						"nicknames": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The nicknames of the person.",
						},
					},
				},
				Description: "Nested map of locale → name fields.",
			},
			"unions": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of union IDs the profile belongs to.",
			},
			"projects": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of project IDs the profile is a member of.",
			},
			"birth":             eventSchema("Birth event information."),
			"baptism":           eventSchema("Baptism event information."),
			"death":             eventSchema("Death event information."),
			"burial":            eventSchema("Burial event information."),
			"cause_of_death":    schema.StringAttribute{Computed: true, Description: "Profile's death cause."},
			"current_residence": locationSchema("Profile's current residence."),
			"about": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Profile's about-me section, keyed by locale.",
			},
			"public":      schema.BoolAttribute{Computed: true, Description: "Profile's public visibility."},
			"alive":       schema.BoolAttribute{Computed: true, Description: "Profile's alive status."},
			"deleted":     schema.BoolAttribute{Computed: true, Description: "Profile's deleted status."},
			"merged_into": schema.StringAttribute{Computed: true, Description: "The ID of the profile this profile was merged into, if any."},
			"created_at":  schema.StringAttribute{Computed: true, Description: "The Unix epoch time in seconds when the profile was created."},
		},
	}
}

func eventSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"name":        schema.StringAttribute{Computed: true, Description: "Event's name."},
			"description": schema.StringAttribute{Computed: true, Description: "Event's description."},
			"date":        dateRangeSchema("Event's date."),
			"location":    locationSchema("Event's location."),
		},
		Description: description,
	}
}

func locationSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"city":            schema.StringAttribute{Computed: true, Description: "City name."},
			"country":         schema.StringAttribute{Computed: true, Description: "Country name."},
			"county":          schema.StringAttribute{Computed: true, Description: "County name."},
			"latitude":        schema.Float64Attribute{Computed: true, Description: "Latitude coordinate."},
			"longitude":       schema.Float64Attribute{Computed: true, Description: "Longitude coordinate."},
			"place_name":      schema.StringAttribute{Computed: true, Description: "Place name."},
			"state":           schema.StringAttribute{Computed: true, Description: "State name."},
			"street_address1": schema.StringAttribute{Computed: true, Description: "First line of the street address."},
			"street_address2": schema.StringAttribute{Computed: true, Description: "Second line of the street address."},
			"street_address3": schema.StringAttribute{Computed: true, Description: "Third line of the street address."},
		},
		Description: description,
	}
}

func dateRangeSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"range":     schema.StringAttribute{Computed: true, Description: "Range (before, after, or between)."},
			"circa":     schema.BoolAttribute{Computed: true, Description: "Indicates whether the date is an approximation."},
			"day":       schema.Int32Attribute{Computed: true, Description: "Day of the month."},
			"month":     schema.Int32Attribute{Computed: true, Description: "Month of the year."},
			"year":      schema.Int32Attribute{Computed: true, Description: "Date's year."},
			"end_circa": schema.BoolAttribute{Computed: true, Description: "Indicates whether the end date is an approximation."},
			"end_day":   schema.Int32Attribute{Computed: true, Description: "Date's end day (only valid if range is between)."},
			"end_month": schema.Int32Attribute{Computed: true, Description: "Date's end month (only valid if range is between)."},
			"end_year":  schema.Int32Attribute{Computed: true, Description: "Date's end year (only valid if range is between)."},
		},
		Description: description,
	}
}
