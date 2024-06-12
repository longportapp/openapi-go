package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/trade"
	"github.com/shopspring/decimal"
)

func main() {
	// create trade context from environment variables
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	tradeContext, err := trade.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer tradeContext.Close()
	ctx := context.Background()

	// subscribe order status
	tradeContext.OnTrade(func(ev *trade.PushEvent) {
		// handle order changing event
		log.Printf("order event: %+v\n", ev)
	})

	// submit order
	order := &trade.SubmitOrder{
		Symbol:            "700.HK",
		OrderType:         trade.OrderTypeLO,
		Side:              trade.OrderSideBuy,
		SubmittedQuantity: 200,
		TimeInForce:       trade.TimeTypeDay,
		SubmittedPrice:    decimal.NewFromFloat(12),
	}
	orderId, err := tradeContext.SubmitOrder(ctx, order)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("orderId: %v\n", orderId)
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
