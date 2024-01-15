package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type domainsCreateResult struct {
	Domain        string `xml:"Domain,attr"`
	Registered    bool   `xml:"Registered,attr"`
	ChargedAmount string `xml:"ChargedAmount,attr"`
}

type domainsCreateCommandResponse struct {
	Result *domainsCreateResult `xml:"DomainCreateResult"`
}

type domainsCreateResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *domainsCreateCommandResponse `xml:"CommandResponse"`
}

func DomainsCreate(client *namecheap.Client, domainName string, years string, nameservers string, info *UserAddrGetInfoCommandResponse) (*domainsCreateCommandResponse, error) {
	var response domainsCreateResponse

	params := map[string]string{
		"Command":    "namecheap.domains.create",
		"DomainName": domainName,

		"Years":                   years,
		"RegistrantFirstName":     info.Result.FirstName,
		"RegistrantLastName":      info.Result.LastName,
		"RegistrantAddress1":      info.Result.Address1,
		"RegistrantAddress2":      info.Result.Address2,
		"RegistrantCity":          info.Result.City,
		"RegistrantStateProvince": info.Result.StateProvince,
		"RegistrantPostalCode":    info.Result.PostalCode,
		"RegistrantCountry":       info.Result.Country,
		"RegistrantPhone":         info.Result.Phone,
		"RegistrantEmailAddress":  info.Result.EmailAddress,

		"TechFirstName":     info.Result.FirstName,
		"TechLastName":      info.Result.LastName,
		"TechAddress1":      info.Result.Address1,
		"TechAddress2":      info.Result.Address2,
		"TechCity":          info.Result.City,
		"TechStateProvince": info.Result.StateProvince,
		"TechPostalCode":    info.Result.PostalCode,
		"TechCountry":       info.Result.Country,
		"TechPhone":         info.Result.Phone,
		"TechEmailAddress":  info.Result.EmailAddress,

		"AdminFirstName":     info.Result.FirstName,
		"AdminLastName":      info.Result.LastName,
		"AdminAddress1":      info.Result.Address1,
		"AdminAddress2":      info.Result.Address2,
		"AdminCity":          info.Result.City,
		"AdminStateProvince": info.Result.StateProvince,
		"AdminPostalCode":    info.Result.PostalCode,
		"AdminCountry":       info.Result.Country,
		"AdminPhone":         info.Result.Phone,
		"AdminEmailAddress":  info.Result.EmailAddress,

		"AuxBillingFirstName":     info.Result.FirstName,
		"AuxBillingLastName":      info.Result.LastName,
		"AuxBillingAddress1":      info.Result.Address1,
		"AuxBillingAddress2":      info.Result.Address2,
		"AuxBillingCity":          info.Result.City,
		"AuxBillingStateProvince": info.Result.StateProvince,
		"AuxBillingPostalCode":    info.Result.PostalCode,
		"AuxBillingCountry":       info.Result.Country,
		"AuxBillingPhone":         info.Result.Phone,
		"AuxBillingEmailAddress":  info.Result.EmailAddress,

		"Extended attributes": "",
		"Nameservers":         nameservers,
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
