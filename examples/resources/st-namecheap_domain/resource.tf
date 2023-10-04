resource "st-namecheap_domain" "domain-com" {
  provider = st-namecheap
  domain = "sige-test11.com"
  purchase_years = 1
  min_days_remaining = 90
}
