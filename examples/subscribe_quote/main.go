package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/quote"
)

func main() {
	// create quote context from environment variables
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	quoteContext, err := quote.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer quoteContext.Close()
	ctx := context.Background()
	quoteContext.OnQuote(func(pe *quote.PushQuote) {
		bytes, _ := json.Marshal(pe)
		fmt.Println(string(bytes))
	})
	// Subscribe some symbols
	err = quoteContext.Subscribe(ctx, []string{"700.HK"}, []quote.SubType{quote.SubTypeBrokers, quote.SubTypeDepth, quote.SubTypeTrade, quote.SubTypeQuote}, true)
	if err != nil {
		log.Fatal(err)
		return
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
