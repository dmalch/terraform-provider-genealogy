package geni

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// about_me	String	Profile's about me section (cf. detail_strings) (must be requested)
// baptism 	Event	Profile's baptism event info
// big_tree	Boolean	True if the profile is attached to the big tree
// birth 	Event	Profile's birth event info
// block_exists	Boolean	Indicates whether the profile is blocked
// burial 	Event	Profile's burial event info
// cause_of_death 	String	Profile's death cause
// claimed	Boolean	True if the profile is claimed by a user
// created_at 	String	Timestamp of when the profile was created
// created_by	String	URL (or id) of the profile who added this profile to the tree
// creator 	String	URL (or id) of the user who added this profile to the tree
// curator	String	Profile's curator's url (or id)
// current_residence	Location	Profile's current address
// death 	Event	Profile's death event info
// detail_strings 	Hash	Nested maps of locales to details fields (eg. about me) to values (must be requested)
// display_name	String	Profile's display name
// documents_updated_at 	String	Timestamp of the last document updated/added to the profile. Will not be return if no documents exist.
// email	String	Profile's email address
// events 	Array of Events	Events associated with this profile (must be requested)
// first_name	String	Profile's first name
// gender	String	Profile's gender
// get_email	Boolean	Indicates whether the profile can receive email
// guid	String	Profile's globally unique identifier
// id	String	Profile's node id
// is_alive	Boolean	True if the profile is living
// language	String	Язык профиля
// last_name	String	Profile's last name
// locked	Boolean	True if the profile has been locked down by a curator
// maiden_name	String	Profile's maiden name
// managers	Array of Strings	URLs (or ids) of profile(s) currently managing this profile
// master_profile	Boolean	Indicates whether the profile is a master profile
// merge_note	Array or String	Note explaining the profile's merge status
// merge_pending	Boolean	Indicates whether the profile has a pending merge
// merged_into	String	URL (or id) of the profile this profile is currently merged into
// middle_name	String	Profile's middle name
// mugshot_urls	PhotoImageSizeMap	All sizes of the profile's main photo
// name	String	Profile's name as it appears on the site to the current user
// names 	Hash	Nested maps of locales to name fields to values.
// Example: {"de": {"last_name": "Smith"}}
// nicknames 	Array of Strings	Also known as. Returned as an array, but can be set as a comma delimited list.
// occupation 	String	Профессия профиля
// phone_numbers	Array of PhoneNumbers	Profile's phone numbers
// photos_updated_at 	String	Timestamp of the last photo updated/added to the profile. Will not be return if no photos exist.
// premium_start_date	String	Дата перехода на подписку Pro
// profile_url 	String	URL to access profile in a browser
// public	Boolean	True если профиль общедоступный
// relationship	String	Profile's relationship to the current user (if any)
// requested_merges	Array of Strings	URLs (or ids) of the profile(s) requested to be merged into this one
// suffix	String	Profile's suffix
// unions	Array of Strings	URLs to unions
// updated_at 	String	Timestamp of when the profile was last updated
// url	String	URL to access profile through the API
// videos_updated_at 	String	Timestamp of the last video updated/added to the profile. Will not be return if no videos exist.
type ProfileResponse struct {
	// Id is the profile's node id
	Id string `json:"id"`
	// Guid is the profile's globally unique identifier
	Guid string `json:"guid"`
	// FirstName is the profile's first name
	FirstName string `json:"first_name"`
	// LastName is the profile's last name
	LastName string `json:"last_name"`
	// Gender is the profile's gender
	Gender string `json:"gender"`
	// Names is the name info
	Names map[string]NameResponse `json:"names"`
	// Birth is the birth event info
	Birth *EventResponse `json:"birth"`
	// Baptism is the baptism event info
	Baptism *EventResponse `json:"baptism"`
	// Death is the death event info
	Death *EventResponse `json:"death"`
	// Burial is the burial event info
	Burial *EventResponse `json:"burial"`
	// Events is the events associated with this profile
	Events []EventResponse `json:"events"`
	// IsAlive is a boolean that indicates whether the profile is living
	IsAlive bool `json:"is_alive"`
	// CreatedAt is the timestamp of when the profile was created
	CreatedAt string `json:"created_at"`
}

type NameResponse struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

// EventResponse is the response for an event
type EventResponse struct {
	Date        DateResponse `json:"date"`
	Description string       `json:"description"`
	Location    string       `json:"location"`
	Name        string       `json:"name"`
}

// DateResponse is the response for a date
type DateResponse struct {
	// Circa is a boolean that indicates whether the date is approximate
	Circa bool `json:"circa"`
	// Day is the day of the month
	Day int `json:"day"`
	// EndCirca is a boolean that indicates whether the end date is approximate
	EndCirca bool `json:"end_circa"`
	// EndDay is the end day of the month (only valid if range is between)
	EndDay int `json:"end_day"`
	// EndMonth is the end month of the year (only valid if range is between)
	EndMonth int `json:"end_month"`
	// EndYear is the end year (only valid if range is between)
	EndYear int `json:"end_year"`
	// Month is the month of the year
	Month int `json:"month"`
	// Range is the range (before, after, or between)
	Range string `json:"range"`
	// Year is the year
	Year int `json:"year"`
}

// LocationResponse is the response for a location
type LocationResponse struct {
	// City is the city name
	City string `json:"city"`
	// Country is the country name
	Country string `json:"country"`
	// County is the county name
	County string `json:"county"`
	// Latitude is the latitude
	Latitude float64 `json:"latitude"`
	// Longitude is the longitude
	Longitude float64 `json:"longitude"`
	// PlaceName is the place name
	PlaceName string `json:"place_name"`
	// State is the state name
	State string `json:"state"`
	// StreetAddress1 is the street address line 1
	StreetAddress1 string `json:"street_address1"`
	// StreetAddress2 is the street address line 2
	StreetAddress2 string `json:"street_address2"`
	// StreetAddress3 is the street address line 3
	StreetAddress3 string `json:"street_address3"`
}

func CreateProfile(accessToken, name, description string) (string, error) {
	return "", nil
}

func GetProfile(accessToken, profileId string) (*ProfileResponse, error) {
	// Create a new HTTP request
	requestUrl := geniUrl + "api/" + profileId + "?access_token=" + accessToken + "&api_version=1"
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	body, err := doRequest(req)
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

func UpdateProfile(accessToken, profileId, name, description string) error {
	return nil
}

func DeleteProfile(accessToken, profileId string) error {
	return nil
}
