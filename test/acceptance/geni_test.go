package acceptance

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"testing"

	"github.com/dmalch/terraform-provider-geni/internal"
)

func TestAccExampleWidget_createProfile(t *testing.T) {
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

func TestAccExampleWidget_createUnionWithTwoPartners(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartners(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.husband", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.wife", tfjsonpath.New("first_name"), knownvalue.StringExact("Jane")),
				},
			},
		},
	})
}

func unionWithTwoPartners(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		}

		resource "geni_profile" "husband" {
		  first_name = "John"
		  last_name  = "Doe"
		}
		
		resource "geni_profile" "wife" {
		  first_name = "Jane"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		 partners = [
			geni_profile.husband.id,
			geni_profile.wife.id,
		 ]
		}
		`
}
