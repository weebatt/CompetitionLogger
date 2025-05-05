package report

import (
	"CompetitionLogger/pkg/events"
	"fmt"
)

var eventComments = map[int]string{
	1:  "[%s] The competitor(%d) registered",
	2:  "[%s] The start time for the competitor(%d) was set by a draw to %s",
	3:  "[%s] The competitor(%d) is on the start line",
	4:  "[%s] The competitor(%d) has started",
	5:  "[%s] The competitor(%d) is on the firing range(%s)",
	6:  "[%s] The target(%s) has been hit by competitor(%d)",
	7:  "[%s] The competitor(%d) left the firing range",
	8:  "[%s] The competitor(%d) entered the penalty laps",
	9:  "[%s] The competitor(%d) left the penalty laps",
	10: "[%s] The competitor(%d) ended the main lap",
	11: "[%s] The competitor(%d) can't continue: %s",
	32: "[%s] The competitor(%d) is disqualified",
	33: "[%s] The competitor(%d) has finished",
}

func GenerateLog(event events.Event) string {
	comment, exists := eventComments[event.EventID]
	if !exists {
		return fmt.Sprintf("Unknown event %d for competitor %d", event.EventID, event.CompetitorID)
	}

	switch event.EventID {
	case 2, 5, 11:
		return fmt.Sprintf(comment, event.Time, event.CompetitorID, event.ExtraParams)
	case 6:
		return fmt.Sprintf(comment, event.Time, event.ExtraParams, event.CompetitorID)
	default:
		return fmt.Sprintf(comment, event.Time, event.CompetitorID)
	}
}
