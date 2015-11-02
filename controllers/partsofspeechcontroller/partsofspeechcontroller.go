package partsofspeechcontroller

import (
	"encoding/json"
	"github.com/stinkyfingers/badlibs/models/partsofspeech"
	"io/ioutil"
	"net/http"
)

func FindPartsOfSpeech(w http.ResponseWriter, r *http.Request) {
	var p partsofspeech.PartOfSpeech
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, &p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	ps, err := p.FindMatch()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(ps)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}

func CreatePartOfSpeech(w http.ResponseWriter, r *http.Request) {
	var p partsofspeech.PartOfSpeech
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, &p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = p.Create()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}
