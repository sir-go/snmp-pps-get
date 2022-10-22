package main

import (
	"encoding/json"
	"os"
)

type Cfg struct {
	Oids *[]struct {
		Models *[]string `json:"models"`
		Tx     string    `json:"tx"`
		Rx     string    `json:"rx"`
	} `json:"oids"`
	Default struct {
		Tx string `json:"tx"`
		Rx string `json:"rx"`
	} `json:"default"`
}

func LoadConfig(confpath string) (*Cfg, error) {

	conf := new(Cfg)
	file, err := os.Open(confpath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}(file)

	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}
