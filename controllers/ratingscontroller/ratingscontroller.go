package ratingscontroller

import (
	"encoding/json"
	"github.com/stinkyfingers/badlibs/models/ratings"
	"io/ioutil"
	"net/http"
)

func FindRatings(w http.ResponseWriter, r *http.Request) {
	var rate ratings.Rating
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.Unmarshal(requestBody, &rate)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	rs, err := rate.FindMatch()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	j, err := json.Marshal(rs)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}
