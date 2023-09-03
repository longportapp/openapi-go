package config

type ConfigType int32

const (
	ConfigTypeEnv ConfigType = iota
	ConfigTypeYAML
	ConfigTypeTOML
)
