package dex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/httputil"
)

func GetKline(tickerid string, period time.Duration) *KLine {
	if strings.Contains(tickerid, "_usdt") {
		v := url.Values{}
		v.Add("tickerid", tickerid)
		v.Add("period", strconv.Itoa(int(period.Minutes())))
		v.Add("reach", strconv.Itoa(int(time.Now().Unix())))
		v.Add("since", "")
		v.Add("utc", "0")
		v.Add("webp", "1")

		url := "https://dncapi.xxbitcoin.xyz/api/v1/kline/market?" + v.Encode()

		kline := &KLine{}

		bytes, err := httputil.Request(http.MethodGet, url, nil, nil)
		if err != nil {
			log.Errors(err, "response", string(bytes))
		}
		err = json.Unmarshal(bytes, kline)
		if err != nil {
			log.Errors(err)
		}
		return kline
	} else {
		return nil
	}
}

func GetPrice(tickerid string) (float64, error) {
	if strings.Contains(tickerid, "_usdt") {
		return GetBinancePrice(tickerid)
	} else {
		return 0, fmt.Errorf("only support usdt tickerid")
	}
}

func BinanceFutures(n int) chan string {
	tar := make(chan string, n)
	wg := sync.WaitGroup{}
	for i := 1; (i-1)*100 <= n; i++ {
		wg.Add(1)
		go page(i, tar, &wg)
	}
	wg.Wait()
	log.Infos("symbols", len(tar))
	return tar
}

func page(i int, dst chan string, wg *sync.WaitGroup) {
	value := url.Values{}
	value.Add("page", strconv.Itoa(i))
	value.Add("exchangecode", "binance")
	value.Add("pagesize", "100")
	value.Add("webp", "1")
	url := "https://dncapi.xxbitcoin.xyz/api/v3/futuresmarket/exchange/futures/market?" + value.Encode()
	bytes, err := httputil.Get(url, nil)
	if err != nil {
		panic(err)
	}
	var res PriceResponse
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		panic(err)
	}
	for _, list := range res.Data.List {
		for _, m := range list.Markets {
			if m.Pair2 == "USDT" && strings.Contains(m.Tickerid, "_swap_") {
				dst <- m.Tickerid
			}
		}
	}
	wg.Done()
}
