package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGCPAuth(t *testing.T) {
	t.Skip("live test only")
	g := &GCP{}
	err := g.Authorize(context.Background(), "TODO")
	require.Nil(t, err)
}
