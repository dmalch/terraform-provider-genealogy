package acceptance

import (
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

func TestAccDocument_createTextDocument(t *testing.T) {
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccDocument_createTextDocumentWithDetails(t *testing.T) {
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
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

func TestAccDocument_updateTextDocumentWithDetails(t *testing.T) {
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
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

func TestAccDocument_createPngDocument(t *testing.T) {
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
					  file = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name = "cs-white-fff.png"
					  content_type = "image/png"
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("image/png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.StringExact("cs-white-fff.png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
				},
			},
		},
	})
}

func TestAccDocument_createPngDocumentWithDetails(t *testing.T) {
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
					  file = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name = "cs-white-fff.png"
					  content_type = "image/png"	
					  description = "This is a test document description."
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("image/png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.StringExact("cs-white-fff.png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
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

func TestAccDocument_updatePngDocumentWithDetails(t *testing.T) {
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
					  file = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name = "cs-white-fff.png"
					  content_type = "image/png"
					  labels = [
						"Military",
					  ]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("image/png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.StringExact("cs-white-fff.png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
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
					  file = filebase64("${path.module}/assets/cs-white-fff.png")
					  file_name = "cs-white-fff.png"
					  content_type = "image/png"
					  description = "This is a test document description."
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("image/png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.StringExact("cs-white-fff.png")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.Null()),
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

func TestAccDocument_createUrlDocument(t *testing.T) {
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
					  source_url = "https://example.com"
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/html")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.StringExact("https://example.com")),
				},
			},
		},
	})
}

func TestAccDocument_createUrlDocumentWithDetails(t *testing.T) {
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
					  source_url = "https://example.com"
					  description = "This is a test document description."
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/html")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.StringExact("https://example.com")),
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

func TestAccDocument_createUrlDocumentWithPerson(t *testing.T) {
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
					}
			
					resource "geni_document" "test" {
					  title = "Test Document"
					  source_url = "https://example.com"
					  profiles = [
						geni_profile.test.id,
					  ]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("first_name"), knownvalue.StringExact("John")),
					statecheck.ExpectKnownValue("geni_profile.test", tfjsonpath.New("last_name"), knownvalue.StringExact("Doe")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/html")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.StringExact("https://example.com")),
					statecheck.CompareValueCollection("geni_document.test", []tfjsonpath.Path{tfjsonpath.New("profiles")},
						"geni_profile.test", tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}

func TestAccDocument_updateUrlDocumentWithDetails(t *testing.T) {
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
					  source_url = "https://example.com"
					  labels = [
						"Military",
					  ]
					}
					`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("title"), knownvalue.StringExact("Test Document")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/html")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.StringExact("https://example.com")),
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
					  source_url = "https://example.com"
					  description = "This is a test document description."
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
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("content_type"), knownvalue.StringExact("text/html")),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("text"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("file_name"), knownvalue.Null()),
					statecheck.ExpectKnownValue("geni_document.test", tfjsonpath.New("source_url"), knownvalue.StringExact("https://example.com")),
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
