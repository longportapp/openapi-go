package quote

import (
	"github.com/longbridgeapp/openapi-go/http"
	"github.com/longbridgeapp/openapi-go/longbridge"
)

const (
	DefaultQuoteUrl         = "wss://openapi-quote.longbridgeapp.com/v2"
)

// Options for quote context
type Options struct {
	quoteURL   string
	httpClient *http.Client
	lbOpts     *longbridge.Options
	logLevel   string
}

// Option
type Option func(*Options)

// WithQuoteURL to set url for quote context
func WithQuoteURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.quoteURL = url
		}
	}
}

// WithHttpClient to set http client for quote context
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
		quoteURL: DefaultQuoteUrl,
		lbOpts:   longbridge.NewOptions(),
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
