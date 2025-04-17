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
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("about"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("alive"), knownvalue.Bool(false)),
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
		  alive = false
		  public = true
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
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.StringExact("Smith")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.StringExact("John Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("about"), knownvalue.StringExact("This is a test profile")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("alive"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").AtMapKey("name"), knownvalue.StringExact("Birth of John Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").AtMapKey("date").AtMapKey("year"), knownvalue.Int32Exact(1980)),
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
		  middle_name = "Johnson"
		  birth_last_name = "Smith"
		  display_name = "John Doe"
		  about = "This is a test profile"
		  public = true
		  alive = false
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
			  latitude = 55.8948313,
			  longitude = 44.0386238,
			  //postal_code = "606302",
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

func TestAccProfile_createProfileWithDeathDetails(t *testing.T) {
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
					  last_name  = "Doe"
					  middle_name = "Johnson"
					  birth_last_name = "Smith"
					  display_name = "John Doe"
					  about = "This is a test profile"
					  public = true
					  alive = false
					  gender= "male"
					  cause_of_death = "natural"
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
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.StringExact("Smith")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.StringExact("John Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("about"), knownvalue.StringExact("This is a test profile")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("alive"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("death").AtMapKey("name"), knownvalue.StringExact("Death of John Doe at Hospital")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("death").AtMapKey("date").AtMapKey("year"), knownvalue.Int32Exact(1999)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("cause_of_death"), knownvalue.StringExact("natural")),
				},
			},
		},
	})
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
		  alive = false
		  public = true
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

func TestAccProfile_createProfileWithEmptyBirthLocation(t *testing.T) {
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
					  last_name  = "Doe"
					  gender     = "male"
					  alive = false
					  public = true
					  birth      = {
						location = {
						}
					  }
					}
				`,
				ExpectError: regexp.MustCompile(`birth.location.city" must be specified when "birth.location" is\s*specified`),
			},
		},
	})
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
		  alive = false
		  public = true
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
		  alive = false
		  public = true
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
		  alive = false
		  public = true
		  names = {
			"en-US" = {
				first_name = "John"
			}
			"ru" = {
				first_name = "Иван"
				middle_name = "Иванович"
				last_name = "Иванов"
				birth_last_name = "Петров"
				display_name = "Иван Иванович Иванов"
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
		  alive = false
		  public = true
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
		  alive = false
		  public = true
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

func TestAccProfile_updateProfileAliveStatus(t *testing.T) {
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
						last_name  = "Doe"
						alive = true
		  				public = false
					}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("alive"), knownvalue.Bool(true)),
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
						last_name  = "Doe"
						alive = false
		  				public = true
					}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("alive"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccProfile_updateProfilePublicStatus(t *testing.T) {
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
						last_name  = "Doe"
						alive = false
						public = true
					}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("public"), knownvalue.Bool(true)),
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
						last_name  = "Doe"
						alive = false
						public = false
					}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("public"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccProfile_createProfileWithMiddleNameAndRemoveIt1(t *testing.T) {
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
				  alive = false
				  public = true
				  names = {
					"en-US" = {
					  first_name = "John"
					  middle_name = "Johnson"
					  last_name = "Doe"
					  birth_last_name = "Smith"
					  display_name = "John Doe"
					}
				  }
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.StringExact("Smith")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.StringExact("John Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("birth_last_name"), knownvalue.StringExact("Smith")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("display_name"), knownvalue.StringExact("John Doe")),
				},
			},
			{
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}
		
				resource "geni_profile" "test" {
				  alive = false
				  public = true
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
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("birth_last_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("display_name"), knownvalue.Null()),
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
				  birth_last_name = "Smith"
				  display_name = "John Doe"
				  alive = false
				  public = true
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.StringExact("Johnson")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.StringExact("Smith")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.StringExact("John Doe")),
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
				  last_name = "Doee"
				  alive = false
				  public = true
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("middle_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doee")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth_last_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("display_name"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccProfile_updateProfileLocation(t *testing.T) {
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
					  last_name  = "Doe"
					  gender     = "male"
					  alive = false
					  public = true
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
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("city"), knownvalue.StringExact("New York")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("country"), knownvalue.StringExact("USA")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("place_name"), knownvalue.StringExact("Hospital")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("state"), knownvalue.StringExact("New York")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("street_address1"), knownvalue.StringExact("123 Main St")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("street_address2"), knownvalue.StringExact("Apt 1")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("street_address3"), knownvalue.StringExact("Floor 2")),
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
					  last_name  = "Doe"
					  gender     = "male"
					  alive = false
					  public = true
					  birth      = {
						name = "Birth of John Doe"
						date = {
						  year = 1980
						  month = 1
						  day = 1
						}
						location = {
						  place_name = "New Place"
						}
					  }
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("city"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("country"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("place_name"), knownvalue.StringExact("New Place")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("state"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("street_address1"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("street_address2"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("location").AtMapKey("street_address3"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccProfile_updateProfileDate(t *testing.T) {
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
					  last_name  = "Doe"
					  gender     = "male"
					  alive = false
					  public = true
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
					  }
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("range"), knownvalue.StringExact("between")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("year"), knownvalue.Int32Exact(1980)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("month"), knownvalue.Int32Exact(1)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("day"), knownvalue.Int32Exact(1)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("circa"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_year"), knownvalue.Int32Exact(1980)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_month"), knownvalue.Int32Exact(1)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_day"), knownvalue.Int32Exact(1)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_circa"), knownvalue.Bool(true)),
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
					  last_name  = "Doe"
					  gender     = "male"
					  alive = false
					  public = true
					  birth      = {
						name = "Birth of John Doe"
						date = {
						  year = 1980
						}
					  }
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("range"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("year"), knownvalue.Int32Exact(1980)),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("month"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("day"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("circa"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_year"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_month"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_day"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("birth").
						AtMapKey("date").AtMapKey("end_circa"), knownvalue.Null()),
				},
			},
		},
	})
}
