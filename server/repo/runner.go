package repo

import (
	"time"

	"github.com/google/uuid"
)

type Runner struct {
	tableName struct{}  `pg:"runner"`
	MapID     uuid.UUID `json:"map_id"`
	UserID    string    `json:"user_id"`
	Chron     int       `json:"chron"`
	LastRun   time.Time `json:"last_run"`
	Enabled   bool      `json:"enabled"`
}
