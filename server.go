package main

import (
	"errors"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

func (bp *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cf, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	bp.Data[id] = cf
	renderJSON(w, id)
}

func (ts *Service) getAllConfig(w http.ResponseWriter, req *http.Request) {
	allConf := []*Config{}
	for _, v := range ts.Data {
		for i := 0; i < len(v); i++ {
			allConf = append(allConf, v[i])
		}
	}

	renderJSON(w, allConf)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, ok := ts.Data[id]

	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	if value, ok := ts.Data[id]; ok {

		delete(ts.Data, id)
		renderJSON(w, value)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *Service) addConfigToExistingGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cf, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(req)["id"]
	task, ok := ts.Data[id]

	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	for _, c := range cf {
		for i := 0; i < len(cf); i++ {
			task = append(task, c)
		}
	}
	ts.Data[id] = task
	renderJSON(w, task)

}
