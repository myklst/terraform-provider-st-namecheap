package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	namecheap "terraform-provider-st-namecheap/namecheap"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: namecheap.Provider,
	})
}
