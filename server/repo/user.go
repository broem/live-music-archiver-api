package repo

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	InstallDate time.Time `json:"install_date"`
}

type Token struct {
	UserID              string `json:"user_id,omitempty"`
	Name                string `json:"name,omitempty"`
	Email               string `json:"email,omitempty"`
	*jwt.StandardClaims `json:"standard_claims,omitempty"`
}
