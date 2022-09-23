package trade

import "github.com/longbridgeapp/openapi-go/http"

const DefaultTradeUrl = "wss://openapi-trade.longbridgeapp.com/v2"

// Options for quote context
type Options struct {
	TradeURL   string
	HttpClient *http.Client
}

// Option
type Option func(*Options)

// WithTradeURL to set trade url for trade context
func WithTradeURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.TradeURL = url
		}
	}
}

// WithHttpClient to set http client for trade context
func WithHttpClient(client *http.Client) Option {
	return func(o *Options) {
		if client != nil {
			o.HttpClient = client
		}
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		TradeURL: DefaultTradeUrl,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
