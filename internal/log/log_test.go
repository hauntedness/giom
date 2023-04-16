package log

import (
	"fmt"
	"testing"
)

func TestInfos(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "key is string, value is any",
			args: args{
				args: []any{"name", "zhang san", "age", 18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			Infos(tt.args.args...)
			Errors(fmt.Errorf("pretend some error"), tt.args.args...)
		})
	}
}
