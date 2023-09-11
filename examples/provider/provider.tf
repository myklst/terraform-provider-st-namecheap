terraform {
  required_providers {
    st-namecheap = {
      source = "myklst/st-namecheap"
      version = "= 0.1.0"
    }
  }
}


provider "st-namecheap" {
  user_name   = "XXX"
  api_user    = "XXX"
  api_key     = "XXXX"
  client_ip   = "180.255.72.50"
  use_sandbox = false

}
