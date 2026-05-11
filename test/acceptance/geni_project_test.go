package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccProject_getProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_project" "test" {
					  id = "project-8"
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
					  projects = [data.geni_project.test.id]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("projects").AtSliceIndex(0), knownvalue.StringExact("project-8")),
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
					data "geni_project" "test-8" {
					  id = "project-8"
					}

					data "geni_project" "test-9" {
					  id = "project-9"
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
					  projects = [data.geni_project.test-8.id,data.geni_project.test-9.id]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("names").AtMapKey("en-US").AtMapKey("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("projects"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("project-8"),
						knownvalue.StringExact("project-9"),
					})),
				},
			},
		},
	})
}

func TestAccProject_importProfilePopulatesProjects(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "geni_project" "test" {
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
					  projects = [data.geni_project.test.id]
					}
					`,
			},
			{
				ResourceName:    "geni_profile.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					attrs := states[0].Attributes
					if attrs["projects.#"] != "1" {
						return fmt.Errorf("expected projects.# = 1, got %q", attrs["projects.#"])
					}
					if attrs["projects.0"] != "project-8" {
						return fmt.Errorf("expected projects.0 = project-8, got %q", attrs["projects.0"])
					}
					return nil
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
					  id = "project-8"
					}

					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is a test document."
					  projects = [data.geni_project.test.id]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("projects").AtSliceIndex(0), knownvalue.StringExact("project-8")),
				},
			},
		},
	})
}
