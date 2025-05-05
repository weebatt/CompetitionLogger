package events

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         string
	EventID      int
	CompetitorID int
	ExtraParams  string
}

func LoadEvents(pathToEvents string) (*os.File, error) {
	eventsFile, err := os.Open(pathToEvents)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	return eventsFile, nil
}

func ParseEvents(eventsFile *os.File) (map[int][]Event, error) {
	eventsByCompetitor := make(map[int][]Event)
	scanner := bufio.NewScanner(eventsFile)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		event, err := ParseEvent(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing event line: %v", err)
		}

		competitorID := event.CompetitorID
		eventsByCompetitor[competitorID] = append(eventsByCompetitor[competitorID], event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning events: %v", err)
	}

	defer eventsFile.Close()
	return eventsByCompetitor, nil
}

func ParseEvent(line string) (Event, error) {
	if !strings.HasPrefix(line, "[") {
		return Event{}, fmt.Errorf("incorrect time's format")
	}

	endTimeIdx := strings.Index(line, "]")
	if endTimeIdx == -1 {
		return Event{}, fmt.Errorf("incorrect time's format: missing close bracket")
	}

	timeStr := line[1:endTimeIdx]
	layout := "15:04:05.000"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return Event{}, fmt.Errorf("error in time parsing %q: %v", timeStr, err)
	}

	formattedTime := parsedTime.Format(layout)

	rest := strings.TrimSpace(line[endTimeIdx+1:])
	parts := strings.Fields(rest)
	if len(parts) < 2 {
		return Event{}, fmt.Errorf("less than two parts in %q: %v", rest, err)
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return Event{}, fmt.Errorf("error in event id parsing: %v", err)
	}

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		return Event{}, fmt.Errorf("error in competitor id parsing: %v", err)
	}

	extraParams := ""
	if len(parts) > 2 {
		extraParams = strings.Join(parts[2:], " ")
	}

	return Event{
		Time:         formattedTime,
		EventID:      eventID,
		CompetitorID: competitorID,
		ExtraParams:  extraParams,
	}, nil
}
