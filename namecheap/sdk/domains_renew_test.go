package sdk

import (
	"testing"
)

func TestDomainsRenew(t *testing.T) {
	// Because NameCheap sandbox payment isn't working, cannot process domain
	// renewal.

	// client := namecheap.NewClient(&namecheap.ClientOptions{
	// 	UserName:   os.Getenv("NAMECHEAP_USER_NAME"),
	// 	ApiUser:    os.Getenv("NAMECHEAP_API_USER"),
	// 	ApiKey:     os.Getenv("NAMECHEAP_API_KEY"),
	// 	ClientIp:   os.Getenv("NAMECHEAP_CLIENT_IP"),
	// 	UseSandbox: os.Getenv("NAMECHEAP_USE_SANDBOX") == "true",
	// })

	// if _, err := DomainsRenew(client, "example.com", "1"); err != nil {
	// 	t.Error(err)
	// }
}
