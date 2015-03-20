package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/stinkyfingers/badlibs/controllers/application"
	"github.com/stinkyfingers/badlibs/controllers/libscontroller"
	"github.com/stinkyfingers/badlibs/controllers/partsofspeechcontroller"
	"github.com/stinkyfingers/badlibs/controllers/ratingscontroller"
)

var (
	port = flag.String("port", ":8080", "Port to run on")
)

func main() {
	flag.Parse()
	fmt.Print("Running. \n")

	//FILES
	rh.AddRoute(regexp.MustCompile("/public/js/"), http.StripPrefix("/public/js/", http.FileServer(http.Dir("public/js"))))
	rh.AddRoute(regexp.MustCompile("/public/templates/"), http.StripPrefix("/public/templates/", http.FileServer(http.Dir("public/templates"))))
	rh.AddRoute(regexp.MustCompile("/public/css/"), http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))

	//API
	rh.AddRoute(regexp.MustCompile("/lib/create"), http.HandlerFunc(libscontroller.CreateLib))
	rh.AddRoute(regexp.MustCompile("/lib/update"), http.HandlerFunc(libscontroller.UpdateLib))
	rh.AddRoute(regexp.MustCompile("/lib/delete"), http.HandlerFunc(libscontroller.DeleteLib))
	rh.AddRoute(regexp.MustCompile("/lib/get"), http.HandlerFunc(libscontroller.GetLib))
	rh.AddRoute(regexp.MustCompile("/lib/find"), http.HandlerFunc(libscontroller.FindLib))
	rh.AddRoute(regexp.MustCompile("/ratings/find"), http.HandlerFunc(ratingscontroller.FindRatings))
	rh.AddRoute(regexp.MustCompile("/partsofspeech/find"), http.HandlerFunc(partsofspeechcontroller.FindPartsOfSpeech))

	//APP
	rh.AddRoute(regexp.MustCompile("/.*"), http.HandlerFunc(application.Application))

	err := http.ListenAndServe(":"+os.Getenv("PORT"), &rh)
	if err != nil {
		log.Print(err)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request) string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fn(rw, r)
	}
}

var rh RegexpHandler

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}
type RegexpHandler struct {
	routes []*route
}

func (rh *RegexpHandler) AddRoute(pattern *regexp.Regexp, handler http.Handler) {
	ro := route{pattern: pattern, handler: handler}
	rh.routes = append(rh.routes, &ro)
}

func (rh *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rh.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
