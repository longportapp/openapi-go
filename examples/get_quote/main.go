package main

import (
	"context"
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
	// close connection
	defer quoteContext.Close()
	ctx := context.Background()
	// Get basic information of securities
	quotes, err := quoteContext.Quote(ctx, []string{"700.HK", "AAPL.US", "TSLA.US", "NFLX.US"})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("quotes: %+v\n", quotes[0])

	warrants, err := quoteContext.WarrantList(ctx, "700.HK", quote.WarrantFilter{
		SortBy:     quote.WarrantVolume,
		SortOrder:  quote.WarrantAsc,
		SortOffset: 0,
		SortCount:  10,
	}, quote.WarrantEN)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("warrants: %+v\n", warrants[0])
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
