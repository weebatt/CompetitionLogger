package main

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/internal/report"
	"CompetitionLogger/pkg/events"
	"fmt"
	"log"
)

func main() {
	// Loading and parsing config.json
	configBytes, err := config.LoadConfig("/Users/macbook/Downloads/CompetitionLogger/sunny_5_skiers/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	race, err := config.ParseConfig(configBytes)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	fmt.Printf("Race: %v\n", race)

	// Loading and parsing events.txt
	eventsFile, err := events.LoadEvents("/Users/macbook/Downloads/CompetitionLogger/sunny_5_skiers/events")
	if err != nil {
		log.Fatalf("Error loading events: %v", err)
	}

	eventsMap, err := events.ParseEvents(eventsFile)
	if err != nil {
		log.Fatalf("Error parsing events: %v", err)
	}

	// Generating race's logs
	keys := events.SortMapByKey(eventsMap)
	for _, key := range keys {
		generatedLog := report.GenerateLog(eventsMap[key])
		fmt.Printf("%v\n", generatedLog)
	}

}
