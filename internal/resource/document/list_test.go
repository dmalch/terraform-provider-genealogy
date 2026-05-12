package document

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/onsi/gomega"

	"github.com/dmalch/go-geni"
)

// listRequest builds a list.ListRequest carrying the live managed-resource
// schemas. The caller must have already registered gomega for the current
// test via RegisterTestingT.
func listRequest(t *testing.T, includeResource bool) list.ListRequest {
	t.Helper()
	r := NewResource()

	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	Expect(schemaResp.Diagnostics.HasError()).To(BeFalse(), "building resource schema")

	withIdentity, ok := r.(resource.ResourceWithIdentity)
	Expect(ok).To(BeTrue(), "document resource must implement ResourceWithIdentity")

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
	t.Run("Returns 'Title (id)' for a document with a title", func(t *testing.T) {
		RegisterTestingT(t)
		got := displayNameFor(&geni.DocumentResponse{Id: "document-1", Title: "Birth Certificate"})
		Expect(got).To(Equal("Birth Certificate (document-1)"))
	})

	t.Run("Falls back to the bare ID when title is empty", func(t *testing.T) {
		RegisterTestingT(t)
		got := displayNameFor(&geni.DocumentResponse{Id: "document-2"})
		Expect(got).To(Equal("document-2"))
	})
}

func TestBuildListResult(t *testing.T) {
	t.Run("Populates Identity with the document ID", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequest(t, false)
		givenResponse := &geni.DocumentResponse{Id: "document-42", Title: "Test Doc"}

		result, ok := buildListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Identity).NotTo(BeNil())

		var identity ResourceIdentityModel
		Expect(result.Identity.Get(t.Context(), &identity).HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("document-42"))
	})

	t.Run("Sets a human-readable DisplayName", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequest(t, false)
		givenResponse := &geni.DocumentResponse{Id: "document-42", Title: "Test Doc"}

		result, ok := buildListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		Expect(result.DisplayName).To(Equal("Test Doc (document-42)"))
	})

	t.Run("Populates Resource via ValueFrom when IncludeResource is true", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequest(t, true)
		givenResponse := &geni.DocumentResponse{
			Id:          "document-43",
			Title:       "Birth Certificate",
			ContentType: ptr("image/png"),
		}

		result, ok := buildListResult(t.Context(), givenResponse, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Resource).NotTo(BeNil())

		var model ResourceModel
		Expect(result.Resource.Get(t.Context(), &model).HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("document-43"))
		Expect(model.Title.ValueString()).To(Equal("Birth Certificate"))
		Expect(model.ContentType.ValueString()).To(Equal("image/png"))
	})
}
