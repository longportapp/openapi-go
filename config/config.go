package config

import (
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"

	"github.com/longportapp/openapi-go/log"
)

type IConfig interface {
	GetConfig(opts *Options) (*Config, error)
}

var configTypeMap = map[ConfigType]IConfig{
	ConfigTypeEnv:  &EnvConfig{},
	ConfigTypeYAML: &YAMLConfig{},
	ConfigTypeTOML: &TOMLConfig{},
}

type Region string

var (
	RegionCN Region = "cn"

	cnHttpUrl  = "https://openapi.longportapp.cn"
	cnQuoteUrl = "wss://openapi-quote.longportapp.cn"
	cnTradeUrl = "wss://openapi-trade.longportapp.cn"
)

// Config store Longbridge config
type Config struct {
	// Client custom http client
	Client *http.Client

	HttpURL     string        `env:"LONGBRIDGE_HTTP_URL,LONGPORT_HTTP_URL" yaml:"LONGBRIDGE_HTTP_URL,LONGPORT_HTTP_URL" toml:"LONGBRIDGE_HTTP_URL,LONGPORT_HTTP_URL"`
	HTTPTimeout time.Duration `env:"LONGBRIDGE_HTTP_TIMEOUT,LONGPORT_HTTP_TIMEOUT" yaml:"LONGBRIDGE_HTTP_TIMEOUT,LONGPORT_HTTP_TIMEOUT" toml:"LONGPORT_HTTP_TIMEOUT"`
	AppKey      string        `env:"LONGBRIDGE_APP_KEY,LONGPORT_APP_KEY" yaml:"LONGBRIDGE_APP_KEY,LONGPORT_APP_KEY" toml:"LONGBRIDGE_APP_KEY,LONGPORT_APP_KEY"`
	AppSecret   string        `env:"LONGBRIDGE_APP_SECRET,LONGPORT_APP_SECRET" yaml:"eONGBRIDGE_APP_SECRET,LONGPORT_APP_SECRET" toml:"LONGBRIDGE_APP_SECRET,LONGPORT_APP_SECRET"`
	AccessToken string        `env:"LONGBRIDGE_ACCESS_TOKEN,LONGPORT_ACCESS_TOKEN" yaml:"LONGBRIDGE_ACCESS_TOKEN,LONGPORT_ACCESS_TOKEN" toml:"LONGBRIDGE_ACCESS_TOKEN,LONGPORT_ACCESS_TOKEN"`
	TradeUrl    string        `env:"LONGBRIDGE_TRADE_URL,LONGPORT_TRADE_URL" yaml:"LONGBRIDGE_TRADE_URL,LONGPORT_TRADE_URL" toml:"LONGBRIDGE_TRADE_URL,LONGPORT_TRADE_URL"`
	QuoteUrl    string        `env:"LONGBRIDGE_QUOTE_URL,LONGPORT_QUOTE_URL" yaml:"LONGBRIDGE_QUOTE_URL,LONGPORT_QUOTE_URL" toml:"LONGBRIDGE_QUOTE_URL,LONGPORT_QUOTE_URL"`

	LogLevel string `env:"LONGBRIDGE_LOG_LEVEL,LONGPORT_LOG_LEVEL" yaml:"LONGBRIDGE_LOG_LEVEL,LONGPORT_LOG_LEVEL" toml:"LONGBRIDGE_LOG_LEVEL,LONGPORT_LOG_LEVEL"`
	logger   log.Logger

	// longbridge protocol config
	AuthTimeout    time.Duration `env:"LONGBRIDGE_AUTH_TIMEOUT,LONGPORT_AUTH_TIMEOUT" yaml:"LONGBRIDGE_AUTH_TIMEOUT,LONGPORT_AUTH_TIMEOUT" toml:"LONGBRIDGE_AUTH_TIMEOUT,LONGPORT_AUTH_TIMEOUT"`
	Timeout        time.Duration `env:"LONGBRIDGE_TIMEOUT,LONGPORT_TIMEOUT" yaml:"LONGBRIDGE_TIMEOUT,LONGPORT_TIMEOUT" toml:"LONGBRIDGE_TIMEOUT,LONGPORT_TIMEOUT"`
	WriteQueueSize int           `env:"LONGBRIDGE_WRITE_QUEUE_SIZE,LONGPORT_WRITE_QUEUE_SIZE" yaml:"LONGBRIDGE_WRITE_QUEUE_SIZE,LONGPORT_WRITE_QUEUE_SIZE" toml:"LONGBRIDGE_WRITE_QUEUE_SIZE,LONGPORT_WRITE_QUEUE_SIZE"`
	ReadQueueSize  int           `env:"LONGBRIDGE_READ_QUEUE_SIZE,LONGPORT_READ_QUEUE_SIZE" yaml:"LONGBRIDGE_READ_QUEUE_SIZE,LONGPORT_READ_QUEUE_SIZE" toml:"LONGBRIDGE_READ_QUEUE_SIZE,LONGPORT_READ_QUEUE_SIZE"`
	ReadBufferSize int           `env:"LONGBRIDGE_READ_BUFFER_SIZE,LONGPORT_READ_BUFFER_SIZE" yaml:"LONGBRIDGE_READ_BUFFER_SIZE,LONGPORT_READ_BUFFER_SIZE" toml:"LONGBRIDGE_READ_BUFFER_SIZE,LONGPORT_READ_BUFFER_SIZE"`
	MinGzipSize    int           `env:"LONGBRIDGE_MIN_GZIP_SIZE,LONGPORT_MIN_GZIP_SIZE" yaml:"LONGBRIDGE_MIN_GZIP_SIZE,LONGPORT_MIN_GZIP_SIZE" toml:"LONGBRIDGE_MIN_GZIP_SIZE,LONGPORT_MIN_GZIP_SIZE"`
	Region         Region        `env:"LONGPORT_REGION" yaml:"LONGPORT_REGION" toml:"LONGPORT_REGION"`
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

	if options.appKey != nil {
		configData.AppKey = *options.appKey
	}
	if options.appSecret != nil {
		configData.AppSecret = *options.appSecret
	}
	if options.accessToken != nil {
		configData.AccessToken = *options.accessToken
	}

	if configData.Region == RegionCN {
		configData.HttpURL = cnHttpUrl
		configData.QuoteUrl = cnQuoteUrl
		configData.TradeUrl = cnTradeUrl
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
