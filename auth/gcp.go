package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const clientId = "520868981613-vpe0s1lild8cl62hblmg9bfl01fplu06.apps.googleusercontent.com" // TODO config

var ErrTokenExpired = errors.New("token is expired")
var ErrBadClientId = errors.New("unknown client id")
var ErrNoPayload = errors.New("no payload")
var ErrMalformedClaim = errors.New("malformed JWT claims")

type GCP struct {
	token *Token
}

type Token struct {
	Name     string `json:"name"`
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Audience string `json:"aud"`
	Expires  string `json:"exp"`
	IssuedAt string `json:"iat"`
}

func (g *GCP) Authorize(ctx context.Context, token string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", token), nil)
	if err != nil {
		return err
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var tok Token
	err = json.NewDecoder(resp.Body).Decode(&tok)
	if err != nil {
		return err
	}
	issuedAt, err := strconv.Atoi(tok.IssuedAt)
	if err != nil {
		return err
	}
	if time.Now().Add(time.Hour*-24).UnixMilli() < int64(issuedAt) {
		return ErrTokenExpired
	}
	g.token = &tok
	return err
}

func (g *GCP) UserId() (string, error) {
	if g == nil || g.token == nil {
		return "", ErrNoPayload
	}
	return g.token.Subject, nil
}

func (g *GCP) UserName() (string, error) {
	if g == nil || g.token == nil {
		return "", ErrNoPayload
	}
	return g.token.Name, nil
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
