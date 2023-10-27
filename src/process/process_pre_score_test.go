package process

import (
	"reflect"
	"testing"

	"github.com/guojia99/my-cubing-core/model"
)

func Test__getProject(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want model.Project
	}{
		{
			name: "ok",
			in:   "333 1,2,3,45,6",
			want: model.Cube333,
		},
		{
			name: "has empty",
			in:   "   555 1, 2,3,4,5",
			want: model.Cube555,
		},
		{
			name: "empty",
			in:   "",
			want: model.Project(""),
		},
		{
			name: "empty 4",
			in:   "    ",
			want: model.Project(""),
		},
		{
			name: "二阶五魔",
			in:   "二阶五魔 1, 2, 3",
			want: model.XCube222Minx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := _getProject(tt.in); got != tt.want {
				t.Errorf("_getProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__getNumbers(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "11",
			args: args{
				in: "12 abcs sd1231 da3213 231 31.1231 12313.1",
			},
			want: []float64{
				12, 1231, 3213, 231, 31.1231, 12313.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := _getNumbers(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_getNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
