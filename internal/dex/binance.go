package dex

import (
	"encoding/json"
	"fmt"

	"github.com/hauntedness/giom/internal/runtime"
	"github.com/hauntedness/httputil"
)

type PriceResponse struct {
	Data Data_  `json:"data"`
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type Data_ struct {
	TotalCount int64  `json:"total_count"`
	TotalPages int64  `json:"total_pages"`
	Page       int64  `json:"page"`
	List       []List `json:"list"`
}

type List struct {
	Code           string   `json:"code"`
	Symbol         string   `json:"symbol"`
	Name           string   `json:"name"`
	Price          float64  `json:"price"`
	Contracttype   string   `json:"contracttype"`
	Contractvalue  string   `json:"contractvalue"`
	Minunitchange  string   `json:"minunitchange"`
	Leveragechoose string   `json:"leveragechoose"`
	Openingcharge  string   `json:"openingcharge"`
	Closingcharge  string   `json:"closingcharge"`
	Tradingtime    string   `json:"tradingtime"`
	Fundsrate      float64  `json:"fundsrate"`
	OpenInterest   float64  `json:"open_interest"`
	Markets        []Market `json:"martkets"`
}

type Market struct {
	Tickerid     string  `json:"tickerid"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Pair1        string  `json:"pair1"`
	Pair2        Pair2   `json:"pair2"`
	Title        string  `json:"title"`
	Logo         string  `json:"logo"`
	Price        float64 `json:"price"`
	Changerate   float64 `json:"changerate"`
	Amount       float64 `json:"amount"`
	OpenInterest float64 `json:"open_interest"`
	Refreshtime  int64   `json:"refreshtime"`
	IsFocus      int64   `json:"is_focus"`
	Showkline    int64   `json:"showkline"`
}

type Pair2 string

const (
	Usd  Pair2 = "USD"
	Usdt Pair2 = "USDT"
)

func GetBinancePrice(tickerid string) (float64, error) {
	url := "https://dncapi.xxbitcoin.xyz/api/v3/futuresmarket/exchange/futures/market?exchangecode=binance&page=1&pageSize=1000&webp=1"
	bytes, err := httputil.Get(url, nil)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", runtime.Source(), err)
	}
	var value PriceResponse
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", runtime.Source(), err)
	}
	for _, v := range value.Data.List {
		for _, v2 := range v.Markets {
			if v2.Tickerid == tickerid {
				return v2.Price, nil
			}
		}
	}
	return 0, fmt.Errorf("%s, price not found %s", runtime.Source(), tickerid)
}
