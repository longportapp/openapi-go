package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type YAMLConfig struct {
}

func (c *YAMLConfig) GetConfig(opts *Options) (data *Config, err error) {
	data = &Config{}
	bytes, err := ioutil.ReadFile(opts.filePath)
	if err != nil {
		err = errors.Wrapf(err, "YAML ReadFile err")
		return
	}
	err = yaml.Unmarshal(bytes, data)
	if err != nil {
		err = errors.Wrapf(err, "YAML GetConfig err")
		return
	}
	return
}
