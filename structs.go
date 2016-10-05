package main

type Payload struct {
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
	Runners []Runner `json:"runners"`
	HasFixed bool `json:"hasFixedOdds"`
}

type Runner struct {
	Name string `json:"runnerName"`
	Number int `json:"runnerNumber"`
	Odds struct{
			Win float32 `json:"returnWin"`
			Status string `json:"bettingStatus"`
		} `json:"fixedOdds"`
}
