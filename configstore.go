package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type ConfigStore struct {
	cli *api.Client
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

	return &ConfigStore{
		cli: client,
	}, nil
}

func (ps *ConfigStore) Get(id string, version string) (*Configs, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKey(id, version), nil)

	if err != nil || pair == nil {
		return nil, err
	}

	configs := &Configs{}
	err = json.Unmarshal(pair.Value, configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (ps *ConfigStore) GetIdemKey(id string) (*Idem, error) {
	kv := ps.cli.KV()

	pair, _, err := kv.Get(id, nil)

	if err != nil || pair == nil {
		return nil, err
	}

	idem := &Idem{}
	err = json.Unmarshal(pair.Value, idem)
	if err != nil {
		return nil, err
	}

	return idem, nil
}

func (ps *ConfigStore) GetAll() ([]*Configs, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	configs := []*Configs{}
	for _, pair := range data {
		config := &Configs{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (ps *ConfigStore) Delete(id string, version string) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructKey(id, version), nil)
	if err != nil {
		return nil, err
	}

	return map[string]string{"Deleted": id}, nil
}

// func (ps *ConfigStore) DeleteVersion(id string, version string) (*Configs, string, error) {
// 	kv := ps.cli.KV()
// 	pair, _, err := kv.Get(constructKey(id), nil)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	configs := &Configs{}
// 	err = json.Unmarshal(pair.Value, configs)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	var configCopy = configs
// 	_, error := kv.Delete(constructKey(id), nil)
// 	if error != nil {
// 		return nil, "", error
// 	}
// 	delete(configCopy.Configs, version)
// 	data, err := json.Marshal(configCopy)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	p := &api.KVPair{Key: id, Value: data}
// 	_, err = kv.Put(p, nil)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return configs, id, nil

// }

func (ps *ConfigStore) Post(configs *Configs, version string) (*Configs, string, error) {
	kv := ps.cli.KV()

	sid, rid := generateKey(version)
	// post.Id = rid

	data, err := json.Marshal(configs)
	if err != nil {
		return nil, "", err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, "", err
	}

	return configs, rid, nil
}

func (ps *ConfigStore) PostIdemKey(idemId string, idem *Idem) error {
	kv := ps.cli.KV()

	data, err := json.Marshal(idem)
	
	if err != nil {
		return err
	}

	p := &api.KVPair{Key: idemId, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}

	return err
}

func (ps *ConfigStore) PostNewVersion(configs *Configs, id string, version string) (*Configs, string, error) {
	kv := ps.cli.KV()

	sid, rid := generateKeyNewVersion(id, version)
	// post.Id = rid

	data, err := json.Marshal(configs)
	if err != nil {
		return nil, "", err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, "", err
	}

	return configs, rid, nil
}

func (ps *ConfigStore) Put(configs *Configs, rid string) (*Configs, error) {
	// kv := ps.cli.KV()

	// // data, err := json.Marshal(configs)
	// // if err != nil {
	// // 	return nil, err
	// // }

	// pair, _, err := kv.Get(constructKey(rid, version), nil)
	// if err != nil {
	// 	return nil, err
	// }

	// post := &Configs{}
	// err = json.Unmarshal(pair.Value, post)
	// if err != nil {
	// 	return nil, err
	// }
	// data, err := json.Marshal(post)
	// if err != nil {
	// 	return nil, err
	// }

	// p := &api.KVPair{Key: rid, Value: data}
	// _, err = kv.Put(p, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// return configs, nil
	return nil, nil
}
