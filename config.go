package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	WordSearchSystemAddress string `json:"wordSearchSystemAddress"`
}

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
