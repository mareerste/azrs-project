package main

import (
	"errors"
	"fmt"
	"mime"
	"net/http"

	// "sort"
	// "strings"

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
	allConf := []*Configs{}
	for _, v := range ts.Data {
		for i := 0; i < len(v); i++ {
			allConf = append(allConf, v[i])
		}
	}

	renderJSON(w, allConf)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	err := errors.New("version not found")
	http.Error(w, err.Error(), http.StatusNotFound)
	return
}

// func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	task, ok := ts.Data[id]

// 	if !ok {
// 		err := errors.New("key not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	renderJSON(w, task)
// }

func (ts *Service) getConfigHandlerVersion(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	fmt.Println(version)
	if len(version) == 0 {
		err := errors.New("version not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, ok := ts.Data[id]
	result, error := task[0].Configs[version]

	if !error {
		err := errors.New("version not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, result)
}

// func (ts *Service) getFilteredConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	labels := mux.Vars(req)["labels"]
// 	labelMap := map[string]string{}
// 	s := strings.Split(labels, ";")
// 	for _, row := range s {
// 		rosParse := strings.Split(row, ":")
// 		labelMap[rosParse[0]] = rosParse[1]
// 	}

// 	task, ok := ts.Data[id]
// 	var newTask []*Configs

// 	for i := 0; i < len(task); i++ {
// 		entries := task[i].Entries
// 		if len(labelMap) == len(task[i].Entries) {
// 			check := false
// 			keys := make([]string, 0, len(entries))
// 			for k := range entries {
// 				keys = append(keys, k)
// 			}

// 			sort.Strings(keys)
// 			for _, k := range keys {
// 				i, ok := labelMap[k]
// 				if ok == false {
// 					check = true
// 					break
// 				} else {
// 					if i != entries[k] {
// 						check = true
// 						break
// 					}
// 				}
// 			}
// 			if check != true {
// 				newTask = append(newTask, task[i])
// 			}
// 		}
// 	}

// 	if !ok {
// 		err := errors.New("key not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	} else if len(newTask) == 0 {
// 		err := errors.New("params not match")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	renderJSON(w, newTask)
// }

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
		task = append(task, c)
	}
	ts.Data[id] = task
	renderJSON(w, task)

}

func (ts *Service) addNewVersionToConfigsHandler(w http.ResponseWriter, req *http.Request) {
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
		task = append(task, c)
	}
	ts.Data[id] = task
	renderJSON(w, task)

}
