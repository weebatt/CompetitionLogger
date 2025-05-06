package events

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

const (
	key = "logger"
)

func TestLoadEvents(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	ctx := context.WithValue(context.Background(), key, zap.NewNop())

	tests := []struct {
		name        string
		path        string
		content     string
		wantNil     bool
		wantLogInfo bool
	}{
		{
			name:        "valid file",
			path:        "events.txt",
			content:     "[09:05:59.867] 1 1\n",
			wantNil:     false,
			wantLogInfo: true,
		},
		{
			name:        "non-existent file",
			path:        "nonexistent.txt",
			content:     "",
			wantNil:     true,
			wantLogInfo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.content != "" {
				tmpfile, err := os.CreateTemp("", tt.path)
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tmpfile.Name())
				if _, err := tmpfile.WriteString(tt.content); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				if err := tmpfile.Close(); err != nil {
					t.Fatalf("Failed to close temp file: %v", err)
				}
				tt.path = tmpfile.Name()
			}

			file := LoadEvents(ctx, tt.path)
			if (file == nil) != tt.wantNil {
				t.Errorf("LoadEvents() file = %v, wantNil %v", file, tt.wantNil)
			}
			if file != nil {
				file.Close()
			}
		})
	}
}

func TestParseEvents(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	ctx := context.WithValue(context.Background(), key, zap.NewNop())

	tests := []struct {
		name       string
		input      string
		wantEvents []Event
		wantByTime map[string]Event
		wantByComp map[int][]Event
		wantNil    bool
	}{
		{
			name: "valid events",
			input: `[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000`,
			wantEvents: []Event{
				{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByTime: map[string]Event{
				"09:05:59.867": {Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
			},
			wantNil: false,
		},
		{
			name:       "empty file",
			input:      "",
			wantEvents: []Event{},
			wantByTime: map[string]Event{},
			wantByComp: map[int][]Event{},
			wantNil:    false,
		},
		{
			name: "multiple competitors",
			input: `[09:05:59.867] 1 1
					[09:15:00.841] 2 1 09:30:00.000
					[09:16:00.000] 1 2
					[09:17:00.000] 4 2`,
			wantEvents: []Event{
				{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				{Time: "09:16:00.000", EventID: 1, CompetitorID: 2, ExtraParams: ""},
				{Time: "09:17:00.000", EventID: 4, CompetitorID: 2, ExtraParams: ""},
			},
			wantByTime: map[string]Event{
				"09:05:59.867": {Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				"09:16:00.000": {Time: "09:16:00.000", EventID: 1, CompetitorID: 2, ExtraParams: ""},
				"09:17:00.000": {Time: "09:17:00.000", EventID: 4, CompetitorID: 2, ExtraParams: ""},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
				2: {
					{Time: "09:16:00.000", EventID: 1, CompetitorID: 2, ExtraParams: ""},
					{Time: "09:17:00.000", EventID: 4, CompetitorID: 2, ExtraParams: ""},
				},
			},
			wantNil: false,
		},
		{
			name: "invalid time format",
			input: `[invalid] 1 1
[09:15:00.841] 2 1 09:30:00.000`,
			wantEvents: []Event{
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByTime: map[string]Event{
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
			},
			wantNil: false,
		},
		{
			name: "invalid event ID",
			input: `[09:05:59.867] invalid 1
[09:15:00.841] 2 1 09:30:00.000`,
			wantEvents: []Event{
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByTime: map[string]Event{
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
			},
			wantNil: false,
		},
		{
			name:       "non-existent file",
			input:      "",
			wantEvents: []Event{},
			wantByTime: map[string]Event{},
			wantByComp: map[int][]Event{},
			wantNil:    false,
		},
		{
			name: "empty lines and spaces",
			input: `
[09:05:59.867] 1 1

[09:15:00.841] 2 1 09:30:00.000
`,
			wantEvents: []Event{
				{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByTime: map[string]Event{
				"09:05:59.867": {Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
			},
			wantNil: false,
		},
		{
			name: "invalid competitor ID",
			input: `[09:05:59.867] 1 invalid
[09:15:00.841] 2 1 09:30:00.000`,
			wantEvents: []Event{
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByTime: map[string]Event{
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
			},
			wantNil: false,
		},
		{
			name: "missing fields",
			input: `[09:05:59.867] 1
[09:15:00.841] 2 1 09:30:00.000`,
			wantEvents: []Event{
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByTime: map[string]Event{
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
				},
			},
			wantNil: false,
		},
		{
			name: "multiple extra params",
			input: `[09:05:59.867] 1 1 some extra params
[09:15:00.841] 2 1 09:30:00.000 another param`,
			wantEvents: []Event{
				{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: "some extra params"},
				{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000 another param"},
			},
			wantByTime: map[string]Event{
				"09:05:59.867": {Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: "some extra params"},
				"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000 another param"},
			},
			wantByComp: map[int][]Event{
				1: {
					{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: "some extra params"},
					{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000 another param"},
				},
			},
			wantNil: false,
		},
		{
			name: "many events",
			input: func() string {
				var builder strings.Builder
				for i := 0; i < 60; i++ {
					fmt.Fprintf(&builder, "[09:%02d:00.000] %d %d\n", i, i+1, (i%2)+1)
				}
				return builder.String()
			}(),
			wantEvents: func() []Event {
				events := make([]Event, 60)
				for i := 0; i < 60; i++ {
					events[i] = Event{
						Time:         fmt.Sprintf("09:%02d:00.000", i),
						EventID:      i + 1,
						CompetitorID: (i % 2) + 1,
						ExtraParams:  "",
					}
				}
				return events
			}(),
			wantByTime: func() map[string]Event {
				byTime := make(map[string]Event)
				for i := 0; i < 60; i++ {
					event := Event{
						Time:         fmt.Sprintf("09:%02d:00.000", i),
						EventID:      i + 1,
						CompetitorID: (i % 2) + 1,
						ExtraParams:  "",
					}
					byTime[event.Time] = event
				}
				return byTime
			}(),
			wantByComp: func() map[int][]Event {
				byComp := make(map[int][]Event)
				for i := 0; i < 60; i++ {
					event := Event{
						Time:         fmt.Sprintf("09:%02d:00.000", i),
						EventID:      i + 1,
						CompetitorID: (i % 2) + 1,
						ExtraParams:  "",
					}
					byComp[event.CompetitorID] = append(byComp[event.CompetitorID], event)
				}
				return byComp
			}(),
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var file *os.File
			if tt.name != "non-existent file" {
				tmpfile, err := os.CreateTemp("", "events*.txt")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				defer os.Remove(tmpfile.Name())
				if _, err := tmpfile.WriteString(tt.input); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				if err := tmpfile.Close(); err != nil {
					t.Fatalf("Failed to close temp file: %v", err)
				}
				file, err = os.Open(tmpfile.Name())
				if err != nil {
					t.Fatalf("Failed to open temp file: %v", err)
				}
				defer file.Close()
			}

			store := ParseEvents(ctx, file)
			if (store == nil) != tt.wantNil {
				t.Errorf("ParseEvents() store = %v, wantNil %v", store, tt.wantNil)
			}
			if store == nil {
				if len(tt.wantEvents) > 0 {
					t.Errorf("ParseEvents() returned nil, but want events: %v", tt.wantEvents)
				}
				return
			}

			if !slices.Equal(store.events, tt.wantEvents) {
				t.Errorf("ParseEvents() events = %v, want %v", store.events, tt.wantEvents)
			}

			byTime := store.ByTime()
			if !reflect.DeepEqual(byTime, tt.wantByTime) {
				t.Errorf("ByTime() = %v, want %v", byTime, tt.wantByTime)
			}

			byComp := store.ByCompetitor()
			if !reflect.DeepEqual(byComp, tt.wantByComp) {
				t.Errorf("ByCompetitor() = %v, want %v", byComp, tt.wantByComp)
			}
		})
	}
}

func TestByTime(t *testing.T) {
	store := &EventStore{
		events: []Event{
			{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
			{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
		},
	}

	want := map[string]Event{
		"09:05:59.867": {Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
		"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
	}

	got := store.ByTime()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ByTime() = %v, want %v", got, want)
	}
}

func TestByCompetitor(t *testing.T) {
	store := &EventStore{
		events: []Event{
			{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
			{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
			{Time: "09:16:00.000", EventID: 1, CompetitorID: 2, ExtraParams: ""},
		},
	}

	want := map[int][]Event{
		1: {
			{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
			{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
		},
		2: {
			{Time: "09:16:00.000", EventID: 1, CompetitorID: 2, ExtraParams: ""},
		},
	}

	got := store.ByCompetitor()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ByCompetitor() = %v, want %v", got, want)
	}
}

func TestSortMapByKey(t *testing.T) {
	eventsMap := map[string]Event{
		"09:15:00.841": {Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000"},
		"09:05:59.867": {Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
	}

	want := []string{"09:05:59.867", "09:15:00.841"}
	got := SortMapByKey(eventsMap)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SortMapByKey() = %v, want %v", got, want)
	}
}

func TestParseEvent(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	ctx := context.WithValue(context.Background(), key, zap.NewNop())

	tests := []struct {
		name string
		line string
		want Event
	}{
		{
			name: "valid event",
			line: "[09:05:59.867] 1 1",
			want: Event{Time: "09:05:59.867", EventID: 1, CompetitorID: 1, ExtraParams: ""},
		},
		{
			name: "valid event with extra params",
			line: "[09:15:00.841] 2 1 09:30:00.000 some text",
			want: Event{Time: "09:15:00.841", EventID: 2, CompetitorID: 1, ExtraParams: "09:30:00.000 some text"},
		},
		{
			name: "invalid time format",
			line: "[invalid] 1 1",
			want: Event{},
		},
		{
			name: "missing closing bracket",
			line: "[09:05:59.867 1 1",
			want: Event{},
		},
		{
			name: "invalid event ID",
			line: "[09:05:59.867] invalid 1",
			want: Event{},
		},
		{
			name: "invalid competitor ID",
			line: "[09:05:59.867] 1 invalid",
			want: Event{},
		},
		{
			name: "missing fields",
			line: "[09:05:59.867] 1",
			want: Event{},
		},
		{
			name: "no prefix",
			line: "09:05:59.867 1 1",
			want: Event{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseEvent(ctx, tt.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
