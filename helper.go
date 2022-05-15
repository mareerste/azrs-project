package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// *RequestPost
func decodeBody(r io.Reader) ([]*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cf []*Config

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

func generateKey() (string, string) {
	id := uuid.New().String()

	return fmt.Sprintf("configs/%s", id), id
}

func constructKey(id string) string {
	return fmt.Sprintf("configs/%s", id)
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{cli: client}, nil
}

// KREIRANJE CONF
// TODO ? vvv valja li prosledjen parametar: conf []*Config ?
func (confStore *ConfigStore) AddConfig(conf *Config) (*Config, error) {
	kv := confStore.cli.KV()
	sid, rid := generateKey()
	conf.Id = rid

	data, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}

	c := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(c, nil)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// VRACANJE CONF
func (confStore *ConfigStore) GetConf(id string) (*Config, error) {
	kv := ps.cli.KV()
	pair, _, err := kv.Get(constructKey(id), nil)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = json.Unmarshal(pair.Value, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// VRACANJE SVIH POSTOVA
// TODO: Valja li kada je primenjeno na nas model (gde nema posebnu
// strukturu ConfGroup, vec konstruisemo niz pri vracanju odgovora umesto toga) ???
const (
	all = "allConfigs"
)

func (confStore *ConfigStore) GetAll() ([]*Config, error) {
	kv := confStore.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	configs := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

// BRISANJE CONFIG-A
func (confStore *ConfigStore) Delete(id string) (map[string]string, error) {
	kv := confStore.cli.KV()
	_, err := kv.Delete(constructKey(id), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}
