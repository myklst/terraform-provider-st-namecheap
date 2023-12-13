package sdk

import (
	"testing"
)

func TestDomainsReactivate(t *testing.T) {
	// Because NameCheap sandbox payment isn't working, cannot process domain
	// reactivation.

	// client := namecheap.NewClient(&namecheap.ClientOptions{
	// 	UserName:   os.Getenv("NAMECHEAP_USER_NAME"),
	// 	ApiUser:    os.Getenv("NAMECHEAP_API_USER"),
	// 	ApiKey:     os.Getenv("NAMECHEAP_API_KEY"),
	// 	ClientIp:   os.Getenv("NAMECHEAP_CLIENT_IP"),
	// 	UseSandbox: os.Getenv("NAMECHEAP_USE_SANDBOX") == "true",
	// })

	// if _, err := DomainsReactivate(client, "example.com", "2"); err != nil {
	// 	t.Error(err)
	// }
}
