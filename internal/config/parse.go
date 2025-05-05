package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Race struct {
	Laps        int
	LapLen      int
	PenaltyLen  int
	FiringLines int
	Start       string
	StartDelta  string
}

func LoadConfig(pathToConfig string) ([]byte, error) {
	configFile, err := os.Open(pathToConfig)
	if err != nil {
		return nil, fmt.Errorf("error opening configFile: %v", err)
	}
	defer configFile.Close()

	configFileBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("error with convert configFile into []byte: %v", err)
	}

	return configFileBytes, nil
}

func ParseConfig(jsonFileBytes []byte) (Race, error) {
	race := Race{}
	err := json.Unmarshal(jsonFileBytes, &race)
	if err != nil {
		return race, fmt.Errorf("error filling race instance: %v", err)
	}

	return race, nil
}
