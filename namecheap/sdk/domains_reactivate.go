package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type domainsReactivateResult struct {
	Domain    string `xml:"Domain,attr"`
	IsSuccess bool   `xml:"IsSuccess,attr"`
}

type domainsReactivateCommandResponse struct {
	Result *domainsReactivateResult `xml:"DomainReactivateResult"`
}

type domainsReactivateResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *domainsReactivateCommandResponse `xml:"CommandResponse"`
}

func DomainsReactivate(client *namecheap.Client, domains string, years string) (*domainsReactivateCommandResponse, error) {
	var response domainsReactivateResponse

	params := map[string]string{
		"Command":    "namecheap.domains.reactivate",
		"DomainName": domains,
		"YearsToAdd": years,
	}
	if _, err := doXmlWithBackoff(client, params, &response); err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", apiErr.Message, apiErr.Number)
	}

	return response.CommandResponse, nil

}
