package main

import (
	"encoding/json"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const url = `https://public.opendatasoft.com/api/records/1.0/search/?dataset=covid-19-germany-landkreise&q=Aachen&facet=last_update&facet=name&facet=rs&facet=bez&facet=bl`

type apiResponse struct {
	Records []struct {
		Fields struct {
			DeathRate          float64   `json:"death_rate"`           // Death rate overall
			Deaths             int       `json:"deaths"`               // Deaths
			CasesPer100K       float64   `json:"cases_per_100k"`       // Cases per 100k population
			CasesPerPopulation float64   `json:"cases_per_population"` // Number of cases in entire population
			Cases7Bl           int       `json:"cases7_bl"`            // Cases over the last 7 days in Bundesland
			LastUpdate         time.Time `json:"last_update"`          // Last updated at
			Ewz                int       `json:"ewz"`                  // Einwohnerzahl
			EwzBl              int       `json:"ewz_bl"`               // Einwohnerzahl Bundesland
			Cases7Per100K      float64   `json:"cases7_per_100k"`      // Cases over the last 7 days per 100k persons
			Bl                 string    `json:"bl"`                   // Bundesland
			Death7Bl           int       `json:"death7_bl"`            // Deaths over the last 7 days in Bundesland
			Cases              int       `json:"cases"`                // Number of cases registered from the beginning
			Name               string    `json:"name"`                 // Name of the Region
		} `json:"fields"`
	} `json:"records"`
}

func getUpdate(cmdHandler *commandHandler, u *telegram.Update, args []string) {
	res, err := cmdHandler.bot.Client.Get(url)
	if err != nil {

	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var response apiResponse

	err = decoder.Decode(&response)
	if err != nil {

	}
}
