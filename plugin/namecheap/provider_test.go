package namecheap

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"namecheap": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	verifyEnvExists(t, "NAMECHEAP_API_KEY")
	verifyEnvExists(t, "NAMECHEAP_API_USER")
	verifyEnvExists(t, "NAMECHEAP_DOMAIN")
}

func verifyEnvExists(t *testing.T, key string) {
	if v := os.Getenv(key); v == "" {
		t.Fatal(fmt.Sprintf("%s must be set for acceptance tests.", key))
	}
}
