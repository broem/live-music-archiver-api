package runner

import (
	"time"

	"github.com/broem/live-music-archiver-api/server/repo"
	"github.com/broem/live-music-archiver-api/server/scraper"
)

type Runner struct {
	Repo *repo.Repository
}

func NewRunner(r *repo.Repository) *Runner {
	return &Runner{
		Repo: r,
	}
}

func (r *Runner) Run() {
	for {
		time.Sleep(time.Second * 60)
		runners, _ := r.Repo.GetRunners()
		for _, runner := range runners {
			// if the last run time is greater than the chron, run the scraper
			if runner.LastRun.After(time.Now().Add(time.Hour * time.Duration(runner.Chron))) {
				// pull the map from the repo
				mapRun, _ := r.Repo.GetEventMapperByID(runner.MapID.String())
				// // run the scraper
				scraper := scraper.NewScraper()
				scraper.ScrapeEvent(mapRun, false)
				// // update the last run time
				runner.LastRun = time.Now()
				r.Repo.UpsertRunner(runner)
			}
		}
	}
}
