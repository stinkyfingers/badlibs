package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/stinkyfingers/badlibs/controllers/libscontroller"
	"github.com/stinkyfingers/badlibs/controllers/partsofspeechcontroller"
	"github.com/stinkyfingers/badlibs/controllers/ratingscontroller"
)

// var (
// 	port = flag.String("port", ":8080", "Port to run on")
// )

func main() {
	flag.Parse()
	fmt.Print("Running. \n")

	//API
	rh.AddRoute(regexp.MustCompile("/lib/create"), middleware(http.HandlerFunc(libscontroller.CreateLib)))
	rh.AddRoute(regexp.MustCompile("/lib/update"), middleware(http.HandlerFunc(libscontroller.UpdateLib)))
	rh.AddRoute(regexp.MustCompile("/lib/delete"), middleware(http.HandlerFunc(libscontroller.DeleteLib)))
	rh.AddRoute(regexp.MustCompile("/lib/get"), middleware(http.HandlerFunc(libscontroller.GetLib)))
	rh.AddRoute(regexp.MustCompile("/lib/find"), middleware(http.HandlerFunc(libscontroller.FindLib)))
	rh.AddRoute(regexp.MustCompile("/ratings/find"), middleware(http.HandlerFunc(ratingscontroller.FindRatings)))
	rh.AddRoute(regexp.MustCompile("/partsofspeech/find"), middleware(http.HandlerFunc(partsofspeechcontroller.FindPartsOfSpeech)))
	rh.AddRoute(regexp.MustCompile("/partsofspeech/create"), middleware(http.HandlerFunc(partsofspeechcontroller.CreatePartOfSpeech)))

	//openshift env var
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	if bind == ":" {
		bind = ":8081"
	}
	err := http.ListenAndServe(bind, &rh)
	if err != nil {
		log.Print("BIND ", bind)
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

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
