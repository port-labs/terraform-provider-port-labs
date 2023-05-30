package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/port-labs/terraform-provider-port-labs/port"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		port.New,
		providerserver.ServeOpts{
			Address: "registry.terraform.io/<namespace>/<provider_name>",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

}

//ProviderFunc: func() *schema.Provider {
// return port.Provider()
// },
