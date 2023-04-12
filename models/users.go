package libs

import (
	"time"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Auth struct {
	User      User      `json:"user"`
	OIDCToken string    `json:"token"`
	Expires   time.Time `json:"expires"`
}

type AuthStorer interface {
	GetAuth(token string) (*Auth, error)
	UpsertAuth(a *Auth) (*Auth, error)
}
