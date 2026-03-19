package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccProject_getProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_project" "test" {
					  id = "project-6"
					}
					`,
				PlanOnly: true,
			},
		},
	})
}

func TestAccProject_addProfileToProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_project" "test" {
					  id = "project-6"
					}

					resource "geni_profile" "test" {
					  names = {
						"en-US" = {
							first_name = "John"
							last_name = "Doe"
						}
					  }
					  alive = false
					  public = true
					  projects = [data.geni_project.test.id]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("projects").AtSliceIndex(0), knownvalue.StringExact("project-6")),
				},
			},
		},
	})
}

func TestAccProject_addProfileToTwoProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_project" "test-6" {
					  id = "project-6"
					}

					data "geni_project" "test-8" {
					  id = "project-8"
					}

					resource "geni_profile" "test" {
					  names = {
						"en-US" = {
							first_name = "John"
							last_name = "Doe"
						}
					  }
					  alive = false
					  public = true
					  projects = [data.geni_project.test-6.id,data.geni_project.test-8.id]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("projects").AtSliceIndex(0), knownvalue.StringExact("project-6")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("projects").AtSliceIndex(1), knownvalue.StringExact("project-8")),
				},
			},
		},
	})
}

func TestAccProject_addDocumentToProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_project" "test" {
					  id = "project-6"
					}

					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is a test document."
					  projects = [data.geni_project.test.id]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("projects").AtSliceIndex(0), knownvalue.StringExact("project-6")),
				},
			},
		},
	})
}
