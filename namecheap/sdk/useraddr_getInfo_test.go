package sdk

import (
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"testing"
)

func TestUseraddrGetInfo(t *testing.T) {
	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   "haker0032",
		ApiUser:    "haker0032",
		ApiKey:     "34a51f53d7294358b213f0f65a7f9da7",
		ClientIp:   "180.255.72.50",
		UseSandbox: false,
	})

	r, err := UseraddrGetInfo(client, "0")
	if err != nil {
		t.Error(err)
	}

	t.Log(*r.Result.FirstName)
	t.Log(*r.Result.LastName)
	t.Log(*r.Result.PostalCode)

}
