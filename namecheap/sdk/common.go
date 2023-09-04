package sdk

import (
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"net/http"
)

const RetryTimes int = 3

func DoXmlWithRetry(client *namecheap.Client, body map[string]string, obj interface{}) (*http.Response, error) {
	var requestResponse *http.Response

	var err error
	for t := 0; t < RetryTimes; t++ {
		requestResponse, err = client.DoXML(body, obj)
		if err == nil {
			return requestResponse, nil
		}
	}

	return nil, err
}
