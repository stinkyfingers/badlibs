package libscontroller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"context"

	libs "github.com/stinkyfingers/badlibs/models"
)

type Server struct {
	Storage libs.LibStorer
}

func NewServer(storage libs.LibStorer) *Server {
	return &Server{
		Storage: storage,
	}
}

func (s *Server) GetLib(w http.ResponseWriter, r *http.Request) {
	log.Print("GET")
	id := r.URL.Query().Get("id")
	l, err := s.Storage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func (s *Server) CreateLib(w http.ResponseWriter, r *http.Request) {
	var l libs.Lib
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, &l)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	newLib, err := s.Storage.Create(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(newLib)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func (s *Server) DeleteLib(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	l, err := s.Storage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if ok := authorize(r.Context(), *l); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = s.Storage.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
	return
}

func (s *Server) UpdateLib(w http.ResponseWriter, r *http.Request) {
	var l libs.Lib
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, &l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if ok := authorize(r.Context(), l); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	updatedLib, err := s.Storage.Update(&l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(updatedLib)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func (s *Server) AllLibs(w http.ResponseWriter, r *http.Request) {
	log.Print("ALL")
	filter, err := getFilter(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	ls, err := s.Storage.All(filter)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	j, err := json.Marshal(ls)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func getFilter(values url.Values) (*libs.Lib, error) {
	if len(values) == 0 {
		return nil, nil
	}
	var lib libs.Lib
	if values.Has("id") {
		lib.ID = values.Get("id")
	}
	if values.Has("title") {
		lib.Title = values.Get("title")
	}
	if values.Has("rating") {
		lib.Rating = values.Get("rating")
	}
	if values.Has("user") {
		lib.User.ID = values.Get("userId")
	}
	if values.Has("created") {
		createdStr := values.Get("created")
		created, err := time.Parse(createdStr, "2006-01-02")
		if err != nil {
			return nil, err
		}
		lib.Created = &created
	}
	return &lib, nil
}

func authorize(ctx context.Context, lib libs.Lib) bool {
	if ctx.Value("userId") != lib.User.ID {
		return false
	}
	return true
}