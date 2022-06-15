package trade

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/longbridgeapp/openapi-go"
	"github.com/longbridgeapp/openapi-go/internal/util"
)

type GetHistoryExecutions struct {
	Symbol  string
	StartAt time.Time
	EndAt   time.Time
}

func (req GetHistoryExecutions) Values() url.Values {
	p := &params{}
	p.Add("symbol", req.Symbol)
	p.AddDate("start_at", req.StartAt)
	p.AddDate("end_at", req.EndAt)
	return p.Values()
}

type GetTodayExecutions struct {
	Symbol  string
	OrderId string
}

func (req GetTodayExecutions) Values() url.Values {
	p := &params{}
	p.Add("symbol", req.Symbol)
	p.Add("order_id", req.OrderId)
	return p.Values()
}

type GetHistoryOrders struct {
	Symbol  string
	Status  []OrderStatus
	Side    OrderSide
	Market  openapi.Market
	StartAt int64
	EndAt   int64
}

func (r *GetHistoryOrders) Values() url.Values {
	p := &params{}
	p.Add("symbol", string(r.Symbol))
	p.Add("side", string(r.Side))
	p.Add("market", string(r.Market))
	p.AddInt("start_at", r.StartAt)
	p.AddInt("end_at", r.EndAt)
	vals := p.Values()
	for _, s := range r.Status {
		vals.Add("status", string(s))
	}
	return vals
}

type GetTodayOrders struct {
	Symbol string
	Status []OrderStatus
	Side   OrderSide
	Market openapi.Market
}

func (r *GetTodayOrders) Values() url.Values {
	p := &params{}
	p.Add("symbol", string(r.Symbol))
	p.Add("side", string(r.Side))
	p.Add("market", string(r.Market))
	vals := p.Values()
	for _, s := range r.Status {
		vals.Add("status", string(s))
	}
	return vals
}

type ReplaceOrder struct {
	OrderId         string `json:"order_id"`
	Quantity        int64  `json:"quantity"`
	Price           string `json:"price"`
	TriggerPrice    string `json:"trigger_price"`
	LimitOffset     string `json:"limit_offset"`
	TrailingAmount  string `json:"trailing_ammount"`
	TrailingPercent string `json:"trailing_percent"`
	Remark          string `json:"remark"`
}

type SubmitOrder struct {
	Symbol            string
	OrderType         OrderType
	Side              OrderSide
	SubmittedQuantity int64
	SubmittedPrice    string
	TriggerPrice      string
	LimitOffset       string
	TrailingAmount    string
	TrailingPercent   string
	ExpireDate        time.Time
	OutsideRTH        OutsideRTH
	Remark            string
	TimeInForce       TimeType
}

func (r *SubmitOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Symbol            string `json:"symbol"`
		OrderType         string `json:"order_type"`
		Side              string `json:"side"`
		SubmittedQuantity int64  `json:"submitted_quantity"`
		SubmittedPrice    string `json:"submitted_price"`
		TriggerPrice      string `json:"trigger_price"`
		LimitOffset       string `json:"limit_offset"`
		TrailingAmount    string `json:"trailing_amount"`
		TrailingPercent   string `json:"trailing_percent"`
		ExpireDate        string `json:"expire_date"`
		OutsideRTH        string `json:"outside_rth"`
		Remark            string `json:"remark"`
		TimeInForce       string `json:"time_in_force"`
	}{
		Symbol:            r.Symbol,
		OrderType:         string(r.OrderType),
		Side:              string(r.Side),
		SubmittedQuantity: r.SubmittedQuantity,
		SubmittedPrice:    r.SubmittedPrice,
		TriggerPrice:      r.TriggerPrice,
		LimitOffset:       r.LimitOffset,
		TrailingAmount:    r.TrailingAmount,
		TrailingPercent:   r.TrailingPercent,
		ExpireDate:        util.FormatDate(r.ExpireDate),
		OutsideRTH:        string(r.OutsideRTH),
		Remark:            r.Remark,
		TimeInForce:       string(r.TimeInForce),
	})
}

type GetFundPositions struct {
	Symbols []string
}

func (r *GetFundPositions) Values() url.Values {
	vals := url.Values{}
	for _, s := range r.Symbols {
		vals.Add("symbols", string(s))
	}
	return vals
}

type GetStockPositions struct {
	Symbols []string
}

func (r *GetStockPositions) Values() url.Values {
	vals := url.Values{}
	for _, s := range r.Symbols {
		vals.Add("symbols", string(s))
	}
	return vals
}

type GetCashFlow struct {
	StartAt      time.Time
	EndAt        time.Time
	BusinessType *BalanceType
	Symbol       string
	Page         *int64
	Size         *int64
}

func (r *GetCashFlow) Values() url.Values {
	p := &params{}
	p.Add("symbol", string(r.Symbol))
	p.AddInt("start_at", r.StartAt.Unix())
	p.AddInt("end_at", r.EndAt.Unix())
	p.AddOptInt("page", r.Page)
	if r.BusinessType != nil {
		p.AddInt("business_type", int64(*r.BusinessType))
	}
	return p.Values()
}
