package union

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/gomega"

	geniprofile "github.com/dmalch/go-geni/profile"
)

func TestUnionIDFrom(t *testing.T) {
	t.Run("Keeps an already-known union id", func(t *testing.T) {
		RegisterTestingT(t)
		got := unionIDFrom(types.StringValue("union-1"), &geniprofile.Profile{Unions: []string{"union-9"}})
		Expect(got.ValueString()).To(Equal("union-1"))
	})

	t.Run("Adopts the temp profile's union when the id is null", func(t *testing.T) {
		RegisterTestingT(t)
		got := unionIDFrom(types.StringNull(), &geniprofile.Profile{Unions: []string{"union-9"}})
		Expect(got.ValueString()).To(Equal("union-9"))
	})

	t.Run("Adopts the temp profile's union when the id is unknown", func(t *testing.T) {
		RegisterTestingT(t)
		got := unionIDFrom(types.StringUnknown(), &geniprofile.Profile{Unions: []string{"union-9"}})
		Expect(got.ValueString()).To(Equal("union-9"))
	})

	t.Run("Stays null when there is no temp profile", func(t *testing.T) {
		RegisterTestingT(t)
		got := unionIDFrom(types.StringNull(), nil)
		Expect(got.IsNull()).To(BeTrue())
	})

	t.Run("Stays null when the temp profile has no unions", func(t *testing.T) {
		RegisterTestingT(t)
		got := unionIDFrom(types.StringNull(), &geniprofile.Profile{})
		Expect(got.IsNull()).To(BeTrue())
	})
}
