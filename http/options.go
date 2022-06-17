package http

import "time"

const DefaultHttpUrl = "https://openapi.longbridgeapp.com"
const DefaultTimeout = 5 * time.Second

// Options for http client
type Options struct {
	URL         string
	AppKey      string
	AppSecret   string
	AccessToken string
	Timeout     time.Duration
}

// Option for http client
type Option func(*Options)

// WithURL to set url
func WithURL(url string) Option {
	return func(opts *Options) {
		if url != "" {
			opts.URL = url
		}
	}
}

// WithAppKey to set app key
func WithAppKey(appKey string) Option {
	return func(opts *Options) {
		if appKey != "" {
			opts.AppKey = appKey
		}
	}
}

// WithAppSecret to set app secret
func WithAppSecret(appSecret string) Option {
	return func(opts *Options) {
		if appSecret != "" {
			opts.AppSecret = appSecret
		}
	}
}

// WithAccessToken to set access token
func WithAccessToken(accessToken string) Option {
	return func(opts *Options) {
		if accessToken != "" {
			opts.AccessToken = accessToken
		}
	}
}

// WithTimeout to set timeout
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		if timeout > 0 {
			opts.Timeout = timeout
		}
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		Timeout: DefaultTimeout,
		URL:     DefaultHttpUrl,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
