# Terraform Provider for Geni.com

## Overview

This provider enables managing data on Geni.com through Terraform. It exposes configuration and resources that help
automate genealogical information.

## Disclaimer

This application uses the Geni API but is not endorsed, operated, or sponsored by Geni.com.

## Usage

```hcl
terraform {
  required_providers {
    geni = {
      source  = "dmalch/genealogy"
      version = "~> 0.11"
    }
  }
}

provider "geni" {
}
```

## Configuration

* `access_token`: (Optional) The access token used to authenticate against Geni.com. If not provided, the provider will
  attempt to do a client-side OAuth flow to obtain one.
* `use_sandbox_env`: (Optional) Use the Geni sandbox environment. Default is `false`.
* `use_profile_cache` (Optional) Whether to use the profile cache for faster lookups. It preloads all profiles managed
  by the current user, which may be slow for those with many profiles. Not recommended for use with the `-target` flag.
* `use_document_cache` (Optional) Whether to use the document cache for faster lookups. It preloads all documents
  uploaded by the current user, which may be slow for those with many documents. Not recommended for use with the
  `-target` flag.

## Resources

Below is a brief example of adding these resources in the Terraform configuration, demonstrating how to define a
`geni_profile` and reference it in a `geni_union`:

```hcl
resource "geni_profile" "mother" {
  names = {
    "en-US" = {
      first_name      = "Jane"
      last_name       = "Doe"
      birth_last_name = "Brown"
    }
  }
}

resource "geni_profile" "father" {
  names = {
    "en-US" = {
      first_name  = "John"
      middle_name = "Smith"
      last_name   = "Doe"
    }
    "ru" = {
      first_name  = "Иван"
      middle_name = "Иванович"
      last_name   = "Иванов"
    }
    "es" = {
      first_name  = "Juan"
      middle_name = "Pérez"
      last_name   = "García"
    }
  }
}

resource "geni_profile" "child1" {
  names = {
    "en-US" = {
      first_name   = "Alice"
      last_name    = "Doe"
      display_name = "Alicia"
    }
  }
}

resource "geni_profile" "child2" {
  names = {
    "en-US" = {
      first_name   = "Bob"
      last_name    = "Doe"
      display_name = "Bobby"
    }
  }
}

resource "geni_union" "doe_family" {
  partners = [
    geni_profile.mother.id,
    geni_profile.father.id,
  ]

  children = [
    geni_profile.child1.id,
    geni_profile.child2.id,
  ]
}

resource "geni_document" "example" {
  title       = "Birth Certificate"
  description = "This is a birth certificate document."
  source_url  = "https://example.com/document.pdf"
  profiles = [
    geni_profile.child1.id
  ]
}
```

## Data Sources

```hcl
data "geni_project" "example" {
  # Provide a valid project ID or lookup
  # You can reference this data source in other resources
  project_id = "project-67890"
}
```

## Documentation

For the full provider documentation, see
the [Terraform Registry Documentation](https://registry.terraform.io/providers/dmalch/genealogy/latest/docs).

## Contributing

Pull requests and issues are welcome. Ensure tests pass by running `go test ./...` before submitting changes.

## License

This project is released under a permissive license. Refer to the `LICENSE` file for details.
