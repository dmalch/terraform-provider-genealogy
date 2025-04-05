package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/dmalch/terraform-provider-genealogy/internal"
)

func TestAccDocument_createDocument(t *testing.T) {
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
			
					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is a test document."
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/plain")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.StringExact("This is a test document.")),
				},
			},
		},
	})
}

func TestAccDocument_createDocumentWithDetails(t *testing.T) {
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
			
					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is a test document."
					  description = "This is a test document description."
					  content_type = "text/plain"
					  date = {
						  year = 1980
						  month = 1
						  day = 1
						  circa = true
					  }
					  location = {
						  city = "New York"
						  county = "Alameda"
						  country = "USA"
						  place_name = "Hospital"
						  state = "New York"
						  street_address1 = "123 Main St"
						  street_address2 = "Apt 1"
						  street_address3 = "Floor 2"
					  }
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/plain")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.StringExact("This is a test document.")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("description"), knownvalue.StringExact("This is a test document description.")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("date"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"year":  knownvalue.Int32Exact(1980),
						"month": knownvalue.Int32Exact(1),
						"day":   knownvalue.Int32Exact(1),
						"circa": knownvalue.Bool(true),
					})),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("location"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"city":            knownvalue.StringExact("New York"),
						"country":         knownvalue.StringExact("USA"),
						"place_name":      knownvalue.StringExact("Hospital"),
						"state":           knownvalue.StringExact("New York"),
						"street_address1": knownvalue.StringExact("123 Main St"),
						"street_address2": knownvalue.StringExact("Apt 1"),
						"street_address3": knownvalue.StringExact("Floor 2"),
						"county":          knownvalue.StringExact("Alameda"),
						"latitude":        knownvalue.Null(),
						"longitude":       knownvalue.Null(),
					})),
				},
			},
		},
	})
}

func TestAccDocument_updateDocumentWithDetails(t *testing.T) {
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
			
					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is a test document."
					  labels = [
						"Military",
					  ]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/plain")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.StringExact("This is a test document.")),
				},
			},
			{
				Config: `
					provider "geni" {
					  access_token = "` + testAccessToken + `"
					  use_sandbox_env = true
					}
			
					resource "geni_document" "test" {
					  title = "Test Document"
					  text = "This is an updated test document."
					  description = "This is a test document description."
					  content_type = "text/plain"
					  date = {
						  year = 1980
						  month = 1
						  day = 1
						  circa = true
					  }
					  location = {
						  city = "New York"
						  county = "Alameda"
						  country = "USA"
						  place_name = "Hospital"
						  state = "New York"
						  street_address1 = "123 Main St"
						  street_address2 = "Apt 1"
						  street_address3 = "Floor 2"
					  }
					  labels = [
						"Census",
						"Military",
					  ]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/plain")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.StringExact("This is an updated test document.")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("description"), knownvalue.StringExact("This is a test document description.")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("date"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"year":  knownvalue.Int32Exact(1980),
						"month": knownvalue.Int32Exact(1),
						"day":   knownvalue.Int32Exact(1),
						"circa": knownvalue.Bool(true),
					})),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("location"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"city":            knownvalue.StringExact("New York"),
						"country":         knownvalue.StringExact("USA"),
						"place_name":      knownvalue.StringExact("Hospital"),
						"state":           knownvalue.StringExact("New York"),
						"street_address1": knownvalue.StringExact("123 Main St"),
						"street_address2": knownvalue.StringExact("Apt 1"),
						"street_address3": knownvalue.StringExact("Floor 2"),
						"county":          knownvalue.StringExact("Alameda"),
						"latitude":        knownvalue.Null(),
						"longitude":       knownvalue.Null(),
					})),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("labels"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("Census"),
						knownvalue.StringExact("Military"),
					})),
				},
			},
		},
	})
}
