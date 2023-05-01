package eth

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
)

var (
	Wei   = decimal.NewFromInt(params.Wei)
	GWei  = decimal.NewFromInt(params.GWei)
	Ether = decimal.NewFromInt(params.Ether)
)

func WeiToEther(wei decimal.Decimal) decimal.Decimal {
	return wei.DivRound(Ether, 20)
}

func EtherToWei(ether decimal.Decimal) decimal.Decimal {
	return ether.Mul(Ether)
}
