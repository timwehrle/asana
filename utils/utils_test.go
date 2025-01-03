package utils_test

import (
	"github.com/timwehrle/asana-go"
	"github.com/timwehrle/asana/utils"
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	t.Run("Empty Date", func(t *testing.T) {
		result := utils.FormatDate(nil)
		if result != "None" {
			t.Errorf("Expected 'None', got '%s'", result)
		}
	})

	t.Run("Today", func(t *testing.T) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		date := asana.Date(today)
		result := utils.FormatDate(&date)
		if result != "Today" {
			t.Errorf("Expected 'Today', got '%s'", result)
		}
	})

	t.Run("Tomorrow", func(t *testing.T) {
		tomorrow := time.Now().Add(24 * time.Hour)
		date := asana.Date(tomorrow)
		result := utils.FormatDate(&date)
		if result != "Tomorrow" {
			t.Errorf("Expected 'Tomorrow', got '%s'", result)
		}
	})

	t.Run("Date Within a Week", func(t *testing.T) {
		date := time.Now().Add(3 * 24 * time.Hour)
		expected := date.Format("Mon")
		asanaDate := asana.Date(date)
		result := utils.FormatDate(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date After a Week", func(t *testing.T) {
		futureDate := time.Now().Add(8 * 24 * time.Hour)
		expected := futureDate.Format("Jan 02, 2006")
		asanaDate := asana.Date(futureDate)
		result := utils.FormatDate(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("Date Before Today", func(t *testing.T) {
		pastDate := time.Now().Add(8 * (-24) * time.Hour)
		expected := pastDate.Format("Jan 02, 2006")
		asanaDate := asana.Date(pastDate)
		result := utils.FormatDate(&asanaDate)
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

func TestStringToDate(t *testing.T) {
	tests := []struct {
		name    string
		dateStr string
		layout  string
		want    *asana.Date
		wantErr bool
	}{
		{
			name:    "Valid date with 2006-01-02 layout",
			dateStr: "2024-01-15",
			layout:  "2006-01-02",
			want: func() *asana.Date {
				loc := time.Now().Location()
				d := asana.Date(time.Date(2024, 1, 15, 0, 0, 0, 0, loc))
				return &d
			}(),
			wantErr: false,
		},
		{
			name:    "Valid date with different layout",
			dateStr: "15/01/2024",
			layout:  "02/01/2006",
			want: func() *asana.Date {
				loc := time.Now().Location()
				d := asana.Date(time.Date(2024, 1, 15, 0, 0, 0, 0, loc))
				return &d
			}(),
			wantErr: false,
		},
		{
			name:    "None string returns nil",
			dateStr: "None",
			layout:  "2006-01-02",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Invalid date format",
			dateStr: "2024-13-45",
			layout:  "2006-01-02",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty string",
			dateStr: "",
			layout:  "2006-01-02",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid layout",
			dateStr: "2024-01-15",
			layout:  "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.StringToDate(tt.dateStr, tt.layout)

			if (err != nil) != tt.wantErr {
				t.Errorf("StringToDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if (got == nil) != (tt.want == nil) {
				t.Errorf("StringToDate() got = %v, want %v", got, tt.want)
				return
			}

			if got == nil && tt.want == nil {
				return
			}

			gotTime := time.Time(*got)
			wantTime := time.Time(*tt.want)

			if !gotTime.Equal(wantTime) {
				t.Errorf("StringToDate() got = %v, want %v", gotTime, wantTime)
			}

			if gotTime.Hour() != 0 || gotTime.Minute() != 0 ||
				gotTime.Second() != 0 || gotTime.Nanosecond() != 0 {
				t.Errorf("Time components not zeroed: got = %v", gotTime)
			}

			if gotTime.Location().String() != time.Now().Location().String() {
				t.Errorf("Timezone mismatch: got = %v, want local", gotTime.Location())
			}
		})
	}
}
