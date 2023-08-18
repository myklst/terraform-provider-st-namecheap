package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"os"
	namecheap_provider "terraform-provider-st-namecheap/namecheap"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

/*
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: namecheap.Provider,
	})
}*/

func main() {
	providerAddress := os.Getenv("PROVIDER_LOCAL_PATH")
	if providerAddress == "" {
		providerAddress = "registry.terraform.io/myklst/st-namecheap"
	}

	providerserver.Serve(context.Background(), namecheap_provider.New, providerserver.ServeOpts{
		Address: providerAddress,
	})
}
