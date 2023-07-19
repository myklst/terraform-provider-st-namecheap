terraform {
  required_providers {
    namecheap = {
      source = "tao/namecheap"
      version = "1.0.0"
    }
  }
}


provider "namecheap" {
  baseurl = "https://api.namecheap.com/xml.response?"
  key = "34a51f53d7294358b213f0f65a7f9da7"
  user = "haker0032"
  clientIP = "180.255.72.50"
}

resource "namecheap_domain_record" "gd-fancy-domain" {
  domain   = "hohojiang.com"

}

