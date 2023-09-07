package sdk

import (
	"net/http"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"

	"github.com/cenkalti/backoff/v4"
)

func DoXmlWithBackoff(client *namecheap.Client, body map[string]string, obj interface{}) (*http.Response, error) {
	var requestResponse *http.Response

	operation := func() error {
		var err error
		requestResponse, err = client.DoXML(body, obj)
		return err
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {

		return nil, err
	}

	return requestResponse, nil
}
