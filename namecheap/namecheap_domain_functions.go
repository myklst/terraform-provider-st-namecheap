package namecheap_provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"strings"
	"terraform-provider-st-namecheap/namecheap/sdk"
)

func fixAddressEndWithDot(address *string) *string {
	if !strings.HasSuffix(*address, ".") {
		return namecheap.String(*address + ".")
	}
	return address
}

func renewDomain(ctx context.Context, domain string, years string, client *namecheap.Client) diag.Diagnostics {

	resp, err := sdk.DomainsRenew(client, domain, years)

	if err != nil || *resp.Result.Renew == false {
		log(ctx, "renew domain %s failed, exit", domain)
		log(ctx, "reason:", err.Error())
		return diag.Errorf("renew domain failed", domain)
	}

	return nil
}

func reactivateDomain(ctx context.Context, domain string, years string, client *namecheap.Client) diag.Diagnostics {

	resp, err := sdk.DomainsReactivate(client, domain, years)

	if err != nil || *resp.Result.IsSuccess == false {
		log(ctx, "reactivate domain %s failed, exit", domain)
		log(ctx, "reason:", err.Error())
		return diag.Errorf("reactivate domain failed", domain)
	}

	return nil
}

func createDomainIfNonexist(ctx context.Context, domain string, years string, client *namecheap.Client) diag.Diagnostics {
	//get domain info
	_, err := client.Domains.GetInfo(domain)

	//if domain does not exist, then create
	if err != nil {
		log(ctx, "Can not Get Domain Info:%s", domain)

		//log.Println("Can not Get Domain Info, Creating:%s", domain)
		resp, err := sdk.DomainsAvailable(client, domain)
		if err == nil && *resp.Result.Available == true {
			// no err and available, create
			log(ctx, "Can not Get Domain Info, Creating %s", domain)
			_, err = sdk.DomainsCreate(client, domain, years, _info)

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
