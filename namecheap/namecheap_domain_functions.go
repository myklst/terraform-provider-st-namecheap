package namecheap_provider

import (
	"context"
	"fmt"
	"github.com/agent-tao/go-namecheap-sdk/v2/namecheap"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"
)

func fixAddressEndWithDot(address *string) *string {
	if !strings.HasSuffix(*address, ".") {
		return namecheap.String(*address + ".")
	}
	return address
}

func createDomainIfNonexist(ctx context.Context, domain string, client *namecheap.Client) diag.Diagnostics {
	//get domain info
	_, err := client.Domains.GetInfo(domain)

	//if domain does not exist, then create
	if err != nil {
		log(ctx, "Can not Get Domain Info:%s", domain)

		//log.Println("Can not Get Domain Info, Creating:%s", domain)
		resp, err := client.Domains.DomainsAvailable(domain)
		if err == nil && *resp.Result.Available == true {
			// no err and available, create
			log(ctx, "Can not Get Domain Info, Creating %s", domain)
			_, err = client.Domains.DomainsCreate(domain, _info)

			if err != nil {
				log(ctx, "create domain %s failed, exit", domain)
				log(ctx, "reason:", err.Error())
				return diag.Errorf("create domain failed", domain)
			}

		} else {
			log(ctx, "domain %s is not available, exiting!", domain)
			return diag.Errorf("domain is not available to register, you need to change to another domain", domain)
		}
	} else {
		//skip, do nothing
		tflog.Info(ctx, "Domain %s exist, then do record config", domain)
	}
	return nil

}

func log(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a)
	tflog.Info(ctx, msg)

}
