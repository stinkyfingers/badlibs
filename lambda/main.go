package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stinkyfingers/badlibs/server"
	"github.com/stinkyfingers/lambdify"
)

func main() {
	lambda.Start(lambdify.Lambdify(server.NewMux()))
}
