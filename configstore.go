package main

import (
	"azrs-project/tracer"
	"context"
	"encoding/json"
	"errors"
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

func (ps *ConfigStore) Get(ctx context.Context, id string, version string) (*Configs, error) {
	span := tracer.StartSpanFromContext(ctx, "GetDB")
	defer span.Finish()

	span.LogFields(tracer.LogString("GetDB", "configStore Get"))

	kv := ps.cli.KV()

	pair, _, err := kv.Get(constructKey(id, version), nil)

	if err != nil || pair == nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := &Configs{}
	err = json.Unmarshal(pair.Value, configs)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return configs, nil
}

func (ps *ConfigStore) GetIdemKey(ctx context.Context, id string) (*Idem, error) {

	span := tracer.StartSpanFromContext(ctx, "GetIdemKey")
	defer span.Finish()

	kv := ps.cli.KV()

	pair, _, err := kv.Get(id, nil)

	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	} else if pair == nil {
		// tracer.LogError(span, err) //todo
		span.LogFields(tracer.LogString("idemkey", fmt.Sprintf("Pair doesnt exist")))
		return nil, err
	}

	idem := &Idem{}
	err = json.Unmarshal(pair.Value, idem)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	return idem, nil
}

func (ps *ConfigStore) GetAll(ctx context.Context) ([]*Configs, error) {

	span := tracer.StartSpanFromContext(ctx, "GetAllDB")
	defer span.Finish()

	kv := ps.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}

	configs := []*Configs{}
	for _, pair := range data {
		config := &Configs{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			tracer.LogError(span, err)
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (ps *ConfigStore) Delete(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteDB")
	defer span.Finish()

	span.LogFields(tracer.LogString("configstore", "Delete Conf"))

	kv := ps.cli.KV()
	_, err := kv.Delete(constructKey(id, version), nil)
	if err != nil {
		tracer.LogError(span, err)
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

func (ps *ConfigStore) Post(ctx context.Context, configs *Configs, version string) (*Configs, string, error) {
	span := tracer.StartSpanFromContext(ctx, "postRequest")
	defer span.Finish()

	span.LogFields(tracer.LogString("configstore", "Post"))

	kv := ps.cli.KV()

	sid, rid := generateKey(version)
	// post.Id = rid

	data, err := json.Marshal(configs)
	if err != nil {
		/**/ tracer.LogError(span, errors.New("Cannot marshal"))
		return nil, "", err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, "", err
	}

	return configs, rid, nil
}

func (ps *ConfigStore) PostIdemKey(ctx context.Context, idemId string, idem *Idem) error {
	span := tracer.StartSpanFromContext(ctx, "postIdemKey")
	defer span.Finish()

	span.LogFields(tracer.LogString("configstore", "PostIdemKey"))

	kv := ps.cli.KV()

	data, err := json.Marshal(idem)

	if err != nil {
		tracer.LogError(span, errors.New("Cannot marshal"))
		return err
	}

	p := &api.KVPair{Key: idemId, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return err
	}

	return err
}

func (ps *ConfigStore) PostNewVersion(ctx context.Context, configs *Configs, id string, version string) (*Configs, string, error) {
	span := tracer.StartSpanFromContext(ctx, "PostNewVersion")
	defer span.Finish()

	span.LogFields(tracer.LogString("configstore", "PostNewVersion"))

	cntx := tracer.ContextWithSpan(context.Background(), span)

	kv := ps.cli.KV()

	sid, rid := generateKeyNewVersion(cntx, id, version)
	// post.Id = rid

	data, err := json.Marshal(configs)
	if err != nil {
		tracer.LogError(span, err)
		return nil, "", err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
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
