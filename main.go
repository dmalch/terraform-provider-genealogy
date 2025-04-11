package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/dmalch/terraform-provider-genealogy/internal"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/dmalch/genealogy",
		Debug:   debug,
	})

	if err != nil {
		log.Fatal(err)
	}
}
