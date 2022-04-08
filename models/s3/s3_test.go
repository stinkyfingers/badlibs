package libs

import (
	"bytes"
	"testing"
	"io"

	"github.com/stretchr/testify/require"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3"

	libs "github.com/stinkyfingers/badlibs/models"
)

type mockS3Client struct {
	body io.ReadCloser
	s3iface.S3API
}

func(s *mockS3Client) GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{Body: s.body}, nil
}

func TestAll(t *testing.T) {
	expected := []libs.Lib{{
		Title: "test",
		ID: "one",
	}}
	data := `{"one":{"_id":"one","title":"test"}}`
	mock := &mockS3Client{
		body: io.NopCloser(bytes.NewBuffer([]byte(data))),
	}
	s := S3Storage{
		client: mock,
	}
	ls, err := s.All()
	require.Nil(t, err)
	require.Equal(t, expected, ls)
}