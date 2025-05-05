package config

import (
	"CompetitionLogger/pkg/logger"
	"context"
	"encoding/json"
	"go.uber.org/zap"
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

func LoadConfig(ctx context.Context, pathToConfig string) []byte {
	configFile, err := os.Open(pathToConfig)
	if err != nil {
		logger.GetFromContext(ctx).Error("error opening configFile: ", zap.Error(err))
		return nil
	}
	defer configFile.Close()

	configFileBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		logger.GetFromContext(ctx).Error("error with convert configFile into []byte: ", zap.Error(err))
		return nil
	}

	return configFileBytes
}

func ParseConfig(ctx context.Context, jsonFileBytes []byte) Race {
	race := Race{}
	err := json.Unmarshal(jsonFileBytes, &race)
	if err != nil {
		logger.GetFromContext(ctx).Error("error unmarshal jsonFileBytes: ", zap.Error(err))
		return race
	}

	return race
}
