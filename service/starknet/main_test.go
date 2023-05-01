package starknet

import (
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/hauntedness/giom/internal/log"
	"github.com/hauntedness/httputil"
	"github.com/shopspring/decimal"
)

func TestGetTransaction(t *testing.T) {
	url := `https://alpha-mainnet.starknet.io/feeder_gateway/get_transaction?transactionHash=0x555b5dd91013af5d318e6377f783cc991f9ac86b1e3a6f0f0cbf392285affc2`
	data, err := httputil.Get(url, nil)
	if err != nil {
		t.Error(err)
		return
	}
	log.Infos("data", string(data))
}

func TestGetBlock(t *testing.T) {
	url := `https://alpha-mainnet.starknet.io/feeder_gateway/get_block?blockHash=0x6cf3d163c08f4d72c3bd31105ed976afbbcace1928772ae6c0673d708f5dc01`
	data, err := httputil.Get(url, nil)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(string(data))
}

func TestTransactions(t *testing.T) {
	url := "https://api.starkscan.co/api/v0/transactions"

	data, err := httputil.Get(url, httputil.H{"x-api-key": "docs-starkscan-co-api"})
	// req.Header.Add("accept", "application/json")
	if err != nil {
		t.Error(err)
		return
	}
	response := &GetTransactionsResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		t.Error(err)
		return
	}
	for _, d := range response.Data {
		if d.ActualFee == nil {
			log.Info("none_fee", "hash", d.TransactionHash, "type", d.TransactionType, "status", d.TransactionStatus)
			continue
		}
		var selector strings.Builder
		for i, ac := range d.AccountCalls {
			selector.WriteString(ac.SelectorName)
			if i < len(d.AccountCalls)-1 {
				selector.WriteString(",")
			}
		}
		select_name := selector.String()
		log.Info("fee", "actual_fee_in_wei", d.ActualFee, "actual_fee_in_eth", d.GetActualGasInEth(), "selector_name", select_name)
	}
}

func TestConvert(t *testing.T) {
	value := "509186408520024"
	_, ok := big.NewInt(0).SetString(value, 10)
	if !ok {
		t.Errorf("could not get numeric wei from value: %s", value)
		return
	}
	number := decimal.RequireFromString(value)
	unitEth := decimal.NewFromFloat(params.Ether)
	log.Info(number.String())
	ether := number.Div(unitEth)
	log.Info(ether.String())
}

func TestUnmarshal(t *testing.T) {
	data, err := os.ReadFile("testdata/transactions_response.json")
	if err != nil {
		t.Error(err)
		return
	}
	response := &GetTransactionsResponse{}
	err = json.Unmarshal(data, response)
	if err != nil {
		t.Error(err)
		return
	}
	for _, d := range response.Data {
		if d.ActualFee == nil {
			log.Info("none_fee", "hash", d.TransactionHash, "type", d.TransactionType, "status", d.TransactionStatus)
			continue
		}
		var selector strings.Builder
		for i, ac := range d.AccountCalls {
			selector.WriteString(ac.SelectorName)
			if i < len(d.AccountCalls)-1 {
				selector.WriteString(",")
			}
		}
		select_name := selector.String()
		log.Info("fee", "actual_fee_in_wei", d.ActualFee, "actual_fee_in_eth", d.GetActualGasInEth(), "selector_name", select_name)
	}
}
