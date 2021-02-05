package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Token   string `json:"token"`
	Debug   bool   `json:"debug"`
	Timeout int    `json:"timeout"`
}

func LoadConfig() *Config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	config := Config{}

	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
