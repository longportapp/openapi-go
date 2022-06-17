package config

import (
	"os"

	"github.com/pkg/errors"
)

// Config store Longbridge config
type Config struct {
	HttpURL     string
	AppKey      string
	AppSecret   string
	AccessToken string
	TradeUrl    string
	QuoteUrl    string
}

// NewFormEnv to create config with enviromente variables
func NewFormEnv() (*Config, error) {
	accessToken := GetAccessTokenFromEnv()
	if accessToken == "" {
		return nil, errors.New("Don't has accessToken. Please set access token on LONGBRIDGE_ACCESS_TOKEN env")
	}
	appKey := GetAppKeyFromEnv()
	if appKey == "" {
		return nil, errors.New("Don't has appKey. Please set app key on LONGBRIDGE_APP_KEY env")
	}
	appSecret := GetAppSecretFromEnv()
	if appSecret == "" {
		return nil, errors.New("Don't has appSecret. Please set app secret on LONGBRIDGE_APP_Secret env")
	}
	conf := &Config{
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		TradeUrl:    GetTradeUrlFromEnv(),
		QuoteUrl:    GetQuoteUrlFromEnv(),
		HttpURL:     GetHttpUrlFromEnv(),
	}
	return conf, nil
}

// GetAccessTokenFromEnv
func GetAccessTokenFromEnv() string {
	return os.Getenv("LONGBRIDGE_ACCESS_TOKEN")
}

// GetAppKeyFromEnv
func GetAppKeyFromEnv() string {
	return os.Getenv("LONGBRIDGE_APP_KEY")
}

// GetAppSecretFromEnv
func GetAppSecretFromEnv() string {
	return os.Getenv("LONGBRIDGE_APP_SECRET")
}

// GetTradeUrlFromEnv
func GetTradeUrlFromEnv() string {
	return os.Getenv("LONGBRIDGE_TRADE_URL")
}

// GetQuoteUrlFromEnv
func GetQuoteUrlFromEnv() string {
	return os.Getenv("LONGBRIDGE_QUOTE_URL")
}

// GetHttpUrlFromEnv
func GetHttpUrlFromEnv() string {
	return os.Getenv("LONGBRIDGE_HTTP_URL")
}

// GetLogLevelFromEnv
func GetLogLevelFromEnv() string {
	return os.Getenv("LONGBRIDGE_LOG_LEVEL")
}
