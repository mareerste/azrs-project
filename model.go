// package configstore
package main

// type Service struct {
// 	cf *ConfigStore
// }

// type ConfigStore struct {
// 	cli *api.Client
// }

// type Configs struct {
// 	Configs map[string][]*Config `json:"configs"`
// 	Id      string               `json:"id"`
// 	//version
// }

// type Config struct {
// 	Entries map[string]string `json:"entries"`
// }

type Configs struct {
	// Configs map[string][]*Config `json:"configs"`
	Configs []*Config `json:"configs"`
}

type Config struct {
	Entries map[string]string `json:"entries"`
}

type Idem struct {
	Config_id string `json:"config_id"`
	Status    string `json:"status"`
}
