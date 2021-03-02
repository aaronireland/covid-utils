package riteaid

import (
	"github.com/aaronireland/covid-utils/pkg/api"
)

// BaseURLV2 is the base url for the Rite-Aid v2 services api
const (
	BaseURLV2 = "https://www.riteaid.com/services/ext/v2"
	Success   = "SUCCESS"
)

// API is the RiteAid API client for the v2 services REST API
type API struct {
	*api.Client
}

// NewClient is the constructor to instantiate a new instance of the Rite-Aid API Client
func NewClient() *API {
	return &API{api.NewClient(BaseURLV2, api.JSONClient)}
}
