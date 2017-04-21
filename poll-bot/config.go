package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config holds Bot configuration (read from ./config.json)
type Config struct {
	Token        string   `json:"token"`
	Choices      []string `json:"choices"`
	VotesPerUser int      `json:"votesPerUser"`
	Guild        string   `json:"serverID"`
	Role         string   `json:"roleID"`
	RulesURL     string   `json:"rulesURL"`
}

func getConfig(filename string) (*Config, error) {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf Config

	err = json.Unmarshal(jsonData, &conf)

	return &conf, err
}
