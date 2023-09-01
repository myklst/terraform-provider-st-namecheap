terraform {
  required_providers {
    st-namecheap = {
      source  = "myklst/st-namecheap"
      version = "= 2.2.0"
    }
  }
}

provider "st-namecheap" {
  user_name   = "haker0032"
  api_user    = "haker0032"
  api_key     = "34a51f53d7294358b213f0f65a7f9da7"
  client_ip   = "180.255.72.50"
  use_sandbox = false
}
