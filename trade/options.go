package trade

import "github.com/longbridgeapp/openapi-go/http"

const DefaultTradeUrl = "wss://openapi-trade.longbridgeapp.com"

type Options struct {
	TradeURL   string
	HttpClient *http.Client
}

type Option func(*Options)

func WithTradeURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.TradeURL = url
		}
	}
}

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
