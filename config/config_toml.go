package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type TOMLConfig struct {
}

func (c *TOMLConfig) GetConfig(opts *Options) (*Config, error) {
	parseData := &parseConfig{}
	_, err := toml.DecodeFile(opts.filePath, parseData)
	if err != nil {
		err = errors.Wrapf(err, "TOML GetConfig err")
		return nil, err
	}
	if parseData.Longport == nil {
		return nil, errors.New("Longport config is not exist in toml file")
	}
	return parseData.Longport, nil
}
