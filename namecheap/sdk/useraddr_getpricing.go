package sdk

import (
	"encoding/xml"
	"fmt"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

type register struct {
	Price []*struct {
		Duration string `xml:"Duration,attr"`
		Price    string `xml:"YourPrice,attr"`
	} `xml:"Price"`
}

type userGetPricingResult struct {
	ProductCategory register `xml:"ProductType>ProductCategory>Product"`
}

type userGetPricingCommandResponse struct {
	Result *userGetPricingResult `xml:"UserGetPricingResult"`
}

type userGetPricingResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *userGetPricingCommandResponse `xml:"CommandResponse"`
}

func UserGetPricing(client *namecheap.Client, action string, product string) (*userGetPricingCommandResponse, error) {
	var response userGetPricingResponse

	params := map[string]string{
		"Command": "namecheap.users.getPricing",
		"ProductType": "DOMAIN",
		"ActionName" : action,
		"ProductName": product,
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
