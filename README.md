terraform-provider-st-namecheap
===============================

A Terraform Provider for NameCheap domain management.

## Prerequisites

First you'll need to apply for API access to NameCheap. You can do that on
this [API admin page](https://ap.www.namecheap.com/settings/tools/apiaccess/).

Next, find out your IP address and add that IP (or any other IPs accessing this
API) to this [whitelist admin page](https://ap.www.namecheap.com/settings/tools/apiaccess/whitelisted-ips) on NameCheap.

Once you've done that, make note of the API key, your IP address, and your
username to fill into our `provider` block.


Supported Versions
------------------

| Terraform version | minimum provider version |maxmimum provider version
| ---- |--------------------------| ----|
| >= 1.3.x	| 2.2.0	                   | latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.3.x
-	[Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Local Installation
------------------

1. Run make file `make install-local-custom-provider` to install the provider under ~/.terraform.d/plugins.

2. The provider source should be change to the path that configured in the *Makefile*:

    ```
    terraform {
      required_providers {
        st-alicloud = {
          source = "example.local/myklst/st-namecheap"
        }
      }
    }

    provider "st-namecheap" {
        user_name   = "XXX"
        api_user    = "XXX"
        api_key     = "XXXX"
        client_ip   = "X.X.X.X"
        use_sandbox = false
       
    }
    ```


export TF_LOG=DEBUG
