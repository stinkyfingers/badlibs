package libscontroller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/stinkyfingers/badlibs/models/libs"
)

func GetLib(w http.ResponseWriter, r *http.Request) {
	var l libs.Lib
	l.ID = r.URL.Query().Get("id")

	err := l.Get()
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

func CreateLib(w http.ResponseWriter, r *http.Request) {
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
	ti := time.Now()
	l.Created = &ti

	err = l.Create()
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

func DeleteLib(w http.ResponseWriter, r *http.Request) {
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

	err = l.Delete()
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

func UpdateLib(w http.ResponseWriter, r *http.Request) {
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

	err = l.Update()
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

func FindLib(w http.ResponseWriter, r *http.Request) {
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

	ls, err := l.Find()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(ls)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}
