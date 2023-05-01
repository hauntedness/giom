package main

import (
	"context"
	"strings"
	"time"

	"github.com/hauntedness/giom/service/starknet"
	"github.com/hauntedness/wechatbot/wechat"
)

func main() {
	Tick()
	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
		Tick()
	}
}

func Tick() {
	fees, err := starknet.EstimateFees()
	if err != nil {
		_ = wechat.Send(err.Error(), nil, context.Background())
		return
	}
	var sb strings.Builder
	for k, v := range fees {
		sb.WriteString(k)
		sb.WriteRune(':')
		sb.WriteString(v.String())
		sb.WriteRune('\n')
	}
	_ = wechat.Send(sb.String(), nil, context.TODO())
}
