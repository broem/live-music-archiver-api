package repo

import (
	"time"

	"github.com/google/uuid"
)

type Runner struct {
	MapID     uuid.UUID `json:"map_id" pg:"map_id,pk,type:uuid"`
	UserID    string    `json:"user_id" pg:"user_id,pk"`
	Chron     int       `json:"chron" pg:"chron"`
	LastRun   time.Time `json:"last_run" pg:"last_run"`
	Enabled   bool      `json:"enabled" pg:"enabled"`
}
