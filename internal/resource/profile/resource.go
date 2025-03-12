package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/dmalch/terraform-provider-geni/internal/config"
	"github.com/dmalch/terraform-provider-geni/internal/geni"
)

type Resource struct {
	resource.ResourceWithConfigure
	apiKey types.String
}

func NewProfileResource() resource.Resource {
	return &Resource{}
}

// Metadata provides the resource type name
func (r *Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "geni_profile"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(*config.GeniProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.GeniProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.apiKey = cfg.ApiKey
}

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
type ResourceModel struct {
	ID        types.String `tfsdk:"id"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	CreatedAt types.String `tfsdk:"created_at"`
	Gender    types.String `tfsdk:"gender"`
}

// Create creates the resource
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := geni.CreateProfile(r.apiKey.ValueString(), plan.FirstName.ValueString(), plan.LastName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating profile", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read reads the resource
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := geni.GetProfile(r.apiKey.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading profile", err.Error())
		return
	}

	if profile.FirstName != "" {
		state.FirstName = types.StringValue(profile.FirstName)
	}
	if profile.LastName != "" {
		state.LastName = types.StringValue(profile.LastName)
	}
	if profile.Id != "" {
		state.ID = types.StringValue(profile.Id)
	}
	if profile.Gender != "" {
		state.Gender = types.StringValue(profile.Gender)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := geni.UpdateProfile(r.apiKey.ValueString(), plan.ID.ValueString(), plan.FirstName.ValueString(), plan.LastName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating profile", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := geni.DeleteProfile(r.apiKey.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting profile", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
