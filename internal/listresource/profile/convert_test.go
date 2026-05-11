package profilelist

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	profileresource "github.com/dmalch/terraform-provider-genealogy/internal/resource/profile"
)

func ptr[T any](v T) *T { return &v }

// listRequestForProfile constructs a list.ListRequest carrying the live
// managed-resource schemas, so tests exercise the same schemas the framework
// would hand the list resource at runtime.
func listRequestForProfile(t *testing.T, includeResource bool) list.ListRequest {
	t.Helper()
	r := profileresource.NewProfileResource()

	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}

	var idResp resource.IdentitySchemaResponse
	withIdentity, ok := r.(resource.ResourceWithIdentity)
	if !ok {
		t.Fatalf("profile resource %T does not implement ResourceWithIdentity", r)
	}
	withIdentity.IdentitySchema(t.Context(), resource.IdentitySchemaRequest{}, &idResp)
	if idResp.Diagnostics.HasError() {
		t.Fatalf("failed to build identity schema: %v", idResp.Diagnostics)
	}

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
		resp := &geni.ProfileResponse{Id: "profile-42", FirstName: ptr("John"), LastName: ptr("Doe")}

		result, ok := buildProfileListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Identity).NotTo(BeNil())

		var identity profileresource.ResourceIdentityModel
		diags := result.Identity.Get(context.Background(), &identity)
		Expect(diags.HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("profile-42"))
	})

	t.Run("Sets a human-readable DisplayName", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, false)
		resp := &geni.ProfileResponse{Id: "profile-42", FirstName: ptr("John"), LastName: ptr("Doe")}

		result, ok := buildProfileListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		Expect(result.DisplayName).To(Equal("John Doe (profile-42)"))
	})

	t.Run("Populates Resource via profile.ValueFrom when IncludeResource is true", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, true)
		resp := &geni.ProfileResponse{
			Id:        "profile-43",
			Public:    true,
			FirstName: ptr("John"),
			LastName:  ptr("Doe"),
		}

		result, ok := buildProfileListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Resource).NotTo(BeNil())

		var model profileresource.ResourceModel
		diags := result.Resource.Get(context.Background(), &model)
		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("profile-43"))

		var names map[string]profileresource.NameModel
		Expect(model.Names.ElementsAs(context.Background(), &names, false).HasError()).To(BeFalse())
		Expect(names).To(HaveKey("en-US"))
		Expect(names["en-US"].FirstName.ValueString()).To(Equal("John"))
	})

	t.Run("Leaves Resource at its schema-null default when IncludeResource is false", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForProfile(t, false)
		resp := &geni.ProfileResponse{Id: "profile-44", FirstName: ptr("John"), LastName: ptr("Doe")}

		result, ok := buildProfileListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		// When IncludeResource is false the framework still pre-populates a
		// null-schema Resource via req.NewListResult; ID is null.
		var model profileresource.ResourceModel
		_ = result.Resource.Get(context.Background(), &model)
		Expect(model.ID.IsNull()).To(BeTrue())
	})
}

// Compile-time check: NameModel must be visible from the profile package.
var _ = types.StringNull
