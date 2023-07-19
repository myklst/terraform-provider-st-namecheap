package namecheap

import (
	"fmt"
	"log"
	"terraform-provider-namecheap/api"
)

// Config provides the provider's configuration
type Config struct {
	Key      string
	Name     string
	BaseURL  string
	ClientIP string
}

// Client returns a new client for accessing NameCheap.
func (c *Config) Client() (*api.Client, error) {
	client, err := api.NewClient(c.BaseURL, c.Key, c.Name, c.ClientIP)

	if err != nil {
		return nil, fmt.Errorf("error setting up client: %s", err)
	}

	log.Print("Namecheap Client configured")

	return client, nil
}
