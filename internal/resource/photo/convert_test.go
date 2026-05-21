package photo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	geniphoto "github.com/dmalch/go-geni/photo"
)

func TestValueFrom(t *testing.T) {
	t.Run("Populates the model from a full photo response", func(t *testing.T) {
		RegisterTestingT(t)
		response := &geniphoto.Photo{
			ID:          "photo-1",
			Guid:        "abc123",
			AlbumId:     "album-7",
			Title:       "Wedding",
			Description: "A lovely day",
			Date:        "1 Jan 1990",
			Attribution: "Family archive",
			ContentType: "image/png",
			Url:         "https://api.geni.com/photo-1",
			Tags:        []string{"profile-1", "profile-2"},
			Sizes:       map[string]string{"small": "https://img/s.png"},
			CreatedAt:   "1000",
			UpdatedAt:   "2000",
		}

		model := &ResourceModel{}
		diags := ValueFrom(t.Context(), response, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("photo-1"))
		Expect(model.Title.ValueString()).To(Equal("Wedding"))
		Expect(model.Description.ValueString()).To(Equal("A lovely day"))
		Expect(model.Date.ValueString()).To(Equal("1 Jan 1990"))
		Expect(model.Album.ValueString()).To(Equal("album-7"))
		Expect(model.Guid.ValueString()).To(Equal("abc123"))
		Expect(model.ContentType.ValueString()).To(Equal("image/png"))
		Expect(model.Profiles.Elements()).To(HaveLen(2))
		Expect(model.Sizes.Elements()).To(HaveKey("small"))
		Expect(model.CreatedAt.ValueString()).To(Equal("1000"))
		Expect(model.UpdatedAt.ValueString()).To(Equal("2000"))
	})

	t.Run("Maps empty optional fields to null", func(t *testing.T) {
		RegisterTestingT(t)
		model := &ResourceModel{}
		diags := ValueFrom(t.Context(), &geniphoto.Photo{ID: "photo-2", Title: "Untitled"}, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.Description.IsNull()).To(BeTrue())
		Expect(model.Date.IsNull()).To(BeTrue())
		Expect(model.Album.IsNull()).To(BeTrue())
		Expect(model.Profiles.Elements()).To(BeEmpty())
	})

	t.Run("Leaves file and file_name untouched", func(t *testing.T) {
		RegisterTestingT(t)
		model := &ResourceModel{
			File:     types.StringValue("base64data"),
			FileName: types.StringValue("photo.png"),
		}
		diags := ValueFrom(t.Context(), &geniphoto.Photo{ID: "photo-3", Title: "x"}, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.File.ValueString()).To(Equal("base64data"))
		Expect(model.FileName.ValueString()).To(Equal("photo.png"))
	})
}

func TestRequestFrom(t *testing.T) {
	t.Run("Builds the update request from the model", func(t *testing.T) {
		RegisterTestingT(t)
		request := RequestFrom(ResourceModel{
			Title:       types.StringValue("Wedding"),
			Description: types.StringValue("A lovely day"),
			Date:        types.StringValue("1990"),
		})

		Expect(request.Title).To(Equal("Wedding"))
		Expect(request.Description).To(Equal("A lovely day"))
		Expect(request.Date).To(Equal("1990"))
	})
}

func TestUpdateComputedFields(t *testing.T) {
	t.Run("Fills computed fields and keeps configured values", func(t *testing.T) {
		RegisterTestingT(t)
		response := &geniphoto.Photo{
			ID: "photo-9", Guid: "g9", ContentType: "image/png",
			Url: "https://api/photo-9", CreatedAt: "10", UpdatedAt: "20",
			Tags: []string{"profile-5"},
		}
		model := &ResourceModel{
			Title:    types.StringValue("Kept title"),
			Profiles: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("profile-1")}),
		}
		diags := UpdateComputedFields(t.Context(), response, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.ID.ValueString()).To(Equal("photo-9"))
		Expect(model.Guid.ValueString()).To(Equal("g9"))
		Expect(model.Title.ValueString()).To(Equal("Kept title"))
		// Profiles were already configured — preserved, not overwritten.
		Expect(model.Profiles.Elements()).To(HaveLen(1))
	})

	t.Run("Adopts response tags when profiles are unset", func(t *testing.T) {
		RegisterTestingT(t)
		model := &ResourceModel{Profiles: types.SetNull(types.StringType)}
		diags := UpdateComputedFields(t.Context(), &geniphoto.Photo{ID: "photo-1", Tags: []string{"profile-7"}}, model)

		Expect(diags.HasError()).To(BeFalse())
		Expect(model.Profiles.Elements()).To(HaveLen(1))
	})
}

func TestStringOrNull(t *testing.T) {
	t.Run("Returns null for an empty string", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(stringOrNull("").IsNull()).To(BeTrue())
	})

	t.Run("Returns the value for a non-empty string", func(t *testing.T) {
		RegisterTestingT(t)
		Expect(stringOrNull("x").ValueString()).To(Equal("x"))
	})
}
