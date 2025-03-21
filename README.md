# terraform-provider-genealogy

## Terraform Provider for Geni.com
![img.png](docs/img.png)

## Overview
This provider enables managing data on Geni.com through Terraform. It exposes configuration and resources that help automate genealogical information.

## Disclaimer
This application uses the Geni API but is not endorsed, operated, or sponsored by Geni.com.

## Installation
1. Clone or download this repository.  
2. Run `go build` or `go install` inside the repository directory to build the provider binary.  
3. Move the binary into your Terraform plugins directory.  

## Usage
```hcl
terraform {
  required_providers {
    geni = {
      source  = "dmalch/genealogy"
      version = "~> 0.1"
    }
  }
}

provider "geni" {
  access_token = "your_geni_access_token"
}
```

## Configuration
* `access_token`: The access token used to authenticate against Geni.com.

## Resources

Below is a brief example of adding these resources in the Terraform configuration, demonstrating how to define a `geni_profile` and reference it in a `geni_union`:

```hcl
resource "geni_profile" "mother" {
  first_name = "Jane"
  last_name  = "Doe"
}

resource "geni_profile" "father" {
  first_name = "John"
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
    geni_profile.mother.id,
    geni_profile.father.id,
  ]

  children = [
    geni_profile.child1.id,
    geni_profile.child2.id,
  ]
}
```

## Contributing
Pull requests and issues are welcome. Ensure tests pass by running `go test ./...` before submitting changes.

## License
This project is released under a permissive license. Refer to the `LICENSE` file for details.
