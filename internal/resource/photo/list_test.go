package photo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/onsi/gomega"

	geniphoto "github.com/dmalch/go-geni/photo"
)

// listRequest builds a list.ListRequest carrying the live managed-resource
// schemas. The caller must have already registered gomega for the current test.
func listRequest(t *testing.T, includeResource bool) list.ListRequest {
	t.Helper()
	r := NewResource()

	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	Expect(schemaResp.Diagnostics.HasError()).To(BeFalse(), "building resource schema")

	withIdentity, ok := r.(resource.ResourceWithIdentity)
	Expect(ok).To(BeTrue(), "photo resource must implement ResourceWithIdentity")

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
	t.Run("Returns 'Title (id)' for a photo with a title", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(displayNameFor(&geniphoto.Photo{ID: "photo-1", Title: "Wedding"})).To(Equal("Wedding (photo-1)"))
	})

	t.Run("Falls back to the bare ID when title is empty", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(displayNameFor(&geniphoto.Photo{ID: "photo-2"})).To(Equal("photo-2"))
	})
}

func TestBuildListResult(t *testing.T) {
	t.Run("Populates Identity with the photo ID", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequest(t, false)

		result, ok := buildListResult(t.Context(), &geniphoto.Photo{ID: "photo-42", Title: "Test"}, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())

		var identity ResourceIdentityModel
		Expect(result.Identity.Get(t.Context(), &identity).HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("photo-42"))
	})

	t.Run("Sets a human-readable DisplayName", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequest(t, false)

		result, ok := buildListResult(t.Context(), &geniphoto.Photo{ID: "photo-42", Title: "Test"}, req)

		Expect(ok).To(BeTrue())
		Expect(result.DisplayName).To(Equal("Test (photo-42)"))
	})

	t.Run("Populates Resource via ValueFrom when IncludeResource is true", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequest(t, true)

		result, ok := buildListResult(t.Context(), &geniphoto.Photo{ID: "photo-43", Title: "Wedding", ContentType: "image/png"}, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())

		var model ResourceModel
		Expect(result.Resource.Get(t.Context(), &model).HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("photo-43"))
		Expect(model.Title.ValueString()).To(Equal("Wedding"))
		Expect(model.ContentType.ValueString()).To(Equal("image/png"))
	})
}
