package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/dmalch/terraform-provider-genealogy/internal"
)

func TestAccProfile_createProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profile(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("maiden_name"), knownvalue.Null()),
				},
			},
		},
	})
}

func profile(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  first_name = "John"
		  last_name  = "Doe"
		}
		`
}

func TestAccProfile_createProfileWithDetails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profileWithDetails(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
				},
			},
		},
	})
}

func profileWithDetails(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  first_name = "John"
		  last_name  = "Doe"
		  gender     = "male"
		  birth      = {
			name = "Birth of John Doe"
			date = {
			  range = "between"
			  year = 1980
			  month = 1
			  day = 1
			  circa = true
			  end_year = 1980
			  end_month = 1
			  end_day = 1
			  end_circa = true
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "Hospital"
			  state = "New York"
			  street_address1 = "123 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }
		  baptism = {
			name = "Baptism of John Doe"
			date = {
			  range = "between"
			  year = 1980
			  month = 1
			  day = 1
			  circa = true
			  end_year = 1980
			  end_month = 1
			  end_day = 1
			  end_circa = true
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "Church"
			  state = "New York"
			  street_address1 = "456 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }
		  death = {
			name = "Death of John Doe at Hospital"
			date = {
			  range = "between"
			  year = 1999
			  month = 1
			  day = 1
			  circa = true
			  end_year = 1999
			  end_month = 2
			  end_day = 1
			  end_circa = true
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "Hospital"
			  state = "New York"
			  street_address1 = "123 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }
		  burial = {
			name = "Burial of John Doe"
			date = {
			  range = "between"
			  year = 1999
			  month = 1
			  day = 1
			  circa = true
			  end_year = 1999
			  end_month = 2
			  end_day = 1
			  end_circa = true
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "Cemetery"
			  state = "New York"
			  street_address1 = "111 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }
		}
		`
}

func TestAccProfile_createProfileWithFixedBithDate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profileWithFixedBirthDate(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
				},
			},
		},
	})
}

func profileWithFixedBirthDate(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  first_name = "John"
		  last_name  = "Doe"
		  gender     = "male"
		  birth      = {
			name = "Birth of John Doe"
			date = {
			  year = 1980
			  month = 1
			  day = 1
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "Hospital"
			  state = "New York"
			  street_address1 = "123 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }
		}
	`
}

func TestAccProfile_createProfileWithNamesInOtherLanguages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profileWithNamesInOtherLanguages(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
				},
			},
		},
	})
}

func profileWithNamesInOtherLanguages(accessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + accessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  gender     = "male"
		  names = {
			"en-US" = {
				first_name = "John"
				last_name = "Doe"
			}
			"ru" = {
				first_name = "Иван"
			}
			"he" = {
				first_name = "יוחנן"
			}
			"ar" = {
				first_name = "يوحنا"
			}
		  }
		}
		`
}

func TestAccProfile_createProfileAndAddNamesInOtherLanguages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profile(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
				},
			},
			{
				Config: profileWithNamesInOtherLanguages(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
				},
			},
		},
	})
}

func TestAccProfile_failToCreateProfileWithBothFirstLastNameAndNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config:      profileWithFirstLastNameAndNames(testAccessToken),
				ExpectError: regexp.MustCompile(`Attribute "names\[\\"en-US\\"]" cannot be specified when "first_name" is\s*specified`),
			},
		},
	})
}

func profileWithFirstLastNameAndNames(accessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + accessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  first_name = "John"
		  last_name  = "Doe"
		  names = {
			"en-US" = {
				first_name = "John"
				last_name = "Doe"
			}
		  }
		}
		`
}

func TestAccProfile_createProfileWithDifferentSetOfNamesInDifferentLanguages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profileWithDifferentSetOfNamesInDifferentLanguages(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("ru").AtMapKey("first_name"), knownvalue.StringExact("Иван")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("ru").AtMapKey("middle_name"), knownvalue.StringExact("Иванович")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("ru").AtMapKey("last_name"), knownvalue.StringExact("Иванов")),
				},
			},
		},
	})
}

func profileWithDifferentSetOfNamesInDifferentLanguages(accessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + accessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  names = {
			"en-US" = {
				first_name = "John"
			}
			"ru" = {
				first_name = "Иван"
				middle_name = "Иванович"
				last_name = "Иванов"
			}
		  }
		}
		`
}

func TestAccProfile_updateProfileWithDifferentSetOfNamesInDifferentLanguages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profileWithOneName(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.Null()),
				},
			},
			{
				Config: profileWithDifferentSetOfNamesInDifferentLanguages(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("ru").AtMapKey("first_name"), knownvalue.StringExact("Иван")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("ru").AtMapKey("middle_name"), knownvalue.StringExact("Иванович")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("ru").AtMapKey("last_name"), knownvalue.StringExact("Иванов")),
				},
			},
		},
	})
}

func profileWithOneName(accessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + accessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  names = {
			"en-US" = {
				first_name = "John"
			}
		  }
		}
		`
}

func TestAccProfile_createProfileWithEmptyMiddleName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: profileWithEmptyMiddleName(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
				},
			},
		},
	})
}

func profileWithEmptyMiddleName(accessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + accessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "test" {
		  names = {
			"en-US" = {
				first_name = "John"
				middle_name = null
				last_name = "Doe"
			}
		  }
		}
		`
}

func TestAccProfile_createProfileWithMiddleNameAndRemoveIt(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}
		
				resource "geni_profile" "test" {
				  names = {
					"en-US" = {
						first_name = "John"
						middle_name = "Johnson"
						last_name = "Doe"
					}
				  }
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
				},
			},
			{
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}
		
				resource "geni_profile" "test" {
				  names = {
					"en-US" = {
						first_name = "John"
						last_name = "Doe"
					}
				  }
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
				},
			},
		},
	})
}

func TestAccProfile_createProfileWithMiddleNameAndRemoveIt2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}
		
				resource "geni_profile" "test" {
				  first_name = "John"
				  middle_name = "Johnson"
				  last_name = "Doe"
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
				},
			},
			{
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}
		
				resource "geni_profile" "test" {
				  first_name = "John"
				  middle_name = null
				  last_name = "Doee"
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doee")),
				},
			},
		},
	})
}
