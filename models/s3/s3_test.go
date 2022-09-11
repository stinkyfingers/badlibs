package libs

import (
	"bytes"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/require"

	libs "github.com/stinkyfingers/badlibs/models"
)

type mockS3Client struct {
	body io.ReadCloser
	s3iface.S3API
}

func (s *mockS3Client) GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return &s3.GetObjectOutput{Body: s.body}, nil
}

func TestAll(t *testing.T) {
	expected := []libs.Lib{{
		Title: "test",
		ID:    "one",
		User:  "123",
	}, {
		Title: "test",
		ID:    "two",
		User:  "456",
	}}
	data := `{"one":{"_id":"one","title":"test","user":"123"},"two":{"_id":"two","title":"test","user":"456"}}`
	mock := &mockS3Client{
		body: io.NopCloser(bytes.NewBuffer([]byte(data))),
	}
	s := S3Storage{
		client: mock,
	}
	tests := []struct {
		description string
		expected    []libs.Lib
		filter      *libs.Lib
	}{
		{
			description: "all",
			expected:    expected,
		},
		{
			description: "user 123",
			expected:    expected[0:1],
			filter:      &libs.Lib{User: "123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ls, err := s.All(tt.filter)
			require.Nil(t, err)
			require.ElementsMatch(t, tt.expected, ls)
		})
	}
}
