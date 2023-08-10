
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
  user_name = "haker0032"
  api_user = "haker0032"
  api_key = "XXX"
  client_ip = "180.255.72.50"
  use_sandbox = false
}

resource "namecheap_domain_records" "records" {
  provider = namecheap
  domain = "sige-test11.com"
}


provider "st-namecheap" {
  user_name = "haker0032"
  api_user = "haker0032"
  api_key = "XXX"
  client_ip = "180.255.72.50"
  use_sandbox = false

}

resource "namecheap_domain" "domain-com" {
  provider = st-namecheap
  domain = "sige-test11.com"
  mode = "create"
  years = 1
}