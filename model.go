package main

type Service struct {
	Data map[string][]*Config
}

type Config struct {
	Entries map[string]string `json:"entries"`
}
