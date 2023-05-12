package config

import (
	"net/http"
	"time"

	"github.com/longbridgeapp/openapi-go/log"

	env "github.com/Netflix/go-env"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
)

// Config store Longbridge config
type Config struct {
	// Client custom http client
	Client *http.Client

	HttpURL     string        `env:"LONGBRIDGE_HTTP_URL"`
	HTTPTimeout time.Duration `env:"LONGBRIDGE_HTTP_TIMEOUT`
	AppKey      string        `env:"LONGBRIDGE_APP_KEY"`
	AppSecret   string        `env:"LONGBRIDGE_APP_SECRET"`
	AccessToken string        `env:"LONGBRIDGE_ACCESS_TOKEN"`
	TradeUrl    string        `env:"LONGBRIDGE_TRADE_URL"`
	QuoteUrl    string        `env:"LONGBRIDGE_QUOTE_URL"`

	LogLevel string `env:"LONGBRIDGE_LOG_LEVEL"`
	logger   log.Logger

	// longbridge protocol config
	AuthTimeout    time.Duration `env:"LONGBRIDGE_AUTH_TIMEOUT"`
	Timeout        time.Duration `env:"LONGBRIDGE_TIMEOUT"`
	WriteQueueSize int           `env:"LONGBRIDGE_WRITE_QUEUE_SIZE"`
	ReadQueueSize  int           `env:"LONGBRIDGE_READ_QUEUE_SIZE"`
	ReadBufferSize int           `env:"LONGBRIDGE_READ_BUFFER_SIZE"`
	MinGzipSize    int           `env:"LONGBRIDGE_MIN_GZIP_SIZE"`
}

func (c *Config) SetLogger(l log.Logger) {
	if l != nil {
		l.SetLevel(c.LogLevel)
		c.logger = l
		log.SetLogger(l)
	}
}

func (c *Config) Logger() log.Logger {
	return c.logger
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

	log.SetLevel(conf.LogLevel)

	return conf, nil
}
