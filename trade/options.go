package trade

import "github.com/longbridgeapp/openapi-go/http"

type Options struct {
	TradeURL string
	HttpClient *http.Client
}

type Option func(*Options)

func WithTradeURL(url string) Option {
	return func(o *Options) {
		o.TradeURL = url
	}
}

func WithHttpClient(client *http.Client) Option {
	return func(o *Options) {
		o.HttpClient = client
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
