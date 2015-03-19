package libscontroller

import (
	"encoding/json"
	"github.com/stinkyfingers/badlibs/models/libs"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

func GetLib(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var l libs.Lib
	if ok := bson.IsObjectIdHex(id); !ok {
		http.Error(w, "Not valid ID.", 400)
		return
	}
	l.ID = bson.ObjectIdHex(id)
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
