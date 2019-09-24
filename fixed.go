package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func pprint(s interface{}) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func getJSON(url string, target interface{}) {
	r, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(target); err != nil {
		// responseData,_ := ioutil.ReadAll(r.Body)
		// fmt.Printf("%s", responseData)
		panic(err)

	}
}

func folderPath(date time.Time) string {
	return fmt.Sprintf("/Volumes/Racing/Feeds/Fixies/%s", date.Format("Mon 2 Jan, 2006"))
	// return fmt.Sprintf("/go/src/github.com/user/myProject/app/files/%s", date.Format("Mon 2 Jan, 2006"))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	// initialise defaultDay to be used if no date is specified at run-time
	defaultDay := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	api := "TAB,UBET"
	// rtype := ""
	dayPtr := flag.String("d", defaultDay, "day to fetch, YYYY-MM-DD")
	apiPtr := flag.String("api", api, "API's to fetch from (Tab, Ubet, Ladbrokes")
	// raceTypePtr := flag.String("type", rtype, "Race types: (T)horoughbred, (H)arnes, (G)reyhounds")

	flag.Parse()

	requestedDay, err := time.Parse("2006-01-02", *dayPtr)
	if err != nil {
		panic(err)
	}

	apiArr := strings.Split(*apiPtr, ",")
	// fmt.Printf("%s\n", raceType)

	// LADBROKES
	if stringInSlice("LAD", apiArr) {

		Meetings := make(map[string]LadMeetings)
		Runners := make(map[string]eventRaceDetails)
		finalMeetings := make(map[string]Meeting)

		date := requestedDay.Format("2006-01-02")
		url := fmt.Sprintf("https://api-affiliates.ladbrokes.com.au/racing/meetings?date_from=%s&date_to=%s&country=AUS", date, date)
		// fmt.Println(url)
		getJSON(url, &Meetings)

		for _, meeting := range Meetings["data"].Meetings {
			var finalRaces []Race

			for _, race := range meeting.Races {
				url := fmt.Sprintf("https://api-affiliates.ladbrokes.com.au/racing/events/%s", race.ID)

				getJSON(url, &Runners)
				var finalRunners []Runner

				for _, runner := range Runners["data"].Runners {

					if runner.Odds.FixedWin < 1 {
						continue
					}

					runnerTemp := Runner{
						Name:  runner.Name,
						Price: runner.Odds.FixedWin,
					}
					finalRunners = append(finalRunners, runnerTemp)

				}

				raceTemp := Race{
					Name:       "Race " + strconv.Itoa(race.RaceNumber),
					RaceNumber: race.RaceNumber,
					Runners:    finalRunners,
				}
				if len(raceTemp.Runners) > 0 {
					finalRaces = append(finalRaces, raceTemp)
				}

			}

			meetingid := meeting.Name + meeting.Category

			finalMeetings[meetingid] = Meeting{
				Name:     meeting.Name,
				RaceType: meeting.Category,
				Date:     meeting.Date.Format("2006-01-02"),
				Races:    finalRaces,
			}
		}

		for _, meeting := range finalMeetings {
			if len(meeting.Races) < 1 {
				continue
			}
			createFile := createCompiler("LADBROKES", folderPath(requestedDay))
			createFile(meeting)
		}

	}

	// UBET
	// if stringInSlice("UBET", apiArr) {

	// 	addresses := []struct{ Address, RaceType string }{
	// 		{Address: "https://tab.ubet.com/api/sports/8/102", RaceType: "R"},
	// 		{Address: "https://tab.ubet.com/api/sports/18/101", RaceType: "H"},
	// 		{Address: "https://tab.ubet.com/api/sports/19/144", RaceType: "G"},
	// 	}

	// 	for _, address := range addresses {
	// 		ubet := UBETPayload{}

	// 		getJSON(address.Address, &ubet)

	// 		ubetMeetings := make(map[string]Meeting)

	// 		for _, meeting := range ubet.Meetings {

	// 			dateFormat, err := time.Parse(time.RFC3339, meeting.Time)
	// 			if err != nil {
	// 				panic(err)
	// 			}

	// 			if dateFormat.Format("2006-01-02") == requestedDay.Format("2006-01-02") {
	// 				meetingId := fmt.Sprintf("%s%s", meeting.Id, dateFormat.Format("Mon 2 Jan, 2006"))

	// 				if _, ok := ubetMeetings[meetingId]; !ok {
	// 					venue := strings.Split(meeting.Venue, " - ")

	// 					ubetMeetings[meetingId] = Meeting{
	// 						Name:       strings.Replace(venue[0], "/", " ", -1),
	// 						Date:       dateFormat.Format("2006-01-02"),
	// 						DateFormat: dateFormat,
	// 						RaceType:   address.RaceType,
	// 					}
	// 				}

	// 				val := ubetMeetings[meetingId]

	// 				for _, race := range meeting.Events {
	// 					runners := []Runner{}

	// 					thisRaces := []UBETSubEventJSON{}

	// 					getJSON(fmt.Sprintf("https://ubet.com/api/bettingblock?subEventIds=%d", race.Id), &thisRaces)

	// 					for _, thisRace := range thisRaces {

	// 						if thisRace.HasFixed {
	// 							for _, offer := range thisRace.Offers {
	// 								runners = append(runners, Runner{
	// 									Name:  offer.Name,
	// 									Price: offer.Price,
	// 								})
	// 							}

	// 							val.Races = append(val.Races, Race{
	// 								RaceNumber: 1,
	// 								Runners:    runners,
	// 								Name:       meeting.Name,
	// 							})
	// 						}
	// 					}
	// 				}

	// 				ubetMeetings[meetingId] = val
	// 			}
	// 		}

	// 		for _, meeting := range ubetMeetings {
	// 			createFile := createCompiler("UBET", folderPath(meeting.DateFormat))
	// 			createFile(meeting)
	// 		}
	// 	}
	// }

	//TAB
	if stringInSlice("TAB", apiArr) {
		m := TABPayloadJSON{}
		finalMeetings := make(map[string]Meeting)

		getJSON(fmt.Sprintf("https://api.beta.tab.com.au/v1/tab-info-service/racing/dates/%s/meetings?jurisdiction=VIC", *dayPtr), &m)

		for _, meeting := range m.Meetings {
			races := TABRacesJSON{}

			// races are stored differently depending on the state of the meeting
			if meeting.Mnemonic != "" {
				getJSON(meeting.Links.Races, &races)
			} else {
				races.Races = meeting.Races
			}

			fixedOddsRaces := []TABRaceJSON{}

			for _, race := range races.Races {
				if race.HasFixed && race.Link.Self != "" {
					fixedOddsRaces = append(fixedOddsRaces, race)
				}
			}

			meeting.Races = []TABRaceJSON{}
			var finalRaces []Race

			if len(fixedOddsRaces) > 0 {
				for _, race := range fixedOddsRaces {

					var finalRunners []Runner

					getJSON(race.Link.Self, &race)

					for _, runner := range race.TABRunners {
						if (runner.Odds.Win != 0.00) && (runner.Odds.Status == "Open") {
							finalRunners = append(finalRunners, Runner{
								Name:  runner.Name,
								Price: runner.Odds.Win,
							})
						}
					}

					finalRaces = append(finalRaces, Race{
						Name:       fmt.Sprintf("Race %d", race.Number),
						RaceNumber: race.Number,
						Runners:    finalRunners,
					})
				}
			}

			meetingid := fmt.Sprintf("%s%s", meeting.Name, meeting.RaceType)

			dateFormat, err := time.Parse("2006-01-02", meeting.Date)
			if err != nil {
				panic(err)
			}

			finalMeetings[meetingid] = Meeting{
				Name:       meeting.Name,
				RaceType:   meeting.RaceType,
				Date:       meeting.Date,
				DateFormat: dateFormat,
				Races:      finalRaces,
			}
		}

		for _, meeting := range finalMeetings {
			if len(meeting.Races) > 0 {
				createFile := createCompiler("TAB", folderPath(meeting.DateFormat))
				createFile(meeting)
			}
		}
	}
}

func createCompiler(agency string, folderPath string) func(Meeting) {
	return func(meeting Meeting) {
		r := regexp.MustCompile(" \\(.+")
		in := []byte(meeting.Name)
		out := r.ReplaceAll(in, []byte(""))

		var buffer bytes.Buffer

		meetingName := strings.Title(strings.ToLower(string(out))) // convert uppercase string to lowercase, then titlecase that -- titlecasing uppercase text does not work

		buffer.WriteString(fmt.Sprintf("\n%s %s %s\n", meetingName, meeting.RaceType, meeting.Date))

		races := meeting.Races

		for _, race := range races {
			buffer.WriteString(fmt.Sprintf("%s\n", race.Name))
			getFixedOdds(race.Runners, &buffer)
		}

		fileName := fmt.Sprintf("%s %s - %s Fixed Odds - %s.txt", meetingName, meeting.RaceType, agency, meeting.Date)

		os.Mkdir(fmt.Sprintf("%s", folderPath), 0644)

		filePath := fmt.Sprintf("%s/%s", folderPath, fileName)

		// fmt.Printf("%s\n", filePath)

		if err := ioutil.WriteFile(filePath, buffer.Bytes(), 0644); err != nil {
			panic(err)
		}
	}
}

func odds(p1, p2 *Runner) bool {
	return p1.Price < p2.Price
}

func getFixedOdds(runners []Runner, buffer *bytes.Buffer) {
	By(odds).Sort(runners)

	for _, runner := range runners {
		// remove anything in brackets, eg: (Emergency)
		r := regexp.MustCompile(" \\(.+")
		in := []byte(runner.Name)
		out := r.ReplaceAll(in, []byte(""))

		name := strings.Title(strings.ToLower(string(out))) // convert uppercase string to lowercase, then titlecase that -- titlecasing uppercase text does not work

		// lowercase any letter immediately following an apostrophe
		// this is largely redundant since everything will be capped for comparison later on
		r = regexp.MustCompile("'[A-Z]")
		in = []byte(name)
		out = r.ReplaceAllFunc(in, bytes.ToLower)

		if runner.Price > 1.00 {
			buffer.WriteString(fmt.Sprintf("%6.2f %s\n", runner.Price, string(out)))
		}
	}

	buffer.WriteString(fmt.Sprintf("\n"))
}
