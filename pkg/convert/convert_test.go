package convert

import (
	"github.com/timwehrle/asana-api"
	"testing"
	"time"
)

func TestToDate(t *testing.T) {
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
			layout:  time.DateOnly,
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
			layout:  time.DateOnly,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Invalid date format",
			dateStr: "2024-13-45",
			layout:  time.DateOnly,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty string",
			dateStr: "",
			layout:  time.DateOnly,
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
			got, err := ToDate(tt.dateStr, tt.layout)

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
