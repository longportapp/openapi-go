package quote

import "github.com/longbridgeapp/openapi-go/http"

type Options struct {
	QuoteURL string
	HttpClient *http.Client
}

type Option func(*Options)

func WithQuoteURL(url string) Option {
	return func(o *Options) {
		o.QuoteURL = url
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
