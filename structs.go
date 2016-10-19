package main

import "time"

type UBETPayload struct {
	Meetings []MainEvent `json:"MainEvents"`
}

type MainEvent struct {
	Id string `json:"MeetingId"`
	Venue string `json:"MeetingName"`
	Time string `json:"EventStartTime"`
	Name string `json:"EventName"`
	Events []SubEvent `json:"SubEvents"`
}

type SubEvent struct {
	Id int32 `json:"SubEventId"`
	HasFixed bool `json:"IsFixedPriceRacing"`
	Offers []Offer `json:"Offers"`	
}

type Offer struct {
	Name string `json:"LongDisplayName"`
	Price float32 `json:"WinReturn"`
}

type TABPayload struct {
	Meetings []Meeting `json:"meetings"`
}

type Races struct {
	Races []Race `json:"races"`
}

type Meeting struct {
	Name string `json:"meetingName"`
	RaceType string `json:"raceType"`
	Date string `json:"meetingDate"`
	Races []Race `json:"races"`
	Mnemonic	string `json:"venueMnemonic"`
	Links struct {
		Races string `json:"races"`
	} `json:"_links"`
	DateFormat time.Time
}

type Race struct {
	Name string
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
