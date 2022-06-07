package main

import (
	"azrs-project/tracer"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	name = "config_service"
)

type Service struct {
	cf     *ConfigStore
	tracer opentracing.Tracer
	closer io.Closer
}

func NewService() (*Service, error) {
	cf, err := New()
	if err != nil {
		return nil, err
	}

	tracer, closer := tracer.Init(name)
	// fmt.Println("OVDE SAM")
	// fmt.Println(tracer)
	// fmt.Println(closer)
	opentracing.SetGlobalTracer(tracer)
	return &Service{
		cf:     cf,
		tracer: tracer,
		closer: closer,
	}, nil
}

func (s *Service) GetTracer() opentracing.Tracer {
	return s.tracer
}

func (s *Service) GetCloser() io.Closer {
	return s.closer
}

func (s *Service) CloseTracer() error {
	return s.closer.Close()
}

func (bp *Service) createConfigHandler(w http.ResponseWriter, req *http.Request) {

	// fmt.Println(bp)
	span := tracer.StartSpanFromRequest("createConfigHandler", bp.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling config creation at %s\n", req.URL.Path)))

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("Version is required")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ctx := tracer.ContextWithSpan(context.Background(), span)
	cf, err := decodeBodyConfigs(ctx, req.Body)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req_idempotency_key := req.Header.Get("x-idempotency-key")

	if len(req_idempotency_key) == 0 {
		err := errors.New("x-idempotency-key missing")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	idem, err := bp.cf.GetIdemKey(ctx, req_idempotency_key)

	if idem == nil && err == nil {

		_, rid, error := bp.cf.Post(ctx, cf, version)
		if error != nil {
			tracer.LogError(span, error)
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		idemNew := &Idem{rid, "created"}
		err = bp.cf.PostIdemKey(ctx, req_idempotency_key, idemNew)
		//
		if err != nil {
			err := errors.New("failed to insert idem key into database")
			tracer.LogError(span, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		renderJSON(ctx, w, rid)

	} else {
		// err != nil || idem != nil {
		err := errors.New("this request already exist")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (ts *Service) getAllConfig(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getAllConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling get all configs at %s\n", req.URL.Path)))

	// allConf := []*Configs{}
	// for _, v := range ts.Data {
	// 	for i := 0; i < len(v); i++ {
	// 		allConf = append(allConf, v[i])
	// 	}
	// }
	ctx := tracer.ContextWithSpan(context.Background(), span)

	// renderJSON( w, allConf)

	allConfigs, err := ts.cf.GetAll(ctx)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, allConfigs)
}

func (ts *Service) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling get config at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

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
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	configs, err := ts.cf.Get(ctx, id, version)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, configs)
}

func (ts *Service) createNewVersionHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling get config at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	if len(version) == 0 {
		err := errors.New("can not create new config without version")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	configs, err := ts.cf.Get(ctx, id, version)

	if configs != nil || err != nil {
		err := errors.New("version already exist")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cf, err := decodeBodyConfigs(ctx, req.Body)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, rid, err := ts.cf.PostNewVersion(ctx, cf, id, version)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, rid)
}

func (ts *Service) getFilteredConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getFilteredConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling get filtered config at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("version not found")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	task, ok := ts.cf.Get(ctx, id, version)
	if task == nil {
		err := errors.New("config not found")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	labels := mux.Vars(req)["labels"]
	labelMap := map[string]string{}
	s := strings.Split(labels, ";")
	for _, row := range s {
		rosParse := strings.Split(row, ":")
		labelMap[rosParse[0]] = rosParse[1]
	}

	var newTask []*Config

	for i := 0; i < len(task.Configs); i++ {
		entries := task.Configs[i].Entries
		if len(labelMap) == len(task.Configs[i].Entries) {
			check := false
			keys := make([]string, 0, len(entries))
			for k := range entries {
				keys = append(keys, k)
			}

			sort.Strings(keys)
			for _, k := range keys {
				i, ok := labelMap[k]
				if ok == false {
					check = true
					break
				} else {
					if i != entries[k] {
						check = true
						break
					}
				}
			}
			if check != true {
				newTask = append(newTask, task.Configs[i])
			}
		}
	}

	if ok != nil {
		err := errors.New("key not found")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if len(newTask) == 0 {
		err := errors.New("params not match")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(ctx, w, newTask)
}

func (ts *Service) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling delete config at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("version not found")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	msg, err := ts.cf.Delete(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, msg)
}

func (ts *Service) addConfigToExistingGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("addConfigToExistingGroup", ts.tracer, req)
	defer span.Finish()

	span.LogFields(tracer.LogString("handler", fmt.Sprintf("handling add config to a group at %s\n", req.URL.Path)))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	cf, err := decodeBodyConfig(ctx, req.Body)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]

	if len(version) == 0 {
		err := errors.New("version not found")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	task, ok := ts.cf.Get(ctx, id, version)
	if task == nil {
		err := errors.New("key not found")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if ok != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, c := range cf {
		task.Configs = append(task.Configs, c)
	}

	result, _, err := ts.cf.PostNewVersion(ctx, task, id, version)
	if err != nil {
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if result == nil {
		err := errors.New("Update error")
		tracer.LogError(span, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(ctx, w, result)
}
