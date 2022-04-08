package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/stinkyfingers/badlibs/server"
)

const (
	port = ":8088"
)

func main() {
	flag.Parse()
	fmt.Print("Running. \n")
	rh, err := server.NewMux()
	if err != nil {
	    panic(err)
	}

	err = http.ListenAndServe(port, rh)
	if err != nil {
		log.Print(err)
	}

}
