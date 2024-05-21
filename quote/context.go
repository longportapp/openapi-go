package quote

import (
	"context"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/longportapp/openapi-go"
	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/http"
	"github.com/longportapp/openapi-go/internal/util"
	"github.com/longportapp/openapi-go/longbridge"
	"github.com/longportapp/openapi-go/quote/jsontypes"
)

// QuoteContext is a client for interacting with Longbridge Quote OpenAPI
// Longbrige Quote OpenAPI document is https://open.longportapp.com/en/docs/quote/overview
type QuoteContext struct {
	opts *Options
	core *core
}

// OnQuote set callback function which will be called when server push quote events.
func (c *QuoteContext) OnQuote(f func(*PushQuote)) {
	c.core.SetQuoteHandler(f)
}

// OnTrade set callback function which will be called when server push trade events.
func (c *QuoteContext) OnTrade(f func(*PushTrade)) {
	c.core.SetTradeHandler(f)
}

// OnDepth set callback function which will be called when server push depth events.
func (c *QuoteContext) OnDepth(f func(*PushDepth)) {
	c.core.SetDepthHandler(f)
}

// OnBrokers set callback function which will be called when server push brokers events.
func (c *QuoteContext) OnBrokers(f func(*PushBrokers)) {
	c.core.SetBrokersHandler(f)
}

// Subscribe quote
// Reference: https://open.longportapp.com/en/docs/quote/subscribe/subscribe
func (c *QuoteContext) Subscribe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	return c.core.Subscribe(ctx, symbols, subTypes, isFirstPush)
}

// Unsubscribe quote
// Reference: https://open.longportapp.com/en/docs/quote/subscribe/unsubscribe
func (c *QuoteContext) Unsubscribe(ctx context.Context, unSubAll bool, symbols []string, subTypes []SubType) (err error) {
	return c.core.Unsubscribe(ctx, unSubAll, symbols, subTypes)
}

// Subscriptions obtain the subscription information.
// Reference: https://open.longportapp.com/en/docs/quote/subscribe/subscription
func (c *QuoteContext) Subscriptions(ctx context.Context) (subscriptions map[string][]SubType, err error) {
	return c.core.Subscriptions(ctx)
}

// StaticInfo obtain the basic information of securities.
// Reference: https://open.longportapp.com/en/docs/quote/pull/static
func (c *QuoteContext) StaticInfo(ctx context.Context, symbols []string) (staticInfos []*StaticInfo, err error) {
	return c.core.StaticInfo(ctx, symbols)
}

// Quote obtain the real-time quotes of securities, and supports all types of securities.
// Reference: https://open.longportapp.com/en/docs/quote/pull/quote
func (c *QuoteContext) Quote(ctx context.Context, symbols []string) (quotes []*SecurityQuote, err error) {
	return c.core.Quote(ctx, symbols)
}

// OptionQuote obtain the real-time quotes of US stock options, including the option-specific data.
// Reference: https://open.longportapp.com/en/docs/quote/pull/option-quote
func (c *QuoteContext) OptionQuote(ctx context.Context, symbols []string) (optionQuotes []*OptionQuote, err error) {
	return c.core.OptionQuote(ctx, symbols)
}

// WarrantQuote obtain the real-time quotes of HK warrants, including the warrant-specific data.
// Reference: https://open.longportapp.com/en/docs/quote/pull/warrant-quote
func (c *QuoteContext) WarrantQuote(ctx context.Context, symbols []string) (warrantQuotes []*WarrantQuote, err error) {
	return c.core.WarrantQuote(ctx, symbols)
}

// Depth obtain the depth data of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/depth
func (c *QuoteContext) Depth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	return c.core.Depth(ctx, symbol)
}

// Brokers obtain the real-time broker queue data of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/brokers
func (c *QuoteContext) Brokers(ctx context.Context, symbol string) (securityBrokers *SecurityBrokers, err error) {
	return c.core.Brokers(ctx, symbol)
}

// Participants obtain participant IDs data (which can be synchronized once a day).
// Reference: https://open.longportapp.com/en/docs/quote/pull/broker-ids
func (c *QuoteContext) Participants(ctx context.Context) (infos []*ParticipantInfo, err error) {
	return c.core.Participants(ctx)
}

// Trades obtain the trades data of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/trade
func (c *QuoteContext) Trades(ctx context.Context, symbol string, count int32) (trades []*Trade, err error) {
	return c.core.Trades(ctx, symbol, count)
}

// Intraday obtain the intraday data of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/intraday
func (c *QuoteContext) Intraday(ctx context.Context, symbol string) (lines []*IntradayLine, err error) {
	return c.core.Intraday(ctx, symbol)
}

// Candlesticks obtain the candlestick data of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/candlestick
func (c *QuoteContext) Candlesticks(ctx context.Context, symbol string, period Period, count int32, adjustType AdjustType) (sticks []*Candlestick, err error) {
	return c.core.Candlesticks(ctx, symbol, period, count, adjustType)
}

// OptionChainExpiryDateList obtain the the list of expiration dates of option chain
// Reference: https://open.longportapp.com/en/docs/quote/pull/optionchain-date
func (c *QuoteContext) OptionChainExpiryDateList(ctx context.Context, symbol string) (times []time.Time, err error) {
	return c.core.OptionChainExpiryDateList(ctx, symbol)
}

// OptionChainInfoByDate obtain a list of option securities by the option chain expiry date.
// Reference: https://open.longportapp.com/en/docs/quote/pull/optionchain-date-strike
func (c *QuoteContext) OptionChainInfoByDate(ctx context.Context, symbol string, expiryDate *time.Time) (strikePriceInfos []*StrikePriceInfo, err error) {
	return c.core.OptionChainInfoByDate(ctx, symbol, expiryDate)
}

// WarrantIssuers obtain the warrant issuer IDs data (which can be synchronized once a day).
// Reference: https://open.longportapp.com/en/docs/quote/pull/issuer
func (c *QuoteContext) WarrantIssuers(ctx context.Context) (infos []*IssuerInfo, err error) {
	return c.core.WarrantIssuers(ctx)
}

// WarrantIssuers obtain the quotes of HK warrants, and supports sorting and filtering.
// Reference: https://open.longportapp.com/en/docs/quote/pull/warrant-filter
func (c *QuoteContext) WarrantList(ctx context.Context, symbol string, config WarrantFilter, lang WarrantLanguage) (infos []*WarrantInfo, err error) {
	return c.core.WarrantList(ctx, symbol, config, lang)
}

// TradingSession obtain the daily trading hours of each market.
// Reference: https://open.longportapp.com/en/docs/quote/pull/trade-session
func (c *QuoteContext) TradingSession(ctx context.Context) (sessions []*MarketTradingSession, err error) {
	return c.core.TradingSession(ctx)
}

// TradingDays obtain the trading days of the market.
// Reference: https://open.longportapp.com/en/docs/quote/pull/trade-day
func (c *QuoteContext) TradingDays(ctx context.Context, market openapi.Market, begin *time.Time, end *time.Time) (days *MarketTradingDay, err error) {
	return c.core.TradingDays(ctx, market, begin, end)
}

// CapitalDistribution is used to obtain the daily capital distribution of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/capital-distribution
func (c *QuoteContext) CapitalDistribution(ctx context.Context, symbol string) (capitalDib CapitalDistribution, err error) {
	return c.core.CapitalDistribution(ctx, symbol)
}

// CapitalFlow is used to obtain the daily capital flow intraday of security.
// Reference: https://open.longportapp.com/en/docs/quote/pull/capital-flow-intraday
func (c *QuoteContext) CapitalFlow(ctx context.Context, symbol string) (capitalFlowLines []CapitalFlowLine, err error) {
	return c.core.CapitalFlow(ctx, symbol)
}

// CalcIndex is used to obtain the calculate indexes of securities.
// Reference: https://open.longportapp.com/en/docs/quote/pull/calc-index
func (c *QuoteContext) CalcIndex(ctx context.Context, symbols []string, indexes []CalcIndex) (calcIndexes []*SecurityCalcIndex, err error) {
	return c.core.CalcIndex(ctx, symbols, indexes)
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

// CreateWatchlistGroup use to create watchlist group. Doc: https://open.longportapp.com/en/docs/quote/individual/watchlist_create_group
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.CreateWatchlistGroup(context.Background(), "test", []string{"AAPL.US"})

func (c *QuoteContext) CreateWatchlistGroup(ctx context.Context, name string, symbols []string) (gid int64, err error) {
	var resp struct {
		ID int64 `json:"id,string"`
	}
	err = c.opts.httpClient.Post(ctx, "/v1/watchlist/groups", map[string]interface{}{
		"name":       name,
		"securities": symbols,
	}, &resp)

	gid = resp.ID
	return
}

// DeleteWatchlistGroup use to delete watchlist group.
// If `purge` set to true, the securities in the group will be unfollowed even in the *other* groups.
// If `purge` set to false, the securities in the group will remain in the *other* groups.
// Doc: https://open.longportapp.com/en/docs/quote/individual/watchlist_delete_group
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.DeleteWatchlistGroup(context.Background(), 1, false)

func (c *QuoteContext) DeleteWatchlistGroup(ctx context.Context, id int64, purge bool) (err error) {
	var resp struct{}
	err = c.opts.httpClient.Delete(ctx, "/v1/watchlist/groups", nil, &resp, http.WithBody(map[string]interface{}{
		"id":    id,
		"purge": purge,
	}))
	return
}

// UpdateWatchlistGroup use to update watchlist group.
// Doc: https://open.longportapp.com/en/docs/quote/individual/watchlist_update_group
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.UpdateWatchlistGroup(context.Background(), 1, "test", []string{"AAPL.US"}, quote.AddWatchlist)

func (c *QuoteContext) UpdateWatchlistGroup(ctx context.Context, id int64, name string, symbols []string, mode WatchlistUpdateMode) (err error) {
	var resp struct{}
	err = c.opts.httpClient.Put(ctx, "/v1/watchlist/groups", map[string]interface{}{
		"id":         id,
		"name":       name,
		"securities": symbols,
		"mode":       mode,
	}, &resp)
	return
}

// WatchedGroups to get all watchlist groups.
// Reference: https://open.longportapp.com/en/docs/quote/individual/watchlist_groups
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.WatchedGroups(context.Background())

func (c *QuoteContext) WatchedGroups(ctx context.Context) (groupList []*WatchedGroup, err error) {
	var resp jsontypes.WatchedGroupList
	err = c.opts.httpClient.Get(ctx, "/v1/watchlist/groups", nil, &resp)
	if err != nil {
		return
	}

	err = util.Copy(&groupList, resp.Groups)
	return
}

// SecurityList used to list securities. Doc: https://open.longportapp.com/en/docs/quote/security/security_list
func (c *QuoteContext) SecurityList(ctx context.Context, market openapi.Market, category SecurityListCategory) (list []*Security, err error) {
	var resp jsontypes.SecurityList
	values := url.Values{}
	values.Add("market", string(market))
	values.Add("category", string(category))

	err = c.opts.httpClient.Get(ctx, "/v1/quote/get_security_list", values, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&list, resp.List)
	return
}

// Close
func (c *QuoteContext) Close() error {
	return c.core.Close()
}

// Deprecated: NewFromEnv return QuoteContext, use NewFromCfg plz
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
	lbOpts := longbridge.NewOptions(
		longbridge.WithAuthTimeout(cfg.AuthTimeout),
		longbridge.WithTimeout(cfg.Timeout),
		longbridge.WithReadBufferSize(cfg.ReadBufferSize),
		longbridge.WithReadQueueSize(cfg.ReadQueueSize),
		longbridge.WithWriteQueueSize(cfg.WriteQueueSize),
		longbridge.WithMinGzipSize(cfg.MinGzipSize),
	)
	return New(
		WithQuoteURL(cfg.QuoteUrl),
		WithHttpClient(httpClient),
		WithLogLevel(cfg.LogLevel),
		WithLogger(cfg.Logger()),
		WithLbOptions(lbOpts),
	)
}

// New return QuoteContext with option.
// A connection will be created with quote server.
func New(opt ...Option) (*QuoteContext, error) {
	opts := newOptions(opt...)
	core, err := newCore(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &QuoteContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}
