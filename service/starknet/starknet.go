package starknet

import (
	"encoding/json"
	"strings"

	"github.com/hauntedness/giom/internal/eth"
	"github.com/hauntedness/httputil"
	"github.com/shopspring/decimal"
)

type Fees map[string]decimal.Decimal

func EstimateFees() (Fees, error) {
	url := "https://api.starkscan.co/api/v0/transactions"

	data, err := httputil.Get(url, httputil.H{"x-api-key": "docs-starkscan-co-api"})
	// req.Header.Add("accept", "application/json")
	if err != nil {
		return nil, err
	}
	response := &GetTransactionsResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}
	var sum Fees = make(Fees)
	var count map[string]int64 = make(map[string]int64)
	for _, d := range response.Data {
		// unit in wei
		if d.ActualFee == nil || d.ActualFee.Cmp(decimal.Zero) == 0 {
			continue
		}
		var sb strings.Builder
		sb.WriteRune('[')
		for i, ac := range d.AccountCalls {
			sb.WriteString(ac.SelectorName)
			if i < len(d.AccountCalls)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteRune(']')
		selectors := sb.String()
		count[selectors] = count[selectors] + 1
		sum[selectors] = sum[selectors].Add(*d.ActualFee)

	}
	for k, v := range sum {
		wei := v.Div(decimal.NewFromInt(count[k]))
		sum[k] = eth.WeiToEther(wei)
	}
	return sum, nil
}
