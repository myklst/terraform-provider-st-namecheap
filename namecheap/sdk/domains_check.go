package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type DomainsCheckResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsCheckCommandResponse `xml:"CommandResponse"`
}

type DomainsCheckCommandResponse struct {
	Result *DomainsCheckResult `xml:"DomainCheckResult"`
}

type DomainsCheckResult struct {
	Domain    *string `xml:"Domain,attr"`
	Available *bool   `xml:"Available,attr"`
}

func DomainsAvailable(client *namecheap.Client, domains string) (*DomainsCheckCommandResponse, error) {

	var response DomainsCheckResponse

	params := map[string]string{
		"Command":    "namecheap.domains.check",
		"DomainList": domains,
	}

	_, err := client.DoXML(params, &response)

	if err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil

}
