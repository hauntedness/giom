package page

import (
	"reflect"
	"testing"
)

func TestFrom(t *testing.T) {
	type args struct {
		text   string
		offset int
		size   int
	}
	tests := []struct {
		name string
		args args
		want Page
	}{
		{
			name: "normal",
			args: args{
				text:   "12\n34\n56\n78\n9a\nbc\nde\nf",
				offset: 3,
				size:   2,
			},
			want: Page{
				Lines:  []string{"78", "9a"},
				Number: 1,
				Total:  8,
				Size:   2,
			},
		},
		{
			name: "normal",
			args: args{
				//
				text:   "12\n34\n56\n78\n9a\nbc\nde\nf",
				offset: 5,
				size:   2,
			},
			want: Page{
				Lines:  []string{"bc", "de"},
				Number: 2,
				Total:  8,
				Size:   2,
			},
		},
		{
			name: "normal",
			args: args{
				//
				text:   "12\n34\n56\n78\n9a\nbc\nde\nf",
				offset: 4,
				size:   5,
			},
			want: Page{
				Lines:  []string{"9a", "bc", "de", "f"},
				Number: 0,
				Total:  8,
				Size:   5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := From(tt.args.text, tt.args.offset, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("From() = %v, want %v", got, tt.want)
			}
		})
	}
}
