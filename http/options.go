package http

import "time"

const DefaultHttpUrl = "https://openapi.longbridgeapp.com"
const DefaultTimeout = 5 * time.Second

type Options struct {
	URL         string
	AppKey      string
	AppSecret   string
	AccessToken string
	Timeout     time.Duration
}

type Option func(*Options)

func WithURL(url string) Option {
	return func(opts *Options) {
		if url != "" {
			opts.URL = url
		}
	}
}

func WithAppKey(appKey string) Option {
	return func(opts *Options) {
		if appKey != "" {
			opts.AppKey = appKey
		}
	}
}

func WithAppSecret(appSecret string) Option {
	return func(opts *Options) {
		if appSecret != "" {
			opts.AppSecret = appSecret
		}
	}
}

func WithAccessToken(accessToken string) Option {
	return func(opts *Options) {
		if accessToken != "" {
			opts.AccessToken = accessToken
		}
	}
}

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
