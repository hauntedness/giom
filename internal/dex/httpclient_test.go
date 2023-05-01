package dex

import (
	"testing"
	"time"

	"github.com/hauntedness/giom/internal/log"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

func TestGetTopN(t *testing.T) {
	k := BinanceFutures(2000)
	if k == nil {
		t.Fail()
	}
	set := make(map[string]struct{})
	for {
		select {
		case tickerid := <-k:
			if _, ok := set[tickerid]; !ok {
				t.Log(tickerid)
				set[tickerid] = struct{}{}
			} else {
				t.Fail()
			}
		default:
			return
		}
	}
}

func TestGetPrice(t *testing.T) {
	tickerid := "binance_fil_swap_usdt"
	f, err := GetPrice(tickerid)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("get price of %v %v", tickerid, f)
	if f <= 0 {
		t.Fail()
	}
}

func TestGetKline0(t *testing.T) {
	type args struct {
		tickerid string
		period   time.Duration
	}
	tests := []struct {
		name string
		args args
		want *KLine
	}{
		{
			name: "test",
			args: args{
				tickerid: "binance_fil_swap_usdt",
				period:   time.Minute,
			},
			want: &KLine{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetKline(tt.args.tickerid, tt.args.period)
			unix := got.Data.Kline[0][0]
			t.Logf("start time:%s", time.Unix(int64(unix), 0))
			t.Logf("length of %v:%v", tt.args.tickerid, len(got.Data.Kline))
			df, err := got.ToDataFrame()
			if err != nil {
				t.Fail()
			}
			// y = alpha + beta * x
			alpha, beta := LinearRegression(df.Close())
			t.Logf("alpha %f, beta:%f", alpha, beta)
			t.Log(df.Close())
		})
	}
}

func TestGetKline(t *testing.T) {
	type args struct {
		tickerid string
		period   time.Duration
	}
	type TestItem struct {
		name string
		args args
	}
	tests := []TestItem{
		{
			name: "1m",
			args: args{
				tickerid: "binance_btc_swap_usdt",
				period:   time.Minute,
			},
		},
		{
			name: "15m",
			args: args{
				tickerid: "binance_btc_swap_usdt",
				period:   time.Minute * 15,
			},
		},
		{
			name: "30m",
			args: args{
				tickerid: "binance_btc_swap_usdt",
				period:   time.Hour / 2,
			},
		},
		{
			name: "1h",
			args: args{
				tickerid: "binance_btc_swap_usdt",
				period:   time.Hour,
			},
		},

		{
			name: "1 day",
			args: args{
				tickerid: "binance_btc_swap_usdt",
				period:   time.Hour * 24,
			},
		},
		{
			name: "1 month",
			args: args{
				tickerid: "binance_btc_swap_usdt",
				period:   time.Hour * 24 * 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetKline(tt.args.tickerid, tt.args.period)
			if got == nil || len(got.Data.Kline) == 0 {
				t.Errorf("GetKline() = %v", got)
			} else if df, err := got.ToDataFrame(); err == nil {
				t.Logf("start date: %v", time.Unix(int64(df.UnixAt(0)), 0))
			} else {
				t.Error(err)
			}
		})
	}
}

// volume, 线性回归测量斜率beta,
//
// x := [0,1,2,3,4,5,6...n]
//
// weight := [0,0,0,0,0...0]
func LinearRegression(y []float64) (alpha, beta float64) {
	x := make([]float64, 0, len(y))
	for i := range y {
		x = append(x, float64(i))
	}
	alpha, beta = stat.LinearRegression(x, y, nil, false)
	return
}

func TestBtcAndEth(t *testing.T) {
	k_btc := GetKline("binance_btc_swap_usdt", time.Minute)
	df_btc, err := k_btc.ToDataFrame()
	if err != nil {
		t.Error(err)
		return
	}
	k_eth := GetKline("binance_eth_swap_usdt", time.Minute)
	df_eth, err := k_eth.ToDataFrame()
	if err != nil {
		t.Error(err)
		return
	}
	result := floats.DivTo(make([]float64, df_btc.Len()), df_btc.Close(), df_eth.Close())
	log.Infos("result", result)
}
