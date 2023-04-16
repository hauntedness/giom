package runtime

import (
	"fmt"
	"testing"
)

func TestSource(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			err := fmt.Errorf("%s: %w", Source(), fmt.Errorf("something went wrong"))
			fmt.Println(err)
			err = fmt.Errorf("%s: %w", Source(), err)
			fmt.Println(err)
		})
	}
}
