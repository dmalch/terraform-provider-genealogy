package project

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
)

func TestValueFrom(t *testing.T) {
	t.Run("Happy path, when a fully defined project response is passed", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.ProjectResponse{
			Id:          "project-123",
			Name:        "Test Project",
			Description: ptr("This is a test project"),
			UpdatedAt:   "1719709400",
			CreatedAt:   "1719709300",
		}

		model := &Model{}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("project-123"))
		Expect(model.Name.ValueString()).To(Equal("Test Project"))
		Expect(model.Description.ValueString()).To(Equal("This is a test project"))
		Expect(model.UpdatedAt.ValueString()).To(Equal("1719709400"))
		Expect(model.CreatedAt.ValueString()).To(Equal("1719709300"))
	})

	t.Run("When description is nil", func(t *testing.T) {
		RegisterTestingT(t)
		givenResponse := &geni.ProjectResponse{
			Id:        "project-456",
			Name:      "Minimal Project",
			UpdatedAt: "1719709400",
			CreatedAt: "1719709300",
		}

		model := &Model{}
		diags := ValueFrom(t.Context(), givenResponse, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("project-456"))
		Expect(model.Name.ValueString()).To(Equal("Minimal Project"))
		Expect(model.Description.IsNull()).To(BeTrue())
	})
}

func ptr[T any](s T) *T {
	return &s
}
