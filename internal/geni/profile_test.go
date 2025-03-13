package geni

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetProfile(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	profileId := "profile-5955"
	profile, err := GetProfile(testAccessToken, profileId)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.Id).To(BeEquivalentTo(profileId))
	Expect(profile.FirstName).To(BeEquivalentTo("D"))
	Expect(profile.LastName).To(BeEquivalentTo("M"))
	Expect(profile.Gender).To(BeEquivalentTo("male"))
	Expect(profile.Names).To(HaveKeyWithValue("en-US", NameResponse{
		FirstName: "D",
		LastName:  "M",
	}))
	Expect(profile.Names).To(HaveKeyWithValue("ru", NameResponse{
		FirstName:  "Д",
		LastName:   "М",
		MiddleName: "В",
	}))
}
