# Namecheap Terraform ST Provider

A Terraform Provider for Namecheap domain configuration.

## Prerequisites

First you'll need to apply for API access to Namecheap. You can do that on
this [API admin page](https://ap.www.namecheap.com/settings/tools/apiaccess/).

Next, find out your IP address and add that IP (or any other IPs accessing this API) to
this [whitelist admin page](https://ap.www.namecheap.com/settings/tools/apiaccess/whitelisted-ips) on Namecheap.

Once you've done that, make note of the API key, your IP address, and your username to fill into our `provider` block.

## Usage Example

Make sure your API details are correct in the provider block.

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    st-namecheap = {
      source = "st/namecheap"
      version = "= 2.2.0"
    }

    namecheap = {
      source = "namecheap/namecheap"
      version = "= 2.1.0"
    }
  }
}

provider "namecheap" {
  user_name = "your_username"
  api_user = "your_username"
  api_key = "your_api_key"
  client_ip = "your.ip.address.here"
  use_sandbox = false
}

provider "st-namecheap" {
  user_name = "your_username"
  api_user = "your_username"
  api_key = "your_api_key"
  client_ip = "your.ip.address.here"
  use_sandbox = false

  registrant_firstname = "John"
  registrant_lastname = "Smith"
  registrant_address1 = "Street Ave. 666"
  registrant_city      ="New City"
  registrant_state_province ="CA"
  registrant_postal_code    ="90045"
  registrant_country       ="US"
  registrant_phone         ="+1.6613102107"
  registrant_email_address ="john@gmail.com"

  tech_firstname = "John"
  tech_lastname = "Smith"
  tech_address1 = "Street Ave. 666"
  tech_city = "New City"
  tech_state_province = "CA"
  tech_postal_code = "90045"
  tech_country = "US"
  tech_phone = "+1.6613102107"
  tech_email_address = "john@gmail.com"

  admin_firstname = "John"
  admin_lastname = "Smith"
  admin_address1 = "Street Ave. 666"
  admin_city = "New City"
  admin_state_province = "CA"
  admin_postal_code = "90045"
  admin_country = "US"
  admin_phone = "+1.6613102107"
  admin_email_address = "john@gmail.com"

  aux_billing_firstname = "John"
  aux_billing_lastname = "Smith"
  aux_billing_address1 = "Street Ave. 666"
  aux_billing_city = "New City"
  aux_billing_state_province = "CA"
  aux_billing_postal_code = "90045"
  aux_billing_country = "US"
  aux_billing_phone = "+1.6613102107"
  aux_billing_email_address = "john@gmail.com"
}

resource "namecheap_domain" "domain-com" {
  provider = st-namecheap
  
  domain = "domain.com"
  mode = "create"
  years = 1
}

resource "namecheap_domain_records" "domain-com" {
  provider = namecheap
  
  domain = "domain.com"
  mode = "OVERWRITE"

  record {
    hostname = "dev"
    type = "A"
    address = "10.12.14.19"
  }
}

```

### Contributing

To contribute, please read our [contributing](CONTRIBUTING.md) docs.  
