// package configstore
package main

import (
	"github.com/hashicorp/consul/api"
)

type Service struct {
	cf *ConfigStore
}

type ConfigStore struct {
	cli *api.Client
}

type Configs struct {
	Configs map[string][]*Config `json:"configs"`
	Id      string               `json:"id"`
	//version
}

type Config struct {
	Entries map[string]string `json:"entries"`
}
