package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	//"github.com/gorilla/mux"
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

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := bp.cf.AddConfig(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, config.Id)
}

func (bp *Service) getAllConfig(w http.ResponseWriter, _ *http.Request) {

	allTasks, err := bp.cf.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}

func (bp *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	task, _ := bp.cf.GetConf(id) //task je configs, u configs imam mapu u kojoj imamo niz pokazivaca na configuraciju ili grupu

	if task == nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, task)
}

func (bp *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	//id := mux.Vars(req)["id"]
	//
	//if value, ok := ts.Data[id]; ok {
	//
	//	delete(ts.Data, id)
	//	renderJSON(w, value)
	//} else {
	//	err := errors.New("key not found")
	//	http.Error(w, err.Error(), http.StatusNotFound)
	//}
}

func (bp *Service) addConfigToExistingGroupHandler(w http.ResponseWriter, req *http.Request) {
	//contentType := req.Header.Get("Content-Type")
	//mediatype, _, err := mime.ParseMediaType(contentType)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	//if mediatype != "application/json" {
	//	err := errors.New("Expect application/json Content-Type")
	//	http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	//	return
	//}
	//cf, err := decodeBody(req.Body)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	//id := mux.Vars(req)["id"]
	//task, ok := ts.Data[id]
	//
	//if !ok {
	//	err := errors.New("key not found")
	//	http.Error(w, err.Error(), http.StatusNotFound)
	//	return
	//}
	//
	//for _, c := range cf {
	//	task = append(task, c)
	//}
	//ts.Data[id] = task
	//renderJSON(w, task)

}
