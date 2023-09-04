package sdk

import (
	"encoding/xml"
	"fmt"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type DomainsRenewResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsRenewCommandResponse `xml:"CommandResponse"`
}

type DomainsRenewCommandResponse struct {
	Result *DomainsRenewResult `xml:"DomainRenewResult"`
}

type DomainsRenewResult struct {
	DomainName *string `xml:"DomainName,attr"`
	Renew      *bool   `xml:"Renew,attr"`
}

func DomainsRenew(client *namecheap.Client, domains string, years string) (*DomainsRenewCommandResponse, error) {

	var response DomainsRenewResponse

	params := map[string]string{
		"Command":    "namecheap.domains.renew",
		"DomainName": domains,
		"Years":      years,
	}

	_, err := DoXmlWithRetry(client, params, &response)

	if err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil

}
