package worker

import (
	"CompetitionLogger/internal/config"
	"CompetitionLogger/pkg/events"
	"math"
	"testing"
	"time"
)

func TestSubtractTimes(t *testing.T) {
	type content struct {
		name     string
		t1       string
		t2       string
		expected string
	}
	tests := []content{
		{
			name:     "normal subtraction",
			t1:       "09:30:00.000",
			t2:       "09:59:03.872",
			expected: "00:29:03.872",
		},
		{
			name:     "empty t1",
			t1:       "",
			t2:       "09:59:03.872",
			expected: "00:00:00.000",
		},
		{
			name:     "empty t2",
			t1:       "09:30:00.000",
			t2:       "",
			expected: "00:00:00.000",
		},
		{
			name:     "negative duration",
			t1:       "09:59:03.872",
			t2:       "09:30:00.000",
			expected: "00:00:00.000",
		},
		{
			name:     "same time",
			t1:       "09:30:00.000",
			t2:       "09:30:00.000",
			expected: "00:00:00.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := subtractTimes(tt.t2, tt.t1)
			if result != tt.expected {
				t.Errorf("subtractTimes(%q, %q) = %q, want %q", tt.t2, tt.t1, result, tt.expected)
			}
		})
	}
}

func TestTimeToSeconds(t *testing.T) {
	type content struct {
		name     string
		timeStr  string
		expected float64
	}
	tests := []content{
		{
			name:     "normal time",
			timeStr:  "00:29:03.872",
			expected: 1743.872,
		},
		{
			name:     "zero time",
			timeStr:  "00:00:00.000",
			expected: 0.0,
		},
		{
			name:     "only milliseconds",
			timeStr:  "00:00:00.500",
			expected: 0.5,
		},
		{
			name:     "invalid format",
			timeStr:  "invalid",
			expected: 0.0,
		},
		{
			name:     "missing milliseconds",
			timeStr:  "00:00:00",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimeToSeconds(tt.timeStr)
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("timeToSeconds(%q) = %v, want %v", tt.timeStr, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	type content struct {
		name     string
		duration time.Duration
		expected string
	}

	tests := []content{
		{
			name:     "normal duration",
			duration: 1743872 * time.Millisecond,
			expected: "00:29:03.872",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00:00.000",
		},
		{
			name:     "negative duration",
			duration: -1 * time.Second,
			expected: "00:00:00.000",
		},
		{
			name:     "only milliseconds",
			duration: 500 * time.Millisecond,
			expected: "00:00:00.500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestProcessCompetitor(t *testing.T) {
	racecConfig := config.Race{
		Laps:        2,
		LapLen:      3651,
		PenaltyLen:  50,
		FiringLines: 1,
	}

	type content struct {
		name     string
		events   []events.Event
		expected CompetitorReport
	}

	tests := []content{
		{
			name: "not finished competitor",
			events: []events.Event{
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				{Time: "09:30:01.005", EventID: 4, CompetitorID: 1},
				{Time: "09:49:31.659", EventID: 5, CompetitorID: 1, ExtraParams: "1"},
				{Time: "09:49:33.123", EventID: 6, CompetitorID: 1, ExtraParams: "1"},
				{Time: "09:49:34.650", EventID: 6, CompetitorID: 1, ExtraParams: "2"},
				{Time: "09:49:35.937", EventID: 6, CompetitorID: 1, ExtraParams: "4"},
				{Time: "09:49:37.364", EventID: 6, CompetitorID: 1, ExtraParams: "5"},
				{Time: "09:49:55.915", EventID: 8, CompetitorID: 1},
				{Time: "09:51:48.391", EventID: 9, CompetitorID: 1},
				{Time: "09:59:03.872", EventID: 10, CompetitorID: 1},
				{Time: "09:59:03.872", EventID: 11, CompetitorID: 1, ExtraParams: "Заблудился в лесу"},
			},
			expected: CompetitorReport{
				CompetitorID: 1,
				Status:       "NotFinished",
				TotalTime:    "00:29:03.872",
				Laps: []LapInfo{
					{Time: "00:29:02.867", Speed: 3651.0 / 1742.867},
					{Time: "", Speed: 0.0},
				},
				Penalty: PenaltyInfo{
					Time:  "00:01:52.476",
					Speed: 50.0 / 112.476,
				},
				HitsShots: "4/5",
			},
		},
		{
			name: "not started competitor",
			events: []events.Event{
				{Time: "09:05:59.867", EventID: 1, CompetitorID: 2},
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 2, ExtraParams: "09:30:30.000"},
			},
			expected: CompetitorReport{
				CompetitorID: 2,
				Status:       "NotStarted",
				TotalTime:    "00:00:00.000",
				Laps:         []LapInfo{{Time: "", Speed: 0.0}, {Time: "", Speed: 0.0}},
				Penalty:      PenaltyInfo{Time: "", Speed: 0.0},
				HitsShots:    "0/0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessCompetitor(racecConfig, tt.expected.CompetitorID, tt.events)

			if result.Status != tt.expected.Status {
				t.Errorf("Status = %q, want %q", result.Status, tt.expected.Status)
			}
			if result.TotalTime != tt.expected.TotalTime {
				t.Errorf("TotalTime = %q, want %q", result.TotalTime, tt.expected.TotalTime)
			}
			if result.HitsShots != tt.expected.HitsShots {
				t.Errorf("HitsShots = %q, want %q", result.HitsShots, tt.expected.HitsShots)
			}

			if len(result.Laps) != len(tt.expected.Laps) {
				t.Errorf("Laps length = %d, want %d", len(result.Laps), len(tt.expected.Laps))
			}
			for i, lap := range result.Laps {
				if lap.Time != tt.expected.Laps[i].Time {
					t.Errorf("Lap[%d].Time = %q, want %q", i, lap.Time, tt.expected.Laps[i].Time)
				}
				if math.Abs(lap.Speed-tt.expected.Laps[i].Speed) > 0.001 {
					t.Errorf("Lap[%d].Speed = %v, want %v", i, lap.Speed, tt.expected.Laps[i].Speed)
				}
			}

			if result.Penalty.Time != tt.expected.Penalty.Time {
				t.Errorf("Penalty.Time = %q, want %q", result.Penalty.Time, tt.expected.Penalty.Time)
			}
			if math.Abs(result.Penalty.Speed-tt.expected.Penalty.Speed) > 0.001 {
				t.Errorf("Penalty.Speed = %v, want %v", result.Penalty.Speed, tt.expected.Penalty.Speed)
			}
		})
	}
}
