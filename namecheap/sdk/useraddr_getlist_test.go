package sdk

import (
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"testing"
)

func TestUseraddrGetList(t *testing.T) {

	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   "haker0032",
		ApiUser:    "haker0032",
		ApiKey:     "34a51f53d7294358b213f0f65a7f9da7",
		ClientIp:   "180.255.72.50",
		UseSandbox: false,
	})

	r, err := UseraddrGetList(client)
	if err != nil {
		t.Error(err)
	}

	v := len(*r.Result.List)
	t.Log(v)
	addrId := *(*r.Result.List)[0].AddressId
	t.Log(addrId)

}
