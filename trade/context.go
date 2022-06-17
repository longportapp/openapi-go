//  Package trade provide TradeContext
package trade

import (
	"context"
	"net/url"

	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/http"
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
	resp := &Executions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/execution/history", params.Values(), &resp)
	if err != nil {
		return
	}
	trades = resp.Trades
	return
}

// TodayExecutions will return today's executions
// Reference: https://open.longbridgeapp.com/en/docs/trade/execution/today_executions
func (c *TradeContext) TodayExecutions(ctx context.Context, params *GetTodayExecutions) (trades []*Execution, err error) {
	resp := &Executions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/execution/today", params.Values(), resp)
	if err != nil {
		return
	}
	trades = resp.Trades
	return
}

// HistoryOrders will return history orders
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/history_orders
func (c *TradeContext) HistoryOrders(ctx context.Context, params *GetHistoryOrders) (orders []*Order, hasMore bool, err error) {
	resp := &Orders{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/order/history", params.Values(), resp)
	if err != nil {
		return
	}
	orders = resp.Orders
	hasMore = resp.HasMore
	return
}

// TodayOrders will return today orders
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/today_orders
func (c *TradeContext) TodayOrders(ctx context.Context, params *GetTodayOrders) (orders []*Order, err error) {
	resp := &Orders{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/order/today", params.Values(), resp)
	if err != nil {
		return
	}
	orders = resp.Orders
	return
}

// ReplaceOrder modify quantity or price
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/replace
func (c *TradeContext) ReplaceOrder(ctx context.Context, params *ReplaceOrder) (err error) {
	err = c.opts.HttpClient.Put(ctx, "/v1/trade/order", params, nil)
	return
}

// SubmitOrder HK and US stocks, warrant and option
// Reference: https://open.longbridgeapp.com/en/docs/trade/order/submit
func (c *TradeContext) SubmitOrder(ctx context.Context, params *SubmitOrder) (orderId string, err error) {
	resp := &submitOrderResponse{}
	err = c.opts.HttpClient.Post(ctx, "/v1/trade/order", params, resp)
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
	err = c.opts.HttpClient.Delete(ctx, "/v1/trade/order", values, nil)
	return
}

// AccountBalance to obtain the available, desirable, frozen, to-be-settled, and in-transit funds (fund purchase and redemption) information for each currency of the user.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/account
func (c *TradeContext) AccountBalance(ctx context.Context) (accounts []*AccountBalance, err error) {
	var resp AccountBalances
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/account", nil, &resp)
	if err != nil {
		return
	}
	accounts = resp.List
	return
}

// CashFlow to obtain capital inflow/outflow direction, capital type, capital amount, occurrence time, associated stock code and capital flow description information.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/cashflow
func (c *TradeContext) CashFlow(ctx context.Context, params *GetCashFlow) (cashflows []*CashFlow, err error) {
	var resp CashFlows
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/cashflow", params.Values(), &resp)
	if err != nil {
		return
	}
	cashflows = resp.List
	return
}

// FundPositions to obtain fund position information including account, fund code, holding share, cost net worth, current net worth, and currency.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/fund
func (c *TradeContext) FundPositions(ctx context.Context, symbols []string) (fundPositionChannels []*FundPositionChannel, err error) {
	params := &GetFundPositions{
		Symbols: symbols,
	}
	var resp FundPositions
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/fund", params.Values(), &resp)
	fundPositionChannels = resp.List
	return

}

// StockPositions to obtain stock position information including account, stock code, number of shares held, number of available shares, average position price (calculated according to account settings), and currency.
// Reference: https://open.longbridgeapp.com/en/docs/trade/asset/stock
func (c *TradeContext) StockPositions(ctx context.Context, symbols []string) (stockPositionChannels []*StockPositionChannel, err error) {
	params := &GetStockPositions{
		Symbols: symbols,
	}
	var resp StockPositions
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/stock", params.Values(), &resp)
	stockPositionChannels = resp.List
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
	return New(WithTradeURL(cfg.TradeUrl), WithHttpClient(httpClient))
}

// New return TradeContext with option.
// A connection will be created with Trade server.
func New(opt ...Option) (*TradeContext, error) {
	opts := newOptions(opt...)
	otp, err := opts.HttpClient.GetOTP(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get otp")
	}
	core, err := newCore(opts.TradeURL, otp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &TradeContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}
