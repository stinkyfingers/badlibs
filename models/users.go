package libs

import (
	"time"
)

type Auth struct {
	User      User      `json:"user"`
	OIDCToken string    `json:"token"`
	Expires   time.Time `json:"expires"`
}

type AuthStorer interface {
	GetAuth(token string) (*Auth, error)
	UpsertAuth(a *Auth) (*Auth, error)
}
