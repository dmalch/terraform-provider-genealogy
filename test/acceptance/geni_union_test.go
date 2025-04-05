package acceptance

import (
	"math/big"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/dmalch/terraform-provider-genealogy/internal"
)

func TestAccUnion_createUnionWithTwoPartners(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "husband" {
		  first_name = "John"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "wife" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.husband.id,
			geni_profile.wife.id,
		  ]
		}
		`
}

func TestAccUnion_createUnionWithTwoPartnersAndChild(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "husband" {
		  first_name = "John"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "wife" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "child" {
		  first_name = "Alice"
		  last_name  = "Doe"
		  alive = false
		  public = true
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

func TestAccUnion_createUnionWithTwoPartnersTwoChildren(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "husband" {
		  first_name = "John"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "wife" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "child1" {
		  first_name = "Alice"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "child2" {
		  first_name = "Bob"
		  last_name  = "Doe"
		  alive = false
		  public = true
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

func TestAccUnion_createUnionWithOneParentAndChild(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "mother" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "child" {
		  first_name = "Alice"
		  last_name  = "Doe"
		  alive = false
		  public = true
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

func TestAccUnion_createUnionWithTwoSiblingsWithoutParents(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "sibling1" {
		  first_name = "Alice"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "sibling2" {
		  first_name = "Bob"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  children = [
			geni_profile.sibling1.id,
			geni_profile.sibling2.id,
		  ]
		}
		`
}

func TestAccUnion_createUnionWithTwoSiblingsAndAddParentsInTheSecondStep(t *testing.T) {
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
				  use_sandbox_env = true
				}

				resource "geni_profile" "sibling1" {
				  first_name = "Alice"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
				
				resource "geni_profile" "sibling2" {
				  first_name = "Bob"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "mother" {
				  first_name = "Jane"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}

				resource "geni_profile" "father" {
				  first_name = "John"
				  last_name  = "Doe"
				  alive = false
				  public = true
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

func TestAccUnion_failToCreateUnionWithOneParent(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "mother" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.mother.id,
		  ]
		}
		`
}

func TestAccUnion_failToCreateUnionWithOneChild(t *testing.T) {
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
		  use_sandbox_env = true
		}

		resource "geni_profile" "child" {
		  first_name = "Alice"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  children = [
			geni_profile.child.id,
		  ]
		}
		`
}

func TestAccUnion_failToCreateUnionWithThreePartners(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config:      unionWithThreePartners(testAccessToken),
				ExpectError: regexp.MustCompile(`Too Many Partners`),
			},
		},
	})
}

func unionWithThreePartners(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "partner1" {
		  first_name = "John"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "partner2" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "partner3" {
		  first_name = "Alice"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.partner1.id,
			geni_profile.partner2.id,
			geni_profile.partner3.id,
		  ]
		}
		`
}

func TestAccUnion_failToAddThirdPartnerToUnion(t *testing.T) {
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
			{
				// Try to add a third partner to the union
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}

				resource "geni_profile" "husband" {
				  first_name = "John"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
			
				resource "geni_profile" "wife" {
				  first_name = "Jane"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "partner3" {
				  first_name = "Alice"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
				
				resource "geni_union" "doe_family" {
				  partners = [
					geni_profile.husband.id,
					geni_profile.wife.id,
					geni_profile.partner3.id,
				  ]
				}
				`,
				ExpectError: regexp.MustCompile(`Too Many Partners`),
			},
			{
				// Revert back to the original state
				Config: unionWithTwoPartners(testAccessToken),
			},
		},
	})
}

func TestAccUnion_failToRemovePartnerFromUnion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersAndChild(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(1)),
				},
			},
			{
				// Try to remove a partner from the union
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}

				resource "geni_profile" "husband" {
				  first_name = "John"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "wife" {
				  first_name = "Jane"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "child" {
				  first_name = "Alice"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
				
				resource "geni_union" "doe_family" {
				  partners = [
					geni_profile.husband.id,
				  ]
				
				  children = [
					geni_profile.child.id,
				  ]
				}
				`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccUnion_failToRemoveChildrenFromUnion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersAndChild(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(1)),
				},
			},
			{
				// Try to remove a child from the union
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}

				resource "geni_profile" "husband" {
				  first_name = "John"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "wife" {
				  first_name = "Jane"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "child" {
				  first_name = "Alice"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
				
				resource "geni_union" "doe_family" {
				  partners = [
					geni_profile.husband.id,
					geni_profile.wife.id,
				  ]
				}
				`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccUnion_addAnotherChildToUnion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersAndChild(testAccessToken),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(1)),
				},
			},
			{
				// Add another child to the union
				Config: `
				provider "geni" {
				  access_token = "` + testAccessToken + `"
				  use_sandbox_env = true
				}

				resource "geni_profile" "husband" {
				  first_name = "John"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "wife" {
				  first_name = "Jane"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "child1" {
				  first_name = "Alice"
				  last_name  = "Doe"
				  alive = false
				  public = true
				}
		
				resource "geni_profile" "child2" {
				  first_name = "Bob"
				  last_name  = "Doe"
				  alive = false
				  public = true
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
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("children"), knownvalue.SetSizeExact(2)),
				},
			},
		},
	})
}

func TestAccUnion_createUnionWithTwoPartnersAndDetails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersAndDetails(testAccessToken),
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

func unionWithTwoPartnersAndDetails(testAccessToken string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "husband" {
		  first_name = "John"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "wife" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.husband.id,
			geni_profile.wife.id,
		  ]

		  marriage = {
			name = "Marriage of John and Jane Doe"
			date = {
			  range = "between"
			  year = 1980
			  month = 1
			  day = 1
			  circa = true
			  end_year = 1980
			  end_month = 1
			  end_day = 2
			  end_circa = true
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "City Hall"
			  state = "New York"
			  street_address1 = "123 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }

		  divorce = {
			name = "Divorce of John and Jane Doe"
			date = {
			  range = "between"
			  year = 1980
			  month = 1
			  day = 1
			  circa = true
			  end_year = 1980
			  end_month = 1
			  end_day = 2
			  end_circa = true
			}
			location = {
			  city = "New York"
			  country = "USA"
			  place_name = "City Hall"
			  state = "New York"
			  street_address1 = "123 Main St"
			  street_address2 = "Apt 1"
			  street_address3 = "Floor 2"
			}
		  }
		}
		`
}

func TestAccUnion_updateUnionDetails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		//IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"geni": providerserver.NewProtocol6WithError(internal.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: unionWithTwoPartnersAndMarriageDetails(testAccessToken, "1980"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.husband", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.wife", tfjsonpath.New("first_name"), knownvalue.StringExact("Jane")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.husband", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.wife", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("marriage").AtMapKey("date").AtMapKey("year"),
						knownvalue.NumberExact(big.NewFloat(1980))),
				},
			},
			{
				Config: unionWithTwoPartnersAndMarriageDetails(testAccessToken, "1981"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.husband", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.wife", tfjsonpath.New("first_name"), knownvalue.StringExact("Jane")),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("partners"), knownvalue.SetSizeExact(2)),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.husband", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValueCollection("geni_union.doe_family", []tfjsonpath.Path{tfjsonpath.New("partners")},
						"geni_profile.wife", tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue("geni_union.doe_family", tfjsonpath.New("marriage").AtMapKey("date").AtMapKey("year"),
						knownvalue.NumberExact(big.NewFloat(1981))),
				},
			},
		},
	})
}

func unionWithTwoPartnersAndMarriageDetails(testAccessToken string, year string) string {
	return `
		provider "geni" {
		  access_token = "` + testAccessToken + `"
		  use_sandbox_env = true
		}

		resource "geni_profile" "husband" {
		  first_name = "John"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_profile" "wife" {
		  first_name = "Jane"
		  last_name  = "Doe"
		  alive = false
		  public = true
		}
		
		resource "geni_union" "doe_family" {
		  partners = [
			geni_profile.husband.id,
			geni_profile.wife.id,
		  ]

		  marriage = {
			date = {
			  year = ` + year + `
			  month = 1
			  day = 1
			}
		  }
		}
		`
}
