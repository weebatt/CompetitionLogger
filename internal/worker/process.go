package worker

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/pkg/events"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CompetitorReport struct {
	CompetitorID int
	Status       string
	TotalTime    string
	Laps         []LapInfo
	Penalty      PenaltyInfo
	HitsShots    string
}

type LapInfo struct {
	Time  string
	Speed float64
}

type PenaltyInfo struct {
	Time  string
	Speed float64
}

func ProcessCompetitor(config config.Race, competitorID int, events []events.Event) CompetitorReport {
	var reportTable CompetitorReport
	reportTable.CompetitorID = competitorID
	reportTable.Status = "NotStarted"

	var plannedStart, actualStart, finishTime, lastEventTime string
	var lapTimes []string
	var penaltyStart, penaltyEnd string
	hits := 0
	firingLineVisits := 0

	for _, event := range events {
		switch event.EventID {
		case 2:
			plannedStart = event.ExtraParams
			reportTable.Status = "NotStarted"
		case 4:
			actualStart = event.Time
			reportTable.Status = "Started"
		case 5:
			firingLineVisits++
		case 6:
			hits++
		case 8:
			penaltyStart = event.Time
		case 9:
			penaltyEnd = event.Time
		case 10:
			lapTimes = append(lapTimes, event.Time)
		case 11:
			reportTable.Status = "NotFinished"
			lastEventTime = event.Time
		case 32:
			reportTable.Status = "NotFinished"
			lastEventTime = event.Time
		case 33:
			finishTime = event.Time
			reportTable.Status = "Finished"
		}
	}

	shots := firingLineVisits * 5

	if reportTable.Status == "NotStarted" {
		reportTable.TotalTime = "00:00:00.000"
	} else if reportTable.Status == "NotFinished" {
		reportTable.TotalTime = subtractTimes(lastEventTime, plannedStart)
	} else if reportTable.Status == "Finished" {
		reportTable.TotalTime = subtractTimes(finishTime, plannedStart)
	} else {
		if actualStart != "" && lastEventTime != "" {
			reportTable.TotalTime = subtractTimes(lastEventTime, plannedStart)
		}
	}

	for i := 0; i < config.Laps; i++ {
		var lapInfo LapInfo
		if i < len(lapTimes) {
			start := actualStart
			if i > 0 {
				start = lapTimes[i-1]
			}
			lapInfo.Time = subtractTimes(lapTimes[i], start)
			lapSeconds, _ := timeToSeconds(lapInfo.Time)
			if lapSeconds > 0 {
				lapInfo.Speed = float64(config.LapLen) / lapSeconds
			}
		}
		reportTable.Laps = append(reportTable.Laps, lapInfo)
	}

	if penaltyStart != "" && penaltyEnd != "" {
		reportTable.Penalty.Time = subtractTimes(penaltyEnd, penaltyStart)
		penaltySeconds, _ := timeToSeconds(reportTable.Penalty.Time)
		if penaltySeconds > 0 {
			misses := shots - hits
			penaltyDistance := misses * config.PenaltyLen
			reportTable.Penalty.Speed = float64(penaltyDistance) / penaltySeconds
		}
	}

	reportTable.HitsShots = fmt.Sprintf("%d/%d", hits, shots)

	return reportTable
}

func subtractTimes(t2, t1 string) string {
	layout := "15:04:05.000"
	time1, _ := time.Parse(layout, t1)
	time2, _ := time.Parse(layout, t2)
	duration := time2.Sub(time1)
	return formatDuration(duration)
}

func timeToSeconds(timeStr string) (float64, error) {
	hoursMinutesSeconds := strings.Split(timeStr, ":")
	if len(hoursMinutesSeconds) != 3 {
		return 0, fmt.Errorf("wrong HH:MM:SS format: %s", timeStr)
	}

	secondsMilliseconds := strings.Split(hoursMinutesSeconds[2], ".")
	if len(secondsMilliseconds) != 2 {
		return 0, fmt.Errorf("wrong SS:sss format: %s", timeStr)
	}

	hours, _ := strconv.Atoi(hoursMinutesSeconds[0])
	minutes, _ := strconv.Atoi(hoursMinutesSeconds[1])
	seconds, _ := strconv.Atoi(secondsMilliseconds[0])
	milliseconds, _ := strconv.Atoi(secondsMilliseconds[1])
	return float64(hours*3600+minutes*60+seconds) + float64(milliseconds)/1000.0, nil
}

func formatDuration(d time.Duration) string {
	hours := int(d / time.Hour)
	minutes := int(d/time.Minute) % 60
	seconds := int(d/time.Second) % 60
	milliseconds := int(d/time.Millisecond) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}
