package config

import (
	"net/http"
	"time"

	"github.com/longbridgeapp/openapi-go/log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
)

type IConfig interface {
	GetConfig(opts *Options) (*Config, error)
}

var configTypeMap = map[ConfigType]IConfig{
	ConfigTypeEnv:  &EnvConfig{},
	ConfigTypeYAML: &YAMLConfig{},
	ConfigTypeTOML: &TOMLConfig{},
}

// Config store Longbridge config
type Config struct {
	// Client custom http client
	Client *http.Client

	HttpURL     string        `env:"LONGBRIDGE_HTTP_URL" yaml:"LONGBRIDGE_HTTP_URL" toml:"LONGBRIDGE_HTTP_URL"`
	HTTPTimeout time.Duration `env:"LONGBRIDGE_HTTP_TIMEOUT" yaml:"LONGBRIDGE_HTTP_TIMEOUT" toml:"LONGBRIDGE_HTTP_TIMEOUT"`
	AppKey      string        `env:"LONGBRIDGE_APP_KEY" yaml:"LONGBRIDGE_APP_KEY" toml:"LONGBRIDGE_APP_KEY"`
	AppSecret   string        `env:"LONGBRIDGE_APP_SECRET" yaml:"LONGBRIDGE_APP_SECRET" toml:"LONGBRIDGE_APP_SECRET"`
	AccessToken string        `env:"LONGBRIDGE_ACCESS_TOKEN" yaml:"LONGBRIDGE_ACCESS_TOKEN" toml:"LONGBRIDGE_ACCESS_TOKEN"`
	TradeUrl    string        `env:"LONGBRIDGE_TRADE_URL" yaml:"LONGBRIDGE_TRADE_URL" toml:"LONGBRIDGE_TRADE_URL"`
	QuoteUrl    string        `env:"LONGBRIDGE_QUOTE_URL" yaml:"LONGBRIDGE_QUOTE_URL" toml:"LONGBRIDGE_QUOTE_URL"`

	LogLevel string `env:"LONGBRIDGE_LOG_LEVEL" yaml:"LONGBRIDGE_LOG_LEVEL" toml:"LONGBRIDGE_LOG_LEVEL"`
	logger   log.Logger

	// longbridge protocol config
	AuthTimeout    time.Duration `env:"LONGBRIDGE_AUTH_TIMEOUT" yaml:"LONGBRIDGE_AUTH_TIMEOUT"toml:"LONGBRIDGE_AUTH_TIMEOUT"`
	Timeout        time.Duration `env:"LONGBRIDGE_TIMEOUT" yaml:"LONGBRIDGE_TIMEOUT" toml:"LONGBRIDGE_TIMEOUT"`
	WriteQueueSize int           `env:"LONGBRIDGE_WRITE_QUEUE_SIZE" yaml:"LONGBRIDGE_WRITE_QUEUE_SIZE" toml:"LONGBRIDGE_WRITE_QUEUE_SIZE"`
	ReadQueueSize  int           `env:"LONGBRIDGE_READ_QUEUE_SIZE" yaml:"LONGBRIDGE_READ_QUEUE_SIZE" toml:"LONGBRIDGE_READ_QUEUE_SIZE"`
	ReadBufferSize int           `env:"LONGBRIDGE_READ_BUFFER_SIZE" yaml:"LONGBRIDGE_READ_BUFFER_SIZE" toml:"LONGBRIDGE_READ_BUFFER_SIZE"`
	MinGzipSize    int           `env:"LONGBRIDGE_MIN_GZIP_SIZE" yaml:"LONGBRIDGE_MIN_GZIP_SIZE" toml:"LONGBRIDGE_MIN_GZIP_SIZE"`
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

func New(opts ...Option) (configData *Config, err error) {
	options := newOptions(opts...)
	conf, exist := configTypeMap[options.tp]
	if !exist {
		err = errors.Errorf("config type:%+v not support", options.tp)
		return
	}
	configData, err = conf.GetConfig(options)
	if err != nil {
		err = errors.Wrapf(err, "GetConfig err")
		return
	}
	err = configData.check()
	if err != nil {
		err = errors.Wrapf(err, "New config check err")
		return
	}
	log.SetLevel(configData.LogLevel)
	return
}

func (c *Config) check() (err error) {
	if c.AccessToken == "" {
		err = errors.New("Don't has accessToken. Please set access token on LONGBRIDGE_ACCESS_TOKEN env")
		return
	}
	if c.AppKey == "" {
		err = errors.New("Don't has appKey. Please set app key on LONGBRIDGE_APP_KEY env")
		return
	}
	if c.AppSecret == "" {
		err = errors.New("Don't has appSecret. Please set app secret on LONGBRIDGE_APP_SECRET env")
		return
	}
	return
}

// Deprecated: NewFormEnv to create config with enviromente variables
func NewFormEnv() (*Config, error) {
	return New()
}
