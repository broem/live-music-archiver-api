package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/broem/live-music-archiver-api/server/repo"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

type Scraper struct {
	collector *colly.Collector
}

func NewScraper() *Scraper {
	coll := colly.NewCollector(colly.AllowURLRevisit())
	return &Scraper{
		collector: coll,
	}
}

func (s *Scraper) ScrapeEvent(eventMap *repo.EventMapper, single bool) []*repo.Event {
	events := []*repo.Event{}
	s.collector.OnScraped(func(r *colly.Response) {
		fmt.Println("Scraped", r.Request.URL)
	})

	// if the full event selector doesnt pull back anything try modifying it a bit.
	for {
		events = []*repo.Event{}
		// clear the collector and visit the venue event page
		s.collector.OnHTMLDetach(eventMap.FullEventSelector)
		// remove last . separated part of the selector
		s.collector.OnHTML(eventMap.FullEventSelector, func(e *colly.HTMLElement) {
			evt := &repo.Event{}
			evt.ID = uuid.New().String()
			evt.MapID = eventMap.MapID
			evt.UserID = eventMap.UserID
			evt.CaptureDate = time.Now().UTC()
			evt.Cbsa = eventMap.Cbsa
			evt.StateFips = eventMap.StateFips
			evt.CountyFips = eventMap.CountyFips
			mapEvents(e, eventMap, evt)
			events = append(events, evt)
		})

		if err := s.collector.Visit(parseWebsite(eventMap.VenueBaseURL, true)); err != nil {
			// we dont need to print the error here because the collector will print it
			// fmt.Println(err)
		}

		if len(events) <= 1 {
			s.collector.OnHTMLDetach(eventMap.FullEventSelector)
			splits := strings.Split(eventMap.FullEventSelector, ".")
			if len(splits) > 2 {
				eventMap.FullEventSelector = strings.Join(splits[:len(splits)-2], ".")
				continue
			} else {
				break
			}
		}

		break
	}

	// visit each event description page
	for _, evt := range events {
		s.collector.OnHTML("html", func(e *colly.HTMLElement) {
			mapEvents(e, eventMap, evt)
		})

		if err := s.collector.Visit(parseWebsite(evt.DescriptionURL, true)); err != nil {
			fmt.Println(err)
		}
	}
	s.collector.Wait()

	return events
}

// VisitCollect is a wrapper around colly.Collector.Visit that allows for a selector to be used to scrape a collection of events.
func (s *Scraper) VisitCollect(eventMap *repo.EventMapper, selector string) []*repo.Event {
	s.collector.OnHTMLDetach(selector)
	events := []*repo.Event{}

	s.collector.OnHTML(selector, func(e *colly.HTMLElement) {
		evt := &repo.Event{}
		evt.ID = uuid.New().String()
		evt.MapID = eventMap.MapID
		evt.UserID = eventMap.UserID
		mapEvents(e, eventMap, evt)
		events = append(events, evt)
	})

	if err := s.collector.Visit(parseWebsite(eventMap.VenueBaseURL, true)); err != nil {
		fmt.Println(err)
	}
	s.collector.Wait()

	return events
}


func mapEvents(e *colly.HTMLElement, mapper *repo.EventMapper, evt *repo.Event) {
	// TODO: might be worth looking into.
	// e.DOM.Find(mapper.TitleSelector).Each(func(i int, s *goquery.Selection) {
	// 	evt.Name = s.Find(mapper.EventNameSelector).Text()
	// 	evt.Description = s.Find(mapper.EventDescriptionSelector).Text()
	// 	evt.DescriptionURL = s.Find(mapper.EventDescriptionURLSelector).AttrOr("href", "")
	// 	evt.Date = s.Find(mapper.EventDateSelector).Text()
	// 	evt.Time = s.Find(mapper.EventTimeSelector).Text()
	// })

	// get the title
	e.ForEach(mapper.TitleSelector, func(_ int, el *colly.HTMLElement) {
		evt.Title = strings.TrimSpace(el.Text)
	})

	// get the date
	e.ForEach(mapper.DateSelector, func(_ int, el *colly.HTMLElement) {
		evt.Date = strings.TrimSpace(el.Text)
	})

	// get the time
	e.ForEach(mapper.TimeSelector, func(_ int, el *colly.HTMLElement) {
		evt.Time = strings.TrimSpace(el.Text)
	})

	// get the ticket cost
	e.ForEach(mapper.TicketCostSelector, func(_ int, el *colly.HTMLElement) {
		evt.TicketCost = strings.TrimSpace(el.Text)
	})

	// get the venue
	e.ForEach(mapper.VenueNameSelector, func(_ int, el *colly.HTMLElement) {
		evt.Venue = strings.TrimSpace(el.Text)
	})

	// get the images
	for _, imgSelector := range mapper.ImagesSelector {
		e.ForEach(imgSelector, func(_ int, el *colly.HTMLElement) {
			evt.Images = append(evt.Images, el.Attr("src"))
		})
	}

	// get the event description URL
	e.ForEach(mapper.DescriptionURLSelector, func(_ int, el *colly.HTMLElement) {
		evt.DescriptionURL = el.Attr("href")
		// evt.URL = evt.DescriptionURL
	})

	// get the event description
	e.ForEach(mapper.DescriptionSelector, func(_ int, el *colly.HTMLElement) {
		evt.Description = strings.TrimSpace(el.Text)
	})

	fmt.Printf("%+v\n", e)
}

func parseWebsite(url string, secure bool) string {
	url = trimProtocol(url)

	if secure {
		return "https://" + url
	}

	return "http://" + url
}
