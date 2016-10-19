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
)

func getJson(url string, target interface{}) {
    r, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer r.Body.Close()

    if err = json.NewDecoder(r.Body).Decode(target); err != nil {
    	panic(err)
    }
}

func folderPath(date time.Time) string {
	return fmt.Sprintf("/Volumes/Racing/Feeds/AFL/%s", date.Format("Mon 2 Jan, 2006"))
} 

func main() {
	// initialise defaultDay to be used if no date is specified at run-time
	defaultDay := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	dayPtr := flag.String("d", defaultDay, "day to fetch, YYYY-MM-DD")
	flag.Parse()

	requestedDay, err := time.Parse("2006-01-02", *dayPtr)
	if err != nil {
		panic(err)
	}

	// dayString := dayTime.Format("Mon 2 Jan, 2006")

// UBET

if true {

	// https://ubet.com/api/sports/8/102
	// https://ubet.com/api/sports/18/101
	// https://ubet.com/api/sports/19/144

	addresses := []struct{Address, RaceType string}{
		{ Address: "https://ubet.com/api/sports/8/102", RaceType: "R" },
		{ Address: "https://ubet.com/api/sports/18/101", RaceType: "H" },
		{ Address: "https://ubet.com/api/sports/19/144", RaceType: "G" },
	}

	for _, address := range addresses {
		ubet := UBETPayload{}

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

					thisRaces := []SubEvent{}

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
								Number: 1,
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
			// fmt.Printf("%s", meeting)
			createFile := createCompiler("UBET", folderPath(meeting.DateFormat))
			createFile(meeting)
		}
	}
}

// TAB
	if true {
		m := TABPayload{}
		getJson(fmt.Sprintf("https://api.beta.tab.com.au/v1/tab-info-service/racing/dates/%s/meetings?jurisdiction=VIC", *dayPtr), &m)

		for _, meeting := range m.Meetings {
			races := Races{}

			dateFormat, err := time.Parse("2006-01-02", meeting.Date)
			if err != nil {
				panic(err)
			}

			// races are stored differently depending on the state of the meeting
			if (meeting.Mnemonic != "") {
				getJson(meeting.Links.Races, &races)
			} else {
				races.Races = meeting.Races
			}

			fixedOddsRaces := []Race{}

			for _, race := range races.Races {
				if (race.HasFixed && race.Link.Self != "") {
					fixedOddsRaces = append(fixedOddsRaces, race)
				}
			}

			meeting.Races = []Race{}

			if (len(fixedOddsRaces) > 0) {
				for _, race := range fixedOddsRaces {
					getJson(race.Link.Self, &race)

					for _, runner := range race.TABRunners {
						if ((runner.Odds.Win != 0.00) && (runner.Odds.Status == "Open")) {
							race.Runners = append(race.Runners, Runner{
								Name: runner.Name,
								Price: runner.Odds.Win,	
							})
						}
					}

					race.Name = fmt.Sprintf("Race %d", race.Number)

					meeting.Races = append(meeting.Races, race)
				}

				createFile := createCompiler("TAB", folderPath(dateFormat))
				createFile(meeting)
			}
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

		fmt.Printf("%s\n", filePath)

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
