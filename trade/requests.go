package trade

import (
	"time"

	"github.com/longbridgeapp/openapi-go"
)

type GetHistoryExecutions struct {
	Symbol  string
	StartAt *time.Time
	EndAt   *time.Time
}

type GetTodayExecutions struct {
	Symbol  string
	OrderId string
}

type GetHistoryOrders struct {
	Symbol  string
	Status  []OrderStatus
	Side    []OrderSide
	Market  []openapi.Market
	StartAt *time.Time
	EndAt   *time.Time
}

type GetTodayOrders struct {
	Symbol string
	Status string
	Side   string
	Market string
}

type ReplaceOrder struct {
	OrderId         string
	Quantity        int64
	Price           string
	TriggerPrice    string
	LimitOffset     string
	TrailingAmount  string
	TrailingPercent string
	Remark          string
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
	ExpireDate        *time.Time
	OutsideRTH        *OutsideRTH
	Remark            string
}

type GetFundPositions struct {
	Symbols []string
}

type GetStockPositions struct {
	Symbols []string
}

type GetCashFlow struct {
	StartAt      time.Time
	EndAt        time.Time
	BusinessType *BalanceType
	Symbol       string
	Page         *int
	Size         *int
}
