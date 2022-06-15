package trade

import (
	"context"
	"net/url"

	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/http"
	"github.com/pkg/errors"
)

type TradeContext struct {
	opts *Options
	core *Core
}

func (c *TradeContext) SetOnTrade(f func(*PushEvent)) {
	c.core.SetHandler(f)
}

func (c *TradeContext) Subscribe(ctx context.Context, symbols []string) (subRes *SubResponse, err error) {
	return c.core.Subscribe(ctx, symbols)
}

func (c *TradeContext) Unsubscribe(ctx context.Context, symbols []string) (unsubRes *UnsubResponse, err error) {
	return c.core.Unsubscribe(ctx, symbols)
}

func (c *TradeContext) HistoryExecutions(ctx context.Context, params GetHistoryExecutions) (trades []*Execution, err error) {
	resp := &Executions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/execution/history", params.Values(), trades)
	if err != nil {
		return
	}
	trades = resp.Trades
	return
}

func (c *TradeContext) TodayExecutions(ctx context.Context, params GetTodayExecutions) (trades []*Execution, err error) {
	resp := &Executions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/execution/today", params.Values(), resp)
	if err != nil {
		return
	}
	trades = resp.Trades
	return
}

func (c *TradeContext) HistoryOrders(ctx context.Context, params GetHistoryOrders) (orders []*Order, err error) {
	resp := &Orders{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/order/history", params.Values(), resp)
	if err != nil {
		return
	}
	orders = resp.Orders
	return
}

func (c *TradeContext) TodayOrders(ctx context.Context, params GetTodayOrders) (orders []*Order, err error) {
	resp := &Orders{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/order/today", params.Values(), resp)
	if err != nil {
		return
	}
	orders = resp.Orders
	return
}

func (c *TradeContext) ReplaceOrder(ctx context.Context, params *ReplaceOrder) (err error) {
	err = c.opts.HttpClient.Put(ctx, "/v1/trade/order", params, nil)
	return
}

func (c *TradeContext) SubmitOrder(ctx context.Context, params *SubmitOrder) (orderId string, err error) {
	resp := &SubmitOrderResponse{}
	err = c.opts.HttpClient.Post(ctx, "/v1/trade/order", params, resp)
	if err != nil {
		return
	}
	return resp.OrderId, nil
}

func (c *TradeContext) WithdrawOrder(ctx context.Context, orderId string) (err error) {
	var values url.Values
	values.Add("order_id", orderId)
	err = c.opts.HttpClient.Delete(ctx, "/v1/trade/order", values, nil)
	return
}

func (c *TradeContext) AccountBalance(ctx context.Context) (accounts []*AccountBalance, err error) {
	resp := &AccountBalances{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/account", nil, resp)
	if err != nil {
		return
	}
	accounts = resp.List
	return
}

func (c *TradeContext) CashFlow(ctx context.Context, params GetCashFlow) (cashflows []*CashFlow, err error) {
	resp := &CashFlows{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/cashflow", params.Values(), cashflows)
	if err != nil {
		return
	}
	cashflows = resp.List
	return
}

func (c *TradeContext) FundPositions(ctx context.Context, symbols []string) (fundPositionChannels []*FundPositionChannel, err error) {
	params := &GetFundPositions{
		Symbols: symbols,
	}
	resp := &FundPositions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/fund", params.Values(), resp)
	fundPositionChannels = resp.List
	return

}

func (c *TradeContext) StockPositions(ctx context.Context, symbols []string) (stockPositions *StockPositions, err error) {
	params := &GetStockPositions{
		Symbols: symbols,
	}
	stockPositions = &StockPositions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/stock", params.Values(), stockPositions)
	return
}

func (c *TradeContext) Close() error {
	return c.core.Close()
}

func NewFormEnv() (*TradeContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, err
	}
	return NewFromCfg(cfg)
}

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

func New(opt ...Option) (*TradeContext, error) {
	opts := newOptions(opt...)
	otp, err := opts.HttpClient.GetOTP(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get otp")
	}
	core, err := NewCore(opts.TradeURL, otp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &TradeContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}
