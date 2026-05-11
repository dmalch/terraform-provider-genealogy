package profilelist

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	profileresource "github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
)

func ptr[T any](v T) *T { return &v }

// listRequestForProfile builds a list.ListRequest carrying the live
// managed-resource schemas, so tests exercise the same schemas the framework
// would hand the list resource at runtime. The caller must have already
// registered gomega for the current test via RegisterTestingT.
func listRequestForProfile(t *testing.T, includeResource bool) list.ListRequest {
	t.Helper()
	r := profileresource.NewProfileResource()

	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	Expect(schemaResp.Diagnostics.HasError()).To(BeFalse(), "building resource schema")

	withIdentity, ok := r.(resource.ResourceWithIdentity)
	Expect(ok).To(BeTrue(), "profile resource must implement ResourceWithIdentity")

	var idResp resource.IdentitySchemaResponse
	withIdentity.IdentitySchema(t.Context(), resource.IdentitySchemaRequest{}, &idResp)
	Expect(idResp.Diagnostics.HasError()).To(BeFalse(), "building identity schema")

	return list.ListRequest{
		ResourceSchema:         schemaResp.Schema,
		ResourceIdentitySchema: idResp.IdentitySchema,
		IncludeResource:        includeResource,
	}
}

func TestDisplayNameFor(t *testing.T) {
	t.Run("Returns 'First Last (id)' for an en-US profile", func(t *testing.T) {
		RegisterTestingT(t)
		got := displayNameFor(&geni.ProfileResponse{
			Id:        "profile-1",
			FirstName: ptr("John"),
			LastName:  ptr("Doe"),
		})
		Expect(got).To(Equal("John Doe (profile-1)"))
	})

	t.Run("Prefers the en-US locale entry from a localized names map", func(t *testing.T) {
		RegisterTestingT(t)
		got := displayNameFor(&geni.ProfileResponse{
			Id: "profile-2",
			Names: map[string]geni.NameElement{
				"fr":    {FirstName: ptr("Jean"), LastName: ptr("Dupont")},
				"en-US": {FirstName: ptr("John"), LastName: ptr("Doe")},
			},
		})
		Expect(got).To(Equal("John Doe (profile-2)"))
	})

	t.Run("Falls back to the bare ID when no name fields are populated", func(t *testing.T) {
		RegisterTestingT(t)
		got := displayNameFor(&geni.ProfileResponse{Id: "profile-3"})
		Expect(got).To(Equal("profile-3"))
	})
}

func TestBuildProfileListResult(t *testing.T) {
	t.Run("Populates Identity with the profile ID", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, false)
		givenResponse := &geni.ProfileResponse{Id: "profile-42", FirstName: ptr("John"), LastName: ptr("Doe")}

		result, ok := buildProfileListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Identity).NotTo(BeNil())

		var identity profileresource.ResourceIdentityModel
		Expect(result.Identity.Get(t.Context(), &identity).HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("profile-42"))
	})

	t.Run("Sets a human-readable DisplayName", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, false)
		givenResponse := &geni.ProfileResponse{Id: "profile-42", FirstName: ptr("John"), LastName: ptr("Doe")}

		result, ok := buildProfileListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		Expect(result.DisplayName).To(Equal("John Doe (profile-42)"))
	})

	t.Run("Populates Resource via profile.ValueFrom when IncludeResource is true", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, true)
		givenResponse := &geni.ProfileResponse{
			Id:        "profile-43",
			Public:    true,
			FirstName: ptr("John"),
			LastName:  ptr("Doe"),
		}

		result, ok := buildProfileListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Resource).NotTo(BeNil())

		var model profileresource.ResourceModel
		Expect(result.Resource.Get(t.Context(), &model).HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("profile-43"))

		var names map[string]profileresource.NameModel
		Expect(model.Names.ElementsAs(t.Context(), &names, false).HasError()).To(BeFalse())
		Expect(names).To(HaveKey("en-US"))
		Expect(names["en-US"].FirstName.ValueString()).To(Equal("John"))
	})

	t.Run("Leaves Resource at its schema-null default when IncludeResource is false", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, false)
		givenResponse := &geni.ProfileResponse{Id: "profile-44", FirstName: ptr("John"), LastName: ptr("Doe")}

		result, ok := buildProfileListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		// req.NewListResult pre-populates Resource with a null Raw; with
		// IncludeResource=false buildProfileListResult must not overwrite it.
		Expect(result.Resource.Raw.IsNull()).To(BeTrue())
	})
}
