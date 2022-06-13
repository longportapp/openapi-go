package trade

import (
	"context"
	"net/url"

	"github.com/longbridgeapp/openapi-protocol/go/client"
	"github.com/longbridgeapp/openapi-go/config"
)

type TradeContext struct {
	opts         Options
	streamClient client.Client
	core         *Core
}

func (c *TradeContext) Subscribe(ctx context.Context, symbols []string) (subRes *SubResponse, err error) {
	return c.core.Subscribe(ctx, symbols)
}

func (c *TradeContext) Unsubscribe(ctx context.Context, symbols []string) (unsubRes *UnsubResponse, err error) {
	return c.core.Unsubscribe(ctx, symbols)
}

func (c *TradeContext) HistoryExecutions(ctx context.Context, params GetHistoryExecutions) (trades []Execution, err error) {
	type Response struct {
		Trades []Execution
	}
	resp := &Response{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/execution/history", params, trades)
	if err != nil {
		return
	}
	trades = resp.Trades
	return
}

func (c *TradeContext) TodayExecutions(ctx context.Context, params GetTodayExecutions) (trades []Execution, err error) {
	type Response struct {
		Trades []Execution
	}
	resp := &Response{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/execution/today", params, resp)
	if err != nil {
		return
	}
	trades = resp.Trades
	return
}

func (c *TradeContext) HistoryOrders(ctx context.Context, params GetHistoryOrders) (orders []Order, err error) {
	type Response struct {
		Orders []Order
	}
	resp := &Response{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/order/history", params, resp)
	if err != nil {
		return
	}
	orders = resp.Orders
	return
}

func (c *TradeContext) TodayOrders(ctx context.Context, params GetTodayOrders) (orders []Order, err error) {
	type Response struct {
		Orders []Order
	}
	resp := &Response{}
	err = c.opts.HttpClient.Get(ctx, "/v1/trade/order/today", params, resp)
	if err != nil {
		return
	}
	orders = resp.Orders
	return
}

func (c *TradeContext) ReplaceOrder(ctx context.Context, params ReplaceOrder) (err error) {
	err = c.opts.HttpClient.Put(ctx, "/v1/trade/order", params, nil)
	return
}

func (c *TradeContext) SubmitOrder(ctx context.Context, params SubmitOrder) (orderId string, err error) {
	resp := &SubmitOrderResponse{}
	err = c.opts.HttpClient.Post(ctx, "/v1/trade/order", params, resp)
	return resp.OrderId, nil
}

func (c *TradeContext) WithdrawOrder(ctx context.Context, orderId string) (err error) {
	var params url.Values
	params.Add("order_id", orderId)
	err = c.opts.HttpClient.Delete(ctx, "/v1/trade/order", params, nil)
	return
}

func (c *TradeContext) AccountBalance(ctx context.Context) (accounts []AccountBalance, err error) {
	type Response struct {
		List []AccountBalance
	}
	resp := &Response{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/account", nil, resp)
	if err != nil {
		return
	}
	accounts = resp.List
	return
}

func (c *TradeContext) CashFlow(ctx context.Context, params GetCashFlow) (cashflows []CashFlow, err error) {
	type Response struct {
		List []CashFlow
	}
	resp := &Response{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/cashflow", params, cashflows)
	if err != nil {
		return
	}
	cashflows = resp.List
	return
}

func (c *TradeContext) FundPositions(ctx context.Context, symbols []string) (fundPositions *FundPositions, err error) {
	params := &GetFundPositions{
		Symbols: symbols,
	}
	fundPositions = &FundPositions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/fund", params, fundPositions)
	return
}

func (c *TradeContext) StockPositions(ctx context.Context, params GetStockPositions) (stockPositions *StockPositions, err error) {
	stockPositions = &StockPositions{}
	err = c.opts.HttpClient.Get(ctx, "/v1/asset/stock", params, stockPositions)
	return
}

func NewFromCfg(config.Config) *TradeContext {
	return nil
}

func New(opts ...Option) {

}
