package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var regexpExactlyOne = regexp.MustCompile(`(?s)Invalid Attribute Combination`)

func TestAccDataSourceProfile_byID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_profile" "fixture" {
					  names = {
						"en-US" = {
							first_name = "John"
							last_name  = "Doe"
						}
					  }
					  alive  = false
					  public = true
					}

					data "geni_profile" "by_id" {
					  id = geni_profile.fixture.id
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.geni_profile.by_id", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("data.geni_profile.by_id", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("data.geni_profile.by_id", tfjsonpath.New("public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.geni_profile.by_id", tfjsonpath.New("alive"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccDataSourceProfile_byGUID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "geni_profile" "fixture" {
					  names = {
						"en-US" = {
							first_name = "Jane"
							last_name  = "Roe"
						}
					  }
					  alive  = false
					  public = true
					}

					data "geni_profile" "by_guid" {
					  guid = geni_profile.fixture.guid
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.geni_profile.by_guid", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("Jane")),
					statecheck.ExpectKnownValue("data.geni_profile.by_guid", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Roe")),
				},
			},
		},
	})
}

func TestAccDataSourceProfile_rejectsBothIDAndGUID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_profile" "both" {
					  id   = "profile-1"
					  guid = "abcdef0123456789abcdef0123456789"
					}
					`,
				ExpectError: regexpExactlyOne,
			},
		},
	})
}

func TestAccDataSourceProfile_rejectsNeitherIDNorGUID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_profile" "neither" {}
					`,
				ExpectError: regexpExactlyOne,
			},
		},
	})
}
