package namecheap_provider

import (
	"context"
	"errors"
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
		return diag.Errorf("renew domain [%s] failed", domain)
	}

	log(ctx, "renew domain [%s] success", domain)
	return nil
}

func reactivateDomain(ctx context.Context, domain string, years string, client *namecheap.Client) diag.Diagnostics {

	resp, err := sdk.DomainsReactivate(client, domain, years)

	if err != nil || *resp.Result.IsSuccess == false {
		log(ctx, "reactivate domain %s failed, exit", domain)
		log(ctx, "reason:", err.Error())
		return diag.Errorf("reactivate domain [%s] failed", domain)
	}

	log(ctx, "reactivate domain [%s] success", domain)
	return nil
}

func createDomain(ctx context.Context, domain string, years string, client *namecheap.Client) diag.Diagnostics {
	//get domain info
	_, err := client.Domains.GetInfo(domain)
	if err == nil {
		return diag.Errorf("domain [%s] has been created in this account", domain)
	}

	//if domain does not exist, then create

	//log.Println("Can not Get Domain Info, Creating:%s", domain)
	resp, err := sdk.DomainsAvailable(client, domain)
	if err == nil && *resp.Result.Available == true {
		// no err and available, create
		log(ctx, "Domain [%s] is available, Creating...", domain)
		r, err := getUserAccountContact(client)
		if err != nil {
			log(ctx, "get user Contacts failed, exit")
			log(ctx, "reason:", err.Error())
			return diag.Errorf("get user contacts failed, please check the contacts in the account manually", domain)
		}

		log(ctx, "debug domain create", r)

		_, err = sdk.DomainsCreate(client, domain, years, r)

		if err != nil {
			log(ctx, "create domain [%s] failed, exit", domain)
			log(ctx, "reason:", err.Error())
			return diag.Errorf("create domain [%s] failed", domain)
		}

	} else {
		log(ctx, "domain [%s] is not available, exiting!", domain)
		return diag.Errorf("domain [%s] is not available to register, you need to change to another domain", domain)
	}

	return nil

}

func getUserAccountContact(client *namecheap.Client) (*sdk.UseraddrGetInfoCommandResponse, error) {

	r1, err := sdk.UseraddrGetList(client)
	if err != nil {
		return nil, err
	}
	addrlen := len(*r1.Result.List)
	if addrlen == 0 {
		return nil, errors.New("UseraddrGetList returns 0, please add user contact info to this account")
	}

	addrId := *(*r1.Result.List)[0].AddressId

	r2, err := sdk.UseraddrGetInfo(client, addrId)
	if err != nil {
		return nil, err
	}

	return r2, nil

}

func log(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a)
	tflog.Info(ctx, msg)

}
