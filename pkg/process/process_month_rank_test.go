package process

import (
	"reflect"
	"testing"
	"time"
)

func TestMRank_Do(t *testing.T) {
	tests := []struct {
		input string
		want  time.Time
	}{
		{
			input: "20230101",
			want:  time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input: "2023",
			want:  time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input: "09",
			want:  time.Date(time.Now().Year(), 9, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input: "202305",
			want:  time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.input, func(t *testing.T) {
				if got := getStringTime(tt.input); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getStringTime() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
