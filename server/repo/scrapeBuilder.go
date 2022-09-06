package repo

import "strings"

type ScrapeBuilder struct {
	UserID           string                `json:"userId,omitempty"`
	UserEmail        string                `json:"userEmail,omitempty"`
	Event            ScrapeBuilderContents `json:"event"`
	VenueName        ScrapeBuilderContents `json:"venueName"`
	VenueAddress     ScrapeBuilderContents `json:"venueAddress"`
	VenueContactInfo ScrapeBuilderContents `json:"venueContactInfo"`
	EventTitle       ScrapeBuilderContents `json:"eventTitle"`
	EventDescURL     ScrapeBuilderContents `json:"eventDescURL"`
	EventDesc        ScrapeBuilderContents `json:"eventDesc"`
	Images           ScrapeBuilderContents `json:"images"`
	StartDate        ScrapeBuilderContents `json:"startDate"`
	EndDate          ScrapeBuilderContents `json:"endDate"`
	DoorTime         ScrapeBuilderContents `json:"doorTime"`
	TicketCost       ScrapeBuilderContents `json:"ticketCost"`
	TicketURL        ScrapeBuilderContents `json:"ticketURLS"`
	OtherPerformers  ScrapeBuilderContents `json:"otherPerformers"`
	EventURLS        ScrapeBuilderContents `json:"eventURLS"`
	AgeRequired      ScrapeBuilderContents `json:"ageRequired"`
	FacebookURL      ScrapeBuilderContents `json:"facebookURL"`
	TwitterURL       ScrapeBuilderContents `json:"twitterURL"`
	Misc             ScrapeBuilderContents `json:"misc"`
	Frequency        string                `json:"frequency,omitempty"`
	Cbsa             string                `json:"cbsa,omitempty"`
	StateFips        string                `json:"stateFips,omitempty"`
	CountyFips       string                `json:"countyFips,omitempty"`
}

type ScrapeBuilderContents struct {
	TextContent string `json:"textContent"`
	InnerHTML   string `json:"innerHTML"`
	InnerText   string `json:"innerText"`
	ClassName   string `json:"className"`
	TagName     string `json:"tagName"`
	URL         string `json:"url"`
}

// stringer interface for ScrapeBuilderContents
func (s ScrapeBuilderContents) String() string {
	if s.TagName == "" {
		s.TagName = "div."
	}
	return s.TagName + "." + strings.ReplaceAll(s.ClassName, " ", ".")
}

type Builder struct {
	UserID     string `json:"user_id,omitempty"`
	BuilderMap string `json:"builder_map"`
}
