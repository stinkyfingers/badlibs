package server

import (
	"encoding/json"
	"fmt"
	"github.com/stinkyfingers/badlibs/auth"
	"log"
	"net/http"
	"os"

	"github.com/stinkyfingers/badlibs/controllers/libscontroller"
	libs "github.com/stinkyfingers/badlibs/models"
	filelibs "github.com/stinkyfingers/badlibs/models/file"
	s3libs "github.com/stinkyfingers/badlibs/models/s3"
)

// NewMux returns the router
func NewMux() (http.Handler, error) {
	authentication, err := getAuthentication("GCP")
	if err != nil {
	    return nil, err
	}

	storage, err := getStorage()
	if err != nil {
		return nil, err
	}
	s := libscontroller.NewServer(storage)

	mux := http.NewServeMux()
	mux.Handle("/lib/create", cors(authentication.Middleware(s.CreateLib)))
	mux.Handle("/lib/update", cors(authentication.Middleware(s.UpdateLib)))
	mux.Handle("/lib/delete", cors(authentication.Middleware(s.DeleteLib)))
	mux.Handle("/lib/get", cors(s.GetLib))
	mux.Handle("/lib/all", cors(s.AllLibs))
	mux.Handle("/health", cors(status))
	return mux, nil
}

func isPermittedOrigin(origin string) string {
	var permittedOrigins = []string{
		"http://localhost:3000",
		"https://radlibs.john-shenk.com",
	}
	for _, permittedOrigin := range permittedOrigins {
		if permittedOrigin == origin {
			return origin
		}
	}
	return ""
}

func cors(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permittedOrigin := isPermittedOrigin(r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", permittedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next := http.HandlerFunc(handler)
		next.ServeHTTP(w, r)
	})
}

func status(w http.ResponseWriter, r *http.Request) {
	log.Print("HEALTH")
	status := struct {
		Health string `json:"health"`
	}{
		"healthy",
	}
	j, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(j)
}

func getStorage() (libs.LibStorer, error) {
	switch os.Getenv("STORAGE") {
	case "file":
		filename := "foo.json"
		if envfile := os.Getenv("FILE"); envfile != "" {
			filename = envfile
		}
		return filelibs.NewFileStorage(filename)
	case "s3":
		return s3libs.NewS3Storage(os.Getenv("PROFILE"))
	default:
		return nil, fmt.Errorf("specify env vars for STORAGE (and PROFILE, FILE")
	}
}

func getAuthentication(kind string) (auth.Auth, error) {
	switch kind {
	case "GCP":
		return &auth.GCP{}, nil
	default:
		return nil, fmt.Errorf("%s has not been implemented", kind)
	}
}