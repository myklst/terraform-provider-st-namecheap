package sdk

import (
	"os"
	"testing"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func TestDomainsGetContacts(t *testing.T) {
	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   os.Getenv("NAMECHEAP_USER_NAME"),
		ApiUser:    os.Getenv("NAMECHEAP_API_USER"),
		ApiKey:     os.Getenv("NAMECHEAP_API_KEY"),
		ClientIp:   os.Getenv("NAMECHEAP_CLIENT_IP"),
		UseSandbox: os.Getenv("NAMECHEAP_USE_SANDBOX") == "true",
	})

	if _, err := DomainsGetContacts(client); err != nil {
		// "no purchased domains" error is expected.
		if err.Error() != "no purchased domains" {
			t.Error(err)
		}
	}
}
