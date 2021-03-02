package twilio

import (
	"github.com/aaronireland/covid-utils/pkg/api"
)

const (
	BaseURL = "https://api.twilio.com"
)

type APIClient struct {
	*api.Client
	fromPhoneNumber string
	routes          Routes
	accountSid      string
	authToken       string
}

func NewClient(accountSid, authToken string) (*APIClient, error) {
	client := &APIClient{Client: api.NewClient(BaseURL, api.XMLClient)}
	err := client.Authenticate(accountSid, authToken)
	if err != nil {
		return nil, err
	}
	err = client.SetPhoneNumber()
	if err != nil {
		return nil, err
	}

	return client, nil
}
