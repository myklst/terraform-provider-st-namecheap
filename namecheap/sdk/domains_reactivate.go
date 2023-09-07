package sdk

import (
	"encoding/xml"
	"fmt"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type DomainsReactivateResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsReactivateCommandResponse `xml:"CommandResponse"`
}

type DomainsReactivateCommandResponse struct {
	Result *DomainsReactivateResult `xml:"DomainReactivateResult"`
}

type DomainsReactivateResult struct {
	Domain    *string `xml:"Domain,attr"`
	IsSuccess *bool   `xml:"IsSuccess,attr"`
}

func DomainsReactivate(client *namecheap.Client, domains string, years string) (*DomainsReactivateCommandResponse, error) {

	var response DomainsReactivateResponse

	params := map[string]string{
		"Command":    "namecheap.domains.reactivate",
		"DomainName": domains,
		"YearsToAdd": years,
	}

	_, err := DoXmlWithBackoff(client, params, &response)

	if err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil

}
