package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type TOMLConfig struct {
}

func (c *TOMLConfig) GetConfig(opts *Options) (data *Config, err error) {
	data = &Config{}
	_, err = toml.Decode(opts.filePath, data)
	if err != nil {
		err = errors.Wrapf(err, "TOML GetConfig err")
		return
	}
	return
}
