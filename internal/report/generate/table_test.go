package generate

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/internal/worker"
	"CompetitionLogger/pkg/events"
	"reflect"
	"sort"
	"testing"
)

func TestReportTable(t *testing.T) {
	raceConfig := config.Race{
		Laps:        2,
		LapLen:      3651,
		PenaltyLen:  50,
		FiringLines: 1,
	}

	tests := []struct {
		name      string
		eventsMap map[int][]events.Event
		want      []worker.CompetitorReport
	}{
		{
			name: "multiple competitors",
			eventsMap: map[int][]events.Event{
				1: {
					{Time: "09:00:00.000", EventID: 2, CompetitorID: 1, ExtraParams: "09:00:00.000"},
					{Time: "09:00:01.000", EventID: 4, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:05:00.000", EventID: 10, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:10:00.000", EventID: 10, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:15:00.000", EventID: 33, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:02:00.000", EventID: 5, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:02:01.000", EventID: 6, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:02:02.000", EventID: 6, CompetitorID: 1, ExtraParams: ""},
				},
				2: {
					{Time: "09:00:00.000", EventID: 2, CompetitorID: 2, ExtraParams: "09:00:00.000"},
					{Time: "09:00:01.000", EventID: 4, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:06:00.000", EventID: 10, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:12:00.000", EventID: 10, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:18:00.000", EventID: 33, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:03:00.000", EventID: 5, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:03:01.000", EventID: 6, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:08:00.000", EventID: 8, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:08:30.000", EventID: 9, CompetitorID: 2, ExtraParams: ""},
				},
			},
			want: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "Finished",
					TotalTime:    "00:15:00.000",
					Laps: []worker.LapInfo{
						{Time: "00:04:59.000", Speed: 3651.0 / worker.TimeToSeconds("00:04:59.000")},
						{Time: "00:05:00.000", Speed: 3651.0 / worker.TimeToSeconds("00:05:00.000")},
					},
					Penalty:   worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots: "2/5",
				},
				{
					CompetitorID: 2,
					Status:       "Finished",
					TotalTime:    "00:18:00.000",
					Laps: []worker.LapInfo{
						{Time: "00:05:59.000", Speed: 3651.0 / worker.TimeToSeconds("00:05:59.000")},
						{Time: "00:06:00.000", Speed: 3651.0 / worker.TimeToSeconds("00:06:00.000")},
					},
					Penalty: worker.PenaltyInfo{
						Time:  "00:00:30.000",
						Speed: float64(50*4) / worker.TimeToSeconds("00:00:30.000"),
					},
					HitsShots: "1/5",
				},
			},
		},
		{
			name:      "empty events map",
			eventsMap: map[int][]events.Event{},
			want:      []worker.CompetitorReport{},
		},
		{
			name: "single competitor not started",
			eventsMap: map[int][]events.Event{
				1: {
					{Time: "09:00:00.000", EventID: 2, CompetitorID: 1, ExtraParams: "09:00:00.000"},
				},
			},
			want: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "NotStarted",
					TotalTime:    "00:00:00.000",
					Laps: []worker.LapInfo{
						{Time: "", Speed: 0.0},
						{Time: "", Speed: 0.0},
					},
					Penalty:   worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots: "0/0",
				},
			},
		},
		{
			name: "competitor not finished",
			eventsMap: map[int][]events.Event{
				1: {
					{Time: "09:00:00.000", EventID: 2, CompetitorID: 1, ExtraParams: "09:00:00.000"},
					{Time: "09:00:01.000", EventID: 4, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:05:00.000", EventID: 10, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:10:00.000", EventID: 11, CompetitorID: 1, ExtraParams: ""},
				},
			},
			want: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "NotFinished",
					TotalTime:    "00:10:00.000",
					Laps: []worker.LapInfo{
						{Time: "00:04:59.000", Speed: 3651.0 / worker.TimeToSeconds("00:04:59.000")},
						{Time: "", Speed: 0.0},
					},
					Penalty:   worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots: "0/0",
				},
			},
		},
		{
			name: "invalid time format",
			eventsMap: map[int][]events.Event{
				1: {
					{Time: "09:00:00.000", EventID: 2, CompetitorID: 1, ExtraParams: "09:00:00.000"},
					{Time: "09:00:01.000", EventID: 4, CompetitorID: 1, ExtraParams: ""},
					{Time: "invalid", EventID: 10, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:15:00.000", EventID: 33, CompetitorID: 1, ExtraParams: ""},
				},
			},
			want: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "Finished",
					TotalTime:    "00:15:00.000",
					Laps: []worker.LapInfo{
						{Time: "00:00:00.000", Speed: 0.0},
						{Time: "", Speed: 0.0},
					},
					Penalty:   worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots: "0/0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReportTable(raceConfig, tt.eventsMap)
			if len(got) != len(tt.want) {
				t.Errorf("ReportTable() len = %d, want %d", len(got), len(tt.want))
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].CompetitorID < got[j].CompetitorID
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].CompetitorID < tt.want[i].CompetitorID
			})

			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("ReportTable()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFormatReport(t *testing.T) {
	tests := []struct {
		name    string
		reports []worker.CompetitorReport
		want    string
	}{
		{
			name: "multiple competitors",
			reports: []worker.CompetitorReport{
				{
					CompetitorID: 2,
					Status:       "Finished",
					TotalTime:    "00:18:00.000",
					Laps: []worker.LapInfo{
						{Time: "00:05:59.000", Speed: 3651.0 / worker.TimeToSeconds("00:05:59.000")},
						{Time: "00:06:00.000", Speed: 3651.0 / worker.TimeToSeconds("00:06:00.000")},
					},
					Penalty: worker.PenaltyInfo{
						Time:  "00:00:30.000",
						Speed: float64(50*4) / worker.TimeToSeconds("00:00:30.000"),
					},
					HitsShots: "1/5",
				},
				{
					CompetitorID: 1,
					Status:       "Finished",
					TotalTime:    "00:15:00.000",
					Laps: []worker.LapInfo{
						{Time: "00:04:59.000", Speed: 3651.0 / worker.TimeToSeconds("00:04:59.000")},
						{Time: "00:05:00.000", Speed: 3651.0 / worker.TimeToSeconds("00:05:00.000")},
					},
					Penalty:   worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots: "2/5",
				},
			},
			want: `[Finished] 1 [{00:04:59.000, 12.211}, {00:05:00.000, 12.170}] {,} 2/5
[Finished] 2 [{00:05:59.000, 10.170}, {00:06:00.000, 10.142}] {00:00:30.000, 6.667} 1/5
`,
		},
		{
			name:    "empty reports",
			reports: []worker.CompetitorReport{},
			want:    "",
		},
		{
			name: "single competitor not started",
			reports: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "NotStarted",
					TotalTime:    "00:00:00.000",
					Laps: []worker.LapInfo{
						{Time: "", Speed: 0.0},
						{Time: "", Speed: 0.0},
					},
					Penalty:   worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots: "0/0",
				},
			},
			want: `[NotStarted] 1 [{,}, {,}] {,} 0/0
`,
		},
		{
			name: "competitor with penalty only",
			reports: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "Started",
					TotalTime:    "00:00:30.000",
					Laps: []worker.LapInfo{
						{Time: "", Speed: 0.0},
						{Time: "", Speed: 0.0},
					},
					Penalty: worker.PenaltyInfo{
						Time:  "00:00:30.000",
						Speed: float64(50*4) / worker.TimeToSeconds("00:00:30.000"),
					},
					HitsShots: "0/5",
				},
			},
			want: `[Started] 1 [{,}, {,}] {00:00:30.000, 6.667} 0/5
`,
		},
		{
			name: "nil laps",
			reports: []worker.CompetitorReport{
				{
					CompetitorID: 1,
					Status:       "NotStarted",
					TotalTime:    "00:00:00.000",
					Laps:         nil,
					Penalty:      worker.PenaltyInfo{Time: "", Speed: 0.0},
					HitsShots:    "0/0",
				},
			},
			want: `[NotStarted] 1 [] {,} 0/0
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatReport(tt.reports)
			if got != tt.want {
				t.Errorf("FormatReport() = %q, want %q", got, tt.want)
			}
		})
	}
}
