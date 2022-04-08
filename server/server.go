package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/stinkyfingers/badlibs/controllers/libscontroller"
	s3libs "github.com/stinkyfingers/badlibs/models/s3"
)

// NewMux returns the router
func NewMux() (http.Handler, error) {
	s3Storage, err := s3libs.NewS3Storage(os.Getenv("PROFILE"))
	if err != nil {
	    return nil, err
	}
	s := libscontroller.NewServer(s3Storage)

	mux := http.NewServeMux()
	mux.Handle("/lib/create", middleware(s.CreateLib))
	mux.Handle("/lib/update", middleware(s.UpdateLib))
	mux.Handle("/lib/delete", middleware(s.DeleteLib))
	mux.Handle("/lib/get", middleware(s.GetLib))
	mux.Handle("/lib/all", middleware(s.AllLibs))
	mux.Handle("/health", middleware(status))
	return mux, nil
}

func middleware(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
