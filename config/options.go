package config

import (
	"path"
)

type Options struct {
	tp       ConfigType
	filePath string
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

func newOptions(opt ...Option) *Options {
	opts := Options{
		tp: ConfigTypeEnv,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
