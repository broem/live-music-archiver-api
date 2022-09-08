package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/broem/live-music-archiver-api/server/igscraper"
	"github.com/broem/live-music-archiver-api/server/repo"
	"github.com/broem/live-music-archiver-api/server/runner"
	"github.com/broem/live-music-archiver-api/server/scraper"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Api struct {
	r  *repo.Repository
	s  *igscraper.Scraper
	sc *scraper.Scraper
}

type Config struct {
	Address  string `json:"address" yaml:"address"`
	User     string `json:"user" yaml:"user"`
	Pass     string `json:"pass" yaml:"pass"`
	Database string `json:"database" yaml:"database"`
	Pool     int    `json:"pool" yaml:"pool"`
}

func NewApi(cfg *Config) {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://bleachedsolutions.com", "https://www.bleachedsolutions.com", "*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "../client/dist",
		Index:  "index.html",
		Browse: true,
		HTML5:  true,
	}))

	repo, err := repo.NewRepo(
		repo.Config{
			Address:  cfg.Address,
			User:     cfg.User,
			Pass:     cfg.Pass,
			Database: cfg.Database,
			Pool:     cfg.Pool,
		},
	)
	if err != nil {
		panic("no repo response")
	}

	// create database tables
	err = repo.CreateTables()
	if err != nil {
		panic("no repo response")
	}

	a := Api{
		r:  repo,
		s:  igscraper.NewScraper(),
		sc: scraper.NewScraper(),
	}

	// start runner
	go runner.NewRunner(repo).Run()

	// start IgRunner
	go runner.NewIgRunner(repo, a.s).Run()

	g := e.Group("/api")

	g.POST("/scrapeInstagram", a.scrapeInstagram)
	g.POST("/saveUrlSelection", a.saveUrlSelection, TokenMiddleware)
	g.POST("/logout", a.logout, TokenMiddleware)
	g.POST("/scrapeBuilder", a.scrapeBuilder)
	g.POST("/verified", a.verified)
	g.POST("/updateScrape", a.updateScrape)
	g.POST("/deleteScrape", a.deleteScrape)
	g.POST("/updateIGScrape", a.updateIGScrape)
	g.POST("/deleteIGScrape", a.deleteIGScrape)
	g.GET("/scrapeID/:id", a.ScrapeByID)
	g.GET("/myScrapes/:id", a.myScrapes)
	g.GET("/myIgScrapes/:id", a.myIgScrapes)
	g.GET("/getCurrentScrapeEvents/:id", a.getCurrentScrapeEvents)
	g.GET("/getCurrentIGScrapeEvents/:id", a.getCurrentIGScrapeEvents)

	e.Logger.Fatal(e.Start(":3424"))
}

// myScrapes returns all the scrapes for a user
func (a *Api) myScrapes(c echo.Context) error {
	id := c.Param("id")

	events, err := a.r.GetEventsByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = createEventsCSV(events)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.File("events.csv")
}

// myIgScrapes returns all the scrapes for a user
func (a *Api) myIgScrapes(c echo.Context) error {
	id := c.Param("id")

	events, err := a.r.GetIgCapturedByUserID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = createIGEventsTXT(events)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.File("igRawText.txt")
}

func (a *Api) deleteScrape(c echo.Context) error {
	fmt.Println("deleteScrape")
	p := new(CurrentRunners)
	if err := c.Bind(p); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	// delete runner
	err := a.r.DeleteRunner(p.UserID, p.MapID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// delete the map
	err = a.r.DeleteEventMapper(p.MapID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "deleted")
}

func (a *Api) updateScrape(c echo.Context) error {
	p := new(CurrentRunners)
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// make a runner
	runner := &repo.Runner{
		MapID:   uuid.MustParse(p.MapID),
		UserID:  p.UserID,
		Chron:   mapFrequency(p.Frequency),
		Enabled: p.Enabled,
	}

	// update runner
	err := a.r.UpsertRunner(runner)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "updated")
}

func createEventsCSV(events []*repo.Event) error {
	f, err := os.Create("events.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	for _, event := range events {
		// defer w.Flush()
		// event to csv string
		mslice := event.EventStringSlice()
		err = w.Write(mslice)
		if err != nil {
			return err
		}
	}

	return nil
}

// createIGEventsTXT creates a txt file with all the scraped IG data
func createIGEventsTXT(events []*repo.IgCaptured) error {
	f, err := os.Create("igRawText.txt")
	if err != nil {
		return err
	}
	defer f.Close()

	for _, event := range events {
		// event to txt string
		mslice := event.CapturedString()
		_, err = f.WriteString(mslice)
		if err != nil {
			return err
		}
	}

	return nil
}

// ScrapeByID scrapes an event by id
func (a *Api) ScrapeByID(c echo.Context) error {
	id := c.Param("id")
	event, err := a.r.GetEventMapperByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	events := a.sc.ScrapeEvent(event, false)

	return c.JSON(http.StatusOK, events)
}

type CurrentRunners struct {
	MapID     string `json:"mapID"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Enabled   bool   `json:"enabled"`
	UserID    string `json:"userID"`
}

type CurrentIGRunners struct {
	MapID     string `json:"mapID"`
	Profile   string `json:"profile"`
	Frequency string `json:"frequency"`
	Enabled   bool   `json:"enabled"`
	UserID    string `json:"userID"`
}

func (a *Api) getCurrentIGScrapeEvents(c echo.Context) error {
	id := c.Param("id")

	runners, err := a.r.GetIgRunnersByUserID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// get all map id from runners
	var mapIDs []string
	for _, r := range runners {
		mapIDs = append(mapIDs, r.MapID.String())
	}

	// get enabled maps
	maps, err := a.r.GetIGEventMappersByMapID(mapIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// create collection of current runners
	var currentRunners []CurrentIGRunners
	for _, m := range maps {
		// get the runner for the map from the runners
		var r *repo.IgRunner
		for _, runner := range runners {
			if runner.MapID.String() == m.MapID.String() {
				r = runner
			}
		}

		currentRunners = append(currentRunners, CurrentIGRunners{
			MapID:     m.MapID.String(),
			Profile:   m.IgUserName,
			Frequency: unmapFrequency(r.Chron),
			Enabled:   r.Enabled,
			UserID:    r.UserID,
		})
	}

	return c.JSON(http.StatusOK, currentRunners)
}

func (a *Api) updateIGScrape(c echo.Context) error {
	p := new(CurrentIGRunners)
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// make a runner
	runner := &repo.IgRunner{
		MapID:   uuid.MustParse(p.MapID),
		UserID:  p.UserID,
		Chron:   mapFrequency(p.Frequency),
		Enabled: p.Enabled,
	}

	// update runner
	err := a.r.UpsertIgRunner(runner)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "updated")
}

func (a *Api) deleteIGScrape(c echo.Context) error {
	p := new(CurrentRunners)
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// delete runner
	err := a.r.DeleteIgRunner(p.UserID, p.MapID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// delete the map
	err = a.r.DeleteIgMapper(p.MapID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "deleted")
}

func (a *Api) getCurrentScrapeEvents(c echo.Context) error {
	id := c.Param("id")

	/// for a given ID url, frequncy, enabled
	// need runners by id and enabled
	runners, err := a.r.GetRunnersByID(id)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// get all map id from runners
	var mapIDs []string
	for _, r := range runners {
		mapIDs = append(mapIDs, r.MapID.String())
	}

	// get enabled maps
	maps, err := a.r.GetEventMappersByMapID(mapIDs)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// create collection of current runners
	var currentRunners []CurrentRunners
	for _, m := range maps {
		// get the runner for the map from the runners
		var r *repo.Runner
		for _, runner := range runners {
			if runner.MapID.String() == m.MapID.String() {
				r = runner
			}
		}

		currentRunners = append(currentRunners, CurrentRunners{
			MapID:     m.MapID.String(),
			URL:       m.VenueBaseURL,
			Frequency: unmapFrequency(r.Chron),
			Enabled:   r.Enabled,
			UserID:    r.UserID,
		})
	}

	return c.JSON(http.StatusOK, currentRunners)
}

func (a *Api) scrapeBuilder(c echo.Context) error {
	fmt.Println("Building an event scrape")
	p := new(repo.ScrapeBuilder)
	if err := c.Bind(p); err != nil {
		fmt.Println(err)
		return err
	}

	// try to add user to database
	user := &repo.User{
		UserID:      p.UserID,
		Email:       p.UserEmail,
		InstallDate: time.Now().UTC(),
	}

	err := a.r.AddUser(user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = a.r.SaveScrapeBuilder(p)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// convert to EventMapper
	m := &repo.EventMapper{
		MapID:                    uuid.New(),
		VenueBaseURL:             p.Event.URL,
		FullEventSelector:        fmt.Sprint(p.Event),
		UserID:                   p.UserID,
		TitleSelector:            fmt.Sprint(p.EventTitle),
		DescriptionSelector:      fmt.Sprint(p.EventDesc),
		DescriptionURLSelector:   fmt.Sprint(p.EventDescURL),
		DateSelector:             fmt.Sprint(p.StartDate),
		TimeSelector:             fmt.Sprint(p.DoorTime), // for times we should bring in the current time with timezone
		VenueNameSelector:        fmt.Sprint(p.VenueName),
		VenueAddressSelector:     fmt.Sprint(p.VenueAddress),
		VenueContactInfoSelector: fmt.Sprint(p.VenueContactInfo),
		TicketCostSelector:       fmt.Sprint(p.TicketCost),
		TicketURLSelector:        fmt.Sprint(p.TicketURLS),
		OtherPerformersSelector:  fmt.Sprint(p.OtherPerformers),
		AgeRequiredSelector:      fmt.Sprint(p.AgeRequired),
		FacebookURLSelector:      fmt.Sprint(p.FacebookURL),
		TwitterURLSelector:       fmt.Sprint(p.TwitterURL),
		MiscSelector:             fmt.Sprint(p.Misc),
		Cbsa:                     p.Cbsa,
		StateFips:                p.StateFips,
		CountyFips:               p.CountyFips,
		ImagesSelector:           []string{fmt.Sprint(p.Images)},
	}

	// upsert and return the upserted map
	m, err = a.r.UpsertEventMapper(m)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Add runner
	runner := &repo.Runner{
		MapID:   m.MapID,
		UserID:  p.UserID,
		Chron:   mapFrequency(p.Frequency),
		LastRun: time.Now().UTC(),
	}
	err = a.r.UpsertRunner(runner)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// scrape the first event found
	events := a.sc.ScrapeEvent(m, true)
	if len(events) == 0 {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = a.r.SaveEvents(events)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	e := &repo.EventWithCount{
		Event:      events[0],
		EventCount: len(events),
	}

	fmt.Println("Event scrape complete")
	// we need to have user approve the first scraped event
	// send back first event to user for approval
	return c.JSON(http.StatusOK, e)
}

func (a *Api) verified(c echo.Context) error {
	fmt.Println("Verifying event")
	enabled := new(repo.Enabled)
	err := c.Bind(enabled)
	if err != nil {
		return err
	}

	if !enabled.Enabled {
		// remove the associated map, runner, and events
		err = a.r.DeleteEventMapper(enabled.MapID.String())
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		err = a.r.DeleteRunner(enabled.UserID, enabled.MapID.String())
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		err = a.r.DeleteEvents(enabled.MapID.String())
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		// update the runner to enabled
		err = a.r.UpdateRunner(enabled.UserID, enabled.MapID, true)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		// update mappper to enabled
		err = a.r.UpdateEventMapper(enabled.MapID, true)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

	}

	fmt.Printf("Event verified: %v", enabled.Enabled)

	return c.JSON(http.StatusOK, enabled.Enabled)
}

func mapFrequency(frequency string) int {
	switch frequency {
	case "Every Day":
		return 24
	case "Every Other Day":
		return 48
	case "Every Week":
		return 24 * 7
	case "Every Other Week":
		return 24 * 14
	case "Every Month":
		return 24 * 30
	}
	return 24
}

func unmapFrequency(frequency int) string {
	switch frequency {
	case 24:
		return "Every Day"
	case 48:
		return "Every Other Day"
	case 24 * 7:
		return "Every Week"
	case 24 * 14:
		return "Every Other Week"
	case 24 * 30:
		return "Every Month"
	}
	return "Every Day"
}

func (a *Api) scrapeInstagram(c echo.Context) error {
	s := new(repo.IgMapBuilder)
	if err := c.Bind(s); err != nil {
		return err
	}

	// create and save igmapper
	m := &repo.IgMapper{
		MapID:      uuid.New(),
		UserID:     s.UserID,
		UserEmail:  s.UserEmail,
		IgUserName: s.IgUserName,
	}

	err := a.r.SaveIgMapper(m)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// build runner
	runner := &repo.IgRunner{
		MapID:   m.MapID,
		UserID:  s.UserID,
		Chron:   mapFrequency(s.Frequency),
		Enabled: true,
	}

	// save runner
	err = a.r.UpsertIgRunner(runner)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, true)
}

func (a *Api) saveUrlSelection(c echo.Context) error {
	p := new(repo.Message)
	if err := c.Bind(p); err != nil {
		return err
	}
	log.Println(p)
	return c.JSON(http.StatusOK, "You saved the url")
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	log.Println(r)
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *Api) logout(c echo.Context) error {
	return c.JSON(http.StatusOK, "logged out")
}

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := TokenValid(c.Request())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}
		return next(c)
	}
}
