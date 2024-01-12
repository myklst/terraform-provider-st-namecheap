package sdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/myklst/terraform-provider-st-namecheap/namecheap/sdk"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

func TestUserGetPricing(t *testing.T) {
	client := namecheap.NewClient(&namecheap.ClientOptions{
		UserName:   os.Getenv("NAMECHEAP_USER_NAME"),
		ApiUser:    os.Getenv("NAMECHEAP_API_USER"),
		ApiKey:     os.Getenv("NAMECHEAP_API_KEY"),
		ClientIp:   os.Getenv("NAMECHEAP_CLIENT_IP"),
		UseSandbox: os.Getenv("NAMECHEAP_USE_SANDBOX") == "true",
	})

	resp, err := sdk.UserGetPricing(client, "register", "COM")
	if err != nil {
		t.Error(err)
	}else {
		fmt.Printf("%s", resp.Result.ProductCategory.Price[0].Price)
	}
}
