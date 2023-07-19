package namecheap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_API_KEY", nil),
				Description: "NameCheap API Key.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_API_NAME", nil),
				Description: "NameCheap API Secret.",
			},

			"clientIP": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_CLIENT_IP", nil),
				Description: "Client IP.",
			},

			"baseurl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.namecheap.com/xml.response?",
				Description: "NameCheap Base Url(defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"namecheap_domain_record": resourceDomainRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Key:      d.Get("key").(string),
		Name:     d.Get("name").(string),
		BaseURL:  d.Get("baseurl").(string),
		ClientIP: d.Get("clientIP").(string),
	}

	return config.Client()
}
