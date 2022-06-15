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
type CfDirection int32
type TimeType string

const (
	BalanceTypeUnknown BalanceType = 0
	BalanceTypeCash    BalanceType = 1
	BalanceTypeStock   BalanceType = 2
	BalanceTypeFund    BalanceType = 3

	CfDirectionUnkown CfDirection = 0
	CfDirectionOut    CfDirection = 1
	CfDirectionIn     CfDirection = 2

	TimeTypeDay TimeType = "Day"
	TimeTypeGTC TimeType = "GTC"
	TimeTypeGTD TimeType = "GTD"

	OrderTypeLO      OrderType = "LO"
	OrderTypeELO     OrderType = "ELO"
	OrderTypeMO      OrderType = "MO"
	OrderTypeAO      OrderType = "AO"
	OrderTypeALO     OrderType = "ALO"
	OrderTypeODD     OrderType = "ODD"
	OrderTypeLIT     OrderType = "LIT"
	OrderTypeMIT     OrderType = "MIT"
	OrderTypeTSLPAMT OrderType = "TSLPAMT"
	OrderTypeTSLPPCT OrderType = "TSLPPCT"
	OrderTypeTSMAMT  OrderType = "TSMAMT"
	OrderTypeTSMPCT  OrderType = "TSMPCT"

	OrderSideBuy  OrderSide = "Buy"
	OrderSideSell OrderSide = "Sell"
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
	OrderId string
}

type Orders struct {
	Orders []*Order
}

type Order struct {
	OrderId          string
	Status           OrderStatus
	StockName        string
	Quantity         string
	ExecutedQuantity string
	Price            string
	ExecutedPrice    string
	SubmittedAt      string
	Side             OrderSide
	Symbol           string
	OrderType        OrderType
	LastDone         string
	TriggerPrice     string
	Msg              string
	Tag              OrderTag
	TimeInForce      TimeType
	ExpireDate       string
	UpdatedAt        string
	TriggerAt        string
	TrailingAmount   string
	TrailingPercent  string
	LimitOffset      string
	TriggerStatus    TriggerStatus
	Currency         string
	OutsideRth       OutsideRTH
}

type AccountBalances struct {
	List []*AccountBalance
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
	List []*FundPositionChannel
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
	Symbol            string         `json:"symbol"`
	SymbolName        string         `json:"symbol_name"`
	Quantity          int            `json:"quantity"`
	AvailableQuantity int            `json:"available_quantity"`
	Currency          string         `json:"currency"`
	CostPrice         string         `json:"cost_price"`
	Market            openapi.Market `json:"market"`
}

type CashFlows struct {
	List []*CashFlow
}

type CashFlow struct {
	TransactionFlowName string      `json:"transaction_flow_name"`
	Direction           CfDirection `json:"direction"`
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
	orderChanged *PushOrderChanged
}

type PushOrderChanged struct {
	Side             OrderSide     `json:"side"`
	StockName        string        `json:"stock_name"`
	Quantity         string        `json:"quantity"`
	Symbol           string        `json:"symbol"`
	OrderType        OrderType     `json:"order_type"`
	Price            string        `json:"price"`
	ExecutedQuantity int64         `json:"executed_quantity"`
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
