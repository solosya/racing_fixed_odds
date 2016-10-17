package main

import "encoding/xml"

type APIPayload struct {
	Meetings []APIMeeting `json:"data"`
}

type APIMeeting struct {
	Name string `json:"venue"`
	Date string `json:"date"`
	Betting []Betting `json:"betting"`
}

type Betting struct {
	Code string `json:"betcode"`
	Agency string `json:"agency"`
}

type TABPayload struct {
	Meetings []Meeting `json:"meetings"`
}

type UBETPayload struct {
	XMLName xml.Name `xml:"RaceDay"`
	Meeting XMLMeeting
}

type XMLMeeting struct {
	XMLName xml.Name `xml:"Meeting"`
	Name string `xml:"VenueName,attr"`
	Races []XMLRace `xml:"Race"`
}

type XMLRace struct {
	XMLName xml.Name `xml:"Race"`
	Number string `xml:"RaceNo,attr"`
	Runners []XMLRunner `xml:"Runner"`
}

type XMLRunner struct {
	XMLName xml.Name `xml:"Runner"`
	Name string `xml:"RunnerName,attr"`
	Scratched string `xml:"Scratched,attr"`
	Odds XMLFixed `xml:"FixedOdds"`
}

type XMLFixed struct {
	XMLName xml.Name `xml:"FixedOdds"`
	Price float32 `xml:"WinOdds,attr"`
}

type Races struct {
	Races []Race `json:"races"`
}

type Meeting struct {
	Name string `json:"meetingName"`
	Date string `json:"meetingDate"`
	Races []Race `json:"races"`
	Mnemonic	string `json:"venueMnemonic"`
	Links struct {
		Races string `json:"races"`
	} `json:"_links"`
}

type Race struct {
	Number int `json:"raceNumber"`
	Link struct{
			Self string `json:"self"`
		} `json:"_links,omitempty"`
	TABRunners []TABRunner `json:"runners"`
	Runners []Runner
	HasFixed bool `json:"hasFixedOdds"`
}

type TABRunner struct {
	Name string `json:"runnerName"`
	Number int `json:"runnerNumber"`
	Odds struct{
			Win float32 `json:"returnWin"`
			Status string `json:"bettingStatus"`
		} `json:"fixedOdds"`
}

type Runner struct {
	Name string
	Price float32
}
