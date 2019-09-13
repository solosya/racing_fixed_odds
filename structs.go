package main

import "time"

// LADBROKES data structures

// LadMeetings ...
type LadMeetings struct {
	Meetings []*LadMeetingData `json:"meetings"`
}

// LadMeetingData ...
type LadMeetingData struct {
	Meeting      string          `json:"meeting"`
	Name         string          `json:"name"`
	Date         time.Time       `json:"date"`
	Category     string          `json:"category"`
	CategoryName string          `json:"category_name"`
	Country      string          `json:"country"`
	State        string          `json:"state"`
	Races        []*LadRacesJSON `json:"races"`
}

// LadRacesJSON ...
type LadRacesJSON struct {
	ID             string    `json:"id"`
	RaceNumber     int       `json:"race_number"`
	Name           string    `json:"name"`
	StartTime      time.Time `json:"start_time"`
	TrackCondition string    `json:"track_condition"`
	Distance       uint32    `json:"distance"`
	Weather        string    `json:"weather"`
	Country        string    `json:"country"`
	State          string    `json:"state"`
}

// eventRaceDetails ...
type eventRaceDetails struct {
	// Race eventRace `json:"race"`
	// Results   []*eventResult `json:"results,omitempty"`
	Favourite eventEntrant      `json:"favourite"`
	Runners   []*LADEventRunner `json:"runners"`
	// Mover     eventEntrant   `json:"mover"`
	Error string `json:"error"`
}

// eventEntrant ...
type eventEntrant struct {
	EntrantID    string `json:"-"`
	Name         string `json:"name"`
	IsScratched  bool   `json:"is_scratched"`
	ScratchTime  int64  `json:"scratch_time"`
	Barrier      uint32 `json:"barrier"`
	RunnerNumber int    `json:"runner_number"`
	PrizeMoney   string `json:"prize_money"`
	Age          int    `json:"age"`
	Sex          string `json:"sex"`
	Colour       string `json:"colour"`
	SilkColours  string `json:"silk_colours"`
	FormComment  string `json:"form_comment"`
	ClassLevel   string `json:"class_level"`
	Jockey       string `json:"jockey"`
	Country      string `json:"country"`
	TrainerName  string `json:"trainer_name"`
	Weight       struct {
		Allocated string `json:"allocated"`
		Total     string `json:"total"`
	} `json:"weight"`
	Favourite bool `json:"favourite"`
	Mover     bool `json:"mover"`
}

// LADEventRunner ...
type LADEventRunner struct {
	eventEntrant
	Meta  map[string]string `json:"meta"`
	Flucs []float64         `json:"flucs"`
	Odds  struct {
		FixedWin float64 `json:"fixed_win"`
	}
	ScrTime            *time.Time `json:"scr_time"`
	CompetitorID       string     `json:"competitor_id"`
	RideGuideExists    bool       `json:"ride_guide_exists"`
	RideGuideThumbnail string     `json:"ride_guide_thumbnail"`
	RideGuideFile      string     `json:"ride_guide_file"`
	Trainer            string     `json:"trainer"`
}

// UBETPayload ...
type UBETPayload struct {
	Meetings []UBETMainEventJSON `json:"MainEvents"`
}

// UBETMainEventJSON ...
type UBETMainEventJSON struct {
	ID     string             `json:"MeetingId"`
	Venue  string             `json:"MeetingName"`
	Time   string             `json:"EventStartTime"`
	Name   string             `json:"EventName"`
	Events []UBETSubEventJSON `json:"SubEvents"`
}

// UBETSubEventJSON ...
type UBETSubEventJSON struct {
	ID       int32           `json:"SubEventId"`
	HasFixed bool            `json:"IsFixedPriceRacing"`
	Offers   []UBETOfferJSON `json:"Offers"`
}

// UBETOfferJSON ...
type UBETOfferJSON struct {
	Name  string  `json:"LongDisplayName"`
	Price float32 `json:"WinReturn"`
}

// TABPayloadJSON ...
type TABPayloadJSON struct {
	Meetings []TABMeetingJSON `json:"meetings"`
}

// TABRacesJSON ...
type TABRacesJSON struct {
	Races []TABRaceJSON `json:"races"`
}

// TABMeetingJSON ...
type TABMeetingJSON struct {
	Name     string        `json:"meetingName"`
	RaceType string        `json:"raceType"`
	Date     string        `json:"meetingDate"`
	Races    []TABRaceJSON `json:"races"`
	Mnemonic string        `json:"venueMnemonic"`
	Links    struct {
		Races string `json:"races"`
	} `json:"_links"`
	DateFormat time.Time
}

// TABRaceJSON ...
type TABRaceJSON struct {
	Name   string
	Number int `json:"raceNumber"`
	Link   struct {
		Self string `json:"self"`
	} `json:"_links,omitempty"`
	TABRunners []TABRunnerJSON `json:"runners"`
	Runners    []Runner
	HasFixed   bool `json:"hasFixedOdds"`
}

// TABRunnerJSON ...
type TABRunnerJSON struct {
	Name   string `json:"runnerName"`
	Number int    `json:"runnerNumber"`
	Odds   struct {
		Win    float64 `json:"returnWin"`
		Status string  `json:"bettingStatus"`
	} `json:"fixedOdds"`
}

// Meeting ...
type Meeting struct {
	Name       string
	RaceType   string
	Date       string
	DateFormat time.Time
	Races      []Race
}

// Race ...
type Race struct {
	Name       string
	RaceNumber int
	Runners    []Runner
}

// Runner ...
type Runner struct {
	Name  string
	Price float64
}
