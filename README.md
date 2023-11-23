# Longbridge OpenAPI SDK for Go

`longbridge` provides an easy-to-use interface for invokes [`Longbridge OpenAPI`](https://open.longportapp.com/en/).

## Quickstart

_With Go module support , simply add the following import_

```golang
import "github.com/longportapp/openapi-go"
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
    "github.com/longportapp/openapi-go/config"
    "github.com/longportapp/openapi-go/trade"
    "github.com/longportapp/openapi-go/http"
)

func main() {
    c, err := config.New()

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

All envs is listed in the last of [README](#environment-variables)

### Load from file[yaml,toml]
yaml
```golang
conf, err := config.New(config.WithFilePath("./test.yaml"))
```
toml
```golang
conf, err := config.New(config.WithFilePath("./test.toml"))
```

### Init Config manually

Config structure as follow:

```golang
type Config struct {
    HttpURL     string        `env:"LONGBRIDGE_HTTP_URL" yaml:"LONGBRIDGE_HTTP_URL" toml:"LONGBRIDGE_HTTP_URL"`
    HTTPTimeout time.Duration `env:"LONGBRIDGE_HTTP_TIMEOUT" yaml:"LONGBRIDGE_HTTP_TIMEOUT" toml:"LONGBRIDGE_HTTP_TIMEOUT"`
    AppKey      string        `env:"LONGBRIDGE_APP_KEY" yaml:"LONGBRIDGE_APP_KEY" toml:"LONGBRIDGE_APP_KEY"`
    AppSecret   string        `env:"LONGBRIDGE_APP_SECRET" yaml:"LONGBRIDGE_APP_SECRET" toml:"LONGBRIDGE_APP_SECRET"`
    AccessToken string        `env:"LONGBRIDGE_ACCESS_TOKEN" yaml:"LONGBRIDGE_ACCESS_TOKEN" toml:"LONGBRIDGE_ACCESS_TOKEN"`
    TradeUrl    string        `env:"LONGBRIDGE_TRADE_URL" yaml:"LONGBRIDGE_TRADE_URL" toml:"LONGBRIDGE_TRADE_URL"`
    QuoteUrl    string        `env:"LONGBRIDGE_QUOTE_URL" yaml:"LONGBRIDGE_QUOTE_URL" toml:"LONGBRIDGE_QUOTE_URL"`
    
    LogLevel string `env:"LONGBRIDGE_LOG_LEVEL" yaml:"LONGBRIDGE_LOG_LEVEL" toml:"LONGBRIDGE_LOG_LEVEL"`
    // longbridge protocol config
    AuthTimeout    time.Duration `env:"LONGBRIDGE_AUTH_TIMEOUT" yaml:"LONGBRIDGE_AUTH_TIMEOUT"toml:"LONGBRIDGE_AUTH_TIMEOUT"`
    Timeout        time.Duration `env:"LONGBRIDGE_TIMEOUT" yaml:"LONGBRIDGE_TIMEOUT" toml:"LONGBRIDGE_TIMEOUT"`
    WriteQueueSize int           `env:"LONGBRIDGE_WRITE_QUEUE_SIZE" yaml:"LONGBRIDGE_WRITE_QUEUE_SIZE" toml:"LONGBRIDGE_WRITE_QUEUE_SIZE"`
    ReadQueueSize  int           `env:"LONGBRIDGE_READ_QUEUE_SIZE" yaml:"LONGBRIDGE_READ_QUEUE_SIZE" toml:"LONGBRIDGE_READ_QUEUE_SIZE"`
    ReadBufferSize int           `env:"LONGBRIDGE_READ_BUFFER_SIZE" yaml:"LONGBRIDGE_READ_BUFFER_SIZE" toml:"LONGBRIDGE_READ_BUFFER_SIZE"`
    MinGzipSize    int           `env:"LONGBRIDGE_MIN_GZIP_SIZE" yaml:"LONGBRIDGE_MIN_GZIP_SIZE" toml:"LONGBRIDGE_MIN_GZIP_SIZE"`
}

```

set config field manually

```golang
c, err := config.New()
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
c, err := config.New()

l := newOwnLogger()

c.SetLogger(l)

```

### use custom *(net/http).Client

the default http client is initialized simply as follow:

```golang
cli := &http.Client{Timeout: opts.Timeout}
```

we only set timeout here, you can use you own *(net/http).Client.

```golang
c, err := config.New()

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

	"github.com/longportapp/openapi-go/quote"
	"github.com/longportapp/openapi-go/config"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}
	// create quote context from environment variables
	quoteContext, err := quote.NewFromCfg(conf)
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

	"github.com/longportapp/openapi-go/trade"
	"github.com/longportapp/openapi-go/config"
	"github.com/shopspring/decimal"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}
	// create trade context from environment variables
	tradeContext, err := trade.NewFromCfg(conf)
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
| LONGBRIDGE_HTTP_URL         | longbridge rest api url                        | https://openapi.longportapp.com     |         |
| LONGBRIDGE_APP_KEY          | app key                                        |                                       |         |
| LONGBRIDGE_APP_SECRET       | app secret                                     |                                       |         |
| LONGBRIDGE_ACCESS_TOKEN     | access token                                   |                                       |         |
| LONGBRIDGE_TRADE_URL        | longbridge protocol url for trade context      | wss://openapi-trade.longportapp.com |         |
| LONGBRIDGE_QUOTE_URL        | longbridge protocol url for quote context      | wss://openapi-quote.longportapp.com |         |
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
