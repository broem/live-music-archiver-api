package repo

import (
	"time"

	"github.com/google/uuid"
)

type IgMapBuilder struct {
	UserID     string `json:"user_id"`
	UserEmail  string `json:"user_email"`
	IgUserName string `json:"ig_user_name"`
	Frequency  string `json:"frequency"`
}

type IgMapper struct {
	tableName  struct{}  `pg:"ig_mappers"`
	MapID      uuid.UUID `json:"map_id" pg:",pk,type:uuid"`
	UserID     string    `json:"user_id"`
	UserEmail  string    `json:"user_email"`
	IgUserName string    `json:"ig_user_name"`
	Approved   bool      `json:"approved"`
}

type IgRunner struct {
	tableName struct{}  `pg:"ig_runners"`
	MapID     uuid.UUID `json:"map_id" pg:"map_id,pk,type:uuid"`
	UserID    string    `json:"user_id" pg:"user_id,pk"`
	Chron     int       `json:"chron"`
	LastRun   time.Time `json:"last_run"`
	Enabled   bool      `json:"enabled"`
}

type IgCaptured struct {
	tableName         struct{}  `pg:"ig_captured"`
	MapID             uuid.UUID `json:"map_id" pg:",pk,type:uuid"`
	UserID            string    `json:"user_id" pg:",pk"`
	CaptureDate       time.Time `json:"capture_date"`
	IgUsername        string    `json:"ig_username"`
	RawScrapedPayload string    `json:"raw_scraped_payload"`
}

// CapturedStringSlice returns a slice of strings from the RawScrapedPayload
func (i *IgCaptured) CapturedString() string {
	return "Follower: " + i.IgUsername + " Captured: " + i.CaptureDate.String() + " Payload: " + i.RawScrapedPayload + "\n"
}

