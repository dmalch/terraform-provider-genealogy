package geni

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestCreateProfile1(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)
	profileRequest := ProfileRequest{
		Gender: "male",
		Names: map[string]NameElement{
			"en-US": {
				FirstName: "1TestFirstName",
				LastName:  "1TestLastName",
			},
			"ru": {
				FirstName:  "Ф",
				LastName:   "М",
				MiddleName: "Н",
			},
		},
		Birth: &EventElement{
			Date: DateResponse{
				Day:   19,
				Month: 8,
				Year:  1922,
			},
			Location: &LocationElement{
				Country:   "РСФСР",
				County:    "Спасский уезд, Кирилловская волость",
				PlaceName: "село Кирилово",
				State:     "Тамбовская губерния",
			},
		},
		Death: &EventElement{
			Date: DateResponse{
				Day:   25,
				Month: 9,
				Year:  1993,
			},
		},
	}

	profile, err := CreateProfile(testAccessToken, &profileRequest)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.Id).ToNot(BeEmpty())
	Expect(profile.Guid).ToNot(BeEmpty())
	Expect(profile.FirstName).To(BeEquivalentTo("1TestFirstName"))
	Expect(profile.LastName).To(BeEquivalentTo("1TestLastName"))
	Expect(profile.Gender).To(BeEquivalentTo("male"))
	Expect(profile.Names).To(HaveKeyWithValue("en-US", NameElement{
		FirstName: "1TestFirstName",
		LastName:  "1TestLastName",
	}))
	Expect(profile.Names).To(HaveKeyWithValue("ru", NameElement{
		FirstName:  "Ф",
		LastName:   "М",
		MiddleName: "Н",
	}))
	Expect(profile.Birth).To(Equal(&EventElement{
		Date: DateResponse{
			Day:   19,
			Month: 8,
			Year:  1922,
		},
		Location: &LocationElement{
			Country:   "РСФСР",
			County:    "Спасский уезд, Кирилловская волость",
			PlaceName: "село Кирилово",
			State:     "Тамбовская губерния",
		},
		Name: "Birth of 1TestFirstName 1TestLastName",
	}))
	Expect(profile.Death).To(Equal(&EventElement{
		Date: DateResponse{
			Day:   25,
			Month: 9,
			Year:  1993,
		},
		Name: "Death of 1TestFirstName 1TestLastName",
	}))
	Expect(profile.CreatedAt).ToNot(BeEmpty())
}

func TestCreateProfile2(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)
	profileRequest := ProfileRequest{
		//FirstName: `\u0418\u0432\u0430\u043D`,
		FirstName: "Иван",
	}

	profile, err := CreateProfile(testAccessToken, &profileRequest)

	Expect(err).ToNot(HaveOccurred())
	Expect(profile).ToNot(BeNil())
	Expect(profile.Id).ToNot(BeEmpty())
	Expect(profile.Guid).ToNot(BeEmpty())
	Expect(profile.FirstName).To(BeEquivalentTo("Иван"))
}

func TestGetProfile1(t *testing.T) {
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
	Expect(profile.Names).To(HaveKeyWithValue("en-US", NameElement{
		FirstName: "D",
		LastName:  "M",
	}))
	Expect(profile.Names).To(HaveKeyWithValue("ru", NameElement{
		FirstName:  "Д",
		LastName:   "М",
		MiddleName: "В",
	}))
}

func TestGetProfile2(t *testing.T) {
	t.Skip()
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
	Expect(profile.Names).To(HaveKeyWithValue("en-US", NameElement{
		FirstName: "F",
		LastName:  "M",
	}))
	Expect(profile.Names).To(HaveKeyWithValue("ru", NameElement{
		FirstName:  "Ф",
		LastName:   "М",
		MiddleName: "Н",
	}))
	Expect(profile.Birth).To(Equal(&EventElement{
		Date: DateResponse{
			Day:   19,
			Month: 8,
			Year:  1922,
		},
		Location: &LocationElement{
			Country:   "РСФСР",
			County:    "Спасский уезд, Кирилловская волость",
			PlaceName: "село Кирилово",
			State:     "Тамбовская губерния",
		},
		Name: "Birth of F M",
	}))
	Expect(profile.Death).To(Equal(&EventElement{
		Date: DateResponse{
			Day:   25,
			Month: 9,
			Year:  1993,
		},
		Name: "Death of F M",
	}))
	Expect(profile.CreatedAt).To(BeEquivalentTo("1741860385"))
}

func TestDeleteProfile1(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	err := DeleteProfile(testAccessToken, "profile-g599969")

	Expect(err).ToNot(HaveOccurred())
}
