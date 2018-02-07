package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "strings"
    "bytes"
    "regexp"
    "io/ioutil"
    "flag"
    "os"
    "time"
    "strconv"
)

func pprint(s interface{}) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
	    fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func getJson(url string, target interface{}) {
	// fmt.Println(url)
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
	return fmt.Sprintf("/Users/neenanl/Racing/Feeds/%s", date.Format("Mon 2 Jan, 2006"))
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
	api    := "TAB,UBET"
	dayPtr := flag.String("d", defaultDay, "day to fetch, YYYY-MM-DD")
	apiPtr := flag.String("api", api, "API's to fetch from (Tab, Ubet, Ladbrokes")

	flag.Parse()

	requestedDay, err := time.Parse("2006-01-02", *dayPtr)
	if err != nil {
		panic(err)
	}

	api_arr := strings.Split(*apiPtr, ",")

	// fmt.Printf("%+v\n", api_arr)

	// LADBROKES
	if stringInSlice("LAD", api_arr) {
		
		// fmt.Println("doing ubt")
		// fmt.Println(requestedDay)
		
		LADMeetings := make(map[string]LADMeeting)
		races_json  := make(map[string]LadRaces_json)
		
		getJson(fmt.Sprintf("https://www.ladbrokes.com.au/api/feed/racingList?date=%s", requestedDay.Format("2006-01-02")), &races_json)
		itr := 0
		for _, race := range races_json {
			itr++

			time.Sleep(1 * time.Second)

			raceType    := race.Type
			if raceType == "T" {
				raceType = "R"
			}

			raceId 		:= race.Id
			raceNumber 	:= race.RaceNum
			meetingId 	:= race.Meeting + raceType

			if _, ok := LADMeetings[meetingId]; !ok {
				LADMeetings[meetingId] = LADMeeting {
					Name: race.Meeting,
					RaceType: raceType,
					Date: requestedDay.Format("2006-01-02"),
					races: map[int]LADRaces{},
				}
			}



			if _, a_ok := LADMeetings[meetingId].races[raceNumber]; !a_ok {

				LadRunners  := make(map[string] map[string] map[string]LADRunners_json)
				getJson(fmt.Sprintf("https://www.ladbrokes.com.au/api/feed/eventRunners?event_id=%d", raceId), &LadRunners)

				LADMeetings[meetingId].races[raceNumber] = LADRaces {
					Id: 		race.Id,
					name: 		"Race " + strconv.Itoa(race.RaceNum),
					raceNumber:	race.RaceNum,
					runners: 	map[string]Runner{},
				}


				Loop: // break back to here in event of no fixed odds on race (event)
				for _, event := range LadRunners {
					for _, competitors := range event {
						for k, horse := range competitors {
							// fmt.Println(horse.DetailedPricing.HasFixed)
							if horse.DetailedPricing.HasFixed == false || horse.DetailedPricing.Price < 1 {
								
								// If horse has zero price or no fixed odds, delete this race from the result
								delete(LADMeetings[meetingId].races, raceNumber)
								break Loop
							}
							LADMeetings[meetingId].races[raceNumber].runners[k] = Runner {
								Name: horse.Name,
								Price: horse.DetailedPricing.Price,
							}
						}
					}	
				}
			}

			// if itr > 100 {
			// 	break
			// }
		}



		finalMeetings := make(map[string]Meeting)
		createFile := createCompiler("LADBROKE", folderPath(requestedDay))
		
		for _, meeting := range LADMeetings {
			// fmt.Println(meeting.Name)

			if len(meeting.races) == 0 {
				continue
			}

			var finalRaces []Race
			for _, race := range meeting.races {
				
				var finalRunners []Runner
			
				for _, runner := range race.runners {
					runnerTemp := Runner {
						Name: runner.Name,
						Price: runner.Price,
					}

					finalRunners = append(finalRunners, runnerTemp)
				}

				raceTemp := Race {
					Name: race.name,
					RaceNumber: race.raceNumber,
					Runners: finalRunners,
				}

				finalRaces = append(finalRaces, raceTemp)
			}


			meetingid := meeting.Name + meeting.RaceType

			finalMeetings[meetingid] = Meeting {
				Name: meeting.Name,
				RaceType: meeting.RaceType,
				Date: meeting.Date,
				Races: finalRaces,
			}
		}

		for _, meeting := range finalMeetings {
			createFile(meeting)
		}


		// fmt.Printf("%+v\n", finalMeetings)
		// pprint(finalMeetings)
	}




	// UBET
	if stringInSlice("UBET", api_arr) {

		addresses := []struct{Address, RaceType string} {
			{ Address: "https://ubet.com/api/sports/8/102",  RaceType: "R" },
			{ Address: "https://ubet.com/api/sports/18/101", RaceType: "H" },
			{ Address: "https://ubet.com/api/sports/19/144", RaceType: "G" },
		}

		for _, address := range addresses {
			ubet := UBETPayload{}
			fmt.Println(address.Address)

			getJson(address.Address, &ubet)

			ubetMeetings := make(map[string]Meeting)

			for _, meeting := range ubet.Meetings {

				dateFormat, err := time.Parse(time.RFC3339, meeting.Time)
				if err != nil {
					panic(err)
				}


				if (dateFormat.Format("2006-01-02") == requestedDay.Format("2006-01-02")) {
					meetingId := fmt.Sprintf("%s%s", meeting.Id, dateFormat.Format("Mon 2 Jan, 2006"))

					if _, ok := ubetMeetings[meetingId]; !ok {
						venue := strings.Split(meeting.Venue, " - ")

						ubetMeetings[meetingId] = Meeting{
							Name: strings.Replace(venue[0], "/", " ", -1),
							Date: dateFormat.Format("2006-01-02"),
							DateFormat: dateFormat,
							RaceType: address.RaceType,
						}
					}

					val := ubetMeetings[meetingId]

					for _, race := range meeting.Events {
						runners := []Runner{}

						thisRaces := []UBETSubEvent_json{}

						getJson(fmt.Sprintf("https://ubet.com/api/bettingblock?subEventIds=%d", race.Id), &thisRaces)
						

						for _, thisRace := range thisRaces {


							if (thisRace.HasFixed) {
								for _, offer := range thisRace.Offers {
									runners = append(runners, Runner{
										Name: offer.Name,
										Price: offer.Price,	
									})
								}

								val.Races = append(val.Races, Race{
									RaceNumber: 1,
									Runners: runners,
									Name: meeting.Name,
								})
							}
						}
					}

					ubetMeetings[meetingId] = val
				}
			}

			for _, meeting := range ubetMeetings {
				createFile := createCompiler("UBET", folderPath(meeting.DateFormat))
				createFile(meeting)
			}
		}
	}

	//TAB
	if stringInSlice("TAB", api_arr) {
		m := TABPayload_json{}
		finalMeetings := make(map[string]Meeting)

		getJson(fmt.Sprintf("https://api.beta.tab.com.au/v1/tab-info-service/racing/dates/%s/meetings?jurisdiction=VIC", *dayPtr), &m)
		

		for _, meeting := range m.Meetings {
			races := TABRaces_json{}

			// races are stored differently depending on the state of the meeting
			if (meeting.Mnemonic != "") {
				getJson(meeting.Links.Races, &races)
			} else {
				races.Races = meeting.Races
			}

			fixedOddsRaces := []TABRace_json{}

			for _, race := range races.Races {
				if (race.HasFixed && race.Link.Self != "") {
					fixedOddsRaces = append(fixedOddsRaces, race)
				}
			}

			meeting.Races = []TABRace_json{}
			var finalRaces []Race

			if (len(fixedOddsRaces) > 0) {
				for _, race := range fixedOddsRaces {

					var finalRunners []Runner

					getJson(race.Link.Self, &race)

					for _, runner := range race.TABRunners {
						if ((runner.Odds.Win != 0.00) && (runner.Odds.Status == "Open")) {
							finalRunners = append(finalRunners, Runner{
								Name: runner.Name,
								Price: runner.Odds.Win,	
							})
						}
					}

					finalRaces = append(finalRaces, Race {
						Name: fmt.Sprintf("Race %d", race.Number),
						RaceNumber: race.Number,
						Runners: finalRunners,
					})
				}
			}
			
			meetingid := fmt.Sprintf("%s%s", meeting.Name, meeting.RaceType)
			
			dateFormat, err := time.Parse("2006-01-02", meeting.Date)
			if err != nil {
				panic(err)
			}

			finalMeetings[meetingid] = Meeting {
				Name: 		meeting.Name,
				RaceType: 	meeting.RaceType,
				Date: 		meeting.Date,
				DateFormat: dateFormat,
				Races: 		finalRaces,
			}
		}

		for _, meeting := range finalMeetings {
			createFile := createCompiler("TAB", folderPath(meeting.DateFormat))
			createFile(meeting)
		}
	}
}

func createCompiler(agency string, folderPath string) func (Meeting) {
	return func (meeting Meeting) {
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

func odds (p1, p2 *Runner) bool {
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

		if (runner.Price > 1.00) {
			buffer.WriteString(fmt.Sprintf("%6.2f %s\n", runner.Price, string(out)))
		}
	}
	
	buffer.WriteString(fmt.Sprintf("\n"))
}

