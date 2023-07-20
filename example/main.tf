
terraform {
  required_providers {
    namecheap = {
      source = "namecheap/namecheap"
      version = ">= 2.0.0"
    }
  }
}

provider "namecheap" {
  user_name = "haker0032"
  api_user = "haker0032"
  api_key = "34a51f53d7294358b213f0f65a7f9da7"
  client_ip = "180.255.72.50"
  use_sandbox = false
}

resource "namecheap_domain_records" "domain-com" {
  domain = "domain.com"
  mode = "OVERWRITE"

  record {
    hostname = "dev"
    type = "A"
    address = "10.12.14.19"
  }
}