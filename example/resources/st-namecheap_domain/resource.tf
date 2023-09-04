resource "st-namecheap_domain" "domain-com" {
  provider = st-namecheap

  domain = "sige-test11.com"
  mode   = "create"
  years  = 1
}
