package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/badlibs/controllers/libscontroller"
	"github.com/stinkyfingers/badlibs/controllers/ratingscontroller"
)

// NewMux returns the router
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/lib/create", middleware(libscontroller.CreateLib))
	mux.Handle("/lib/update", middleware(libscontroller.UpdateLib))
	mux.Handle("/lib/delete", middleware(libscontroller.DeleteLib))
	mux.Handle("/lib/get", middleware(libscontroller.GetLib))
	mux.Handle("/lib/find", middleware(libscontroller.FindLib))
	mux.Handle("/ratings/find", middleware(ratingscontroller.FindRatings))
	mux.Handle("/health", middleware(status))
	return mux
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
