package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"google.golang.org/api/idtoken"
)

const clientId = "520868981613-vpe0s1lild8cl62hblmg9bfl01fplu06.apps.googleusercontent.com" // TODO config

var ErrTokenExpired = errors.New("token is expired")
var ErrBadClientId = errors.New("unknown client id")
var ErrNoPayload = errors.New("no payload")
var ErrMalformedClaim = errors.New("malformed JWT claims")

type GCP struct {
	payload *idtoken.Payload
}

func (g *GCP) Authorize(ctx context.Context, token string) error {
	payload, err := idtoken.Validate(ctx, token, clientId)
	if err != nil {
		return err
	}
	if time.Unix(payload.Expires, 0).Before(time.Now()) {
		return ErrTokenExpired
	}
	if payload.Audience != clientId {
		return ErrBadClientId
	}
	g.payload = payload
	return nil
}

func (g *GCP) UserId() (string, error) {
	if g == nil || g.payload == nil {
		return "", ErrNoPayload
	}
	return g.payload.Subject, nil
}

func (g *GCP) UserName() (string, error) {
	if g == nil || g.payload == nil {
		return "", ErrNoPayload
	}
	if val, ok := g.payload.Claims["name"]; ok {
		return val.(string), nil
	}
	return "", ErrMalformedClaim
}

func (g *GCP) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message":"missing token"}`))
			return
		}
		if err := g.Authorize(r.Context(), token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"unauthorized"}`))
			return
		}
		userId, err := g.UserId()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"unauthorized"}`))
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

