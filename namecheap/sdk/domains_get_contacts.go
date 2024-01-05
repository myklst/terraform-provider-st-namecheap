package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

var cache *domainsGetContactsCommandResponse

type contact struct {
	OrganizationName    string `xml:"OrganizationName"`
	JobTitle            string `xml:"JobTitle"`
	FirstName           string `xml:"FirstName"`
	LastName            string `xml:"LastName"`
	Address1            string `xml:"Address1"`
	Address2            string `xml:"Address2"`
	City                string `xml:"City"`
	StateProvince       string `xml:"StateProvince"`
	StateProvinceChoice string `xml:"StateProvinceChoice"`
	PostalCode          string `xml:"PostalCode"`
	Country             string `xml:"Country"`
	Phone               string `xml:"Phone"`
	Fax                 string `xml:"Fax"`
	EmailAddress        string `xml:"EmailAddress"`
	PhoneExt            string `xml:"PhoneExt"`
}

type domainsContactsResult struct {
	Domain     string  `xml:"Domain,attr"`
	Registrant *contact `xml:"Registrant"`
	Tech       *contact `xml:"Tech"`
	Admin      *contact `xml:"Admin"`
	AuxBilling *contact `xml:"AuxBilling"`
}

type domainsGetContactsCommandResponse struct {
	Result *domainsContactsResult `xml:"DomainContactsResult"`
}

type domainsGetContactsResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *domainsGetContactsCommandResponse `xml:"CommandResponse"`
}

func DomainsGetContacts(client *namecheap.Client) (*domainsGetContactsCommandResponse, error) {
	if cache == nil {
		r, err := client.Domains.GetList(&namecheap.DomainsGetListArgs{})
		if err != nil {
			return nil, err
		}
		if r == nil || r.Domains == nil || len(*r.Domains) <= 0 {
			return nil, errors.New("no purchased domains")
		}
		domain := (*r.Domains)[0]

		var response domainsGetContactsResponse

		params := map[string]string{
			"Command":    "namecheap.domains.getContacts",
			"DomainName": *domain.Name,
		}
		if _, err = doXmlWithBackoff(client, params, &response); err != nil {
			return nil, err
		}

		if response.Errors != nil && len(*response.Errors) > 0 {
			apiErr := (*response.Errors)[0]
			return nil, fmt.Errorf("%s (%s)", apiErr.Message, apiErr.Number)
		}

		cache = response.CommandResponse
	}

	return cache, nil
}
