package main

import (
	"errors"
	"mime"
	"net/http"
	"reflect"
)

// * CRUD :

func (bp *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {
	// TREBA DA HANDLE-UJE KADA JE POSLATA JSON KONFIGURACIJA (1)
	//  ILI KONFIGURACIONA GRUPA (VISE KONFIGURACIJA)

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
	println(reflect.TypeOf(req.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	// array := []*Config{}
	// array = append(array, cf)
	println(*cf)
	println(cf)
	// if len(*cf) > 1 {

	// }
	bp.Data[id] = *cf
	renderJSON(w, cf)
}

// // VRACA Config sve
// func (ts *Service) getAllConfig(w http.ResponseWriter, req *http.Request) {
// 	// AKO TEST NE VALJA, VRATITI SE NA allConf LINIJU
// 	//  vvv
// 	allConf := []*Config{}
// 	for _, v := range ts.data {
// 		allConf = append(allConf, v)
// 	}

// 	renderJSON(w, allConf)
// }

// // VRACA 1 Config PO "id" (ili Error)
// func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	task, ok := ts.data[id]

// 	if !ok {
// 		err := errors.New("key not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	renderJSON(w, task)
// }

// // BRISE Config ili GRUPU Config-a, sta god se pronadje pod zadatim "ID"
// func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]

// 	if value, ok := ts.data[id]; ok {

// 		delete(ts.data, id)
// 		renderJSON(w, value)
// 	} else {
// 		err := errors.New("key not found")
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 	}
// }

// //  ! ! ! *TODO - Dodati metodu u main kao handler za POST "/addConfigToExistingGroupHandler["id"]
// func (ts *Service) addConfigToExistingGroupHandler(w http.ResponseWriter, req *http.Request) {
// 	// UBACUJE POSLATI JSON Config u POSTOJECU Config GRUPU ZA KOJU JE POSLAT "ID" KROZ URL
// }
