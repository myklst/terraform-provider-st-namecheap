package sdk

import (
	"encoding/xml"
	"fmt"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type DomainsCreateResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsCreateCommandResponse `xml:"CommandResponse"`
}

type DomainsCreateCommandResponse struct {
	Result *DomainsCreateResult `xml:"DomainCreateResult"`
}

type DomainsCreateResult struct {
	Domain        *string `xml:"Domain,attr"`
	Registered    *bool   `xml:"Registered,attr"`
	ChargedAmount *string `xml:"ChargedAmount,attr"`
}

type DomainCreateInfo struct {
	RegistrantFirstName     string
	RegistrantLastName      string
	RegistrantAddress1      string
	RegistrantCity          string
	RegistrantStateProvince string
	RegistrantPostalCode    string
	RegistrantCountry       string
	RegistrantPhone         string
	RegistrantEmailAddress  string

	TechFirstName     string
	TechLastName      string
	TechAddress1      string
	TechCity          string
	TechStateProvince string
	TechPostalCode    string
	TechCountry       string
	TechPhone         string
	TechEmailAddress  string

	AdminFirstName     string
	AdminLastName      string
	AdminAddress1      string
	AdminCity          string
	AdminStateProvince string
	AdminPostalCode    string
	AdminCountry       string
	AdminPhone         string
	AdminEmailAddress  string

	AuxBillingFirstName     string
	AuxBillingLastName      string
	AuxBillingAddress1      string
	AuxBillingCity          string
	AuxBillingStateProvince string
	AuxBillingPostalCode    string
	AuxBillingCountry       string
	AuxBillingPhone         string
	AuxBillingEmailAddress  string
}

func DomainsCreate(client *namecheap.Client, domainName string, years string, info *DomainsGetContactsCommandResponse) (*DomainsCreateCommandResponse, error) {

	var response DomainsCreateResponse

	params := map[string]string{
		"Command":    "namecheap.domains.create",
		"DomainName": domainName,

		"Years":                   years,
		"RegistrantFirstName":     *info.Result.Registrant.FirstName,
		"RegistrantLastName":      *info.Result.Registrant.LastName,
		"RegistrantAddress1":      *info.Result.Registrant.Address1,
		"RegistrantAddress2":      *info.Result.Registrant.Address2,
		"RegistrantCity":          *info.Result.Registrant.City,
		"RegistrantStateProvince": *info.Result.Registrant.StateProvince,
		"RegistrantPostalCode":    *info.Result.Registrant.PostalCode,
		"RegistrantCountry":       *info.Result.Registrant.Country,
		"RegistrantPhone":         *info.Result.Registrant.Phone,
		"RegistrantEmailAddress":  *info.Result.Registrant.EmailAddress,

		"TechFirstName":     *info.Result.Tech.FirstName,
		"TechLastName":      *info.Result.Tech.LastName,
		"TechAddress1":      *info.Result.Tech.Address1,
		"TechAddress2":      *info.Result.Tech.Address2,
		"TechCity":          *info.Result.Tech.City,
		"TechStateProvince": *info.Result.Tech.StateProvince,
		"TechPostalCode":    *info.Result.Tech.PostalCode,
		"TechCountry":       *info.Result.Tech.Country,
		"TechPhone":         *info.Result.Tech.Phone,
		"TechEmailAddress":  *info.Result.Tech.EmailAddress,

		"AdminFirstName":     *info.Result.Admin.FirstName,
		"AdminLastName":      *info.Result.Admin.LastName,
		"AdminAddress1":      *info.Result.Admin.Address1,
		"AdminAddress2":      *info.Result.Admin.Address2,
		"AdminCity":          *info.Result.Admin.City,
		"AdminStateProvince": *info.Result.Admin.StateProvince,
		"AdminPostalCode":    *info.Result.Admin.PostalCode,
		"AdminCountry":       *info.Result.Admin.Country,
		"AdminPhone":         *info.Result.Admin.Phone,
		"AdminEmailAddress":  *info.Result.Admin.EmailAddress,

		"AuxBillingFirstName":     *info.Result.AuxBilling.FirstName,
		"AuxBillingLastName":      *info.Result.AuxBilling.LastName,
		"AuxBillingAddress1":      *info.Result.AuxBilling.Address1,
		"AuxBillingAddress2":      *info.Result.AuxBilling.Address2,
		"AuxBillingCity":          *info.Result.AuxBilling.City,
		"AuxBillingStateProvince": *info.Result.AuxBilling.StateProvince,
		"AuxBillingPostalCode":    *info.Result.AuxBilling.PostalCode,
		"AuxBillingCountry":       *info.Result.AuxBilling.Country,
		"AuxBillingPhone":         *info.Result.AuxBilling.Phone,
		"AuxBillingEmailAddress":  *info.Result.AuxBilling.EmailAddress,

		"Extended attributes": "",
		"Nameservers":         "",
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
