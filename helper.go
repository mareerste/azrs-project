package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// *RequestPost
func decodeBody(r io.Reader) ([]*Configs, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cf []*Configs

	if err := dec.Decode(&cf); err != nil {
		return nil, err
	}
	return cf, nil
}

func decodeBodyConfig(r io.Reader) ([]*Config, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cf []*Configs

	if err := dec.Decode(&cf); err != nil {
		return nil, err
	}
	return cf, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}
