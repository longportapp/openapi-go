package http


type Options struct {
	URL string
	AppKey string
	AppSecret string
	AccessToken string
}


type Option func(*Options)

func WithURL(url string) Option{
	return func(opts *Options) {
		opts.URL = url
	}
}

func WithAppKey(appKey string) Option{
	return func(opts *Options) {
		opts.AppKey = appKey
	}
}

func WithAppSecret(appSecret string) Option {
	return func(opts *Options) {
		opts.AppSecret = appSecret
	}
}

func WithAccessToken(accessToken string) Option {
	return func(opts *Options) {
		opts.AccessToken = accessToken
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
