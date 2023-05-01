package main

import (
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/hauntedness/giom/internal/dex"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/giom/play/common"
	"gonum.org/v1/gonum/floats"
)

type kline struct {
	data []float64
	unix []float64
}

const length = 300

var _kline = kline{data: make([]float64, length), unix: make([]float64, length)}

func init() {
	for i := range _kline.unix {
		_kline.unix[i] = 320.0 * float64(i+1) / length
	}
	k_btc := dex.GetKline("binance_btc_swap_usdt", time.Minute)
	df_btc, err := k_btc.ToDataFrame()
	if err != nil {
		log.Errors(err)
		return
	}
	k_eth := dex.GetKline("binance_eth_swap_usdt", time.Minute)
	df_eth, err := k_eth.ToDataFrame()
	if err != nil {
		log.Errors(err)
		return
	}
	ratio := floats.DivTo(make([]float64, df_btc.Len()), df_btc.Close(), df_eth.Close())[df_btc.Len()-length:]
	ratio = Normalize(ratio, 30)
	go func() {
		for range time.Tick(time.Second * 1) {
			// log.Infos("ratio", ratio)
			floats.Reverse(ratio)
			copy(_kline.data, ratio)
		}
	}()
}

func (k kline) Layout(gtx layout.Context) layout.Dimensions {
	path := clip.Path{}
	path.Begin(gtx.Ops)
	for i := range k.data {
		x := unit.Dp(k.unix[i])
		y := unit.Dp(k.data[i])
		point := f32.Point{
			X: float32(x),
			Y: float32(y),
		}
		path.LineTo(point)
	}
	spec := path.End()
	paint.FillShape(
		gtx.Ops,
		common.ColorRed,
		clip.Stroke{
			Path:  spec,
			Width: 1,
		}.Op(),
	)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func Normalize(data []float64, weight float64) (ret []float64) {
	ret = make([]float64, len(data))
	min := floats.Min(data)
	max := floats.Max(data)
	total := max - min
	for i := range data {
		ret[i] = (data[i] - min) / total * weight
	}
	return
}


