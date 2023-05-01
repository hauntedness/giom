package dex

import (
	"fmt"
	"time"

	"gonum.org/v1/gonum/mat"
)

const (
	Unix = iota // unix second
	High
	Open
	Low
	Close
	Quantity
	Volume
)
const Columns = 7

type DataFrame struct {
	*mat.Dense
}

func (df *DataFrame) Len() int {
	if df.Dense == nil {
		return 0
	}
	r, _ := df.Dims()
	return r
}

func (df *DataFrame) Unix() []float64 {
	return mat.Col(nil, Unix, df)
}

func (df *DataFrame) UnixAt(i int) float64 {
	return df.At(i, Unix)
}

func (df *DataFrame) High() []float64 {
	return mat.Col(nil, High, df)
}

func (df *DataFrame) HighAt(i int) float64 {
	return df.At(i, High)
}

func (df *DataFrame) Open() []float64 {
	return mat.Col(nil, Open, df)
}

func (df *DataFrame) OpenAt(i int) float64 {
	return df.At(i, Open)
}

func (df *DataFrame) Low() []float64 {
	return mat.Col(nil, Low, df)
}

func (df *DataFrame) LowAt(i int) float64 {
	return df.At(i, Low)
}

func (df *DataFrame) Close() []float64 {
	return mat.Col(nil, Close, df)
}

func (df *DataFrame) CloseAt(i int) float64 {
	return df.At(i, Close)
}

func (df *DataFrame) Quantity() []float64 {
	return mat.Col(nil, Quantity, df)
}

func (df *DataFrame) QuantityAt(i int) float64 {
	return df.At(i, Quantity)
}

// 交易额
func (df *DataFrame) Volume() []float64 {
	return mat.Col(nil, Volume, df)
}

func (df *DataFrame) VolumeAt(i int) float64 {
	return df.At(i, Volume)
}

func (df *DataFrame) Print() {
	formatter := mat.Formatted(df.Dense, mat.FormatPython(), mat.Squeeze())
	fmt.Printf("%#10f\n", formatter)
}

//go:generate goparrot -p "Unix;High;Open;Low;Close;Quantity;Volume;"
func (df *DataFrame) SetUnixAt(i int, v float64) {
	df.Set(i, Unix, v)
}

func (df *DataFrame) SetHighAt(i int, v float64) {
	df.Set(i, High, v)
}

func (df *DataFrame) SetOpenAt(i int, v float64) {
	df.Set(i, Open, v)
}

func (df *DataFrame) SetLowAt(i int, v float64) {
	df.Set(i, Low, v)
}

func (df *DataFrame) SetCloseAt(i int, v float64) {
	df.Set(i, Close, v)
}

func (df *DataFrame) SetQuantityAt(i int, v float64) {
	df.Set(i, Quantity, v)
}

func (df *DataFrame) SetVolumeAt(i int, v float64) {
	df.Set(i, Volume, v)
}

func (df *DataFrame) Slice(i, k, j, l int) *DataFrame {
	return &DataFrame{Dense: df.Dense.Slice(i, k, j, l).(*mat.Dense)}
}

func (df *DataFrame) SearchByTime(t time.Time) int {
	var v []float64 = df.Unix()
	if len(v) == 0 {
		panic("no element")
	}
	if len(v) == 1 {
		return 0
	}
	return search(0, len(v), v, t)
}

func search(start, end int, v []float64, t time.Time) int {
	if start == len(v) {
		return start - 1
	}
	if end-start < 2 {
		return start
	}
	var mid int = (end-start)/2 + start
	value := time.Unix(int64(v[mid]), 0)
	if value.Before(t) {
		return search(mid, end, v, t)
	} else {
		return search(start, mid, v, t)
	}
}
