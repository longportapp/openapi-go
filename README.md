# Longport OpenAPI SDK for Go

`Longport` provides an easy-to-use interface for invokes [`Longport OpenAPI`](https://open.longportapp.com/en/).

## Quickstart

_With Go module support , simply add the following import_

```golang
import "github.com/longportapp/openapi-go"
```

_Setting environment variables(MacOS/Linux)_

```bash
export LONGPORT_APP_KEY="App Key get from user center"
export LONGPORT_APP_SECRET="App Secret get from user center"
export LONGPORT_ACCESS_TOKEN="Access Token get from user center"
```

_Setting environment variables(Windows)_

```powershell
setx LONGPORT_APP_KEY "App Key get from user center"
setx LONGPORT_APP_SECRET "App Secret get from user center"
setx LONGPORT_ACCESS_TOKEN "Access Token get from user center"
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

#### yaml example

To load configuration from a YAML file, use the following code snippet:

```golang
conf, err := config.New(config.WithFilePath("./test.yaml"))
```

Here is an example of what the `test.yaml` file might look like:


```yaml
longport:
  app_key: xxxxx
  app_secret: xxxxx 
  access_token: xxxxx 
```

#### toml example

Similarly, to load configuration from a TOML file, use this code snippet:

```golang
conf, err := config.New(config.WithFilePath("./test.toml"))
```

And here is an example of a `test.toml` file:

```toml
[longport]
app_key = "xxxxx"
app_secret = "xxxxx"
access_token = "xxxxx"
```

### Init Config manually

Config structure as follow:

```golang
type Config struct {
    HttpURL     string        `env:"LONGPORT_HTTP_URL" yaml:"http_url" toml:"http_url"`
    HTTPTimeout time.Duration `env:"LONGPORT_HTTP_TIMEOUT" yaml:"http_timeout" toml:"http_timeout"`
    AppKey      string        `env:"LONGPORT_APP_KEY" yaml:"app_key" toml:"app_key"`
    AppSecret   string        `env:"LONGPORT_APP_SECRET" yaml:"app_secret" toml:"app_secret"`
    AccessToken string        `env:"LONGPORT_ACCESS_TOKEN" yaml:"access_token" toml:"access_token"`
    TradeUrl    string        `env:"LONGPORT_TRADE_URL" yaml:"trade_url" toml:"trade_url"`
    QuoteUrl    string        `env:"LONGPORT_QUOTE_URL" yaml:"quote_url" toml:"quote_url"`
    EnableOvernight bool          `env:"LONGPORT_ENABLE_OVERNIGHT" yaml:"enable_overnight" toml:"enable_overnight"`

    LogLevel string `env:"LONGPORT_LOG_LEVEL" yaml:"log_level" toml:"log_level"`
    // LONGPORT protocol config
    AuthTimeout    time.Duration `env:"LONGPORT_AUTH_TIMEOUT" yaml:"auth_timeout" toml:"timeout"`
    Timeout        time.Duration `env:"LONGPORT_TIMEOUT" yaml:"timeout" toml:"timeout"`
    WriteQueueSize int           `env:"LONGPORT_WRITE_QUEUE_SIZE" yaml:"write_queue_size" toml:"write_queue_size"`
    ReadQueueSize  int           `env:"LONGPORT_READ_QUEUE_SIZE" yaml:"read_queue_size" toml:"read_queue_size"`
    ReadBufferSize int           `env:"LONGPORT_READ_BUFFER_SIZE" yaml:"read_buffer_size" toml:"read_buffer_size"`
    MinGzipSize    int           `env:"LONGPORT_MIN_GZIP_SIZE" yaml:"min_gzip_size" toml:"min_gzip_size"`
    Region Region `env:"LONGPORT_REGION" yaml:"region" toml:"region"`
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

### use custom \*(net/http).Client

the default http client is initialized simply as follow:

```golang
cli := &http.Client{Timeout: opts.Timeout}
```

we only set timeout here, you can use you own \*(net/http).Client.

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

| name                      | description                                                                                           | default value                       | example |
| ------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------- | ------- |
| LONGPORT_REGION           | Set access region, if region equals `cn`, sdk will set httpUrl, quoteUrl, tradeUrl to China endpoints | -                                   | cn      |
| LONGPORT_HTTP_URL         | LONGPORT rest api url                                                                                 | <https://openapi.longportapp.com>   |         |
| LONGPORT_APP_KEY          | app key                                                                                               |                                     |         |
| LONGPORT_APP_SECRET       | app secret                                                                                            |                                     |         |
| LONGPORT_ACCESS_TOKEN     | access token                                                                                          |                                     |         |
| LONGPORT_TRADE_URL        | LONGPORT protocol url for trade context                                                               | wss://openapi-trade.longportapp.com |         |
| LONGPORT_QUOTE_URL        | LONGPORT protocol url for quote context                                                               | wss://openapi-quote.longportapp.com |         |
| LONGPORT_LOG_LEVEL        | log level                                                                                             | info                                |         |
| LONGPORT_AUTH_TIMEOUT     | LONGPORT protocol authorize request time out                                                          | 10 second                           | 10s     |
| LONGPORT_TIMEOUT          | LONGPORT protocol dial timeout                                                                        | 5 second                            | 6s      |
| LONGPORT_WRITE_QUEUE_SIZE | longport protocol write queue size                                                                    | 16                                  |         |
| LONGPORT_READ_QUEUE_SIZE  | longport protocol read queue size                                                                     | 16                                  |         |
| LONGPORT_READ_BUFFER_SIZE | longport protocol read buffer size                                                                    | 4096                                |         |
| LONGPORT_MIN_GZIP_SIZE    | longport protocol minimal gzip size                                                                   | 1024                                |         |
| LONGPORT_ENABLE_OVERNIGHT | enable overnight quote subscription feature                                                           | false                               |         |

## License

Licensed under either of

- Apache License, Version 2.0,([LICENSE-APACHE](./LICENSE-APACHE) or <http://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](./LICENSE-MIT) or <http://opensource.org/licenses/MIT>) at your option.
