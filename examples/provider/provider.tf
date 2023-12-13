terraform {
  required_providers {
    st-namecheap = {
      source  = "myklst/st-namecheap"
      version = "~> 0.1"
    }
  }
}

provider "st-namecheap" {
  user_name   = "xxx"
  api_user    = "xxx"
  api_key     = "xxx"
  client_ip   = "xxx.xxx.xxx.xxx"
  use_sandbox = false
}
