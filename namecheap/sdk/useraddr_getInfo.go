package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)
type userAddrGetInfoResult struct {
	Organization        string `xml:"Organization"`
	JobTitle            string `xml:"JobTitle"`
	FirstName           string `xml:"FirstName"`
	LastName            string `xml:"LastName"`
	Address1            string `xml:"Address1"`
	Address2            string `xml:"Address2"`
	City                string `xml:"City"`
	StateProvince       string `xml:"StateProvince"`
	StateProvinceChoice string `xml:"StateProvinceChoice"`
	PostalCode          string `xml:"Zip"`
	Country             string `xml:"Country"`
	Phone               string `xml:"Phone"`
	Fax                 string `xml:"Fax"`
	EmailAddress        string `xml:"EmailAddress"`
	PhoneExt            string `xml:"PhoneExt"`
}

type UserAddrGetInfoCommandResponse struct {
	Result *userAddrGetInfoResult `xml:"GetAddressInfoResult"`
}

type userAddrGetInfoResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *UserAddrGetInfoCommandResponse `xml:"CommandResponse"`
}

func UserAddrGetInfo(client *namecheap.Client, addrId string) (*UserAddrGetInfoCommandResponse, error) {
	var response userAddrGetInfoResponse

	params := map[string]string{
		"Command":   "namecheap.users.address.getInfo",
		"AddressId": addrId,
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
