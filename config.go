package main

import (
	"encoding/json"
	"os"
)

//Config - defines the config parameters which should be exposed to this microservice
type Config struct {
	WordSearchSystemAddress string `json:"wordSearchSystemAddress"`
}

//ParseConfig - reads the json file at configPath and outputs the Config structure
func ParseConfig(configPath string) (*Config, error) {
	var (
		err  error
		file *os.File
	)
	file, err = os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonDec := json.NewDecoder(file)
	var config Config
	err = jsonDec.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
