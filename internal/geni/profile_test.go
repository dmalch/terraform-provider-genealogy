package geni

import (
	"testing"

	. "github.com/onsi/gomega"
)

func ptr[T any](s T) *T {
	return &s
}

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
			Date: &DateElement{
				Day:   19,
				Month: 8,
				Year:  1922,
			},
			Location: &LocationElement{
				Country:   ptr("РСФСР"),
				County:    ptr("Спасский уезд, Кирилловская волость"),
				PlaceName: ptr("село Кирилово"),
				State:     ptr("Тамбовская губерния"),
			},
		},
		Death: &EventElement{
			Date: &DateElement{
				Day:   25,
				Month: 9,
				Year:  1993,
			},
		},
	}
	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	profile, err := client.CreateProfile(&profileRequest)

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
		Date: &DateElement{
			Day:   19,
			Month: 8,
			Year:  1922,
		},
		Location: &LocationElement{
			Country:   ptr("РСФСР"),
			County:    ptr("Спасский уезд, Кирилловская волость"),
			PlaceName: ptr("село Кирилово"),
			State:     ptr("Тамбовская губерния"),
		},
		Name: "Birth of 1TestFirstName 1TestLastName",
	}))
	Expect(profile.Death).To(Equal(&EventElement{
		Date: &DateElement{
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

	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	profile, err := client.CreateProfile(&profileRequest)

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

	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	profile, err := client.GetProfile(profileId)

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
	Expect(profile.Unions).To(ContainElement("union-1837"))
}

func TestGetProfile2(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	profileId := "profile-5957"

	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	profile, err := client.GetProfile(profileId)

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
		Date: &DateElement{
			Day:   19,
			Month: 8,
			Year:  1922,
		},
		Location: &LocationElement{
			Country:   ptr("РСФСР"),
			County:    ptr("Спасский уезд, Кирилловская волость"),
			PlaceName: ptr("село Кирилово"),
			State:     ptr("Тамбовская губерния"),
		},
		Name: "Birth of F M",
	}))
	Expect(profile.Death).To(Equal(&EventElement{
		Date: &DateElement{
			Day:   25,
			Month: 9,
			Year:  1993,
		},
		Name: "Death of F M",
	}))
	Expect(profile.CreatedAt).To(BeEquivalentTo("1741860385"))
}

func TestUpdateProfile1(t *testing.T) {
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
			Date: &DateElement{
				Day:   19,
				Month: 8,
				Year:  1922,
			},
			Location: &LocationElement{
				Country:   ptr("РСФСР"),
				County:    ptr("Спасский уезд, Кирилловская волость"),
				PlaceName: ptr("село Кирилово"),
				State:     ptr("Тамбовская губерния"),
			},
		},
		Death: &EventElement{
			Date: &DateElement{
				Day:   25,
				Month: 9,
				Year:  1993,
			},
		},
	}

	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	profile, err := client.CreateProfile(&profileRequest)
	Expect(err).ToNot(HaveOccurred())

	profileRequest.FirstName = "2TestFirstName"
	updatedProfile, err := client.UpdateProfile(profile.Id, &profileRequest)
	Expect(err).ToNot(HaveOccurred())

	Expect(updatedProfile).ToNot(BeNil())
	Expect(updatedProfile.Id).ToNot(BeEmpty())
	Expect(updatedProfile.Guid).ToNot(BeEmpty())
	Expect(updatedProfile.FirstName).To(BeEquivalentTo("2TestFirstName"))
	Expect(updatedProfile.LastName).To(BeEquivalentTo("1TestLastName"))
	Expect(updatedProfile.Gender).To(BeEquivalentTo("male"))
	Expect(updatedProfile.Names).To(HaveKeyWithValue("en-US", NameElement{
		FirstName: "2TestFirstName",
		LastName:  "1TestLastName",
	}))
	Expect(updatedProfile.Names).To(HaveKeyWithValue("ru", NameElement{
		FirstName:  "Ф",
		LastName:   "М",
		MiddleName: "Н",
	}))
	Expect(updatedProfile.Birth).To(Equal(&EventElement{
		Date: &DateElement{
			Day:   19,
			Month: 8,
			Year:  1922,
		},
		Location: &LocationElement{
			Country:   ptr("РСФСР"),
			County:    ptr("Спасский уезд, Кирилловская волость"),
			PlaceName: ptr("село Кирилово"),
			State:     ptr("Тамбовская губерния"),
		},
		Name: "Birth of 2TestFirstName 1TestLastName",
	}))
	Expect(updatedProfile.Death).To(Equal(&EventElement{
		Date: &DateElement{
			Day:   25,
			Month: 9,
			Year:  1993,
		},
		Name: "Death of 2TestFirstName 1TestLastName",
	}))
	Expect(updatedProfile.CreatedAt).ToNot(BeEmpty())
}

func TestDeleteProfile1(t *testing.T) {
	t.Skip()
	RegisterTestingT(t)

	client, err := NewClient(testAccessToken, true)
	Expect(err).ToNot(HaveOccurred())

	err = client.DeleteProfile("profile-g599969")

	Expect(err).ToNot(HaveOccurred())
}
