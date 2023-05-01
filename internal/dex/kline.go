package dex

import (
	"errors"

	"gonum.org/v1/gonum/mat"
)

/*
	{
	    "data": {
	        "kline": [
	            [
	                1603065600,  //时间
	                39.7777, // 最高价
	                32.3569, // 开盘
	                20.6023, // 最低价
	                38.6663, // 收盘
	                5387902.82,   // 成交量
	                154196745.159039 // 交易额
	            ]
	        ]
	    },
	    "code": 200,
	    "msg": "success"
	}
*/
type KLine struct {
	Data Data   `json:"data"`
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type Data struct {
	Kline [][]float64 `json:"kline"`
}

func (k *KLine) ToDataFrame() (DataFrame, error) {
	data_ := k.Data.Kline
	r := len(data_)
	if r == 0 {
		return DataFrame{}, errors.New("no data")
	}
	c := len(data_[0])
	records := make([]float64, 0, r*c)
	for _, v := range data_ {
		if len(v) == 7 {
			records = append(records, v[0])
			records = append(records, v[1])
			records = append(records, v[2])
			records = append(records, v[3])
			records = append(records, v[4])
			records = append(records, v[5])
			records = append(records, v[6])
		}
	}

	return DataFrame{
		Dense: mat.NewDense(r, c, records),
	}, nil
}

func (k *KLine) IsEmpty() bool {
	return k == nil || len(k.Data.Kline) == 0
}
