package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	Token   string `json:"token"`
	Debug   bool   `json:"debug"`
	Timeout int    `json:"timeout"`
}

func loadConfig() *config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	config := config{}

	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
