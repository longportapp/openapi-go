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

## Config

### Load from env

Support init config from env, and support load env from `.env` file

```golang
import (
    "github.com/longbridgeapp/openapi-go/config"
    "github.com/longbridgeapp/openapi-go/trade"
    "github.com/longbridgeapp/openapi-go/http"
)

func main() {
    c, err := config.NewFromEnv()

    if err != nil {
        // panic
    }

    // init http client from config
    c, err := http.NewFromCfg(c)

    // init trade context from config
    tc, err := trade.NewFromCfg(c)

    // init quote context from config
    qc, err := quote.NewFromCfg(c)
}

```

All envs is listed in the last of [README](#Environment Variables)

### Init Config manually

Config structure as follow:

```golang
type Config struct {
	HttpURL     string `env:"LONGBRIDGE_HTTP_URL"`
	AppKey      string `env:"LONGBRIDGE_APP_KEY"`
	AppSecret   string `env:"LONGBRIDGE_APP_SECRET"`
	AccessToken string `env:"LONGBRIDGE_ACCESS_TOKEN"`
	TradeUrl    string `env:"LONGBRIDGE_TRADE_URL"`
	QuoteUrl    string `env:"LONGBRIDGE_QUOTE_URL"`

	LogLevel string `env:"LONGBRIDGE_LOG_LEVEL"`
	logger   log.Logger

	// trade longbridge protocol config
	TradeLBAuthTimeout    time.Duration `env:"LONGBRIDGE_TRADE_LB_AUTH_TIMEOUT"`
	TradeLBTimeout        time.Duration `env:"LONGBRIDGE_TRADE_LB_TIMEOUT"`
	TradeLBWriteQueueSize int           `env:"LONGBRIDGE_TRADE_LB_WRITE_QUEUE_SIZE"`
	TradeLBReadQueueSize  int           `env:"LONGBRIDGE_TRADE_LB_READ_QUEUE_SIZE"`
	TradeLBReadBufferSize int           `env:"LONGBRIDGE_TRADE_LB_READ_BUFFER_SIZE"`
	TradeLBMinGzipSize    int           `env:"LONGBRIDGE_TRADE_LB_MIN_GZIP_SIZE"`
	// quote longbridge protocol config
	QuoteLBAuthTimeout    time.Duration `env:"LONGBRIDGE_QUOTE_LB_AUTH_TIMEOUT"`
	QuoteLBTimeout        time.Duration `env:"LONGBRIDGE_QUOTE_LB_TIMEOUT"`
	QuoteLBWriteQueueSize int           `env:"LONGBRIDGE_QUOTE_LB_WRITE_QUEUE_SIZE"`
	QuoteLBReadQueueSize  int           `env:"LONGBRIDGE_QUOTE_LB_READ_QUEUE_SIZE"`
	QuoteLBReadBufferSize int           `env:"LONGBRIDGE_QUOTE_LB_READ_BUFFER_SIZE"`
	QuoteLBMinGzipSize    int           `env:"LONGBRIDGE_QUOTE_LB_MIN_GZIP_SIZE"`
}

```

set config field manually

```golang
c, err := config.NewFromEnv()
c.AppKey = "xxx"
c.AppSecret = "xxx"
c.AccessToken = "xxx"

```

### set custom logger

Our logger interface as follow:

```golang
type Logger interface {
	SetLevel(string)
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
}

```

Your can use you own logger by imply the interface


```golang
c, err := config.NewFromEnv()

l := newOwnLogger()

c.SetLogger(l)

```

### use custom *(net/http).Client

the default http client is initialized simply as follow:

```golang
cli := &http.Client{Timeout: opts.Timeout}
```

we only set timeout here, your use you own *(net/http).Client.

```golang
c, err := config.NewFromEnv()

c.Client = &http.Client{
    Transport: ...
}

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

Support load env from `.env` file.

| name                        | description                                    | default value                         | example |
|-----------------------------|------------------------------------------------|---------------------------------------|---------|
| LONGBRIDGE_HTTP_URL         | longbridge rest api url                        | https://openapi.longbridgeapp.com     |         |
| LONGBRIDGE_APP_KEY          | app key                                        |                                       |         |
| LONGBRIDGE_APP_SECRET       | app secret                                     |                                       |         |
| LONGBRIDGE_ACCESS_TOKEN     | access token                                   |                                       |         |
| LONGBRIDGE_TRADE_URL        | longbridge protocol url for trade context      | wss://openapi-trade.longbridgeapp.com |         |
| LONGBRIDGE_QUOTE_URL        | longbridge protocol url for quote context      | wss://openapi-quote.longbridgeapp.com |         |
| LONGBRIDGE_LOG_LEVEL        | log level                                      | info                                  |         |
| LONGBRIDGE_AUTH_TIMEOUT     | longbridge protocol authorize request time out | 10 second                             | 10s     |
| LONGBRIDGE_TIMEOUT          | longbridge protocol dial timeout               | 5 second                              | 6s      |
| LONGBRIDGE_WRITE_QUEUE_SIZE | longbirdge protocol write queue size           | 16                                    |         |
| LONGBRIDGE_READ_QUEUE_SIZE  | longbirdge protocol read queue size            | 16                                    |         |
| LONGBRIDGE_READ_BUFFER_SIZE | longbirdge protocol read buffer size           | 4096                                  |         |
| LONGBRIDGE_MIN_GZIP_SIZE    | longbirdge protocol minimal gzip size          | 1024                                  |         |

## License

Licensed under either of

* Apache License, Version 2.0,([LICENSE-APACHE](./LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](./LICENSE-MIT) or http://opensource.org/licenses/MIT) at your option.
