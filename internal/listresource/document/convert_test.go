package documentlist

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/onsi/gomega"

	documentresource "github.com/dmalch/terraform-provider-genealogy/internal/resource/document"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func ptr[T any](v T) *T { return &v }

func listRequestForDocument(t *testing.T, includeResource bool) list.ListRequest {
	t.Helper()
	r := documentresource.NewResource()

	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}

	var idResp resource.IdentitySchemaResponse
	withIdentity, ok := r.(resource.ResourceWithIdentity)
	if !ok {
		t.Fatalf("document resource %T does not implement ResourceWithIdentity", r)
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

func TestBuildDocumentListResult(t *testing.T) {
	t.Run("Populates Identity with the document ID", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForDocument(t, false)
		resp := &geni.DocumentResponse{Id: "document-42", Title: "Test Doc"}

		result, ok := buildDocumentListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Identity).NotTo(BeNil())

		var identity documentresource.ResourceIdentityModel
		diags := result.Identity.Get(context.Background(), &identity)
		Expect(diags.HasError()).To(BeFalse())
		Expect(identity.ID.ValueString()).To(Equal("document-42"))
	})

	t.Run("Sets a human-readable DisplayName", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForDocument(t, false)
		resp := &geni.DocumentResponse{Id: "document-42", Title: "Test Doc"}

		result, ok := buildDocumentListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		Expect(result.DisplayName).To(Equal("Test Doc (document-42)"))
	})

	t.Run("Populates Resource via document.ValueFrom when IncludeResource is true", func(t *testing.T) {
		RegisterTestingT(t)
		req := listRequestForDocument(t, true)
		resp := &geni.DocumentResponse{
			Id:          "document-43",
			Title:       "Birth Certificate",
			ContentType: ptr("image/png"),
		}

		result, ok := buildDocumentListResult(t.Context(), resp, req)

		Expect(ok).To(BeTrue())
		Expect(result.Diagnostics.HasError()).To(BeFalse())
		Expect(result.Resource).NotTo(BeNil())

		var model documentresource.ResourceModel
		diags := result.Resource.Get(context.Background(), &model)
		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("document-43"))
		Expect(model.Title.ValueString()).To(Equal("Birth Certificate"))
		Expect(model.ContentType.ValueString()).To(Equal("image/png"))
	})
}
