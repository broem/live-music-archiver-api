package runner

import (
	"fmt"
	"os"
	"time"

	"github.com/broem/live-music-archiver-api/server/igscraper"
	"github.com/broem/live-music-archiver-api/server/repo"
)

type IgRunner struct {
	Repo *repo.Repository
	Scraper *igscraper.Scraper
}

func NewIgRunner(r *repo.Repository, s *igscraper.Scraper) *IgRunner {
	return &IgRunner{
		Repo: r,
		Scraper: s,
	}
}

func (r *IgRunner) Run() {
	for {
		time.Sleep(time.Second * 60)
		runners, _ := r.Repo.GetIgRunners()
		for _, runner := range runners {
			// if the last run time is greater than the chron, run the scraper
			if runner.LastRun.After(time.Now().Add(time.Hour * time.Duration(runner.Chron))) {
				// pull the map from the repo
				mapRun, _ := r.Repo.GetIgMapperByID(runner.MapID.String())
				// // run the scraper
				if !r.Scraper.IsLoggedOn() {
					// pull login from env
					user := os.Getenv("IG_USER")
					pass := os.Getenv("IG_PASS")

					r.Scraper.DoLogin(user, pass)
				}
			
				userInfo, err := r.Scraper.GetUserInfo(mapRun.IgUserName)
				if err != nil {
					fmt.Printf("Error getting user info: %v", err)
				}

				following := r.Scraper.GetFollowing(userInfo)
			
				for _, s := range following {
					captured := r.Scraper.GetPosts(s)
					for _, c := range captured {
						// create igCaptured
						igCaptured := &repo.IgCaptured{
							MapID: runner.MapID,
							UserID: runner.UserID,
							IgUsername: s,
							CaptureDate: time.Now().UTC(),
							RawScrapedPayload: c,
						}

						r.Repo.SaveIgCaptured(igCaptured)

				}
				// // update the last run time
				runner.LastRun = time.Now()
				r.Repo.UpsertIgRunner(runner)
			}
		}
	}
}
}
