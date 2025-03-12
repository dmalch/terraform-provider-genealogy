package geni

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
type Profile struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
}

func CreateProfile(apiKey, name, description string) (string, error) {
	return "", nil
}

func GetProfile(apiKey, profileId string) (*Profile, error) {
	return nil, nil
}

func UpdateProfile(apiKey, profileId, name, description string) error {
	return nil
}

func DeleteProfile(apiKey, profileId string) error {
	return nil
}
