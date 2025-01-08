terraform-provider-st-namecheap
===============================

A Terraform Provider for NameCheap domain management.

## Prerequisites

NameCheap provides a sanbox environment for test purposes, please check [HERE](https://www.namecheap.com/support/knowledgebase/article.aspx/763/63/what-is-sandbox/)
for details.

You'll need to apply for API access to NameCheap. You can do that on this
[API Access page](https://ap.www.sandbox.namecheap.com/settings/tools/apiaccess/).

Next, find out your IP address and add that IP (or any other IPs accessing this
API) to this [API Whitelisted IP page](https://ap.www.sandbox.namecheap.com/settings/tools/apiaccess/whitelisted-ips)
on NameCheap.

Once you've done that, make note of the API key, your IP address, and your
username to fill into `provider` block.

Supported Versions
------------------

| Terraform version | minimum provider version |maxmimum provider version
| ---- |--------------------------| ----|
| >= 1.3.x	| 0.1.0	                   | latest |

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
    ```

Why Custom Provider
-------------------

Namecheap does not support managing resources with Terraform.
