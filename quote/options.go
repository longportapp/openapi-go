package quote

import "github.com/longbridgeapp/openapi-go/http"

const DefaultQuoteUrl = "wss://openapi-quote.longbridgeapp.com/v2"

// Options for quote context
type Options struct {
	QuoteURL   string
	HttpClient *http.Client
}

// Option
type Option func(*Options)

// WithQuoteURL to set url for quote context
func WithQuoteURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.QuoteURL = url
		}
	}
}

// WithHttpClient to set http client for quote context
func WithHttpClient(client *http.Client) Option {
	return func(o *Options) {
		if client != nil {
			o.HttpClient = client
		}
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		QuoteURL: DefaultQuoteUrl,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
