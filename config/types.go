package config

type ConfigType string

const (
	ConfigTypeEnv  ConfigType = ".env"
	ConfigTypeYAML ConfigType = ".yaml"
	ConfigTypeTOML ConfigType = ".toml"
)

type Language string

const (
	LanguageZHCN Language = "zh-CN"
	LanguageZHHK Language = "zh-HK"
	LanguageEN   Language = "en"
)
