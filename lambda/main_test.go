package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestLambdaFunction(t *testing.T) {
	ev := events.ALBTargetGroupRequest{
		HTTPMethod:      "POST",
		Path:            "/badlibs/lib/find",
		Headers:         map[string]string{"Content-Type": "application/json"},
		IsBase64Encoded: false,
		Body:            "{}",
	}
	resp, err := lambdaFunction(ev)
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}
