package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/dmalch/terraform-provider-geni/internal"
)

func main() {
	err := providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/dmalch/geni",
	})
	if err != nil {
		log.Fatal(err)
	}
}
