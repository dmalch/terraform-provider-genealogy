package geni

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
)

type ProfileRequest struct {
	// FirstName is the profile's first name
	FirstName string `json:"first_name,omitempty"`
	// LastName is the profile's last name
	LastName string `json:"last_name,omitempty"`
	// MiddleName is the profile's middle name
	MiddleName string `json:"middle_name,omitempty"`
	// MaidenName is the profile's maiden name
	MaidenName string `json:"maiden_name,omitempty"`
	// DisplayName is the profile's display name
	DisplayName string `json:"display_name,omitempty"`
	// Nicknames is the profile's nicknames
	Nicknames []string `json:"nicknames,omitempty"`
	// Gender is the profile's gender
	Gender string `json:"gender,omitempty"`
	// Names is the name info
	Names map[string]NameElement `json:"names,omitempty"`
	// Birth is the birth event info
	Birth *EventElement `json:"birth,omitempty"`
	// Baptism is the baptism event info
	Baptism *EventElement `json:"baptism,omitempty"`
	// Death is the death event info
	Death *EventElement `json:"death,omitempty"`
	// CauseOfDeath is the cause of death
	CauseOfDeath string `json:"cause_of_death,omitempty"`
	// Burial is the burial event info
	Burial *EventElement `json:"burial,omitempty"`
	// IsAlive is a boolean that indicates whether the profile is living
	IsAlive bool `json:"is_alive,omitempty"`
	// Title is the profile's name title
	Title string `json:"title,omitempty"`
	// AboutMe is the profile's about me section
	AboutMe string `json:"about_me,omitempty"`
	// Occupation is the profile's occupation
	Occupation string `json:"occupation,omitempty"`
	// Suffix is the profile's suffix
	Suffix string `json:"suffix,omitempty"`
	// Public is a boolean that indicates whether the profile is public
	Public bool `json:"public,omitempty"`
	// Locked is a boolean that indicates whether the profile is locked down by a curator
	Locked bool `json:"locked,omitempty"`
	// MergeNote is the note explaining the profile's merge status
	MergeNote []string `json:"merge_note,omitempty"`
}

type ProfileResponse struct {
	// Id is the profile's node id
	Id string `json:"id,omitempty"`
	// Guid is the profile's globally unique identifier
	Guid string `json:"guid,omitempty"`
	// FirstName is the profile's first name
	FirstName string `json:"first_name,omitempty"`
	// LastName is the profile's last name
	LastName string `json:"last_name,omitempty"`
	// MiddleName is the profile's middle name
	MiddleName string `json:"middle_name,omitempty"`
	// MaidenName is the profile's maiden name
	MaidenName string `json:"maiden_name,omitempty"`
	// DisplayName is the profile's display name
	DisplayName string `json:"display_name,omitempty"`
	// Nicknames is the profile's nicknames
	Nicknames []string `json:"nicknames,omitempty"`
	// Gender is the profile's gender
	Gender string `json:"gender,omitempty"`
	// Names is the name info
	Names map[string]NameElement `json:"names,omitempty"`
	// Birth is the birth event info
	Birth *EventElement `json:"birth,omitempty"`
	// Baptism is the baptism event info
	Baptism *EventElement `json:"baptism,omitempty"`
	// Death is the death event info
	Death *EventElement `json:"death,omitempty"`
	// CauseOfDeath is the cause of death
	CauseOfDeath string `json:"cause_of_death,omitempty"`
	// Burial is the burial event info
	Burial *EventElement `json:"burial,omitempty"`
	// Events is the events associated with this profile
	Events []EventElement `json:"events,omitempty"`
	// IsAlive is a boolean that indicates whether the profile is living
	IsAlive bool `json:"is_alive,omitempty"`
	// Title is the profile's name title
	Title string `json:"title,omitempty"`
	// AboutMe is the profile's about me section
	AboutMe string `json:"about_me,omitempty"`
	// Occupation is the profile's occupation
	Occupation string `json:"occupation,omitempty"`
	// Suffix is the profile's suffix
	Suffix string `json:"suffix,omitempty"`
	// Public is a boolean that indicates whether the profile is public
	Public bool `json:"public,omitempty"`
	// Locked is a boolean that indicates whether the profile is locked down by a curator
	Locked bool `json:"locked,omitempty"`
	// Language is the profile's language
	Language string `json:"language,omitempty"`
	// ProfileUrl is the URL to access profile in a browser
	ProfileUrl string `json:"profile_url,omitempty"`
	// MergePending is a boolean that indicates whether the profile has a pending merge
	MergePending bool `json:"merge_pending,omitempty"`
	// MergedInto is the URL (or id) of the profile this profile is currently merged into
	MergedInto string `json:"merged_into,omitempty"`
	// MergeNote is the note explaining the profile's merge status
	MergeNote []string `json:"merge_note,omitempty"`
	// Url is the URL to access profile through the API
	Url string `json:"url,omitempty"`
	// Unions is the URLs to unions
	Unions []string `json:"unions,omitempty"`
	// UpdatedAt is the timestamp of when the profile was last updated
	UpdatedAt string `json:"updated_at,omitempty"`
	// CreatedAt is the timestamp of when the profile was created
	CreatedAt string `json:"created_at,omitempty"`
}

// NameElement is the response for a name.
type NameElement struct {
	// FirstName is the profile's first name
	FirstName *string `json:"first_name,omitempty"`
	// LastName is the profile's last name
	LastName *string `json:"last_name,omitempty"`
	// MiddleName is the profile's middle name
	MiddleName *string `json:"middle_name,omitempty"`
}

// EventElement is the response for an event.
type EventElement struct {
	Date        *DateElement     `json:"date,omitempty"`
	Description string           `json:"description,omitempty"`
	Location    *LocationElement `json:"location,omitempty"`
	Name        string           `json:"name,omitempty"`
}

// DateElement is the response for a date.
type DateElement struct {
	// Circa is a boolean that indicates whether the date is approximate
	Circa *bool `json:"circa,omitempty"`
	// Day is the day of the month
	Day *int32 `json:"day,omitempty"`
	// Month is the month of the year
	Month *int32 `json:"month,omitempty"`
	// Year is the year
	Year *int32 `json:"year,omitempty"`
	// EndCirca is a boolean that indicates whether the end date is approximate
	EndCirca *bool `json:"end_circa,omitempty"`
	// EndDay is the end day of the month (only valid if range is between)
	EndDay *int32 `json:"end_day,omitempty"`
	// EndMonth is the end month of the year (only valid if range is between)
	EndMonth *int32 `json:"end_month,omitempty"`
	// EndYear is the end year (only valid if range is between)
	EndYear *int32 `json:"end_year,omitempty"`
	// Range is the range (before, after, or between)
	Range *string `json:"range,omitempty"`
}

// LocationElement is the response for a location.
type LocationElement struct {
	// City is the city name
	City *string `json:"city,omitempty"`
	// Country is the country name
	Country *string `json:"country,omitempty"`
	// County is the county name
	County *string `json:"county,omitempty"`
	// Latitude is the latitude
	Latitude *big.Float `json:"latitude,omitempty"`
	// Longitude is the longitude
	Longitude *big.Float `json:"longitude,omitempty"`
	// PlaceName is the place name
	PlaceName *string `json:"place_name,omitempty"`
	// State is the state name
	State *string `json:"state,omitempty"`
	// StreetAddress1 is the street address line 1
	StreetAddress1 *string `json:"street_address1,omitempty"`
	// StreetAddress2 is the street address line 2
	StreetAddress2 *string `json:"street_address2,omitempty"`
	// StreetAddress3 is the street address line 3
	StreetAddress3 *string `json:"street_address3,omitempty"`
}

func (c *Client) CreateProfile(request *ProfileRequest) (*ProfileResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		slog.Error("Error marshaling request", "error", err)
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/profile/add"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &profile, nil
}

func escapeString(s string) string {
	return escapeStringToUTF(s)
}

func escapeStringToUTF(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r > 127 {
			sb.WriteString(fmt.Sprintf("\\u%04x", r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func (c *Client) GetProfile(profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	c.fixResponse(&profile)

	return &profile, nil
}

func (c *Client) fixResponse(profile *ProfileResponse) {
	//The only_ids flag does not work for the profile endpoint, so we need to remove
	//the url from the Unions field.
	for i, union := range profile.Unions {
		profile.Unions[i] = strings.Replace(union, apiUrl(c.useSandboxEnv), "", 1)
	}
}

func (c *Client) UpdateProfile(profileId string, request *ProfileRequest) (*ProfileResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		slog.Error("Error marshaling request", "error", err)
		return nil, err
	}

	jsonStr := strings.ReplaceAll(string(jsonBody), "\\\\", "\\")
	jsonStr = escapeString(jsonStr)

	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/update"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonStr))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	return &profile, nil
}

type ResultResponse struct {
	Result string `json:"result,omitempty"`
}

func (c *Client) DeleteProfile(profileId string) error {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/delete"
	req, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		slog.Error("Error creating request", "error", err)
		return err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	var result ResultResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return err
	}

	return nil
}

func (c *Client) AddPartner(profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/add-partner"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	c.fixResponse(&profile)

	return &profile, nil
}

func (c *Client) AddChild(profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/add-child"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	c.fixResponse(&profile)

	return &profile, nil
}

func (c *Client) AddSibling(profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/add-sibling"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var profile ProfileResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	c.fixResponse(&profile)

	return &profile, nil
}

func (c *Client) MergeProfiles(profile1Id, profile2Id string) error {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profile1Id + "/merge/" + profile2Id
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return err
	}

	if err := c.addStandardHeadersAndQueryParams(req); err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	var result ResultResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return err
	}

	return nil
}
