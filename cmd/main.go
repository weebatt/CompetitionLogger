package main

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/internal/report/generate"
	"CompetitionLogger/pkg/events"
	"CompetitionLogger/pkg/logger"
	"context"
	"fmt"
)

func main() {
	// Initialize logger
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	// Loading and parsing config.json
	configBytes := config.LoadConfig(ctx, "/Users/macbook/Projects/test_tasks/CompetitionLogger/sunny_5_skiers/config_test.json")
	raceConfig := config.ParseConfig(ctx, configBytes)

	// Loading and parsing events.txt
	eventsFile := events.LoadEvents(ctx, "/Users/macbook/Projects/test_tasks/CompetitionLogger/sunny_5_skiers/events_test")
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
