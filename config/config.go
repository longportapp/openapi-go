package config

import (
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"

	"github.com/longportapp/openapi-go"
	"github.com/longportapp/openapi-go/log"
	"github.com/longportapp/openapi-go/oauth"
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

	stagingHttpUrl  = "https://openapi.longportapp.xyz"
	stagingQuoteUrl = "wss://openapi-quote.longportapp.xyz"
	stagingTradeUrl = "wss://openapi-trade.longportapp.xyz"
)

// IsStaging returns true when LONGPORT_ENV is set to "staging".
func IsStaging() bool {
	return os.Getenv("LONGPORT_ENV") == "staging"
}

// Config store Longport config
type Config struct {
	// Client custom http client
	Client *http.Client

	HttpURL         string           `env:"LONGPORT_HTTP_URL" yaml:"http_url" toml:"http_url"`
	HTTPTimeout     time.Duration    `env:"LONGPORT_HTTP_TIMEOUT" yaml:"http_timeout" toml:"http_timeout"`
	AppKey          string           `env:"LONGPORT_APP_KEY" yaml:"app_key" toml:"app_key"`
	AppSecret       string           `env:"LONGPORT_APP_SECRET" yaml:"app_secret" toml:"app_secret"`
	AccessToken     string           `env:"LONGPORT_ACCESS_TOKEN" yaml:"access_token" toml:"access_token"`
	TradeUrl        string           `env:"LONGPORT_TRADE_URL" yaml:"trade_url" toml:"trade_url"`
	QuoteUrl        string           `env:"LONGPORT_QUOTE_URL" yaml:"quote_url" toml:"quote_url"`
	EnableOvernight bool             `env:"LONGPORT_ENABLE_OVERNIGHT" yaml:"enable_overnight" toml:"enable_overnight"`
	Language        openapi.Language `env:"LONGPORT_LANGUAGE" yaml:"language" toml:"language"`

	LogLevel string `env:"LONGPORT_LOG_LEVEL" yaml:"log_level" toml:"log_level"`
	logger   log.Logger

	// OAuthClient is set when using OAuth 2.0 with auto-refresh (see WithOAuthClient).
	OAuthClient *oauth.OAuth

	// Longport protocol config
	AuthTimeout    time.Duration `env:"LONGPORT_AUTH_TIMEOUT" yaml:"auth_timeout" toml:"auth_timeout"`
	Timeout        time.Duration `env:"LONGPORT_TIMEOUT" yaml:"timeout" toml:"timeout"`
	WriteQueueSize int           `env:"LONGPORT_WRITE_QUEUE_SIZE" yaml:"write_queue_size" toml:"write_queue_size"`
	ReadQueueSize  int           `env:"LONGPORT_READ_QUEUE_SIZE" yaml:"read_queue_size" toml:"read_queue_size"`
	ReadBufferSize int           `env:"LONGPORT_READ_BUFFER_SIZE" yaml:"read_buffer_size" toml:"read_buffer_size"`
	MinGzipSize    int           `env:"LONGPORT_MIN_GZIP_SIZE" yaml:"min_gzip_size" toml:"min_gzip_size"`
	Region         Region        `env:"LONGPORT_REGION" yaml:"region" toml:"region"`

	// ExtraHeaders are additional HTTP headers sent with every request.
	ExtraHeaders map[string]string
}

// parseConfig is a config for toml/yaml
type parseConfig struct {
	Longport *Config `toml:"longport" yaml:"longport"`
}

// WithHeader adds a custom HTTP header that will be sent with every request.
func (c *Config) WithHeader(key, value string) *Config {
	if c.ExtraHeaders == nil {
		c.ExtraHeaders = make(map[string]string)
	}
	c.ExtraHeaders[key] = value
	return c
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
	if options.oauthClient != nil {
		configData.OAuthClient = options.oauthClient
	}

	if IsStaging() {
		configData.HttpURL = stagingHttpUrl
		configData.QuoteUrl = stagingQuoteUrl
		configData.TradeUrl = stagingTradeUrl
	} else if configData.Region == RegionCN {
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
	if c.OAuthClient != nil {
		// OAuth 2.0 with client: token and app key are resolved at request time
		return nil
	}
	if c.AccessToken == "" {
		err = errors.New("missing access token (set LONGPORT_ACCESS_TOKEN or use WithOAuthClient)")
		return
	}
	if c.AppKey == "" {
		err = errors.New("missing app key (set LONGPORT_APP_KEY or use WithOAuthClient)")
		return
	}
	// WithOAuthClient skips this path; here AppSecret is required
	if c.AppSecret == "" {
		err = errors.New("missing app secret (set LONGPORT_APP_SECRET or use WithOAuthClient)")
		return
	}
	return
}

// Deprecated: NewFormEnv to create config from environment variables
func NewFormEnv() (*Config, error) {
	return New()
}
