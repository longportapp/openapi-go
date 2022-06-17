package quote

import (
	"context"
	"time"

	"github.com/longbridgeapp/openapi-go"
	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/http"
	"github.com/pkg/errors"
)

// QuoteContext is a client for interacting with Longbridge Quote OpenAPI
// Longbrige Quote OpenAPI document is https://open.longbridgeapp.com/en/docs/quote/overview
type QuoteContext struct {
	opts *Options
	core *core
}

// OnQuote set callback function which will be called when server push events.
func (c *QuoteContext) OnQuote(f func(*PushEvent)) {
	c.core.SetHandler(f)
}

// Subscribe quote
// Reference: https://open.longbridgeapp.com/en/docs/quote/subscribe/subscribe
func (c *QuoteContext) Subscribe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	return c.core.Subscribe(ctx, symbols, subTypes, isFirstPush)
}

// Unsubscribe quote
// Reference: https://open.longbridgeapp.com/en/docs/quote/subscribe/unsubscribe
func (c *QuoteContext) Unsubscribe(ctx context.Context, unSubAll bool, symbols []string, subTypes []SubType) (err error) {
	return c.core.Unsubscribe(ctx, unSubAll, symbols, subTypes)
}

// Subscriptions obtain the subscription information.
// Reference: https://open.longbridgeapp.com/en/docs/quote/subscribe/subscription
func (c *QuoteContext) Subscriptions(ctx context.Context) (subscriptions map[string][]SubType, err error) {
	return c.core.Subscriptions(ctx)
}

// StaticInfo obtain the basic information of securities.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/static
func (c *QuoteContext) StaticInfo(ctx context.Context, symbols []string) (staticInfos []*StaticInfo, err error) {
	return c.core.StaticInfo(ctx, symbols)
}

// Quote obtain the real-time quotes of securities, and supports all types of securities.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/quote
func (c *QuoteContext) Quote(ctx context.Context, symbols []string) (quotes []*SecurityQuote, err error) {
	return c.core.Quote(ctx, symbols)
}

// OptionQuote obtain the real-time quotes of US stock options, including the option-specific data.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/option-quote
func (c *QuoteContext) OptionQuote(ctx context.Context, symbols []string) (optionQuotes []*OptionQuote, err error) {
	return c.core.OptionQuote(ctx, symbols)
}

// WarrantQuote obtain the real-time quotes of HK warrants, including the warrant-specific data.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/warrant-quote
func (c *QuoteContext) WarrantQuote(ctx context.Context, symbols []string) (warrantQuotes []*WarrantQuote, err error) {
	return c.core.WarrantQuote(ctx, symbols)
}

// Depth obtain the depth data of security.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/depth
func (c *QuoteContext) Depth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	return c.core.Depth(ctx, symbol)
}

// Brokers obtain the real-time broker queue data of security.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/brokers
func (c *QuoteContext) Brokers(ctx context.Context, symbol string) (securityBrokers *SecurityBrokers, err error) {
	return c.core.Brokers(ctx, symbol)
}

// Participants obtain participant IDs data (which can be synchronized once a day).
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/broker-ids
func (c *QuoteContext) Participants(ctx context.Context) (infos []*ParticipantInfo, err error) {
	return c.core.Participants(ctx)
}

// Trades obtain the trades data of security.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/trade
func (c *QuoteContext) Trades(ctx context.Context, symbol string, count int32) (trades []*Trade, err error) {
	return c.core.Trades(ctx, symbol, count)
}

// Intraday obtain the intraday data of security.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/intraday
func (c *QuoteContext) Intraday(ctx context.Context, symbol string) (lines []*IntradayLine, err error) {
	return c.core.Intraday(ctx, symbol)
}

// Candlesticks obtain the candlestick data of security.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/candlestick
func (c *QuoteContext) Candlesticks(ctx context.Context, symbol string, period Period, count int32, adjustType AdjustType) (sticks []*Candlestick, err error) {
	return c.core.Candlesticks(ctx, symbol, period, count, adjustType)
}

// OptionChainExpiryDateList obtain the the list of expiration dates of option chain
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/optionchain-date
func (c *QuoteContext) OptionChainExpiryDateList(ctx context.Context, symbol string) (times []time.Time, err error) {
	return c.core.OptionChainExpiryDateList(ctx, symbol)
}

// OptionChainInfoByDate obtain a list of option securities by the option chain expiry date.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/optionchain-date-strike
func (c *QuoteContext) OptionChainInfoByDate(ctx context.Context, symbol string, expiryDate *time.Time) (times []time.Time, err error) {
	return c.core.OptionChainInfoByDate(ctx, symbol, expiryDate)
}

// WarrantIssuers obtain the warrant issuer IDs data (which can be synchronized once a day).
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/issuer
func (c *QuoteContext) WarrantIssuers(ctx context.Context) (infos []*IssuerInfo, err error) {
	return c.core.WarrantIssuers(ctx)
}

// TradingSession obtain the daily trading hours of each market.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/trade-session
func (c *QuoteContext) TradingSession(ctx context.Context) (sessions []*MarketTradingSession, err error) {
	return c.core.TradingSession(ctx)
}

// TradingDays obtain the trading days of the market.
// Reference: https://open.longbridgeapp.com/en/docs/quote/pull/trade-day
func (c *QuoteContext) TradingDays(ctx context.Context, market openapi.Market, begin *time.Time, end *time.Time) (days *MarketTradingDay, err error) {
	return c.core.TradingDays(ctx, market, begin, end)
}

// RealtimeQuote to get quote infomations on local store
func (c *QuoteContext) RealtimeQuote(ctx context.Context, symbols []string) ([]*Quote, error) {
	return c.core.RealtimeQuote(ctx, symbols)
}

// RealtimeDepth to get depth infomations on local store
func (c *QuoteContext) RealtimeDepth(ctx context.Context, symbol string) (*SecurityDepth, error) {
	return c.core.RealtimeDepth(ctx, symbol)
}

// RealtimeTrades to get trade infomations on local store
func (c *QuoteContext) RealtimeTrades(ctx context.Context, symbol string) ([]*Trade, error) {
	return c.core.RealtimeTrades(ctx, symbol)
}

// RealtimeBrokers to get broker infomations on local store
func (c *QuoteContext) RealtimeBrokers(ctx context.Context, symbol string) (*SecurityBrokers, error) {
	return c.core.RealtimeBrokers(ctx, symbol)
}

// Close
func (c *QuoteContext) Close() error {
	return c.core.Close()
}

// NewFromEnv return QuoteContext with environment variables.
func NewFormEnv() (*QuoteContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, err
	}
	return NewFromCfg(cfg)
}

// NewFromCfg return QuoteContext with config.Config
func NewFromCfg(cfg *config.Config) (*QuoteContext, error) {
	httpClient, err := http.New(
		http.WithAccessToken(cfg.AccessToken),
		http.WithAppKey(cfg.AppKey),
		http.WithAppSecret(cfg.AppSecret),
		http.WithURL(cfg.HttpURL),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return New(WithQuoteURL(cfg.QuoteUrl), WithHttpClient(httpClient))
}

// New return QuoteContext with option.
// A connection will be created with quote server.
func New(opt ...Option) (*QuoteContext, error) {
	opts := newOptions(opt...)
	otp, err := opts.HttpClient.GetOTP(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get otp")
	}
	core, err := newCore(opts.QuoteURL, otp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &QuoteContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}
