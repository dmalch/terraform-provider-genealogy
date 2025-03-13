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

func TestGetProfile2(t *testing.T) {
	//t.Skip()
	RegisterTestingT(t)

	profileId := "profile-5957"
	profile, err := GetProfile(testAccessToken, profileId)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.Id).To(BeEquivalentTo(profileId))
	Expect(profile.Guid).To(BeEquivalentTo("598352"))
	Expect(profile.FirstName).To(BeEquivalentTo("F"))
	Expect(profile.LastName).To(BeEquivalentTo("M"))
	Expect(profile.Gender).To(BeEquivalentTo("male"))
	Expect(profile.Names).To(HaveKeyWithValue("en-US", NameResponse{
		FirstName: "F",
		LastName:  "M",
	}))
	Expect(profile.Names).To(HaveKeyWithValue("ru", NameResponse{
		FirstName:  "Ф",
		LastName:   "М",
		MiddleName: "Н",
	}))
	Expect(profile.Birth).To(Equal(&EventResponse{
		Date: DateResponse{
			Day:   19,
			Month: 8,
			Year:  1922,
		},
		Name: "Birth of F M",
	}))
	Expect(profile.Death).To(Equal(&EventResponse{
		Date: DateResponse{
			Day:   25,
			Month: 9,
			Year:  1993,
		},
		Name: "Death of F M",
	}))
	Expect(profile.CreatedAt).To(BeEquivalentTo("1741860385"))
}
