package twilio

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// ActiveStatus indicates the Twilio account is active
const ActiveStatus = "active"

// AccountsResponse is the XML root
type AccountsResponse struct {
	XMLName  xml.Name `xml:"TwilioResponse"`
	Accounts Accounts `xml:"Accounts"`
}

// Accounts is the XML parent for any/all Twilio accounts under the Sid
type Accounts struct {
	XMLName  xml.Name  `xml:"Accounts"`
	Accounts []Account `xml:"Account"`
}

// Account contains the Twilio account status and URI routes
type Account struct {
	XMLName xml.Name `xml:"Account"`
	Status  string   `xml:"Status"`
	Routes  Routes   `xml:"SubresourceUris"`
}

// Routes are the URI routes for the client Twilio account
type Routes struct {
	XMLName              xml.Name `xml:"SubresourceUris"`
	Messages             string   `xml:"Messages"`
	IncomingPhoneNumbers string   `xml:"IncomingPhoneNumbers"`
}

// Authenticate attempts to call the /Accounts URI and check the account status
func (t *APIClient) Authenticate(accountSid, authToken string) error {
	req, _ := http.NewRequest("GET", t.URI("2010-04-01/Accounts"), nil)
	req.SetBasicAuth(accountSid, authToken)
	var accounts AccountsResponse
	err := t.Do(req, &accounts)
	if err != nil {
		return fmt.Errorf("Unable to authenticate Twilio Account Sid %s: %s", accountSid, err)
	}

	if len(accounts.Accounts.Accounts) > 0 {
		account := accounts.Accounts.Accounts[0]

		// Account is active and authenticated
		if account.Status == ActiveStatus {
			t.routes = account.Routes
			t.accountSid = accountSid
			t.authToken = authToken

			return nil
		}
	}

	return fmt.Errorf("Twilio Account Sid %s is not active", accountSid)
}
