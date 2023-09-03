package config

type Options struct {
	tp       ConfigType
	filePath string
}

type Option func(*Options)

// WithFilePath config path
func WithFilePath(path string) Option {
	return func(o *Options) {
		if path != "" {
			o.filePath = path
		}
	}
}

// WithConfigType config init type
func WithConfigType(tp ConfigType) Option {
	return func(o *Options) {
		o.tp = tp
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
