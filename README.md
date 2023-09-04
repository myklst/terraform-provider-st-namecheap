terraform-provider-st-namecheap
===============================

A Terraform Provider for NameCheap domain management.

## Prerequisites

First you'll need to apply for API access to NameCheap. You can do that on
this [API admin page](https://ap.www.namecheap.com/settings/tools/apiaccess/).

Next, find out your IP address and add that IP (or any other IPs accessing this
API) to this [whitelist admin page](https://ap.www.namecheap.com/settings/tools/apiaccess/whitelisted-ips) on NameCheap.

Once you've done that, make note of the API key, your IP address, and your
username to fill into our `provider` block.

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

export TF_LOG=DEBUG
