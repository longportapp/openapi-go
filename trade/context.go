// Package trade provide TradeContext
package trade

import (
	"context"
	"net/url"

	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/http"
	"github.com/longbridgeapp/openapi-go/internal/util"
	"github.com/longbridgeapp/openapi-go/longbridge"
	"github.com/longbridgeapp/openapi-go/trade/jsontypes"

	"github.com/pkg/errors"
)

// TradeContext is a client for interacting with Longbridge Trade OpenAPI.
// Longbrige Quote OpenAPI document is https://open.longbridgeapp.com/en/docs/trade/trade-overview
type TradeContext struct {
	opts *Options
	core *core
}

// OnQuote set callback function which will be called when server push events.
func (c *TradeContext) OnTrade(f func(*PushEvent)) {
	c.core.SetHandler(f)
}

// Subscribe topics then the handler will receive push event.
// Reference: https://open.longbridgeapp.com/en/docs/trade/trade-push#subscribe
func (c *TradeContext) Subscribe(ctx context.Context, topics []string) (subRes *SubResponse, err error) {
	return c.core.Subscribe(ctx, topics)
}

// Unsubscribe topics then the handler will not receive the symbol's event.
// Reference: https://open.longbridgeapp.com/en/docs/trade/trade-push#cancel-subscribe
func (c *TradeContext) Unsubscribe(ctx context.Context, topics []string) (unsubRes *UnsubResponse, err error) {
	return c.core.Unsubscribe(ctx, topics)
}

// HistoryExecutions will return history executions.
// Reference: https://open.longbridgeapp.com/en/docs/trade/execution/history_executions
func (c *TradeContext) HistoryExecutions(ctx context.Context, params *GetHistoryExecutions) (trades []*Execution, err error) {
	resp := &jsontypes.Executions{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/execution/history", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&trades, resp.Trades)
	return
}

// TodayExecutions will return today's executions
// Reference: https://open.longbridgeapp.com/en/docs/trade/execution/today_executions
func (c *TradeContext) TodayExecutions(ctx context.Context, params *GetTodayExecutions) (trades []*Execution, err error) {
	resp := &jsontypes.Executions{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/execution/today", params.Values(), resp)
	if err != nil {
		return
	}
	err = util.Copy(&trades, resp.Trades)
	return
}

// HistoryOrders will return history orders
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/history_orders
func (c *TradeContext) HistoryOrders(ctx context.Context, params *GetHistoryOrders) (orders []*Order, hasMore bool, err error) {
	resp := &jsontypes.Orders{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/order/history", params.Values(), resp)
	if err != nil {
		return
	}
	hasMore = resp.HasMore
	err = util.Copy(&orders, resp.Orders)
	return
}

// TodayOrders will return today orders
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/today_orders
func (c *TradeContext) TodayOrders(ctx context.Context, params *GetTodayOrders) (orders []*Order, err error) {
	resp := &jsontypes.Orders{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/order/today", params.Values(), resp)
	if err != nil {
		return
	}
	err = util.Copy(&orders, resp.Orders)
	return
}

// ReplaceOrder modify quantity or price
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/replace
func (c *TradeContext) ReplaceOrder(ctx context.Context, params *ReplaceOrder) (err error) {
	var jsonbody jsontypes.ReplaceOrder
	err = util.Copy(&jsonbody, params)
	if err != nil {
		return
	}
	err = c.opts.httpClient.Put(ctx, "/v1/trade/order", jsonbody, nil)
	return
}

// SubmitOrder HK and US stocks, warrant and option
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/submit
func (c *TradeContext) SubmitOrder(ctx context.Context, params *SubmitOrder) (orderId string, err error) {
	var jsonbody jsontypes.SubmitOrder
	err = util.Copy(&jsonbody, params)
	if err != nil {
		return
	}
	resp := &jsontypes.SubmitOrderResponse{}
	err = c.opts.httpClient.Post(ctx, "/v1/trade/order", jsonbody, resp)
	if err != nil {
		return
	}
	return resp.OrderId, nil
}

// WithdrawOrder to close an open order
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/withdraw
func (c *TradeContext) WithdrawOrder(ctx context.Context, orderId string) (err error) {
	values := url.Values{}
	values.Add("order_id", orderId)
	err = c.opts.httpClient.Delete(ctx, "/v1/trade/order", values, nil)
	return
}

// AccountBalance to obtain the available, desirable, frozen, to-be-settled, and in-transit funds (fund purchase and redemption) information for each currency of the user.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/account
func (c *TradeContext) AccountBalance(ctx context.Context) (accounts []*AccountBalance, err error) {
	var resp jsontypes.AccountBalances
	err = c.opts.httpClient.Get(ctx, "/v1/asset/account", nil, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&accounts, resp.List)
	return
}

// CashFlow to obtain capital inflow/outflow direction, capital type, capital amount, occurrence time, associated stock code and capital flow description information.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/cashflow
func (c *TradeContext) CashFlow(ctx context.Context, params *GetCashFlow) (cashflows []*CashFlow, err error) {
	var resp jsontypes.CashFlows
	err = c.opts.httpClient.Get(ctx, "/v1/asset/cashflow", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&cashflows, resp.List)
	return
}

// FundPositions to obtain fund position information including account, fund code, holding share, cost net worth, current net worth, and currency.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/fund
func (c *TradeContext) FundPositions(ctx context.Context, symbols []string) (fundPositionChannels []*FundPositionChannel, err error) {
	params := &GetFundPositions{
		Symbols: symbols,
	}
	var resp jsontypes.FundPositions
	err = c.opts.httpClient.Get(ctx, "/v1/asset/fund", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&fundPositionChannels, resp.List)
	return

}

// StockPositions to obtain stock position information including account, stock code, number of shares held, number of available shares, average position price (calculated according to account settings), and currency.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/stock
func (c *TradeContext) StockPositions(ctx context.Context, symbols []string) (stockPositionChannels []*StockPositionChannel, err error) {
	params := &GetStockPositions{
		Symbols: symbols,
	}
	var resp jsontypes.StockPositions
	err = c.opts.httpClient.Get(ctx, "/v1/asset/stock", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&stockPositionChannels, resp.List)
	return
}

// MarginRatio is used to obtain the initial margin ratio, maintain the margin ratio and strengthen the margin ratio of stocks.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/margin_ratio
func (c *TradeContext) MarginRatio(ctx context.Context, symbol string) (marginRatio MarginRatio, err error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	var resp jsontypes.MarginRatio
	err = c.opts.httpClient.Get(ctx, "/v1/risk/margin-ratio", values, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&marginRatio, resp)
	return
}

// Close
func (c *TradeContext) Close() error {
	return c.core.Close()
}

// NewFormEnv return TradeContext with environment variables.
func NewFormEnv() (*TradeContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, err
	}
	return NewFromCfg(cfg)
}

// NewFromCfg return TradeContext with config.Config.
func NewFromCfg(cfg *config.Config) (*TradeContext, error) {
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
		longbridge.WithAuthTimeout(cfg.TradeLBAuthTimeout),
		longbridge.WithTimeout(cfg.TradeLBTimeout),
		longbridge.WithReadBufferSize(cfg.TradeLBReadBufferSize),
		longbridge.WithReadQueueSize(cfg.TradeLBReadQueueSize),
		longbridge.WithWriteQueueSize(cfg.TradeLBWriteQueueSize),
		longbridge.WithMinGzipSize(cfg.TradeLBMinGzipSize),
	)
	return New(
		WithTradeURL(cfg.TradeUrl),
		WithHttpClient(httpClient),
    WithLbOptions(lbOpts),
		WithLogLevel(cfg.LogLevel),
	)
}

// New return TradeContext with option.
// A connection will be created with Trade server.
func New(opt ...Option) (*TradeContext, error) {
	opts := newOptions(opt...)
	core, err := newCore(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &TradeContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}
