package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridgeapp/openapi-go/trade"
)

func main() {
	tradeContext, err := trade.NewFormEnv()
	if err != nil {
		log.Fatal(err)
	}
	defer tradeContext.Close()
	ctx := context.Background()
	ab, err := tradeContext.AccountBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", ab[0])
}
