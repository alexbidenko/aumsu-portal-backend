package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, data interface{}) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")

	w.Write(jsonString)
}

func ParseRequestBody(w http.ResponseWriter, r *http.Request, data interface{}) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
