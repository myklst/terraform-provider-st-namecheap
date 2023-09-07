package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	provider "github.com/myklst/terraform-provider-st-namecheap/namecheap"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	providerAddress := os.Getenv("PROVIDER_LOCAL_PATH")
	if providerAddress == "" {
		providerAddress = "registry.terraform.io/myklst/st-namecheap"
	}

	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: providerAddress,
	})
}
