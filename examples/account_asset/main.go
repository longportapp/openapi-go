package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/trade"
)

func main() {
	// create trade context from environment variables
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	tradeContext, err := trade.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer tradeContext.Close()
	ctx := context.Background()
	// Get AccountBalance infomation
	ab, err := tradeContext.AccountBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", ab[0])
}
