package main

import (
	"errors"

	"mime"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/gorilla/mux"
)

type Service struct {
	cf *ConfigStore
}

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

	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("Version is required")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cf, err := decodeBodyConfigs(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, rid, err := bp.cf.Post(cf, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, rid)
}

func (ts *Service) getAllConfig(w http.ResponseWriter, req *http.Request) {
	// allConf := []*Configs{}
	// for _, v := range ts.Data {
	// 	for i := 0; i < len(v); i++ {
	// 		allConf = append(allConf, v[i])
	// 	}
	// }

	// renderJSON(w, allConf)

	allConfigs, err := ts.cf.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allConfigs)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	// id := mux.Vars(req)["id"]
	// task, ok := ts.Data[id]

	// if !ok {
	// 	err := errors.New("key not found")
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }
	// renderJSON(w, task)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("version not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	configs, err := ts.cf.Get(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, configs)
}

func (ts *Service) createNewVersionHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	if len(version) == 0 {
		err := errors.New("can not create new config without version")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	configs, err := ts.cf.Get(id, version)

	if configs != nil {
		err := errors.New("version already exist")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

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

	cf, err := decodeBodyConfigs(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, rid, err := ts.cf.PostNewVersion(cf, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, rid)
}

// func (ts *Service) getFilteredConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	version := mux.Vars(req)["version"]

// 	if len(version) == 0 {
// 		err := errors.New("version not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	task, ok := ts.cf.Get(id)
// 	result, error := task.Configs[version]
// 	if !error {
// 		err := errors.New("version not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	labels := mux.Vars(req)["labels"]
// 	labelMap := map[string]string{}
// 	s := strings.Split(labels, ";")
// 	for _, row := range s {
// 		rosParse := strings.Split(row, ":")
// 		labelMap[rosParse[0]] = rosParse[1]
// 	}

// 	var newTask []*Config

// 	for i := 0; i < len(result); i++ {
// 		entries := result[i].Entries
// 		if len(labelMap) == len(result[i].Entries) {
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
// 				newTask = append(newTask, result[i])
// 			}
// 		}
// 	}

// 	if ok != nil {
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
	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("version not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	msg, err := ts.cf.Delete(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, msg)
}

func (ts *Service) addConfigToExistingGroupHandler(w http.ResponseWriter, req *http.Request) {
	// contentType := req.Header.Get("Content-Type")
	// mediatype, _, err := mime.ParseMediaType(contentType)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// if mediatype != "application/json" {
	// 	err := errors.New("Expect application/json Content-Type")
	// 	http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	// 	return
	// }
	// cf, err := decodeBodyConfig(req.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// id := mux.Vars(req)["id"]
	// version := mux.Vars(req)["version"]

	// if len(version) == 0 {
	// 	err := errors.New("version not found")
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }
	// task, ok := ts.cf.Get(id, version)
	// // result, error := task.Configs[version]

	// // if !error {
	// // 	err := errors.New("version not found")
	// // 	http.Error(w, err.Error(), http.StatusNotFound)
	// // 	return
	// // }

	// if ok != nil {
	// 	err := errors.New("key not found")
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	// for _, c := range cf {
	// 	task.Configs = append(task.Configs, c)
	// }

	// ts.cf.Delete(id, version)
	// // delete(task.Configs, version)
	// // task.Configs[version] = result
	// _, newId, _ := ts.cf.Post(task, version)

	// renderJSON(w, newId)

	// // renderJSON(w, task)

}
