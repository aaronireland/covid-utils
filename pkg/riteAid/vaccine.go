package riteaid

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aaronireland/covid-utils/pkg/api"
)

const (
	CovidVaccineService = "PREF-112"
)

type Slots struct {
	FirstDose  bool `json:"1"`
	SecondDose bool `json:"2"`
}

type CheckSlotsResponse struct {
	Status            string `json:"Status"`
	ErrCode           string `json:"ErrCde"`
	ErrMessage        string `json:"ErrMsg"`
	ErrMessageDetails string `json:"ErrMsgDtl"`
	Data              Slots  `json:"slots"`
}

func (c *API) AppointmentsAvailable(ctx context.Context, store Store) (firstDose, secondDose bool, err error) {
	filterStore := api.Param{"storeNumber", strconv.Itoa(store.Number)}

	req, err := http.NewRequest("GET", c.URI("vaccine/checkSlots", filterStore), nil)
	if err != nil {
		return false, false, err
	}

	req = req.WithContext(ctx)

	resp := CheckSlotsResponse{}
	if err := c.Do(req, &resp); err != nil {
		return false, false, err
	}

	if resp.Status != Success {
		err = fmt.Errorf("%s Error: %s - %s", resp.ErrCode, resp.ErrMessage, resp.ErrMessageDetails)
		return false, false, err
	}

	return resp.Data.FirstDose, resp.Data.SecondDose, nil
}
