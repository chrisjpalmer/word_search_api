package main

import (
	"encoding/json"
	"os"
)

//Config - defines the config parameters which should be exposed to this microservice
type Config struct {
	//The address of the word_search_system for gRPC calls
	WordSearchSystemAddress string `json:"wordSearchSystemAddress"`
	//The port to listen on for the http server
	HTTPListenAddress string `json:"httpListenAddress"`
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
