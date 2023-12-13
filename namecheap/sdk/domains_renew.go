package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type domainsRenewResult struct {
	DomainName *string `xml:"DomainName,attr"`
	Renew      *bool   `xml:"Renew,attr"`
}

type domainsRenewCommandResponse struct {
	Result *domainsRenewResult `xml:"DomainRenewResult"`
}

type domainsRenewResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *domainsRenewCommandResponse `xml:"CommandResponse"`
}

func DomainsRenew(client *namecheap.Client, domains string, years string) (*domainsRenewCommandResponse, error) {
	var response domainsRenewResponse

	params := map[string]string{
		"Command":    "namecheap.domains.renew",
		"DomainName": domains,
		"Years":      years,
	}
	if _, err := doXmlWithBackoff(client, params, &response); err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil

}
