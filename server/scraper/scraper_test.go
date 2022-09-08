package scraper

import (
	"testing"

	"github.com/broem/live-music-archiver-api/server/repo"
)

func TestCamelScraper(t *testing.T) {
	s := NewScraper()
	evts := s.ScrapeEvent(&repo.EventMapper{
		FullEventSelector: "DIV.col-12.eventWrapper.rhpSingleEvent.py-4.px-0",
		TitleSelector: "div.col-12.px-0.eventTitleDiv",
		DateSelector: "div.eventDateDetails.mt-md-0.mb-md-2",
		TicketCostSelector: "div.col-12.d-inline-block.eventsColor.eventCost.pt-md-0.mt-md-0.mb-md-2",
		VenueNameSelector: "div.col-12.eventsVenueDiv",
		ImagesSelector: []string{"div.rhp-event-thumb"},
		VenueBaseURL:      "https://thecamel.org/events/",
		DescriptionURLSelector: "div.col-12.mt-2.text-center.eventMoreInfo",
		DescriptionSelector: "div.col-sm-12.px-0.singleEventDescription.emptyDesc",
	}, false)

	if len(evts) == 0 {
		t.Errorf("Expected at least 1 event, got 0")
	}

}
