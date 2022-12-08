package quote

import (
	"time"

	"github.com/longbridgeapp/openapi-go/http"
)

const (
	DefaultQuoteUrl         = "wss://openapi-quote.longbridgeapp.com/v2"
	DefaultLBWriteQueueSize = 16
	DefaultLBReadBufferSize = 4096
	DefaultLBReadQueueSize  = 16
	DefaultLBMinGzipSize    = 1024
	DefaultLBDialTimeout    = time.Second * 5
	DefaultLBAuthTimeout    = time.Second * 10
)

// Options for quote context
type Options struct {
	quoteURL         string
	httpClient       *http.Client
	lbAuthTimeout    time.Duration
	lbTimeout        time.Duration
	lbWriteQueueSize int
	lbReadQueueSize  int
	lbReadBufferSize int
	lbMinGzipSize    int
	logLevel         string
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

func WithLBAuthTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		if timeout > 0 {
			o.lbAuthTimeout = timeout
		}
	}
}

func WithLBTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		if timeout > 0 {
			o.lbTimeout = timeout
		}
	}
}

func WithLBWriteQueueSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.lbWriteQueueSize = size
		}
	}
}

func WithLBReadQueueSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.lbReadQueueSize = size
		}
	}
}

func WithLBReadBufferSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.lbReadBufferSize = size
		}
	}
}

func WithLBMinGzipSize(size int) Option {
	return func(o *Options) {
		if size > 0 {
			o.lbMinGzipSize = size
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
		quoteURL:         DefaultQuoteUrl,
		lbAuthTimeout:    DefaultLBAuthTimeout,
		lbTimeout:        DefaultLBDialTimeout,
		lbWriteQueueSize: DefaultLBWriteQueueSize,
		lbReadQueueSize:  DefaultLBReadQueueSize,
		lbReadBufferSize: DefaultLBReadBufferSize,
		lbMinGzipSize:    DefaultLBMinGzipSize,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
