package config

type ConfigType string

const (
	ConfigTypeEnv  ConfigType = ".env"
	ConfigTypeYAML ConfigType = ".yaml"
	ConfigTypeTOML ConfigType = ".toml"
)
