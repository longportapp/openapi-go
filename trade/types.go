package trade

import (
	"time"

	"github.com/longbridgeapp/openapi-go"
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

type Execution struct {
	OrderId     string
	TradeId     string
	Symbol      string
	TradeDoneAt time.Time
	Quantity    string
	Price       string
}

type Executions struct {
	Trades []*Execution
}

type SubmitOrderResponse struct {
	OrderId string `json:"order_id"`
}

type Orders struct {
	HasMore bool     `json:"has_more"`
	Orders  []*Order `json:"orders"`
}

type Order struct {
	OrderId          string        `json:"order_id"`
	Status           OrderStatus   `json:"status"`
	StockName        string        `json:"stock_name"`
	Quantity         string        `json:"quantity"`
	ExecutedQuantity string        `json:"executed_quantity"`
	Price            string        `json:"price"`
	ExecutedPrice    string        `json:"executed_price"`
	SubmittedAt      string        `json:"submmited_at"`
	Side             OrderSide     `json:"side"`
	Symbol           string        `json:"symbol"`
	OrderType        OrderType     `json:"order_type"`
	LastDone         string        `json:"last_done"`
	TriggerPrice     string        `json:"trigger_price"`
	Msg              string        `json:"msg"`
	Tag              OrderTag      `json:"tag"`
	TimeInForce      TimeType      `json:"time_in_force"`
	ExpireDate       string        `json:"expire_date"`
	UpdatedAt        string        `json:"update_at"`
	TriggerAt        string        `json:"trigger_at"`
	TrailingAmount   string        `json:"trailing_amount"`
	TrailingPercent  string        `json:"trailing_percent"`
	LimitOffset      string        `json:"limit_offset"`
	TriggerStatus    TriggerStatus `json:"trigger_status"`
	Currency         string        `json:"currency"`
	OutsideRth       OutsideRTH    `json:"outside_rth"`
}

type AccountBalances struct {
	List []*AccountBalance `json:"list"`
}

type AccountBalance struct {
	TotalCash              string      `json:"total_cash"`
	MaxFinanceAmount       string      `json:"max_finance_amount"`
	RemainingFinanceAmount string      `json:"remaining_finance_amount"`
	RiskLevel              string      `json:"risk_level"`
	MarginCall             string      `json:"margin_call"`
	Currency               string      `json:"currency"`
	CashInfos              []*CashInfo `json:"cash_infos"`
}

type FundPositions struct {
	List []*FundPositionChannel `json:"list"`
}

type FundPositionChannel struct {
	AccountChannel string          `json:"account_channel"`
	Positions      []*FundPosition `json:"positions"`
}

type FundPosition struct {
	Symbol               string     `json:"symbol"`
	CurrentNetAssetValue string     `json:"current_net_assset_value"`
	NetAssetValueDay     *time.Time `json:"net_asseet_value_day"`
	SymbolName           string     `json:"symbol_name"`
	Currency             string     `json:"currency"`
	CostNetAssetValue    string     `json:"cost_net_asset_value"`
	HoldingUnits         string     `json:"holding_units"`
}

type StockPositions struct {
	List []*StockPositionChannel `json:"list"`
}

type StockPositionChannel struct {
	AccountChannel string           `json:"account_channel"`
	Positions      []*StockPosition `json:"stock_info"`
}

type StockPosition struct {
	Symbol            string         `json:"symbol"`
	SymbolName        string         `json:"symbol_name"`
	Quantity          string         `json:"quantity"`
	AvailableQuantity string         `json:"available_quantity"`
	Currency          string         `json:"currency"`
	CostPrice         string         `json:"cost_price"`
	Market            openapi.Market `json:"market"`
}

type CashFlows struct {
	List []*CashFlow `json:"list"`
}

type CashFlow struct {
	TransactionFlowName string      `json:"transaction_flow_name"`
	Direction           OfDirection `json:"direction"`
	BusinessType        BalanceType `json:"business_type"`
	Balance             string      `json:"balance"`
	Currency            string      `json:"currency"`
	BusinessTime        string      `json:"business_time"`
	Symbol              string      `json:"symbol"`
	Description         string      `json:"description"`
}

type CashInfo struct {
	WithdrawCash  string `json:"withdraw_cash"`
	AvailableCash string `json:"avaliable_cash"`
	FrozenCash    string `json:"frozen_cash"`
	SettlingCash  string `json:"settling_cash"`
	Currency      string `json:"currency"`
}

type PushEvent struct {
	Event string            `json:"event"`
	Data  *PushOrderChanged `json:"data"`
}

type PushOrderChanged struct {
	Side             OrderSide     `json:"side"`
	StockName        string        `json:"stock_name"`
	Quantity         string        `json:"quantity"`
	Symbol           string        `json:"symbol"`
	OrderType        OrderType     `json:"order_type"`
	Price            string        `json:"price"`
	ExecutedQuantity string        `json:"executed_quantity"`
	ExecutedPrice    string        `json:"executed_price"`
	OrderId          string        `json:"order_id"`
	Currency         string        `json:"currency"`
	Status           OrderStatus   `json:"status"`
	SubmittedAt      string        `json:"submitted_at"`
	UpdatedAt        string        `json:"update_at"`
	TriggerPrice     string        `json:"trigger_price"`
	Msg              string        `json:"msg"`
	Tag              OrderTag      `json:"tag"`
	TriggerStatus    TriggerStatus `json:"trigger_status"`
	TriggerAt        string        `json:"trigger_at"`
	TrailingAmount   string        `json:"trailing_amount"`
	TrailingPercent  string        `json:"trailing_percent"`
	LimitOffset      string        `json:"limit_offset"`
	AccountNo        string        `json:"account_no"`
}

type SubResponse struct {
	Success []string
	Fail    []*SubResponseFail
	Current []string
}

type SubResponseFail struct {
	Topic  string
	Reason string
}

type UnsubResponse struct {
	Current []string
}
