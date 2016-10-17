package main

import (
    "net/http"
    "encoding/json"
    "encoding/xml"
    "compress/gzip"
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

func getNumberWord(num int) string {
	return map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
		10: "ten",
		11: "eleven",
		12: "twelve",
	}[num]
}

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

func getXML(url string, target interface{}) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept-Encoding", "gzip")

	r, err := client.Do(req)

	if err != nil{
	    return
	}
    defer r.Body.Close()

    if (r.Status != "404 Not Found") {
	    reader, _ := gzip.NewReader(r.Body)
	    defer reader.Close()

		if err = xml.NewDecoder(reader).Decode(target); err != nil {
			return
		}
	} else {
		return fmt.Errorf("404 Not Found!")
	}

	return
}

func getJsonAuth(url string, target interface{}) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("Rohan", "cham045F-1")
	r, err := client.Do(req)

	if err != nil{
	    panic (err)
	}
    defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(target); err != nil {
		panic(err)
	}
}

func main() {
	// initialise defaultDay to be used if no date is specified at run-time
	defaultDay := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	dayPtr := flag.String("d", defaultDay, "day to fetch, YYYY-MM-DD")
	flag.Parse()

	dayTime, err := time.Parse("2006-01-02", *dayPtr)
	if err != nil {
		panic(err)
	}

	dayString := dayTime.Format("Mon 2 Jan, 2006")

	folderPath := fmt.Sprintf("/Volumes/Racing/Feeds/AFL/%s", dayString)

// UBET
	if true {
		ubetDay := dayTime.Format("2006/01/02")

		api := APIPayload{}
		getJsonAuth(fmt.Sprintf("http://racing-api.pagemasters.com.au/meeting?date=%s", dayTime), &api)

		createFile := createCompiler("UBET", folderPath)

		for _, meeting := range api.Meetings {
			ubetMeeting := Meeting{Date: dayString}

			for _, betting := range meeting.Betting {
				if (betting.Agency == "ubet") {
					m := UBETPayload{}

					meetingPath := fmt.Sprintf("https://tatts.com/pagedata/racing/%s/%s.xml", ubetDay, strings.ToUpper(betting.Code))

					if err := getXML(meetingPath, &m); err == nil {

						fmt.Printf("%s\n", m.Meeting.Name)

						ubetMeeting.Name = m.Meeting.Name

						races := []Race{}

						for _, race := range m.Meeting.Races {
							fmt.Printf("Race!\n")
							r := UBETPayload{}
							getXML(fmt.Sprintf("https://tatts.com/pagedata/racing/%s/%s%s.xml", ubetDay, strings.ToUpper(betting.Code), race.Number), &r)

							runners := []Runner{}

							for _, runner := range r.Meeting.Races[0].Runners {
								fmt.Printf("%s\n", runner)

								if ((runner.Scratched == "N") && (runner.Odds.Price != 0.00)) {
									runners = append(runners, Runner{
										Name: runner.Name,
										Price: runner.Odds.Price,	
									})
								}
							}

							raceNumber, _ := strconv.Atoi(race.Number)

							races = append(races, Race{
								Number: raceNumber,
								Runners: runners,	
							})
						}

						ubetMeeting.Races = races

						createFile(ubetMeeting)
					}
				}
			}
		}
	}

// TAB
	if true {
		m := TABPayload{}
		getJson(fmt.Sprintf("https://api.beta.tab.com.au/v1/tab-info-service/racing/dates/%s/meetings?jurisdiction=VIC", *dayPtr), &m)

		createFile := createCompiler("TAB", folderPath)

		for _, meeting := range m.Meetings {
			races := Races{}

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

					meeting.Races = append(meeting.Races, race)
				}

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

		buffer.WriteString(fmt.Sprintf("\n%s %s\n", meetingName, meeting.Date))
		races := meeting.Races

		for _, race := range races {
			buffer.WriteString(fmt.Sprintf("Race %s - \n", getNumberWord(race.Number)))
			getFixedOdds(race.Runners, &buffer)
		}

		fileName := fmt.Sprintf("%s - %s Fixed Odds.txt", meetingName, agency)

		os.Mkdir(fmt.Sprintf("%s", folderPath), 0644)

		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", folderPath, fileName), buffer.Bytes(), 0644); err != nil {
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

		buffer.WriteString(fmt.Sprintf("%6.2f %s\n", runner.Price, string(out)))
	}
	
	buffer.WriteString(fmt.Sprintf("\n"))
}
