# Longbridge OpenAPI SDK for Go

`longbridge` provides an easy-to-use interface for invokes [`Longbridge OpenAPI`](https://open.longbridgeapp.com/en/).

## Quickstart

_With Go module support , simply add the following import_

```golang
import "github.com/longbridgeapp/openapi-go"
```

_Setting environment variables(MacOS/Linux)_

```bash
export LONGBRIDGE_APP_KEY="App Key get from user center"
export LONGBRIDGE_APP_SECRET="App Secret get from user center"
export LONGBRIDGE_ACCESS_TOKEN="Access Token get from user center"
```

_Setting environment variables(Windows)_

```powershell
setx LONGBRIDGE_APP_KEY "App Key get from user center"
setx LONGBRIDGE_APP_SECRET "App Secret get from user center"
setx LONGBRIDGE_ACCESS_TOKEN "Access Token get from user center"
```

## Quote API (Get basic information of securities)

```golang
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridgeapp/openapi-go/quote"
)

func main() {
	// create quote context from environment variables
	quoteContext, err := quote.NewFormEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer quoteContext.Close()
	ctx := context.Background()
	// Get basic information of securities
	quotes, err := quoteContext.Quote(ctx, []string{"700.HK", "AAPL.US", "TSLA.US", "NFLX.US"})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("quotes: %v", quotes)
}
```

## Trade API (Submit order)

```golang
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridgeapp/openapi-go/trade"
	"github.com/shopspring/decimal"
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
		SubmittedPrice:    decimal.NewFromFloat(12),
	}
	orderId, err := tradeContext.SubmitOrder(ctx, order)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("orderId: %v\n", orderId)
}
```

## Environment Variables

-- support load env from `.env` file

| name                                 | description                                                      | default value                            | example |
|--------------------------------------|------------------------------------------------------------------|------------------------------------------|---------|
| LONGBRIDGE_HTTP_URL                  | longbridge rest api url                                          | https://openapi.longbridgeapp.com        |         |
| LONGBRIDGE_APP_KEY                   | app key                                                          |                                          |         |
| LONGBRIDGE_APP_SECRET                | app secret                                                       |                                          |         |
| LONGBRIDGE_ACCESS_TOKEN              | access token                                                     |                                          |         |
| LONGBRIDGE_TRADE_URL                 | longbridge protocol url for trade context                        | wss://openapi-trade.longbridgeapp.com/v2 |         |
| LONGBRIDGE_QUOTE_URL                 | longbridge protocol url for quote context                        | wss://openapi-quote.longbridgeapp.com/v2 |         |
| LONGBRIDGE_LOG_LEVEL                 | log level                                                        | info                                     |         |
| LONGBRIDGE_TRADE_LB_AUTH_TIMEOUT     | longbridge protocol authorize request time out for trade context | 10 second                                | 10s     |
| LONGBRIDGE_TRADE_LB_TIMEOUT          | longbridge protocol dial timeout for trade context               | 5 second                                 | 6s      |
| LONGBRIDGE_TRADE_LB_WRITE_QUEUE_SIZE | longbirdge protocol write queue size for trade context           | 16                                       |         |
| LONGBRIDGE_TRADE_LB_READ_QUEUE_SIZE  | longbirdge protocol read queue size for trade context            | 16                                       |         |
| LONGBRIDGE_TRADE_LB_READ_BUFFER_SIZE | longbirdge protocol read buffer size for trade context           | 4096                                     |         |
| LONGBRIDGE_TRADE_LB_MIN_GZIP_SIZE    | longbirdge protocol minimal gzip size for trade context          | 1024                                     |         |
| LONGBRIDGE_QUOTE_LB_AUTH_TIMEOUT     | longbridge protocol authorize request time out for quote context | 10 second                                | 10s     |
| LONGBRIDGE_QUOTE_LB_TIMEOUT          | longbridge protocol dial timeout for quote context               | 5 second                                 | 6s      |
| LONGBRIDGE_QUOTE_LB_WRITE_QUEUE_SIZE | longbirdge protocol write queue size for quote context           | 16                                       |         |
| LONGBRIDGE_QUOTE_LB_READ_QUEUE_SIZE  | longbirdge protocol read queue size for quote context            | 16                                       |         |
| LONGBRIDGE_QUOTE_LB_READ_BUFFER_SIZE | longbirdge protocol read buffer size for quote context           | 4096                                     |         |
| LONGBRIDGE_QUOTE_LB_MIN_GZIP_SIZE    | longbirdge protocol minimal gzip size for quote context          | 1024                                     |         |

## License

Licensed under either of

* Apache License, Version 2.0,([LICENSE-APACHE](./LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](./LICENSE-MIT) or http://opensource.org/licenses/MIT) at your option.
