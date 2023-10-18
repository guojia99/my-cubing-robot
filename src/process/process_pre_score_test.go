package process

import (
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
