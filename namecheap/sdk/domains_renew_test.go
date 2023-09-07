package sdk

import (
	"testing"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func TestDomainsRenew(t *testing.T) {
	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   "haker0032",
		ApiUser:    "haker0032",
		ApiKey:     "34a51f53d7294358b213f0f65a7f9da7",
		ClientIp:   "180.255.72.50",
		UseSandbox: false,
	})

	r, err := DomainsRenew(client, "hohojiang.com", "1")
	if err != nil {

		t.Error(err)
	}
	t.Log(*r.Result.DomainName)
	t.Log(*r.Result.Renew)

}
