package config

import (
	"time"

	"github.com/pkg/errors"

	env "github.com/Netflix/go-env"
)

// Config store Longbridge config
type Config struct {
	HttpURL     string `env:"LONGBRIDGE_HTTP_URL"`
	AppKey      string `env:"LONGBRIDGE_APP_KEY"`
	AppSecret   string `env:"LONGBRIDGE_APP_SECRET"`
	AccessToken string `env:"LONGBRIDGE_ACCESS_TOKEN"`
	TradeUrl    string `env:"LONGBRIDGE_TRADE_URL"`
	QuoteUrl    string `env:"LONGBRIDGE_QUOTE_URL"`
	LogLevel    string `env:"LONGBRIDGE_LOG_LEVEL"`

	// trade longbridge protocol config
	TradeLBAuthTimeout    time.Duration `env:"LONGBRIDGE_TRADE_LB_AUTH_TIMEOUT"`
	TradeLBTimeout        time.Duration `env:"LONGBRIDGE_TRADE_LB_TIMEOUT"`
	TradeLBWriteQueueSize int           `env:"LONGBRIDGE_TRADE_LB_WRITE_QUEUE_SIZE"`
	TradeLBReadQueueSize  int           `env:"LONGBRIDGE_TRADE_LB_READ_QUEUE_SIZE"`
	TradeLBReadBufferSize int           `env:"LONGBRIDGE_TRADE_LB_READ_BUFFER_SIZE"`
	TradeLBMinGzipSize    int           `env:"LONGBRIDGE_TRADE_LB_MIN_GZIP_SIZE"`
	// quote longbridge protocol config
	QuoteLBAuthTimeout    time.Duration `env:"LONGBRIDGE_TRADE_LB_AUTH_TIMEOUT"`
	QuoteLBTimeout        time.Duration `env:"LONGBRIDGE_TRADE_LB_TIMEOUT"`
	QuoteLBWriteQueueSize int           `env:"LONGBRIDGE_TRADE_LB_WRITE_QUEUE_SIZE"`
	QuoteLBReadQueueSize  int           `env:"LONGBRIDGE_TRADE_LB_READ_QUEUE_SIZE"`
	QuoteLBReadBufferSize int           `env:"LONGBRIDGE_TRADE_LB_READ_BUFFER_SIZE"`
	QuoteLBMinGzipSize    int           `env:"LONGBRIDGE_TRADE_LB_MIN_GZIP_SIZE"`
}

// NewFormEnv to create config with enviromente variables
func NewFormEnv() (*Config, error) {
	conf := &Config{}
	_, err := env.UnmarshalFromEnviron(conf)
	if err != nil {
		return nil, errors.Wrap(err, "env load error")
	}
	if conf.AccessToken == "" {
		return nil, errors.New("Don't has accessToken. Please set access token on LONGBRIDGE_ACCESS_TOKEN env")
	}
	if conf.AppKey == "" {
		return nil, errors.New("Don't has appKey. Please set app key on LONGBRIDGE_APP_KEY env")
	}
	if conf.AppSecret == "" {
		return nil, errors.New("Don't has appSecret. Please set app secret on LONGBRIDGE_APP_Secret env")
	}
	return conf, nil
}
