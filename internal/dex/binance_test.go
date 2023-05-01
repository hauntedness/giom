package dex

import (
	"testing"
)

func TestGetBinancePrice(t *testing.T) {
	type args struct {
		tickerid string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test",
			args{"binance_hnt_swap_usdt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetBinancePrice(tt.args.tickerid); err != nil || got == 0 {
				t.Errorf("GetBinancePrice() = %v", got)
			}
		})
	}
}
