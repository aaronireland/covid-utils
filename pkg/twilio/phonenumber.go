package twilio

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

// IncomingPhoneNumbers is the XML root response body for the IncomingPhoneNumbers URI
type IncomingPhoneNumbers struct {
	XMLName      xml.Name `xml:"TwilioResponse"`
	PhoneNumbers struct {
		XMLName      xml.Name              `xml:"IncomingPhoneNumbers"`
		PhoneNumbers []IncomingPhoneNumber `xml:"IncomingPhoneNumber"`
	} `xml:"IncomingPhoneNumbers"`
}

// IncomingPhoneNumber is the Twilio phone data used to send SMS/MMS/Voice calls
type IncomingPhoneNumber struct {
	XMLName      xml.Name `xml:"IncomingPhoneNumber"`
	PhoneNumber  string   `xml:"PhoneNumber"`
	Capabilities struct {
		XMLName xml.Name `xml:"Capabilities"`
		Voice   bool     `xml:"Voice"`
		SMS     bool     `xml:"SMS"`
		MMS     bool     `xml:"MMS"`
	} `xml:"Capabilities"`
}

// SetPhoneNumber accesses the IncomingPhoneNumbers URI to find an SMS capable phone number
func (t *APIClient) SetPhoneNumber() error {
	req, _ := http.NewRequest("GET", t.URI(t.routes.IncomingPhoneNumbers), nil)
	req.SetBasicAuth(t.accountSid, t.authToken)
	var phoneNums IncomingPhoneNumbers
	err := t.Do(req, &phoneNums)
	if err != nil {
		return fmt.Errorf("unable to check for phone numbers attached to account: %s", err)
	}
	if len(phoneNums.PhoneNumbers.PhoneNumbers) == 0 {
		return errors.New("no phone numbers associated to account")
	}
	for _, phone := range phoneNums.PhoneNumbers.PhoneNumbers {
		if phone.Capabilities.SMS {
			t.fromPhoneNumber = phone.PhoneNumber
			return nil
		}
	}
	return errors.New("no SMS capable phone number associated to account")
}
