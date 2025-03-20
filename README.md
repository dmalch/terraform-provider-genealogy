# terraform-provider-genealogy

## Terraform Provider for Geni.com
![img.png](docs/img.png)

## Overview
This provider enables managing data on Geni.com through Terraform. It exposes configuration and resources that help automate genealogical information.  

## Installation
1. Clone or download this repository.  
2. Run `go build` or `go install` inside the repository directory to build the provider binary.  
3. Move the binary into your Terraform plugins directory.  

## Usage
```hcl
terraform {
  required_providers {
    genealogy = {
      source  = "dmalch/genealogy"
      version = "~> 0.1"
    }
  }
}

provider "genealogy" {
  access_token = "your_geni_access_token"
}
```

## Configuration
* `access_token`: The access token used to authenticate against Geni.com.

## Contributing
Pull requests and issues are welcome. Ensure tests pass by running `go test ./...` before submitting changes.

## License
This project is released under a permissive license. Refer to the `LICENSE` file for details.
