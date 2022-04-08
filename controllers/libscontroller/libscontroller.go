package libscontroller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	libs "github.com/stinkyfingers/badlibs/models"

)

type Server struct {
	Storage libs.LibStorer
}

func NewServer(storage libs.LibStorer) *Server{
	return &Server{
		Storage: storage,
	}
}

func (s *Server) GetLib(w http.ResponseWriter, r *http.Request) {
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
	var l *libs.Lib
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, l)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	l, err = s.Storage.Create(l)
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

func (s *Server) DeleteLib(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := s.Storage.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	return
}

func (s *Server) UpdateLib(w http.ResponseWriter, r *http.Request) {
	var l *libs.Lib
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, l)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	l, err = s.Storage.Update(l)
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

func (s *Server) AllLibs(w http.ResponseWriter, r *http.Request) {
	ls, err := s.Storage.All()
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
