package main

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/internal/report/generate"
	"CompetitionLogger/pkg/events"
	"fmt"
	"log"
)

func main() {
	// Loading and parsing config.json
	configBytes, err := config.LoadConfig("/Users/macbook/Projects/test_tasks/CompetitionLogger/sunny_5_skiers/config_test.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	raceConfig, err := config.ParseConfig(configBytes)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	fmt.Print(raceConfig)

	// Loading and parsing events.txt
	eventsFile, err := events.LoadEvents("/Users/macbook/Projects/test_tasks/CompetitionLogger/sunny_5_skiers/events_test")
	if err != nil {
		log.Fatalf("Error loading events: %v", err)
	}

	store, err := events.ParseEvents(eventsFile)
	if err != nil {
		log.Fatalf("Error parsing events: %v", err)
	}

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
