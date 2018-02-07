package main

import "time"

// LADBROKES data structures

type LADMeeting struct {
	Name string
	RaceType string
	Date string
	races map[int]LADRaces
}

type LADRaces struct {
	Id int 
	name string
	raceNumber int 
	runners map[string]Runner
}

type LadRaces_json struct {
	Id int `json:"EventID"`
	RaceNum int `json:"RaceNum"`
	Meeting string `json:"Meeting"`
	Type string `json:"RaceType"`
	Runners map[string]LADRunners_json
}

type LADRunners_json struct {
	Name string `json:"Name"`
	DetailedPricing struct {
		HasFixed bool `json:"fixedWin"`
		Price float32 `json:"startingPriceGuarantee"`
	} `json:"DetailedPricing"`
}



// UBET data structures
type UBETPayload struct {
	Meetings []UBETMainEvent_json `json:"MainEvents"`
}


type UBETMainEvent_json struct {
	Id 		string		`json:"MeetingId"`
	Venue 	string		`json:"MeetingName"`
	Time 	string		`json:"EventStartTime"`
	Name 	string		`json:"EventName"`
	Events []UBETSubEvent_json	`json:"SubEvents"`
}

type UBETSubEvent_json struct {
	Id int32 `json:"SubEventId"`
	HasFixed bool `json:"IsFixedPriceRacing"`
	Offers []UBETOffer_json `json:"Offers"`	
}

type UBETOffer_json struct {
	Name string `json:"LongDisplayName"`
	Price float32 `json:"WinReturn"`
}







// TAB data structures
type TABPayload_json struct {
	Meetings []TABMeeting_json `json:"meetings"`
}

type TABRaces_json struct {
	Races []TABRace_json `json:"races"`
}

type TABMeeting_json struct {
	Name string `json:"meetingName"`
	RaceType string `json:"raceType"`
	Date string `json:"meetingDate"`
	Races []TABRace_json `json:"races"`
	Mnemonic	string `json:"venueMnemonic"`
	Links struct {
		Races string `json:"races"`
	} `json:"_links"`
	DateFormat time.Time
}

type TABRace_json struct {
	Name string
	Number int `json:"raceNumber"`
	Link struct{
			Self string `json:"self"`
		} `json:"_links,omitempty"`
	TABRunners []TABRunner_json `json:"runners"`
	Runners []Runner
	HasFixed bool `json:"hasFixedOdds"`
}

type TABRunner_json struct {
	Name string `json:"runnerName"`
	Number int `json:"runnerNumber"`
	Odds struct{
			Win float32 `json:"returnWin"`
			Status string `json:"bettingStatus"`
		} `json:"fixedOdds"`
}




// Final data structures
type Meeting struct {
	Name string
	RaceType string
	Date string
	DateFormat time.Time
	Races []Race
}

type Race struct {
	Name string
	RaceNumber int
	Runners []Runner
}

type Runner struct {
	Name string
	Price float32
}
