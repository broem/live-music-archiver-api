package repo

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	tableName        struct{}  `pg:"captured"`
	UserID           string    `json:"user_id"`
	ID               string    `json:"eventID" pg:"captured_id,pk,type:uuid"`
	Title            string    `json:"eventTitle"`
	DescriptionURL   string    `json:"eventDescURL"`
	Description      string    `json:"eventDesc"`
	URL              string    `json:"eventURL"`
	Date             string    `json:"eventDate"`
	Time             string    `json:"eventTime"`
	Venue            string    `json:"eventVenue"`
	VenueAddress     string    `json:"eventVenueAddress"`
	VenueContactInfo string    `json:"eventVenueContactInfo"`
	TicketCost       string    `json:"eventTicketCost"`
	TicketURL        string    `json:"eventTicketURL"`
	OtherPerformers  string    `json:"eventOtherPerformers"`
	AgeRequired      string    `json:"eventAgeRequired"`
	FacebookURL      string    `json:"eventFacebookURL"`
	TwitterURL       string    `json:"eventTwitterURL"`
	Misc             string    `json:"eventMisc"`
	Images           []string  `json:"eventImages"`
	CaptureDate      time.Time `json:"captureDate"`
	MapID            uuid.UUID `json:"mapId" pg:"type:uuid"`
	Cbsa             string    `json:"cbsa"`
	StateFips        string    `json:"stateFips"`
	CountyFips       string    `json:"countyFips"`
}

// Event to string
func (e Event) String() string {
	// returns all fields comma separated
	ret := e.UserID + "," + e.ID + "," + e.Title + "," + e.DescriptionURL + "," + e.Description + "," + e.URL + "," + e.Date + "," + e.Time + "," + e.Venue + "," + e.VenueAddress + "," + e.VenueContactInfo + "," + e.TicketCost + "," + e.TicketURL + "," + e.OtherPerformers + "," + e.AgeRequired + "," + e.FacebookURL + "," + e.TwitterURL + "," + e.Misc + "," + e.CaptureDate.String()

	// remove commas next to each other
	// ret = strings.ReplaceAll(ret, ",,", ",")
	return ret
}

// EventStringSlice returns a slice of strings from an Event
func (e *Event) EventStringSlice() []string {
	return []string{e.UserID, e.ID, e.Title, e.DescriptionURL, e.Description, e.URL, e.Date, e.Time, e.Venue, e.VenueAddress, e.VenueContactInfo, e.TicketCost, e.TicketURL, e.OtherPerformers, e.AgeRequired, e.FacebookURL, e.TwitterURL, e.Misc, e.CaptureDate.String()}
}

type EventWithCount struct {
	Event      *Event
	EventCount int
}

type EventMapper struct {
	tableName                struct{}  `pg:"mappers"`
	MapID                    uuid.UUID `json:"map_id" pg:"type:uuid,pk"`
	VenueBaseURL             string    `json:"venue_base_url"`
	FullEventSelector        string    `json:"full_event"`
	UserID                   string    `json:"user_id"`
	TitleSelector            string    `json:"eventTitle"`
	DescriptionSelector      string    `json:"eventDesc"`
	DescriptionURLSelector   string    `json:"eventURL"`
	DateSelector             string    `json:"eventDate"`
	TimeSelector             string    `json:"eventTime"`
	VenueNameSelector        string    `json:"eventVenue"`
	VenueAddressSelector     string    `json:"eventVenueAddress"`
	VenueContactInfoSelector string    `json:"eventVenueContactInfo"`
	TicketCostSelector       string    `json:"eventTicketCost"`
	TicketURLSelector        string    `json:"eventTicketURL"`
	OtherPerformersSelector  string    `json:"eventOtherPerformers"`
	AgeRequiredSelector      string    `json:"eventAgeRequired"`
	FacebookURLSelector      string    `json:"eventFacebookURL"`
	TwitterURLSelector       string    `json:"eventTwitterURL"`
	MiscSelector             string    `json:"eventMisc"`
	ImagesSelector           []string  `json:"eventImages"`
	Approved                 bool      `json:"approved"`
	Cbsa                     string    `json:"cbsa"`
	StateFips                string    `json:"state_fips"`
	CountyFips               string    `json:"county_fips"`
}

type Enabled struct {
	Enabled bool      `json:"enabled"`
	MapID   uuid.UUID `json:"mapId"`
	UserID  string    `json:"userId"`
}
