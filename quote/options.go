package quote

import (
	"github.com/longportapp/openapi-go"
	"github.com/longportapp/openapi-go/http"
	"github.com/longportapp/openapi-go/log"
	"github.com/longportapp/openapi-go/longbridge"
	protocol "github.com/longportapp/openapi-protocol/go"
)

const (
	DefaultQuoteUrl = "wss://openapi-quote.longportapp.com/v2"
)

// Options for quote context
type Options struct {
	quoteURL        string
	httpClient      *http.Client
	lbOpts          *longbridge.Options
	logLevel        string
	logger          log.Logger
	enableOvernight bool
	language        openapi.Language
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

func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		if l != nil {
			o.logger = l
		}
	}
}

func WithEnableOvernight(enable bool) Option {
	return func(o *Options) {
		o.enableOvernight = enable
	}
}

func WithLanguage(language openapi.Language) Option {
	return func(o *Options) {
		o.language = language
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		quoteURL: DefaultQuoteUrl,
		lbOpts:   longbridge.NewOptions(),
		logger:   &protocol.DefaultLogger{},
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
