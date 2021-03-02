/*
*  Script finds available COVID-19 vaccination appointment within a radius of given zipcodes. Currently only Rite-Aid
*  locations are supported. The API only returns 10 stores per request so use several addresses and zipcodes.
 */
package main

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aaronireland/covid-utils/pkg/riteaid"
	"github.com/aaronireland/covid-utils/pkg/twilio"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

// RadiusInMiles the radius filter to pass to the Rite-Aid api. A max of 10 stores is returned so keep this small
const RadiusInMiles = 25

var (
	addresses           []string
	debug, sendSMS      bool
	state, phoneNumbers *string
	smsRecipients       []string
	twilioAPI           *twilio.APIClient
)

func init() {
	log.SetHandler(text.New(os.Stderr))

	debugFlag := flag.Bool("debug", false, "Set to test SMS and add verbose logging")
	state = flag.String("state", "", "Only search sites in this state")
	phoneNumbers = flag.String("sms", "", "A comma-separated list of phones to alert via SMS")
	debug = *debugFlag
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		panic("At least one address or zip-code required")
	}

	for _, arg := range args {
		addresses = append(addresses, strings.Split(arg, ",")...)
	}

	smsRecipients = append(smsRecipients, strings.Split(*phoneNumbers, ",")...)

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")

	var err error
	twilioAPI, err = twilio.NewClient(sid, token)
	if err != nil {
		log.WithError(err).Warn("unable to connect to Twilio. SMS alerts disabled")
	}
}

func main() {
	riteAidAPI := riteaid.NewClient()
	if twilioAPI != nil && len(smsRecipients) > 0 {
		sendSMS = true
		log.Debug("SMS is enabled")
	}

	ctx := context.Background()

	var sitesToCheck []riteaid.Store
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	log.Infof("Checking for Rite-Aid stores within %d mile radius of %d addresses", RadiusInMiles, len(addresses))
	for _, address := range addresses {
		if len(address) < 5 {
			continue
		}
		log.Debugf("Checking: %s...", address)
		waitFor := time.Duration(r.Intn(900) + 75)
		time.Sleep(time.Millisecond * waitFor)
		stores, err := riteAidAPI.GetCovidVaccinationSites(ctx, address, RadiusInMiles)

		if err != nil {
			log.Fatalf("Failed to check Rite-Aid covid vaccination sites: %w", err)
		}

		for _, store := range stores {
			var exists bool
			for _, site := range sitesToCheck {
				if site.Number == store.Number {
					exists = true
					break
				}
			}
			if !exists && (store.State == *state || *state == "") {
				sitesToCheck = append(sitesToCheck, store)
			}
		}
	}

	log.Infof("Found %d sites", len(sitesToCheck))
	var sitesWithAvailability []riteaid.Store
	randomSite := r.Intn(len(sitesToCheck))

	for ix, store := range sitesToCheck {

		log.Debugf("Store #%d: %s %s, %s, %s", store.Number, store.Address, store.City, store.State, store.Zipcode)
		waitFor := time.Duration(r.Intn(500) + 75)
		time.Sleep(time.Millisecond * waitFor)
		first, second, err := riteAidAPI.AppointmentsAvailable(ctx, store)

		if err != nil {
			log.Errorf("Failed to check vaccination availability at Rite-Aid Store %d: %s", store.Number, err)
			continue
		}

		if first || second || (debug && ix == randomSite) {
			sitesWithAvailability = append(sitesWithAvailability, store)
		}
	}

	if len(sitesWithAvailability) == 0 {
		log.Info("No vaccination appointment availability...")
	} else {
		s := "sites"
		if len(sitesWithAvailability) == 1 {
			s = "site"
		}

		log.Infof("Vaccination appointment found at %d %s!", len(sitesWithAvailability), s)

		for _, store := range sitesWithAvailability {
			log.Info(store.FullDescription())
			if sendSMS {
				for _, to := range smsRecipients {
					err := twilioAPI.SendSMS(ctx, to, store.Description())
					if err != nil {
						log.WithField("sms_phone", to).WithError(err).Error("SMS notification failed!")
					}

				}
			}
		}
	}
}
