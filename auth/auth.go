package auth

import (
	"context"
	"net/http"
)

type Auth interface {
	Authorize(ctx context.Context, token string) error
	UserId() (string, error)
	UserName() (string, error)
	Middleware(handlerFunc http.HandlerFunc) http.HandlerFunc
}
