package sdk

import (
	"encoding/xml"
	"fmt"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type UseraddrGetListResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *UseraddrGetListCommandResponse `xml:"CommandResponse"`
}

type UseraddrGetListCommandResponse struct {
	Result *UseraddrGetListResult `xml:"AddressGetListResult"`
}

type UseraddrGetListResult struct {
	List *[]struct {
		AddressId   *string `xml:"AddressId,attr"`
		AddressName *string `xml:"AddressName,attr"`
	} `xml:"List"`
}

func UseraddrGetList(client *namecheap.Client) (*UseraddrGetListCommandResponse, error) {

	var response UseraddrGetListResponse

	params := map[string]string{
		"Command": "namecheap.users.address.getList",
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
