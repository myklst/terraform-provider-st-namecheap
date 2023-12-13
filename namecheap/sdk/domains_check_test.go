package sdk

import (
	"os"
	"testing"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func TestDomainsCheck(t *testing.T) {
	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   os.Getenv("NAMECHEAP_USER_NAME"),
		ApiUser:    os.Getenv("NAMECHEAP_API_USER"),
		ApiKey:     os.Getenv("NAMECHEAP_API_KEY"),
		ClientIp:   os.Getenv("NAMECHEAP_CLIENT_IP"),
		UseSandbox: os.Getenv("NAMECHEAP_USE_SANDBOX") == "true",
	})

	r, err := DomainsAvailable(client, "example.com")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(*r.Result.Available)
	}
}
