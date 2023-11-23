package trade

import (
	"github.com/longportapp/openapi-go/http"
	"github.com/longportapp/openapi-go/log"
	"github.com/longportapp/openapi-go/longbridge"
	protocol "github.com/longportapp/openapi-protocol/go"
)

const (
	DefaultTradeUrl = "wss://openapi-trade.longportapp.com/v2"
)

// Options for quote context
type Options struct {
	tradeURL   string
	httpClient *http.Client
	lbOpts     *longbridge.Options
	logLevel   string
	logger     log.Logger
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

// WithLogLevel use to set log level
func WithLogLevel(level string) Option {
	return func(o *Options) {
		if level != "" {
			o.logLevel = level
		}
	}
}

// WithLogger use custom protocol.Logger implementation
func WithLogger(logger log.Logger) Option {
	return func(o *Options) {
		if logger != nil {
			o.logger = logger
		}
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		tradeURL: DefaultTradeUrl,
		lbOpts:   longbridge.NewOptions(),
		logger:   &protocol.DefaultLogger{},
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
