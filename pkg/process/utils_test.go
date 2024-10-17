package process

import (
	"reflect"
	"testing"
)

func TestCutMsgWithFields(t *testing.T) {
	type args struct {
		input string
		cut   string
	}
	tests := []struct {
		name       string
		args       args
		wantHeader string
		wantTitle  string
		wantValues []string
	}{
		{
			name: "1",
			args: args{
				input: "录入[选手] 1 2 3 4 5",
				cut:   " ",
			},
			wantHeader: "录入",
			wantTitle:  "选手",
			wantValues: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "2",
			args: args{
				input: "录入[选手] 1,2,3,4,5",
				cut:   ",",
			},
			wantHeader: "录入",
			wantTitle:  "选手",
			wantValues: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "3",
			args: args{
				input: "录入 1,2,3,4,5",
				cut:   ",",
			},
			wantHeader: "录入",
			wantTitle:  "",
			wantValues: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "4",
			args: args{
				input: "录入  1,2,3,4,5",
				cut:   ",",
			},
			wantHeader: "录入",
			wantTitle:  "",
			wantValues: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "5",
			args: args{
				input: "pk[WCA项目] 徐永浩女装vs女装徐永浩",
				cut:   "vs",
			},
			wantHeader: "pk",
			wantTitle:  "WCA项目",
			wantValues: []string{"徐永浩女装", "女装徐永浩"},
		},
		{
			name: "6",
			args: args{
				input: "录入 clock 10.35 14.41 12.24 9.90 10.69",
				cut:   " ",
			},
			wantHeader: "录入",
			wantTitle:  "",
			wantValues: []string{"clock", "10.35", "14.41", "12.24", "9.90", "10.69"},
		},
		{
			name: "7",
			args: args{
				input: " pk 兔兔vs嘉",
				cut:   "vs",
			},
			wantHeader: "pk",
			wantTitle:  "",
			wantValues: []string{"兔兔", "嘉"},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotHeader, gotTitle, gotValues := CutMsgWithFields(tt.args.input, tt.args.cut)
				if gotHeader != tt.wantHeader {
					t.Errorf("CutMsgWithFields() gotHeader = %v, want %v", gotHeader, tt.wantHeader)
				}
				if gotTitle != tt.wantTitle {
					t.Errorf("CutMsgWithFields() gotTitle = %v, want %v", gotTitle, tt.wantTitle)
				}
				if !reflect.DeepEqual(gotValues, tt.wantValues) {
					t.Errorf("CutMsgWithFields() gotValues = %v, want %v", gotValues, tt.wantValues)
				}
			},
		)
	}
}
