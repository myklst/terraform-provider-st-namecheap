package namecheap_provider

import (
	"context"
	"github.com/agent-tao/go-namecheap-sdk/v2/namecheap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _info namecheap.DomainCreateInfo

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A registered user name for namecheap",
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_USER_NAME", nil),
			},

			"api_user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A registered api user for namecheap",
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_API_USER", nil),
			},

			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The namecheap API key",
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_API_KEY", nil),
			},

			"client_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client IP address",
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_CLIENT_IP", nil),
				Default:     "0.0.0.0",
			},

			"use_sandbox": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use sandbox API endpoints",
				DefaultFunc: schema.EnvDefaultFunc("NAMECHEAP_USE_SANDBOX", false),
			},

			//--------------------------------------------------------------------------------

			"years": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of years to register",
				Default:     "2",
			},

			"registrant_firstname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First name of the Registrant user",
				Default:     "",
			},

			"registrant_lastname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second name of the Registrant user",
				Default:     "",
			},

			"registrant_address1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address1 of the Registrant user",
				Default:     "",
			},

			"registrant_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City of the Registrant user",
				Default:     "",
			},

			"registrant_state_province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State/Province of the Registrant user",
				Default:     "",
			},

			"registrant_postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "PostalCode of the Registrant user",
				Default:     "",
			},

			"registrant_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country of the Registrant user",
				Default:     "",
			},

			"registrant_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number in the format +NNN.NNNNNNNNNN",
				Default:     "",
			},

			"registrant_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the Registrant user",
				Default:     "",
			},

			"tech_firstname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First name of the Tech user",
				Default:     "",
			},

			"tech_lastname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second name of the Tech user",
				Default:     "",
			},

			"tech_address1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address1 of the Tech user",
				Default:     "",
			},

			"tech_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City of the Tech user",
				Default:     "",
			},

			"tech_state_province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State/Province of the Tech user",
				Default:     "",
			},

			"tech_postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "PostalCode of the Tech user",
				Default:     "",
			},

			"tech_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country of the Tech user",
				Default:     "",
			},

			"tech_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number in the format +NNN.NNNNNNNNNN",
				Default:     "",
			},

			"tech_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the Tech user",
				Default:     "",
			},

			"admin_firstname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First name of the Admin user",
				Default:     "",
			},

			"admin_lastname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second name of the Admin user",
				Default:     "",
			},

			"admin_address1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address1 of the Admin user",
				Default:     "",
			},

			"admin_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City of the Admin user",
				Default:     "",
			},

			"admin_state_province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State/Province of the Admin user",
				Default:     "",
			},

			"admin_postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "PostalCode of the Admin user",
				Default:     "",
			},

			"admin_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country of the Admin user",
				Default:     "",
			},

			"admin_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number in the format +NNN.NNNNNNNNNN",
				Default:     "",
			},

			"admin_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the Admin user",
				Default:     "",
			},

			"aux_billing_firstname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "First name of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_lastname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Second name of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_address1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address1 of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_state_province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State/Province of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "PostalCode of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country of the AuxBilling user",
				Default:     "",
			},

			"aux_billing_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone number in the format +NNN.NNNNNNNNNN",
				Default:     "",
			},

			"aux_billing_email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address of the AuxBilling user",
				Default:     "",
			},

			//---------------------------------------------------------------------------------
		},
		ResourcesMap: map[string]*schema.Resource{
			"namecheap_domain": resourceNamecheapDomain(),
		},
		ConfigureContextFunc: configureContext,
	}
}

func configureContext(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {

	userName := data.Get("user_name").(string)
	apiUser := data.Get("api_user").(string)
	apiKey := data.Get("api_key").(string)
	clientIp := data.Get("client_ip").(string)
	useSandbox := data.Get("use_sandbox").(bool)

	_info.Years = data.Get("years").(string)
	_info.RegistrantFirstName = data.Get("registrant_firstname").(string)
	_info.RegistrantLastName = data.Get("registrant_lastname").(string)
	_info.RegistrantAddress1 = data.Get("registrant_address1").(string)
	_info.RegistrantCity = data.Get("registrant_city").(string)
	_info.RegistrantStateProvince = data.Get("registrant_state_province").(string)
	_info.RegistrantPostalCode = data.Get("registrant_postal_code").(string)
	_info.RegistrantCountry = data.Get("registrant_country").(string)
	_info.RegistrantPhone = data.Get("registrant_phone").(string)
	_info.RegistrantEmailAddress = data.Get("registrant_email_address").(string)

	_info.TechFirstName = data.Get("tech_firstname").(string)
	_info.TechLastName = data.Get("tech_lastname").(string)
	_info.TechAddress1 = data.Get("tech_address1").(string)
	_info.TechCity = data.Get("tech_city").(string)
	_info.TechStateProvince = data.Get("tech_state_province").(string)
	_info.TechPostalCode = data.Get("tech_postal_code").(string)
	_info.TechCountry = data.Get("tech_country").(string)
	_info.TechPhone = data.Get("tech_phone").(string)
	_info.TechEmailAddress = data.Get("tech_email_address").(string)

	_info.AdminFirstName = data.Get("admin_firstname").(string)
	_info.AdminLastName = data.Get("admin_lastname").(string)
	_info.AdminAddress1 = data.Get("admin_address1").(string)
	_info.AdminCity = data.Get("admin_city").(string)
	_info.AdminStateProvince = data.Get("admin_state_province").(string)
	_info.AdminPostalCode = data.Get("admin_postal_code").(string)
	_info.AdminCountry = data.Get("admin_country").(string)
	_info.AdminPhone = data.Get("admin_phone").(string)
	_info.AdminEmailAddress = data.Get("admin_email_address").(string)

	_info.AuxBillingFirstName = data.Get("aux_billing_firstname").(string)
	_info.AuxBillingLastName = data.Get("aux_billing_lastname").(string)
	_info.AuxBillingAddress1 = data.Get("aux_billing_address1").(string)
	_info.AuxBillingCity = data.Get("aux_billing_city").(string)
	_info.AuxBillingStateProvince = data.Get("aux_billing_state_province").(string)
	_info.AuxBillingPostalCode = data.Get("aux_billing_postal_code").(string)
	_info.AuxBillingCountry = data.Get("aux_billing_country").(string)
	_info.AuxBillingPhone = data.Get("aux_billing_phone").(string)
	_info.AuxBillingEmailAddress = data.Get("aux_billing_email_address").(string)

	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   userName,
		ApiUser:    apiUser,
		ApiKey:     apiKey,
		ClientIp:   clientIp,
		UseSandbox: useSandbox,
	})

	return client, diag.Diagnostics{}
}
