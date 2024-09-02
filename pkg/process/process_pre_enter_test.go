package process

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/guojia99/my-cubing-core/model"
)

func TestPreEnter_CheckPrefix(t *testing.T) {
	in := "三阶 D D D D D"

	c := &PreEnter{}

	fmt.Println(c.CheckPrefix(in))
}

func Test__getProject(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name  string
		args  args
		want  model.Project
		want1 int
	}{
		{
			name: "1",
			args: args{
				in: "333[2] 1 2 3 4 5",
			},
			want:  model.Cube333,
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, got1 := _getProject(tt.args.in)
				if got != tt.want {
					t.Errorf("_getProject() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("_getProject() got1 = %v, want %v", got1, tt.want1)
				}
			},
		)
	}
}

func Test__getResults(t *testing.T) {
	type args struct {
		in string
		pj model.Project
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "xxx1",
			args: args{
				in: "444[2] 1 2 3 4 5",
				pj: model.Cube444,
			},
			want:    []string{"1", "2", "3", "4", "5"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := _getResults(tt.args.in, tt.args.pj)
				if (err != nil) != tt.wantErr {
					t.Errorf("_getResults() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("_getResults() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
