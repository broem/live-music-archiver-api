package repo

import (
	"time"

	"github.com/google/uuid"
)

type Runner struct {
	tableName struct{}  `pg:"runner"`
	MapID     uuid.UUID `json:"map_id" pg:"map_id,pk"`
	UserID    string    `json:"user_id" pg:"user_id,pk"`
	Chron     int       `json:"chron"`
	LastRun   time.Time `json:"last_run"`
	Enabled   bool      `json:"enabled"`
}
