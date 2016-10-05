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

func main() {
	// initialise defaultDay to be used if no date is specified at run-time
	defaultDay := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	dayPtr := flag.String("d", defaultDay, "day to fetch, YYYY-MM-DD")
	flag.Parse()

	// fetch meeting list for specified day
	m := Payload{}
	getJson(fmt.Sprintf("https://api.beta.tab.com.au/v1/tab-info-service/racing/dates/%s/meetings?jurisdiction=VIC", *dayPtr), &m)

	if (len(m.Meetings) > 0) {
		os.Mkdir(fmt.Sprintf("/Volumes/Racing/Feeds/Fixies/%s", *dayPtr), 0644)

		for _, meeting := range m.Meetings {
			var buffer bytes.Buffer

			buffer.WriteString(fmt.Sprintf("\n%s %s\n", meeting.Name, meeting.Date))
			races := Races{}

			// races are stored differently depending on the state of the meeting
			if (meeting.Mnemonic != "") {
				getJson(meeting.Links.Races, &races)
			} else {
				races.Races = meeting.Races
			}

			for _, race := range races.Races {
				if (race.HasFixed && race.Link.Self != "") {
					getJson(race.Link.Self, &race)

					buffer.WriteString(fmt.Sprintf("Race %s - \n", getNumberWord(race.Number)))
					getFixedOdds(race.Runners, &buffer)
				}
			}

			if err := ioutil.WriteFile(fmt.Sprintf("/Volumes/Racing/Feeds/Fixies/%s/%s.txt", meeting.Date, meeting.Name), buffer.Bytes(), 0644); err != nil {
				panic(err)
			}
		}
	}
}

func odds (p1, p2 *Runner) bool {
	return p1.Odds.Win < p2.Odds.Win
}

func getFixedOdds(runners []Runner, buffer *bytes.Buffer) {
	By(odds).Sort(runners)

	for _, runner := range runners {
		if ((runner.Odds.Win != 0.00) && (runner.Odds.Status == "Open")) {
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

			buffer.WriteString(fmt.Sprintf("%6.2f %s\n", runner.Odds.Win, string(out)))
		}
	}
	
	buffer.WriteString(fmt.Sprintf("\n"))
}
