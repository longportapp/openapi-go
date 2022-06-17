package trade

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/longbridgeapp/openapi-go"
	"github.com/longbridgeapp/openapi-go/internal/util"
)

type GetHistoryExecutions struct {
	Symbol  string    // optional
	StartAt time.Time // optional
	EndAt   time.Time // optional
}

func (req *GetHistoryExecutions) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", req.Symbol)
	p.AddDate("start_at", req.StartAt)
	p.AddDate("end_at", req.EndAt)
	return p.Values()
}

type GetTodayExecutions struct {
	Symbol  string // optional
	OrderId string // optional
}

func (req *GetTodayExecutions) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", req.Symbol)
	p.Add("order_id", req.OrderId)
	return p.Values()
}

type GetHistoryOrders struct {
	Symbol  string         // optional
	Status  []OrderStatus  // optional
	Side    OrderSide      // optional
	Market  openapi.Market // optional
	StartAt int64          // optional
	EndAt   int64          // optional
}

func (r *GetHistoryOrders) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
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
	Symbol string         // optional
	Status []OrderStatus  // optional
	Side   OrderSide      // optional
	Market openapi.Market // optional
}

func (r *GetTodayOrders) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
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
	OrderId         string `json:"order_id"`         // required
	Quantity        uint64 `json:"quantity,string"`  // required
	Price           string `json:"price"`            // LO / ELO / ALO / ODD / LIT Order Required
	TriggerPrice    string `json:"trigger_price"`    // LIT / MIT Order Required
	LimitOffset     string `json:"limit_offset"`     // TSLPAMT / TSLPPCT Order Required
	TrailingAmount  string `json:"trailing_ammount"` // TSLPAMT / TSMAMT Order Required
	TrailingPercent string `json:"trailing_percent"` // TSLPPCT / TSMAPCT Order Required
	Remark          string `json:"remark"`
}

type SubmitOrder struct {
	Symbol            string    // required
	OrderType         OrderType // required
	Side              OrderSide // required
	SubmittedQuantity uint64    // required
	SubmittedPrice    string    // LO / ELO / ALO / ODD / LIT Order Required
	TriggerPrice      string    // LIT / MIT Order Required
	LimitOffset       string    // TSLPAMT / TSLPPCT Order Required
	TrailingAmount    string    // TSLPAMT / TSMAMT Order Required
	TrailingPercent   string    // TSLPPCT / TSMAPCT Order Required
	ExpireDate        time.Time // required when time_in_force is GTD
	OutsideRTH        OutsideRTH
	Remark            string
	TimeInForce       TimeType // required
}

func (r *SubmitOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Symbol            string `json:"symbol"`
		OrderType         string `json:"order_type"`
		Side              string `json:"side"`
		SubmittedQuantity uint64 `json:"submitted_quantity,string"`
		SubmittedPrice    string `json:"submitted_price,omitempty"`
		TriggerPrice      string `json:"trigger_price,omitempty"`
		LimitOffset       string `json:"limit_offset,omitempty"`
		TrailingAmount    string `json:"trailing_amount,omitempty"`
		TrailingPercent   string `json:"trailing_percent,omitempty"`
		ExpireDate        string `json:"expire_date,omitempty"`
		OutsideRTH        string `json:"outside_rth,omitempty"`
		Remark            string `json:"remark,omitempty"`
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
		ExpireDate:        util.FormatDate(&r.ExpireDate),
		OutsideRTH:        string(r.OutsideRTH),
		Remark:            r.Remark,
		TimeInForce:       string(r.TimeInForce),
	})
}

type GetFundPositions struct {
	Symbols []string // optional
}

func (r *GetFundPositions) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	vals := url.Values{}
	for _, s := range r.Symbols {
		vals.Add("symbols", string(s))
	}
	return vals
}

type GetStockPositions struct {
	Symbols []string // optional
}

func (r *GetStockPositions) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	vals := url.Values{}
	for _, s := range r.Symbols {
		vals.Add("symbols", string(s))
	}
	return vals
}

type GetCashFlow struct {
	StartAt      int64 // start timestamp , required
	EndAt        int64 // end timestamp, required
	BusinessType BalanceType
	Symbol       string
	Page         int64
	Size         int64
}

func (r *GetCashFlow) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", r.Symbol)
	p.AddInt("start_at", r.StartAt)
	p.AddInt("end_at", r.EndAt)
	p.AddOptInt("page", r.Page)
	p.AddOptInt("size", r.Size)
	p.AddInt("business_type", int64(r.BusinessType))
	return p.Values()
}
