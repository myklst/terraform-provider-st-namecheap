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
	Years string

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

func DomainsCreate(client *namecheap.Client, domainName string, info DomainCreateInfo) (*DomainsCreateCommandResponse, error) {

	var response DomainsCreateResponse

	params := map[string]string{
		"Command":    "namecheap.domains.create",
		"DomainName": domainName,

		"Years":                   info.Years,
		"RegistrantFirstName":     info.RegistrantFirstName,
		"RegistrantLastName":      info.RegistrantLastName,
		"RegistrantAddress1":      info.RegistrantAddress1,
		"RegistrantCity":          info.RegistrantCity,
		"RegistrantStateProvince": info.RegistrantStateProvince,
		"RegistrantPostalCode":    info.RegistrantPostalCode,
		"RegistrantCountry":       info.RegistrantCountry,
		"RegistrantPhone":         info.RegistrantPhone,
		"RegistrantEmailAddress":  info.RegistrantEmailAddress,

		"TechFirstName":     info.TechFirstName,
		"TechLastName":      info.TechLastName,
		"TechAddress1":      info.TechAddress1,
		"TechCity":          info.TechCity,
		"TechStateProvince": info.TechStateProvince,
		"TechPostalCode":    info.TechPostalCode,
		"TechCountry":       info.TechCountry,
		"TechPhone":         info.TechPhone,
		"TechEmailAddress":  info.TechEmailAddress,

		"AdminFirstName":     info.AdminFirstName,
		"AdminLastName":      info.AdminLastName,
		"AdminAddress1":      info.AdminAddress1,
		"AdminCity":          info.AdminCity,
		"AdminStateProvince": info.AdminStateProvince,
		"AdminPostalCode":    info.AdminPostalCode,
		"AdminCountry":       info.AdminCountry,
		"AdminPhone":         info.AdminPhone,
		"AdminEmailAddress":  info.AdminEmailAddress,

		"AuxBillingFirstName":     info.AuxBillingFirstName,
		"AuxBillingLastName":      info.AuxBillingLastName,
		"AuxBillingAddress1":      info.AuxBillingAddress1,
		"AuxBillingCity":          info.AuxBillingCity,
		"AuxBillingStateProvince": info.AuxBillingStateProvince,
		"AuxBillingPostalCode":    info.AuxBillingPostalCode,
		"AuxBillingCountry":       info.AuxBillingCountry,
		"AuxBillingPhone":         info.AuxBillingPhone,
		"AuxBillingEmailAddress":  info.AuxBillingEmailAddress,

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
