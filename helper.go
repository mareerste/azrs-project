package main

import (
	"azrs-project/tracer"
	"context"
	"encoding/json"
	"fmt"
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

	var cf []*Config

	if err := dec.Decode(&cf); err != nil {
		return nil, err
	}
	return cf, nil
}
func decodeBodyConfigs(ctx context.Context, r io.Reader) (*Configs, error) {
	span := tracer.StartSpanFromContext(ctx, "decodeBodyConfigs")
	defer span.Finish()

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cf *Configs

	if err := dec.Decode(&cf); err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	return cf, nil
}

func renderJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	span := tracer.StartSpanFromContext(ctx, "renderJSON")
	defer span.Finish()
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()  //unused
}

const (
	configs    = "configs/%s/%s"
	configsLab = "configs/%s/%s/%s"
	all        = "configs"
)

// ../config/safasfassafasfasf/version1
// ../config/safasfassafasfasf/version2

// "asffasfasfasf" => safasfas+version1

func generateKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(configs, id, version), id
}

func generateKeyNewVersion(id string, version string) (string, string) {
	return fmt.Sprintf(configs, id, version), id
}

func constructKey(id string, version string) string {
	return fmt.Sprintf(configs, id, version)
}




