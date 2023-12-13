package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type userAddrGetListResult struct {
	List *[]struct {
		AddressId   *string `xml:"AddressId,attr"`
		AddressName *string `xml:"AddressName,attr"`
	} `xml:"List"`
}

type userAddrGetListCommandResponse struct {
	Result *userAddrGetListResult `xml:"AddressGetListResult"`
}

type userAddrGetListResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *userAddrGetListCommandResponse `xml:"CommandResponse"`
}

func UserAddrGetList(client *namecheap.Client) (*userAddrGetListCommandResponse, error) {
	var response userAddrGetListResponse

	params := map[string]string{
		"Command": "namecheap.users.address.getList",
	}
	if _, err := doXmlWithBackoff(client, params, &response); err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil
}
