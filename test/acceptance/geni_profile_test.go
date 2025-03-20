package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/dmalch/terraform-provider-geni/internal"
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
				},
			},
		},
	})
}

func profile(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
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
			name = "Baptism"
			description = "Baptized in the USA"
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
			name = "Death"
			description = "Died in the USA"
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
			name = "Burial"
			description = "Buried in the USA"
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
