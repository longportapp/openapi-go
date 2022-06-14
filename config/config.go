package config

import (
	"os"

	"github.com/pkg/errors"
)

const ()

type Config struct {
	HttpURL     string
	AppKey      string
	AppSecret   string
	AccessToken string
	TradeUrl    string
	QuoteUrl    string
}

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
	tradeUrl := GetTradeUrlFromEnv()
	quoteUrl := GetQuoteUrlFromEnv()
	httpUrl := GetHttpUrlFromEnv()
	conf := &Config{
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: accessToken,
		TradeUrl:    tradeUrl,
		QuoteUrl:    quoteUrl,
		HttpURL:     httpUrl,
	}
	return conf, nil
}

func GetAccessTokenFromEnv() string {
	return os.Getenv("LONGBRIDGE_ACCESS_TOKEN")
}

func GetAppKeyFromEnv() string {
	return os.Getenv("LONGBRIDGE_APP_KEY")
}

func GetAppSecretFromEnv() string {
	return os.Getenv("LONGBRIDGE_APP_SECRET")
}

func GetTradeUrlFromEnv() string {
	return os.Getenv("LONGBRIDGE_TRADE_URL")
}

func GetQuoteUrlFromEnv() string {
	return os.Getenv("LONGBRIDGE_QUOTE_URL")
}

func GetHttpUrlFromEnv() string {
	return os.Getenv("LONGBRIDGE_HTTP_URL")
}
