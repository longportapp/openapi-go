# Longport OpenAPI SDK for Go

`Longport` provides an easy-to-use interface for invoking [Longport OpenAPI](https://open.longportapp.com/).

## Quickstart

_With Go module support , simply add the following import_

```golang
import "github.com/longportapp/openapi-go"
```

## Authentication

Longport OpenAPI supports two authentication methods:

1. **OAuth 2.0 (Recommended)** — Bearer tokens, no HMAC; token is stored and refreshed automatically. See [Config: Using OAuth 2.0](#using-oauth-20-recommended).
2. **Legacy API Key** — App key, secret, and access token via environment variables. See [Config: Using Legacy API Key](#using-legacy-api-key-environment-variables).

If you use OAuth, register an OAuth client first to get your `client_id`:

```bash
curl -X POST https://openapi.longportapp.com/v1/oauth2/client/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Application",
    "redirect_uris": ["http://localhost:60355/callback"],
    "grant_types": ["authorization_code", "refresh_token"]
  }'
```

Response:
```json
{
  "client_id": "your-client-id-here",
  "name": "My Application",
  "redirect_uris": ["http://localhost:60355/callback"]
}
```

Save the `client_id` for use when building config (see Config section below).

## Config

You can create a `Config` in two ways: with OAuth (recommended) or with legacy API key (environment variables). The examples in this README use OAuth.

### Using OAuth 2.0 (Recommended)

**Token storage:** After the authorization flow, the SDK stores the access and refresh tokens under `~/.longport/openapi/tokens/<client_id>` (or `%USERPROFILE%\.Longport\openapi\tokens\<client_id>` on Windows). It loads and refreshes them automatically on later runs, so you typically authorize once per machine.

```golang
import (
    "context"
    "fmt"
    "log"

    "github.com/longportapp/openapi-go/config"
    "github.com/longportapp/openapi-go/oauth"
)

func main() {
    o := oauth.New("your-client-id").
        OnOpenURL(func(url string) {
            fmt.Println("Please visit:", url)
        })
    if err := o.Build(context.Background()); err != nil {
        log.Fatal(err)
    }

    cfg, err := config.New(config.WithOAuthClient(o))
    if err != nil {
        log.Fatal(err)
    }
    // Use cfg with quote.NewFromCfg(cfg), trade.NewFromCfg(cfg), http.NewFromCfg(cfg), etc.
}
```

**Benefits:** No shared secret; no per-request HMAC; token load/refresh/persist is automatic.

### Using Legacy API Key (Environment Variables)

For backward compatibility you can use the traditional three keys. **Load from env** and **Load from file** (see below) only support Legacy API Key (app key, secret, access token); they do not support OAuth. Set env vars (or a `.env` file), then call `config.New()`.

_macOS / Linux_

```bash
export LONGPORT_APP_KEY="App Key get from user center"
export LONGPORT_APP_SECRET="App Secret get from user center"
export LONGPORT_ACCESS_TOKEN="Access Token get from user center"
```

_Windows_

```powershell
setx LONGPORT_APP_KEY "App Key get from user center"
setx LONGPORT_APP_SECRET "App Secret get from user center"
setx LONGPORT_ACCESS_TOKEN "Access Token get from user center"
```


```golang
c, err := config.New()
if err != nil {
    log.Fatal(err)
}
// Use c with NewFromCfg(c) for quote, trade, http.

```

All supported env vars are listed in [Environment Variables](#environment-variables).

### Load From File (YAML, TOML)

Load from file (and load from env above) only supports Legacy API Key. For OAuth, use [Using OAuth 2.0 (Recommended)](#using-oauth-20-recommended).

#### YAML Example

To load configuration from a YAML file, use the following code snippet:

```golang
conf, err := config.New(config.WithFilePath("./test.yaml"))
```

Here is an example of what the `test.yaml` file might look like:


```yaml
Longport:
  app_key: xxxxx
  app_secret: xxxxx 
  access_token: xxxxx 
```

#### TOML Example

Similarly, to load configuration from a TOML file, use this code snippet:

```golang
conf, err := config.New(config.WithFilePath("./test.toml"))
```

And here is an example of a `test.toml` file:

```toml
[Longport]
app_key = "xxxxx"
app_secret = "xxxxx"
access_token = "xxxxx"
```

### Init Config Manually

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
    Language    openapi.Language `env:"LONGPORT_LANGUAGE" yaml:"language" toml:"language"`

    LogLevel string `env:"LONGPORT_LOG_LEVEL" yaml:"log_level" toml:"log_level"`
    // Longport protocol config
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

### Set Custom Logger

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

### Use Custom \*(net/http).Client

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

## Quote API (Get Basic Information of Securities)

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longportapp/openapi-go/config"
    "github.com/longportapp/openapi-go/oauth"
    "github.com/longportapp/openapi-go/quote"
)

func main() {
    o := oauth.New("your-client-id").
        OnOpenURL(func(url string) {
            fmt.Println("Please visit:", url)
        })
    if err := o.Build(context.Background()); err != nil {
        log.Fatal(err)
    }
    cfg, err := config.New(config.WithOAuthClient(o))
    if err != nil {
        log.Fatal(err)
    }
    quoteContext, err := quote.NewFromCfg(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer quoteContext.Close()
    quotes, err := quoteContext.Quote(context.Background(), []string{"700.HK", "AAPL.US", "TSLA.US", "NFLX.US"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("quotes: %v", quotes)
}
```

## Fundamental API (ETF Asset Allocation)

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longportapp/openapi-go/config"
    "github.com/longportapp/openapi-go/fundamental"
)

func main() {
    cfg, err := config.NewFromEnv()
    if err != nil {
        log.Fatal(err)
    }
    fctx, err := fundamental.NewFromCfg(cfg)
    if err != nil {
        log.Fatal(err)
    }
    resp, err := fctx.EtfAssetAllocation(context.Background(), "QQQ.US")
    if err != nil {
        log.Fatal(err)
    }
    for _, group := range resp.Info {
        fmt.Printf("group: type=%d report_date=%s items=%d\n",
            group.AssetType, group.ReportDate, len(group.Lists))
        for _, item := range group.Lists {
            fmt.Printf("  %s (%s) ratio=%s\n", item.Name, item.Symbol, item.PositionRatio)
        }
    }
}
```

## Counter API (Symbol ↔ counter_id)

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longportapp/openapi-go/config"
    "github.com/longportapp/openapi-go/counter"
    "github.com/longportapp/openapi-go/quote"
)

func main() {
    // Pure local conversion (no network), backed by an embedded directory:
    fmt.Println(counter.SymbolToCounterID("QQQ.US")) // ETF/US/QQQ
    fmt.Println(counter.SymbolToCounterID("700.HK"))  // ST/HK/700
    fmt.Println(counter.CounterIDToSymbol("IX/HK/HSI")) // HSI.HK
    fmt.Println(counter.IsETF("SPY.US"))                // true

    // Batch resolution, local-first with a remote fallback (and caching):
    cfg, err := config.NewFromEnv()
    if err != nil {
        log.Fatal(err)
    }
    qctx, err := quote.NewFromCfg(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer qctx.Close()
    ids, err := qctx.ResolveCounterIds(context.Background(), []string{"TSLA.US", "QQQ.US"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ids) // map[QQQ.US:ETF/US/QQQ TSLA.US:ST/US/TSLA]
}
```

## Trade API (Submit Order)

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longportapp/openapi-go/config"
    "github.com/longportapp/openapi-go/oauth"
    "github.com/longportapp/openapi-go/trade"
    "github.com/shopspring/decimal"
)

func main() {
    o := oauth.New("your-client-id").
        OnOpenURL(func(url string) {
            fmt.Println("Please visit:", url)
        })
    if err := o.Build(context.Background()); err != nil {
        log.Fatal(err)
    }
    cfg, err := config.New(config.WithOAuthClient(o))
    if err != nil {
        log.Fatal(err)
    }
    tradeContext, err := trade.NewFromCfg(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer tradeContext.Close()
    order := &trade.SubmitOrder{
        Symbol:            "700.HK",
        OrderType:         trade.OrderTypeLO,
        Side:              trade.OrderSideBuy,
        SubmittedQuantity: 200,
        TimeInForce:       trade.TimeTypeDay,
        SubmittedPrice:    decimal.NewFromFloat(12),
    }
    orderId, err := tradeContext.SubmitOrder(context.Background(), order)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("orderId: %v\n", orderId)
}
```

## Environment Variables

Support load env from `.env` file.

| name                      | description                                                                                           | default value                       | example | optional       |
| ------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------- | ------- | -------------- |
| LONGPORT_REGION           | Set access region, if region equals `cn`, SDK will set httpUrl, quoteUrl, tradeUrl to China endpoints | -                                   | cn      | cn             |
| LONGPORT_HTTP_URL         | Longport REST API URL                                                                               | <https://openapi.longportapp.com>    |         |                |
| LONGPORT_APP_KEY          | app key                                                                                               |                                     |         |                |
| LONGPORT_APP_SECRET       | app secret                                                                                            |                                     |         |                |
| LONGPORT_ACCESS_TOKEN     | access token                                                                                          |                                     |         |                |
| LONGPORT_TRADE_URL        | Longport protocol URL for trade context                                                             | wss://openapi-trade.longportapp.com  |         |                |
| LONGPORT_QUOTE_URL        | Longport protocol URL for quote context                                                             | wss://openapi-quote.longportapp.com  |         |                |
| LONGPORT_LOG_LEVEL        | log level                                                                                             | info                                |         |                |
| LONGPORT_AUTH_TIMEOUT     | Longport protocol authorize request timeout                                                         | 10 second                           | 10s     |                |
| LONGPORT_TIMEOUT          | Longport protocol dial timeout                                                                      | 5 second                            | 6s      |                |
| LONGPORT_WRITE_QUEUE_SIZE | Longport protocol write queue size                                                                  | 16                                  |         |                |
| LONGPORT_READ_QUEUE_SIZE  | Longport protocol read queue size                                                                   | 16                                  |         |                |
| LONGPORT_READ_BUFFER_SIZE | Longport protocol read buffer size                                                                  | 4096                                |         |                |
| LONGPORT_MIN_GZIP_SIZE    | Longport protocol minimal gzip size                                                                 | 1024                                |         |                |
| LONGPORT_ENABLE_OVERNIGHT | enable overnight quote subscription feature                                                           | false                               |         |                |
| LONGPORT_LANGUAGE         | set user language for some information.                                                              | -                                   | en      | en,zh-CN,zh-HK |

## License

Licensed under either of

- Apache License, Version 2.0,([LICENSE-APACHE](./LICENSE-APACHE) or <http://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](./LICENSE-MIT) or <http://opensource.org/licenses/MIT>) at your option.
