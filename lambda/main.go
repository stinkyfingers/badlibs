package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stinkyfingers/badlibs/server"
)

func lambdaFunction(ev events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	log.Print("path", ev.Path)
	var queries string
	for k, v := range ev.QueryStringParameters {
		queries += fmt.Sprintf("&%s=%s", k, v)
	}
	queries = strings.Replace(queries, "&", "?", 1)

	path := strings.TrimPrefix(ev.Path, "/badlibs")

	req, err := http.NewRequest(ev.HTTPMethod, path+queries, strings.NewReader(ev.Body))
	if err != nil {
		return events.ALBTargetGroupResponse{
			StatusCode:        http.StatusInternalServerError,
			StatusDescription: "500 Internal Server Error",
			Body:              err.Error(),
			IsBase64Encoded:   false,
		}, nil
	}
	log.Print("REQ", req)

	rr := httptest.NewRecorder()
	server.NewMux().ServeHTTP(rr, req)
	headers := map[string]string{"Content-Type": "application/json"}
	multiValueHeaders := make(map[string][]string)
	for k, v := range rr.Result().Header {
		if len(v) > 1 {
			multiValueHeaders[k] = v
		} else if len(v) == 1 {
			headers[k] = v[0]
		}
	}
	log.Print(rr.Body.String(), rr)
	return events.ALBTargetGroupResponse{
		Body:              rr.Body.String(),
		IsBase64Encoded:   false,
		StatusCode:        http.StatusOK,
		StatusDescription: "200 OK",
		Headers:           headers,
		MultiValueHeaders: multiValueHeaders,
	}, nil

}

func main() {
	lambda.Start(lambdaFunction)
}
