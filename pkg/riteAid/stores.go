package riteaid

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aaronireland/covid-utils/pkg/api"
)

type Store struct {
	Number              int      `json:"storeNumber"`
	Name                string   `json:"name"`
	Type                string   `json:"storeType"`
	Address             string   `json:"address"`
	City                string   `json:"city"`
	State               string   `json:"state"`
	Zipcode             string   `json:"zipcode"`
	FullZipcode         string   `json:"fullZipCode"`
	TimeZone            string   `json:"timeZone"`
	Latitude            float64  `json:"latitude"`
	Longitude           float64  `json:"longitude"`
	FullPhone           string   `json:"fullPhone"`
	LocationDescription string   `json:"locationDescription"`
	HoursMonday         string   `json:"storeHoursMonday"`
	HoursTuesday        string   `json:"storeHoursTuesday"`
	HoursWednesday      string   `json:"storeHoursWednesday"`
	HoursThursday       string   `json:"storeHoursThursday"`
	HoursFriday         string   `json:"storeHoursFriday"`
	HoursSaturday       string   `json:"storeHoursSaturday"`
	HoursSunday         string   `json:"storeHoursSunday"`
	RxHoursMonday       string   `json:"rxHrsMon"`
	RxHoursTuesday      string   `json:"rxHrsTue"`
	RxHoursWednesday    string   `json:"rxHrsWed"`
	RxHoursThursday     string   `json:"rxHrsThu"`
	RxHoursFriday       string   `json:"rxHrsFri"`
	RxHoursSaturday     string   `json:"rxHrsSat"`
	RxHoursSunday       string   `json:"rxHrsSun"`
	MilesFromCenter     float64  `json:"milesFromCenter"`
	SpecialServiceKeys  []string `json:"specialServiceKeys"`
}

type GetStoresResponse struct {
	Status            string `json:"Status"`
	ErrCode           string `json:"ErrCde"`
	ErrMessage        string `json:"ErrMsg"`
	ErrMessageDetails string `json:"ErrMsgDtl"`
	Data              struct {
		Stores []Store `json:"stores"`
	} `json:"Data"`
}

func (s Store) FullDescription() string {
	var lines []string
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("---------- %s Store #%d ----------", s.Name, s.Number))
	lines = append(lines, "")
	lines = append(lines, s.Address)
	lines = append(lines, fmt.Sprintf("%s, %s %s", s.City, s.State, s.Zipcode))
	lines = append(lines, s.LocationDescription)
	lines = append(lines, fmt.Sprintf("Phone: %s", s.FullPhone))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Monday Hours: %s", s.HoursMonday))
	lines = append(lines, fmt.Sprintf("Tuesday Hours: %s", s.HoursTuesday))
	lines = append(lines, fmt.Sprintf("Wednesday Hours: %s", s.HoursWednesday))
	lines = append(lines, fmt.Sprintf("Thursday Hours: %s", s.HoursThursday))
	lines = append(lines, fmt.Sprintf("Friday Hours: %s", s.HoursFriday))
	lines = append(lines, fmt.Sprintf("Saturday Hours: %s", s.HoursSaturday))
	lines = append(lines, fmt.Sprintf("Sunday Hours: %s", s.HoursSunday))

	return strings.Join(lines, "\n")
}

func (s Store) Description() string {
	var lines []string
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("%s Store #%d", s.Name, s.Number))
	lines = append(lines, "")
	lines = append(lines, s.Address)
	lines = append(lines, fmt.Sprintf("%s, %s %s", s.City, s.State, s.Zipcode))
	lines = append(lines, fmt.Sprintf("Phone: %s", s.FullPhone))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Monday Hours: %s", s.HoursMonday))
	lines = append(lines, fmt.Sprintf("Tuesday Hours: %s", s.HoursTuesday))
	lines = append(lines, fmt.Sprintf("Wednesday Hours: %s", s.HoursWednesday))
	lines = append(lines, fmt.Sprintf("Thursday Hours: %s", s.HoursThursday))
	lines = append(lines, fmt.Sprintf("Friday Hours: %s", s.HoursFriday))
	lines = append(lines, fmt.Sprintf("Saturday Hours: %s", s.HoursSaturday))
	lines = append(lines, fmt.Sprintf("Sunday Hours: %s", s.HoursSunday))

	return strings.Join(lines, "\n")
}
func (c *API) GetCovidVaccinationSites(ctx context.Context, zipCode string, radius int) (stores []Store, err error) {
	filterZip := api.Param{"address", zipCode}
	filterRadius := api.Param{"radius", strconv.Itoa(radius)}
	covidVaccineServiceSites := api.Param{"attrFilter", CovidVaccineService}

	req, err := http.NewRequest("GET", c.URI("stores/getStores", filterZip, filterRadius, covidVaccineServiceSites), nil)
	if err != nil {
		return stores, err
	}

	req = req.WithContext(ctx)

	resp := GetStoresResponse{}
	if err := c.Do(req, &resp); err != nil {
		return stores, err
	}

	if resp.Status != Success {
		err = fmt.Errorf("%s Error: %s - %s", resp.ErrCode, resp.ErrMessage, resp.ErrMessageDetails)
		return stores, err
	}

	return resp.Data.Stores, nil
}
