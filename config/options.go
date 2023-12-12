package config

import (
	"path"
)

type Options struct {
	tp       ConfigType
	filePath string

	appKey      *string
	appSecret   *string
	accessToken *string
}

type Option func(*Options)

// WithFilePath config path
func WithFilePath(filePath string) Option {
	return func(o *Options) {
		if filePath != "" {
			o.filePath = filePath
			fileSuffix := path.Ext(filePath)
			if fileSuffix != "" {
				o.tp = ConfigType(fileSuffix)
			}
		}
	}
}

// WithConfigKey config appKey, appSecret, accessToken
func WithConfigKey(appKey string, appSecret string, accessToken string) Option {
	return func(o *Options) {
		o.appKey = &appKey
		o.appSecret = &appSecret
		o.accessToken = &accessToken
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		tp: ConfigTypeEnv,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
