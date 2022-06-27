package trade

import (
	"time"

	"github.com/longbridgeapp/openapi-go"
	"github.com/shopspring/decimal"
)

type OrderType string
type OrderSide string
type OutsideRTH string // Outside regular trading hours
type OrderStatus string
type Market string
type OrderTag string
type TriggerStatus string
type BalanceType int32
type OfDirection int32
type TimeType string

const (
	// Balance type
	BalanceTypeUnknown BalanceType = 0 // unknown type
	BalanceTypeCash    BalanceType = 1 // cash
	BalanceTypeStock   BalanceType = 2 // stock
	BalanceTypeFund    BalanceType = 3 // fund

	// Outflow direction
	OfDirectionUnkown OfDirection = 0
	OfDirectionOut    OfDirection = 1 // outflow
	OfDirectionIn     OfDirection = 2 // inflow

	// Time force in type
	TimeTypeDay TimeType = "Day" // Day Order
	TimeTypeGTC TimeType = "GTC" // Good Til Canceled Order
	TimeTypeGTD TimeType = "GTD" // Good Til Date Order

	// Order type
	OrderTypeLO      OrderType = "LO"      // Limit Order
	OrderTypeELO     OrderType = "ELO"     // Enhanced Limit Order
	OrderTypeMO      OrderType = "MO"      // Market Order
	OrderTypeAO      OrderType = "AO"      // At-auction Order
	OrderTypeALO     OrderType = "ALO"     // At-auction Limit Order
	OrderTypeODD     OrderType = "ODD"     // Odd Lots Order
	OrderTypeLIT     OrderType = "LIT"     // Limit If Touched
	OrderTypeMIT     OrderType = "MIT"     // Market If Touched
	OrderTypeTSLPAMT OrderType = "TSLPAMT" // Trailing Limit If Touched (Trailing Amount)
	OrderTypeTSLPPCT OrderType = "TSLPPCT" // Trailing Limit If Touched (Trailing Percent)
	OrderTypeTSMAMT  OrderType = "TSMAMT"  // Trailing Market If Touched (Trailing Amount)
	OrderTypeTSMPCT  OrderType = "TSMPCT"  // Trailing Market If Touched (Trailing Percent)

	// Order side
	OrderSideBuy  OrderSide = "Buy"
	OrderSideSell OrderSide = "Sell"

	// Outside RTH
	OutsideRTHOnly    OutsideRTH = "RTH_ONLY"          // Regular trading hour only
	OutsideRTHAny     OutsideRTH = "ANY_TIME"          // Any time
	OutsideRTHUnknown OutsideRTH = "UnknownOutsideRth" // Default is UnknownOutsideRth when the order is not a US stock
)

// Execution is execution details
type Execution struct {
	OrderId     string
	TradeId     string
	Symbol      string
	TradeDoneAt time.Time
	Quantity    string
	Price       *decimal.Decimal
}

// Executions has a Execution list
type Executions struct {
	Trades []*Execution
}

type submitOrderResponse struct {
	OrderId string `json:"order_id"`
}

// Orders has a Order details
type Orders struct {
	HasMore bool     `json:"has_more"`
	Orders  []*Order `json:"orders"`
}

// Order is order details
type Order struct {
	OrderId          string           `json:"order_id"`
	Status           OrderStatus      `json:"status"`
	StockName        string           `json:"stock_name"`
	Quantity         string           `json:"quantity"`
	ExecutedQuantity string           `json:"executed_quantity"`
	Price            *decimal.Decimal `json:"price"`
	ExecutedPrice    *decimal.Decimal `json:"executed_price"`
	SubmittedAt      string           `json:"submmited_at"`
	Side             OrderSide        `json:"side"`
	Symbol           string           `json:"symbol"`
	OrderType        OrderType        `json:"order_type"`
	LastDone         *decimal.Decimal `json:"last_done"`
	TriggerPrice     *decimal.Decimal `json:"trigger_price"`
	Msg              string           `json:"msg"`
	Tag              OrderTag         `json:"tag"`
	TimeInForce      TimeType         `json:"time_in_force"`
	ExpireDate       string           `json:"expire_date"`
	UpdatedAt        string           `json:"update_at"`
	TriggerAt        string           `json:"trigger_at"`
	TrailingAmount   *decimal.Decimal `json:"trailing_amount"`
	TrailingPercent  *decimal.Decimal `json:"trailing_percent"`
	LimitOffset      *decimal.Decimal `json:"limit_offset"`
	TriggerStatus    TriggerStatus    `json:"trigger_status"`
	Currency         string           `json:"currency"`
	OutsideRth       OutsideRTH       `json:"outside_rth"`
}

// AccountBalances has a AccountBalance list
type AccountBalances struct {
	List []*AccountBalance `json:"list"`
}

// AccountBalance is user account balance
type AccountBalance struct {
	TotalCash              *decimal.Decimal `json:"total_cash"`
	MaxFinanceAmount       *decimal.Decimal `json:"max_finance_amount"`
	RemainingFinanceAmount *decimal.Decimal `json:"remaining_finance_amount"`
	RiskLevel              string           `json:"risk_level"`
	MarginCall             *decimal.Decimal `json:"margin_call"`
	Currency               string           `json:"currency"`
	CashInfos              []*CashInfo      `json:"cash_infos"`
}

// FundPositions has a FundPosition list
type FundPositions struct {
	List []*FundPositionChannel `json:"list"`
}

// FundPositionChannel is a account channel's fund position details
type FundPositionChannel struct {
	AccountChannel string          `json:"account_channel"`
	Positions      []*FundPosition `json:"fund_info"`
}

// FundPosition is fund position details
type FundPosition struct {
	Symbol               string `json:"symbol"`
	CurrentNetAssetValue string `json:"current_net_asset_value"`
	NetAssetValueDay     int64  `json:"net_asset_value_day,string"` // timestamp
	SymbolName           string `json:"symbol_name"`
	Currency             string `json:"currency"`
	CostNetAssetValue    string `json:"cost_net_asset_value"`
	HoldingUnits         string `json:"holding_units"`
}

// StockPositions has a StockPosition list
type StockPositions struct {
	List []*StockPositionChannel `json:"list"`
}

// StockPositionChannel is a account channel's stock positions details
type StockPositionChannel struct {
	AccountChannel string           `json:"account_channel"`
	Positions      []*StockPosition `json:"stock_info"`
}

// StockPosition is user stock position details
type StockPosition struct {
	Symbol            string           `json:"symbol"`
	SymbolName        string           `json:"symbol_name"`
	Quantity          string           `json:"quantity"`
	AvailableQuantity string           `json:"available_quantity"`
	Currency          string           `json:"currency"`
	CostPrice         *decimal.Decimal `json:"cost_price"`
	Market            openapi.Market   `json:"market"`
}

// CashFlows has a CashFlow list
type CashFlows struct {
	List []*CashFlow `json:"list"`
}

// CashFlow is cash flow details
type CashFlow struct {
	TransactionFlowName string           `json:"transaction_flow_name"`
	Direction           OfDirection      `json:"direction"`
	BusinessType        BalanceType      `json:"business_type"`
	Balance             *decimal.Decimal `json:"balance"`
	Currency            string           `json:"currency"`
	BusinessTime        string           `json:"business_time"`
	Symbol              string           `json:"symbol"`
	Description         string           `json:"description"`
}

// CashInfo
type CashInfo struct {
	WithdrawCash  *decimal.Decimal `json:"withdraw_cash"`
	AvailableCash *decimal.Decimal `json:"avaliable_cash"`
	FrozenCash    *decimal.Decimal `json:"frozen_cash"`
	SettlingCash  *decimal.Decimal `json:"settling_cash"`
	Currency      string           `json:"currency"`
}

// PushEvent is quote context callback event
type PushEvent struct {
	Event string            `json:"event"`
	Data  *PushOrderChanged `json:"data"`
}

// PushOrderChanged is order change event details
type PushOrderChanged struct {
	Side             OrderSide        `json:"side"`
	StockName        string           `json:"stock_name"`
	Quantity         string           `json:"quantity"`
	Symbol           string           `json:"symbol"`
	OrderType        OrderType        `json:"order_type"`
	Price            *decimal.Decimal `json:"price"`
	ExecutedQuantity string           `json:"executed_quantity"`
	ExecutedPrice    *decimal.Decimal `json:"executed_price"`
	OrderId          string           `json:"order_id"`
	Currency         string           `json:"currency"`
	Status           OrderStatus      `json:"status"`
	SubmittedAt      string           `json:"submitted_at"`
	UpdatedAt        string           `json:"update_at"`
	TriggerPrice     *decimal.Decimal `json:"trigger_price"`
	Msg              string           `json:"msg"`
	Tag              OrderTag         `json:"tag"`
	TriggerStatus    TriggerStatus    `json:"trigger_status"`
	TriggerAt        string           `json:"trigger_at"`
	TrailingAmount   *decimal.Decimal `json:"trailing_amount"`
	TrailingPercent  string           `json:"trailing_percent"`
	LimitOffset      string           `json:"limit_offset"`
	AccountNo        string           `json:"account_no"`
}

// SubResponse is subscribe function response
type SubResponse struct {
	Success []string
	Fail    []*SubResponseFail
	Current []string
}

// SubResponseFail contains subscribe failed reason
type SubResponseFail struct {
	Topic  string
	Reason string
}

type UnsubResponse struct {
	Current []string
}
