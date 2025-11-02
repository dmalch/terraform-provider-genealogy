package geni

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type DetailsString struct {
	// AboutMe is the profile's about me section
	AboutMe *string `json:"about_me"`
}

type ProfileRequest struct {
	// DisplayName is the profile's display name
	DisplayName *string `json:"display_name,omitempty"`
	// Nicknames is the profile's nicknames
	Nicknames []string `json:"nicknames,omitempty"`
	// Gender is the profile's gender
	Gender *string `json:"gender,omitempty"`
	// Names is the name info
	Names map[string]NameElement `json:"names,omitempty"`
	// Birth is the birth event info
	Birth *EventElement `json:"birth,omitempty"`
	// Baptism is the baptism event info
	Baptism *EventElement `json:"baptism,omitempty"`
	// Death is the death event info
	Death *EventElement `json:"death,omitempty"`
	// CauseOfDeath is the cause of death
	CauseOfDeath *string `json:"cause_of_death"`
	// Burial is the burial event info
	Burial *EventElement `json:"burial,omitempty"`
	// IsAlive is a boolean that indicates whether the profile is living
	IsAlive bool `json:"is_alive"`
	// Title is the profile's name title
	Title string `json:"title,omitempty"`
	// CurrentResidence is the profile's current residence
	CurrentResidence *LocationElement `json:"current_residence"`
	// AboutMe is the profile's about me section
	AboutMe *string `json:"about_me"`
	// DetailStrings are nested maps of locales to details fields (e.g. about me) to values
	DetailStrings map[string]DetailsString `json:"detail_strings"`
	// Occupation is the profile's occupation
	Occupation string `json:"occupation,omitempty"`
	// Suffix is the profile's suffix
	Suffix string `json:"suffix,omitempty"`
	// Public is a boolean that indicates whether the profile is public
	Public bool `json:"public"`
	// Locked is a boolean that indicates whether the profile is locked down by a curator
	Locked bool `json:"locked"`
	// MergeNote is the note explaining the profile's merge status
	MergeNote []string `json:"merge_note,omitempty"`
}

type ProfileBulkResponse struct {
	Results    []ProfileResponse `json:"results,omitempty"`
	Page       int               `json:"page,omitempty"`
	TotalCount int               `json:"total_count,omitempty"`
}

type ProfileResponse struct {
	// Id is the profile's node id
	Id string `json:"id,omitempty"`
	// Guid is the profile's globally unique identifier
	Guid string `json:"guid,omitempty"`
	// FirstName is the profile's first name
	FirstName *string `json:"first_name,omitempty"`
	// LastName is the profile's last name
	LastName *string `json:"last_name,omitempty"`
	// MiddleName is the profile's middle name
	MiddleName *string `json:"middle_name,omitempty"`
	// MaidenName is the profile's maiden name
	MaidenName *string `json:"maiden_name,omitempty"`
	// DisplayName is the profile's display name
	DisplayName *string `json:"display_name,omitempty"`
	// Nicknames is the profile's nicknames
	Nicknames []string `json:"nicknames,omitempty"`
	// Gender is the profile's gender
	Gender *string `json:"gender,omitempty"`
	// Names is the name info
	Names map[string]NameElement `json:"names,omitempty"`
	// Birth is the birth event info
	Birth *EventElement `json:"birth,omitempty"`
	// Baptism is the baptism event info
	Baptism *EventElement `json:"baptism,omitempty"`
	// Death is the death event info
	Death *EventElement `json:"death,omitempty"`
	// CauseOfDeath is the cause of death
	CauseOfDeath *string `json:"cause_of_death,omitempty"`
	// Burial is the burial event info
	Burial *EventElement `json:"burial,omitempty"`
	// Events is the events associated with this profile
	Events []EventElement `json:"events,omitempty"`
	// IsAlive is a boolean that indicates whether the profile is living
	IsAlive bool `json:"is_alive"`
	// Title is the profile's name title
	Title string `json:"title,omitempty"`
	// CurrentResidence is the profile's current residence
	CurrentResidence *LocationElement `json:"current_residence"`
	// AboutMe is the profile's about me section
	AboutMe *string `json:"about_me,omitempty"`
	// DetailStrings are nested maps of locales to details fields (e.g. about me) to values
	DetailStrings map[string]DetailsString `json:"detail_strings,omitempty"`
	// Occupation is the profile's occupation
	Occupation string `json:"occupation,omitempty"`
	// Suffix is the profile's suffix
	Suffix string `json:"suffix,omitempty"`
	// Public is a boolean that indicates whether the profile is public
	Public bool `json:"public"`
	// Locked is a boolean that indicates whether the profile is locked down by a curator
	Locked bool `json:"locked"`
	// Language is the profile's language
	Language string `json:"language,omitempty"`
	// ProfileUrl is the URL to access profile in a browser
	ProfileUrl string `json:"profile_url,omitempty"`
	// MergePending is a boolean that indicates whether the profile has a pending merge
	MergePending bool `json:"merge_pending,omitempty"`
	// MergedInto is the ID of the profile this profile is currently merged into
	MergedInto string `json:"merged_into,omitempty"`
	// MergeNote is the note explaining the profile's merge status
	MergeNote []string `json:"merge_note,omitempty"`
	// Url is the URL to access profile through the API
	Url string `json:"url,omitempty"`
	// Unions is the URLs to unions
	Unions []string `json:"unions,omitempty"`
	// Deleted is a boolean that indicates whether the profile is deleted
	Deleted bool `json:"deleted"`
	// UpdatedAt is the timestamp of when the profile was last updated
	UpdatedAt string `json:"updated_at,omitempty"`
	// CreatedAt is the timestamp of when the profile was created
	CreatedAt string `json:"created_at,omitempty"`
}

// NameElement is the response for a name.
type NameElement struct {
	// FirstName is the profile's first name
	FirstName *string `json:"first_name"`
	// LastName is the profile's last name
	LastName *string `json:"last_name"`
	// MiddleName is the profile's middle name
	MiddleName *string `json:"middle_name"`
	// MaidenName is the profile's maiden name
	MaidenName *string `json:"maiden_name"`
	// DisplayName is the profile's display name
	DisplayName *string `json:"display_name"`
	// Nicknames is the profile's comma-separated list of nicknames
	Nicknames *string `json:"nicknames"`
}

// EventElement is the response for an event.
type EventElement struct {
	Date        *DateElement     `json:"date"`
	Description *string          `json:"description,omitempty"`
	Location    *LocationElement `json:"location"`
	Name        string           `json:"name,omitempty"`
}

// DateElement is the response for a date.
type DateElement struct {
	// Circa is a boolean that indicates whether the date is approximate
	Circa *bool `json:"circa"`
	// Day is the day of the month
	Day *int32 `json:"day"`
	// Month is the month of the year
	Month *int32 `json:"month"`
	// Year is the year
	Year *int32 `json:"year"`
	// EndCirca is a boolean that indicates whether the end date is approximate
	EndCirca *bool `json:"end_circa"`
	// EndDay is the end day of the month (only valid if range is between)
	EndDay *int32 `json:"end_day"`
	// EndMonth is the end month of the year (only valid if range is between)
	EndMonth *int32 `json:"end_month"`
	// EndYear is the end year (only valid if range is between)
	EndYear *int32 `json:"end_year"`
	// Range is the range (before, after, or between)
	Range *string `json:"range"`
}

// LocationElement is the response for a location.
type LocationElement struct {
	// City is the city name
	City *string `json:"city"`
	// Country is the country name
	Country *string `json:"country"`
	// County is the county name
	County *string `json:"county"`
	// Latitude is the latitude
	Latitude *float64 `json:"latitude,omitempty"`
	// Longitude is the longitude
	Longitude *float64 `json:"longitude,omitempty"`
	// PlaceName is the place name
	PlaceName *string `json:"place_name"`
	// State is the state name
	State *string `json:"state"`
	// StreetAddress1 is the street address line 1
	StreetAddress1 *string `json:"street_address1"`
	// StreetAddress2 is the street address line 2
	StreetAddress2 *string `json:"street_address2"`
	// StreetAddress3 is the street address line 3
	StreetAddress3 *string `json:"street_address3"`
}

func (c *Client) CreateProfile(ctx context.Context, request *ProfileRequest) (*ProfileResponse, error) {
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

	c.addProfileFieldsQueryParams(req)

	body, err := c.doRequest(ctx, req)
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

func (c *Client) GetProfile(ctx context.Context, profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	c.addProfileFieldsQueryParams(req)

	body, err := c.doRequest(ctx, req,
		WithRequestKey(func() string {
			return profileId
		}),
		WithPrepareBulkRequest(func(req *http.Request, urlMap *sync.Map) {
			// Add a new ids parameter containing IDs of all profiles to be fetched in
			// addition to the current one. First, we need to get the IDs from the map.
			ids := make([]string, 0)

			ids = append(ids, profileId)

			urlMap.Range(func(key, value interface{}) bool {
				if _, ok := value.(context.CancelFunc); ok {
					if keyString, ok := key.(string); ok && strings.Contains(keyString, "profile") {
						ids = append(ids, keyString)
					}
				}
				return true
			})

			if len(ids) > 1 {
				query := req.URL.Query()
				query.Add("ids", strings.Join(ids, ","))
				req.URL.RawQuery = query.Encode()
			}
		}),
		WithParseBulkResponse(func(req *http.Request, body []byte, urlMap *sync.Map) ([]byte, error) {
			// If only one profile is requested, we can skip the bulk response parsing
			if !req.URL.Query().Has("ids") {
				return body, nil
			}

			// Parse the response to get the profile ID
			var response ProfileBulkResponse
			err := json.Unmarshal(body, &response)
			if err != nil {
				slog.Error("Error unmarshaling bulk response", "error", err)
				return nil, err
			}

			var requestedProfileRes []byte

			// Store the response in the map using the profile ID as the key
			for _, profile := range response.Results {

				jsonBody, err := json.Marshal(&profile)
				if err != nil {
					slog.Error("Error marshaling request", "error", err)
					return nil, err
				}

				if profile.Id == profileId {
					requestedProfileRes = jsonBody
					continue
				}

				previous, loaded := urlMap.Swap(profile.Id, jsonBody)
				if loaded {
					// If the previous value is context cancel function, cancel it
					if cancelFunc, ok := previous.(context.CancelFunc); ok {
						cancelFunc()
					}
				}
			}

			return requestedProfileRes, nil
		}))
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

func (c *Client) addProfileFieldsQueryParams(req *http.Request) {
	query := req.URL.Query()
	query.Add("fields", "id,first_name,last_name,middle_name,maiden_name,display_name,nicknames,names,gender,birth,baptism,death,burial,cause_of_death,current_residence,about_me,detail_strings,unions,is_alive,public,deleted,merged_into,updated_at,created_at")
	req.URL.RawQuery = query.Encode()
}

func (c *Client) fixResponse(profile *ProfileResponse) {
	//The only_ids flag does not work for the profile endpoint, so we need to remove
	//the url from the Unions field.
	for i, union := range profile.Unions {
		profile.Unions[i] = strings.Replace(union, apiUrl(c.useSandboxEnv), "", 1)
	}
}

func (c *Client) GetProfiles(ctx context.Context, profileIds []string) (*ProfileBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/profile"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	c.addProfileFieldsQueryParams(req)

	query := req.URL.Query()
	query.Add("ids", strings.Join(profileIds, ","))
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var profiles ProfileBulkResponse
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	for i := range profiles.Results {
		c.fixResponse(&profiles.Results[i])
	}

	return &profiles, nil
}

func (c *Client) GetManagedProfiles(ctx context.Context, page int) (*ProfileBulkResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/user/managed-profiles"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	c.addProfileFieldsQueryParams(req)

	query := req.URL.Query()
	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var profile ProfileBulkResponse
	err = json.Unmarshal(body, &profile)
	if err != nil {
		slog.Error("Error unmarshaling response", "error", err)
		return nil, err
	}

	// Iterate over the profiles and fix the response
	for i := range profile.Results {
		c.fixResponse(&profile.Results[i])
	}

	return &profile, nil
}

func (c *Client) UpdateProfile(ctx context.Context, profileId string, request *ProfileRequest) (*ProfileResponse, error) {
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

	c.addProfileFieldsQueryParams(req)

	body, err := c.doRequest(ctx, req)
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

type ResultResponse struct {
	Result string `json:"result,omitempty"`
}

func (c *Client) DeleteProfile(ctx context.Context, profileId string) error {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/delete"
	req, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		slog.Error("Error creating request", "error", err)
		return err
	}

	body, err := c.doRequest(ctx, req)
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

func (c *Client) AddPartner(ctx context.Context, profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/add-partner"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	c.addProfileFieldsQueryParams(req)

	body, err := c.doRequest(ctx, req)
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

func (c *Client) AddChild(ctx context.Context, profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/add-child"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	c.addProfileFieldsQueryParams(req)

	body, err := c.doRequest(ctx, req)
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

func (c *Client) AddSibling(ctx context.Context, profileId string) (*ProfileResponse, error) {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profileId + "/add-sibling"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	c.addProfileFieldsQueryParams(req)

	body, err := c.doRequest(ctx, req)
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

func (c *Client) MergeProfiles(ctx context.Context, profile1Id, profile2Id string) error {
	url := BaseUrl(c.useSandboxEnv) + "api/" + profile1Id + "/merge/" + profile2Id
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return err
	}

	body, err := c.doRequest(ctx, req)
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
