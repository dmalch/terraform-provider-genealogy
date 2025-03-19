package acceptance

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

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
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.husband", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.wife", tfjsonpath.New("id"), compare.ValuesSame()),
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

func TestAccExampleWidget_createUnionWithTwoPartnersAndChild(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersAndChild(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.husband", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.wife", tfjsonpath.New("first_name"), knownvalue.StringExact("Jane")),
					statecheck.ExpectKnownValue("geni_profile.child", tfjsonpath.New("first_name"), knownvalue.StringExact("Alice")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.husband", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.wife", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(1)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.child", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}

func unionWithTwoPartnersAndChild(testAccessToken string) string {
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
		
		resource "geni_profile" "child" {
		  first_name = "Alice"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.husband.id,
			geni_profile.wife.id,
		  ]
		
		  children = [
			geni_profile.child.id,
		  ]
		}
		`
}

func TestAccExampleWidget_createUnionWithTwoPartnersTwoChildren(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersTwoChildren(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.husband", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.wife", tfjsonpath.New("first_name"), knownvalue.StringExact("Jane")),
					statecheck.ExpectKnownValue("geni_profile.child1", tfjsonpath.New("first_name"), knownvalue.StringExact("Alice")),
					statecheck.ExpectKnownValue("geni_profile.child2", tfjsonpath.New("first_name"), knownvalue.StringExact("Bob")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.husband", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.wife", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.child1", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.child2", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}

func unionWithTwoPartnersTwoChildren(testAccessToken string) string {
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
		
		resource "geni_profile" "child1" {
		  first_name = "Alice"
		  last_name  = "Doe"
		}
		
		resource "geni_profile" "child2" {
		  first_name = "Bob"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.husband.id,
			geni_profile.wife.id,
		  ]
		
		  children = [
			geni_profile.child1.id,
			geni_profile.child2.id,
		  ]
		}
		`
}

func TestAccExampleWidget_createUnionWithOneParentAndChild(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithOneParentAndChild(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.mother", tfjsonpath.New("first_name"), knownvalue.StringExact("Jane")),
					statecheck.ExpectKnownValue("geni_profile.child", tfjsonpath.New("first_name"), knownvalue.StringExact("Alice")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(1)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.mother", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(1)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.child", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}

func unionWithOneParentAndChild(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		}

		resource "geni_profile" "mother" {
		  first_name = "Jane"
		  last_name  = "Doe"
		}
		
		resource "geni_profile" "child" {
		  first_name = "Alice"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.mother.id,
		  ]
		
		  children = [
			geni_profile.child.id,
		  ]
		}
		`
}

func TestAccExampleWidget_createUnionWithTwoSiblingsWithoutParents(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoSiblingsWithoutParents(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.sibling1", tfjsonpath.New("first_name"), knownvalue.StringExact("Alice")),
					statecheck.ExpectKnownValue("geni_profile.sibling2", tfjsonpath.New("first_name"), knownvalue.StringExact("Bob")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.sibling1", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.sibling2", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}

func unionWithTwoSiblingsWithoutParents(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		}

		resource "geni_profile" "sibling1" {
		  first_name = "Alice"
		  last_name  = "Doe"
		}
		
		resource "geni_profile" "sibling2" {
		  first_name = "Bob"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		  children = [
			geni_profile.sibling1.id,
			geni_profile.sibling2.id,
		  ]
		}
		`
}

func TestAccExampleWidget_createUnionWithTwoSiblingsAndAddParentsInTheSecondStep(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoSiblingsWithoutParents(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.sibling1", tfjsonpath.New("first_name"), knownvalue.StringExact("Alice")),
					statecheck.ExpectKnownValue("geni_profile.sibling2", tfjsonpath.New("first_name"), knownvalue.StringExact("Bob")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.sibling1", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("children")},
						"geni_profile.sibling2", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
			{
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				}

				resource "geni_profile" "sibling1" {
				  first_name = "Alice"
				  last_name  = "Doe"
				}
				
				resource "geni_profile" "sibling2" {
				  first_name = "Bob"
				  last_name  = "Doe"
				}
		
				resource "geni_profile" "mother" {
				  first_name = "Jane"
				  last_name  = "Doe"
				}

				resource "geni_profile" "father" {
				  first_name = "John"
				  last_name  = "Doe"
				}
				
				resource "geni_union" "doe_family" {
				  partners = [
					geni_profile.mother.id,
					geni_profile.father.id,
				  ]

				  children = [
					geni_profile.sibling1.id,
					geni_profile.sibling2.id,
				  ]
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.mother", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.father", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(2)),
				},
			},
		},
	})
}

func TestAccExampleWidget_failToCreateUnionWithOneParent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config:      unionWithOneParent(testAccessToken),
				ExpectError: regexp.MustCompile(`Insufficient Attribute Configuration`),
			},
		},
	})
}

func unionWithOneParent(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		}

		resource "geni_profile" "mother" {
		  first_name = "Jane"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.mother.id,
		  ]
		}
		`
}

func TestAccExampleWidget_failToCreateUnionWithOneChild(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config:      unionWithOneChild(testAccessToken),
				ExpectError: regexp.MustCompile(`Insufficient Attribute Configuration`),
			},
		},
	})
}

func unionWithOneChild(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		}

		resource "geni_profile" "child" {
		  first_name = "Alice"
		  last_name  = "Doe"
		}
		
		resource "geni_union" "doe_family" {
		  children = [
			geni_profile.child.id,
		  ]
		}
		`
}
