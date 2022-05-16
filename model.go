package main

type Service struct {
	Data map[string][]*Configs
}

// type Configs struct {
// 	Version string    `json:"version"`
// 	Configs []*Config `json:"configs"`
// }

type Configs struct {
	Configs map[string][]*Config `json:"configs"`
}

type Config struct {
	Entries map[string]string `json:"entries"`
}
