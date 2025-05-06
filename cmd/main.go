package main

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/internal/report/generate"
	"CompetitionLogger/pkg/events"
	"CompetitionLogger/pkg/logger"
	"context"
	"fmt"
	"os"
)

func main() {
	// Initialize logger
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	// Initialize paths to input data
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logger.GetFromContext(ctx).Error("Config path not set or set incorrect")
	}

	eventsPath := os.Getenv("EVENTS_PATH")
	if eventsPath == "" {
		logger.GetFromContext(ctx).Error("Events path not set or set incorrect")
	}

	// Loading and parsing config.json
	configBytes := config.LoadConfig(ctx, configPath)
	raceConfig := config.ParseConfig(ctx, configBytes)

	// Loading and parsing events.txt
	eventsFile := events.LoadEvents(ctx, eventsPath)
	store := events.ParseEvents(ctx, eventsFile)

	// Generating race's logs
	keys := events.SortMapByKey(store.ByTime())
	for _, key := range keys {
		generatedLog := generate.Log(store.ByTime()[key])
		fmt.Printf("%v\n", generatedLog)
	}

	// Generating race's report table
	reports := generate.ReportTable(raceConfig, store.ByCompetitor())
	fmt.Println(generate.FormatReport(reports))
}
