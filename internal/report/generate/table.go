package generate

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/internal/worker"
	"CompetitionLogger/pkg/events"
	"fmt"
	"sort"
	"strings"
)

func ReportTable(config config.Race, eventsMap map[int][]events.Event) []worker.CompetitorReport {
	var reports []worker.CompetitorReport

	for competitorID, es := range eventsMap {
		report := worker.ProcessCompetitor(config, competitorID, es)
		reports = append(reports, report)
	}

	return reports
}

func FormatReport(reports []worker.CompetitorReport) string {
	var result strings.Builder
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].CompetitorID < reports[j].CompetitorID
	})
	for _, r := range reports {
		result.WriteString(fmt.Sprintf("[%s] %d ", r.Status, r.CompetitorID))

		result.WriteString("[")
		for i, lap := range r.Laps {
			if lap.Time == "" {
				result.WriteString("{,}")
			} else {
				result.WriteString(fmt.Sprintf("{%.12s, %.3f}", lap.Time, lap.Speed))
			}
			if i < len(r.Laps)-1 {
				result.WriteString(", ")
			}
		}
		result.WriteString("] ")

		if r.Penalty.Time == "" {
			result.WriteString("{,}")
		} else {
			result.WriteString(fmt.Sprintf("{%.12s, %.3f}", r.Penalty.Time, r.Penalty.Speed))
		}
		result.WriteString(" ")

		result.WriteString(r.HitsShots)

		result.WriteString("\n")
	}
	return result.String()
}
