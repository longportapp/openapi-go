package config

type Config struct {
	httpURL string
	appKey string
	appSecret string
	accessToken string
	tradeWsUrl string
	quoteWsUrl string
}

func NewFormEnv() *Config {
	return nil
}
