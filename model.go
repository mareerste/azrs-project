// package configstore
package main

import (
	"github.com/hashicorp/consul/api"
)

type Service struct {
	cf *ConfigStore
}

type Config struct {
	// *Da li uraditi i version na ConfGroup (entitet/struktura koju
	// nemamo u kodu) i kako prilagoditi handler-e i json body na novi nacin struktuisanja ??
	// TODO uncomment:
	//version string
	Entries map[string]string `json:"entries"`
}

type ConfigStore struct {
	cli *api.Client
}
