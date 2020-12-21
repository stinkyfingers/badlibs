package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stinkyfingers/badlibs/server"
)

func main() {
	flag.Parse()
	fmt.Print("Running. \n")

	rh := server.NewMux()

	//openshift env var
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	if bind == ":" {
		bind = ":8081"
	}
	err := http.ListenAndServe(bind, rh)
	if err != nil {
		log.Print("BIND ", bind)
		log.Print(err)
	}

}
