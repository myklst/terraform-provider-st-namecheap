
terraform {
  required_providers {
    namecheap = {
      source = "tao/namecheap"
      version = "= 2.2.0"
    }
  }
}

provider "namecheap" {
  user_name = "haker0032"
  api_user = "haker0032"
  api_key = "XXX"
  client_ip = "180.255.72.50"
  use_sandbox = false

  years = "1"

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


resource "namecheap_domain_records" "domain-com" {
  domain = "sige-test11.com"

    record {
        hostname = "dev"
        type = "A"
        address = "10.12.14.19"
      }

    record {
        hostname = "dev"
        type = "CNAME"
        address = "hoho.com"
      }

}