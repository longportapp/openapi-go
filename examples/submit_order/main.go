package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridgeapp/openapi-go/trade"
)

func main() {
	// create trade context from environment variables
	tradeContext, err := trade.NewFormEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer tradeContext.Close()
	ctx := context.Background()
	// submit order
	order := &trade.SubmitOrder{
		Symbol:            "700.HK",
		OrderType:         trade.OrderTypeLO,
		Side:              trade.OrderSideBuy,
		SubmittedQuantity: 200,
		TimeInForce:       trade.TimeTypeDay,
		SubmittedPrice:    "10",
	}
	orderId, err := tradeContext.SubmitOrder(ctx, order)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("orderId: %v", orderId)
}
