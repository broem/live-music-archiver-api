package repo

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	tableName   struct{}  `pg:"users"`
	UserID      string    `json:"user_id,omitempty"`
	Email       string    `json:"email,omitempty"`
	InstallDate time.Time `json:"install_date,omitempty"`
}

type Token struct {
	UserID              string `json:"user_id,omitempty"`
	Name                string `json:"name,omitempty"`
	Email               string `json:"email,omitempty"`
	*jwt.StandardClaims `json:"standard_claims,omitempty"`
}
