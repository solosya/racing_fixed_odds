package main

type TABPayload struct {
	Meetings []Meeting `json:"meetings"`
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
