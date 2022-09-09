package auth

import (
	"context"
	"net/http"
	"time"

	libs "github.com/stinkyfingers/badlibs/models"
)

type Internal struct {
	ID      string
	Token   string
	Expires time.Time
	Storage libs.LibStorer
}

func (i *Internal) Authorize(ctx context.Context, token string) error {
	auth, err := i.Storage.GetAuth(token)
	if err != nil {
		return err
	}
	if time.Now().After(auth.Expires) {
		return ErrTokenExpired
	}
	i.Token = auth.OIDCToken
	i.ID = auth.User.ID
	return err
}

func (i *Internal) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message":"missing token"}`))
			return
		}
		if err := i.Authorize(r.Context(), token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message":"unauthorized"}`))
			return
		}
		ctx := context.WithValue(r.Context(), "userId", i.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
