package twilio

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aaronireland/covid-utils/pkg/api"
)

func (t *APIClient) SendSMS(ctx context.Context, to, message string) error {
	params := api.Params{
		api.Param{"To", to},
		api.Param{"From", t.fromPhoneNumber},
		api.Param{"Body", message},
	}

	reqFormData := *strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", t.URI(t.routes.Messages), &reqFormData)
	if err != nil {
		return fmt.Errorf("Malformed POST request: %s", err)
	}
	req = req.WithContext(ctx)
	req.SetBasicAuth(t.accountSid, t.authToken)

	var resp struct{}
	err = t.Do(req, &resp)

	if err != nil {
		return fmt.Errorf("SMS to %s failed: %s", to, err)
	}

	return nil
}
