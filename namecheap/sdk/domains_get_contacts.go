package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type DomainsGetContactsResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsGetContactsCommandResponse `xml:"CommandResponse"`
}

type DomainsGetContactsCommandResponse struct {
	Result *DomainsContactsResult `xml:"DomainContactsResult"`
}

type DomainsContactsResult struct {
	Domain     *string  `xml:"Domain,attr"`
	Registrant *Contact `xml:"Registrant"`
	Tech       *Contact `xml:"Tech"`
	Admin      *Contact `xml:"Admin"`
	AuxBilling *Contact `xml:"AuxBilling"`
}

type Contact struct {
	OrganizationName    *string `xml:"OrganizationName"`
	JobTitle            *string `xml:"JobTitle"`
	FirstName           *string `xml:"FirstName"`
	LastName            *string `xml:"LastName"`
	Address1            *string `xml:"Address1"`
	Address2            *string `xml:"Address2"`
	City                *string `xml:"City"`
	StateProvince       *string `xml:"StateProvince"`
	StateProvinceChoice *string `xml:"StateProvinceChoice"`
	PostalCode          *string `xml:"PostalCode"`
	Country             *string `xml:"Country"`
	Phone               *string `xml:"Phone"`
	Fax                 *string `xml:"Fax"`
	EmailAddress        *string `xml:"EmailAddress"`
	PhoneExt            *string `xml:"PhoneExt"`
}

func DomainsGetContacts(client *namecheap.Client) (*DomainsGetContactsCommandResponse, error) {

	r, err := client.Domains.GetList(&namecheap.DomainsGetListArgs{})
	if err != nil {
		return nil, err
	}

	if len(*r.Domains) < 1 {
		return nil, errors.New("no exist domains in this account,so can't get contacts to create domain ")
	}
	domain := (*r.Domains)[0]

	var response DomainsGetContactsResponse

	params := map[string]string{
		"Command":    "namecheap.domains.getContacts",
		"DomainName": *domain.Name,
	}

	_, err = DoXmlWithBackoff(client, params, &response)

	if err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil

}
