package config

import (
	"github.com/Netflix/go-env"
	"github.com/pkg/errors"
)

type EnvConfig struct {
}

func (c *EnvConfig) GetConfig(_ *Options) (data *Config, err error) {
	data = &Config{}
	_, err = env.UnmarshalFromEnviron(data)
	if err != nil {
		err = errors.Wrapf(err, "Env GetConfig err")
		return
	}
	return
}
