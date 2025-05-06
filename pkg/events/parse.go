package events

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"CompetitionLogger/pkg/logger"
	"go.uber.org/zap"
)

type Event struct {
	Time         string
	EventID      int
	CompetitorID int
	ExtraParams  string
}

type EventStore struct {
	events []Event
}

func LoadEvents(ctx context.Context, pathToEvents string) *os.File {
	eventsFile, err := os.Open(pathToEvents)
	if err != nil {
		logger.GetFromContext(ctx).Debug("error opening file: ", zap.Error(err))
		return nil
	}

	logger.GetFromContext(ctx).Info("success loading events from file", zap.String("path_to_events", pathToEvents))
	return eventsFile
}

func ParseEvents(ctx context.Context, eventsFile *os.File) *EventStore {
	store := &EventStore{}
	if eventsFile == nil {
		logger.GetFromContext(ctx).Warn("Events file is nil")
		return store
	}

	scanner := bufio.NewScanner(eventsFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event := parseEvent(ctx, line)
		if event.Time == "" && event.EventID == 0 && event.CompetitorID == 0 {
			continue
		}
		store.events = append(store.events, event)
	}

	if err := scanner.Err(); err != nil {
		logger.GetFromContext(ctx).Error("error scanning file", zap.Error(err))
		return store
	}

	logger.GetFromContext(ctx).Info("success parsed events")
	return store
}

func parseEvent(ctx context.Context, line string) Event {
	if !strings.HasPrefix(line, "[") {
		logger.GetFromContext(ctx).Error("incorrect time's format", zap.String("line", line))
		return Event{}
	}

	endTimeIdx := strings.Index(line, "]")
	if endTimeIdx == -1 {
		logger.GetFromContext(ctx).Error("incorrect time's format: missing close bracket")
		return Event{}
	}

	timeStr := line[1:endTimeIdx]
	layout := "15:04:05.000"
	_, err := time.Parse(layout, timeStr)
	if err != nil {
		logger.GetFromContext(ctx).Error("error parsing time", zap.String("time", timeStr), zap.Error(err))
		return Event{}
	}

	formattedTime := timeStr

	rest := strings.TrimSpace(line[endTimeIdx+1:])
	parts := strings.Fields(rest)
	if len(parts) < 2 {
		logger.GetFromContext(ctx).Error("incorrect format: missing fields", zap.String("line", line))
		return Event{}
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		logger.GetFromContext(ctx).Error("error parsing event id", zap.String("eventID", parts[0]), zap.Error(err))
		return Event{}
	}

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		logger.GetFromContext(ctx).Error("error parsing competitor id", zap.String("competitorID", parts[1]), zap.Error(err))
		return Event{}
	}

	extraParams := ""
	if len(parts) > 2 {
		extraParams = strings.Join(parts[2:], " ")
	}

	logger.GetFromContext(ctx).Info("success parse line into Event")
	return Event{
		Time:         formattedTime,
		EventID:      eventID,
		CompetitorID: competitorID,
		ExtraParams:  extraParams,
	}
}

// ByTime using time as a key
func (s *EventStore) ByTime() map[string]Event {
	result := make(map[string]Event)
	for _, event := range s.events {
		result[event.Time] = event
	}
	return result
}

// ByCompetitor using CompetitorID as a key
func (s *EventStore) ByCompetitor() map[int][]Event {
	result := make(map[int][]Event)
	for _, event := range s.events {
		result[event.CompetitorID] = append(result[event.CompetitorID], event)
	}
	return result
}

func SortMapByKey(eventsMap map[string]Event) []string {
	keys := make([]string, 0, len(eventsMap))
	for key := range eventsMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
