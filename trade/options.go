package trade

import (
	"github.com/longbridgeapp/openapi-go/http"
	"github.com/longbridgeapp/openapi-go/longbridge"
)

const (
	DefaultTradeUrl = "wss://openapi-trade.longbridgeapp.com/v2"
)

// Options for quote context
type Options struct {
	tradeURL   string
	httpClient *http.Client
	lbOpts     *longbridge.Options
	logLevel   string
}

// Option
type Option func(*Options)

// WithTradeURL to set trade url for trade context
func WithTradeURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.tradeURL = url
		}
	}
}

// WithHttpClient to set http client for trade context
func WithHttpClient(client *http.Client) Option {
	return func(o *Options) {
		if client != nil {
			o.httpClient = client
		}
	}
}

func WithLbOptions(opts *longbridge.Options) Option {
	return func(o *Options) {
		if opts != nil {
			o.lbOpts = opts
		}
	}
}

func WithLogLevel(level string) Option {
	return func(o *Options) {
		if level != "" {
			o.logLevel = level
		}
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		tradeURL: DefaultTradeUrl,
		lbOpts:   longbridge.NewOptions(),
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
