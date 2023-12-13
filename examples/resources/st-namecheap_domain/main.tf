resource "st-namecheap_domain" "domain" {
  domain             = "example.com"
  purchase_years     = 1
  min_days_remaining = 90
}
