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
	OrderTypeSLO     OrderType = "SLO"     // SLO order

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

// Orders has a Order details
type Orders struct {
	HasMore bool
	Orders  []*Order
}

// Order is order details
type Order struct {
	OrderId          string
	Status           OrderStatus
	StockName        string
	Quantity         string
	ExecutedQuantity string
	Price            *decimal.Decimal
	ExecutedPrice    *decimal.Decimal
	SubmittedAt      string
	Side             OrderSide
	Symbol           string
	OrderType        OrderType
	LastDone         *decimal.Decimal
	TriggerPrice     *decimal.Decimal
	Msg              string
	Tag              OrderTag
	TimeInForce      TimeType
	ExpireDate       string
	UpdatedAt        string
	TriggerAt        string
	TrailingAmount   *decimal.Decimal
	TrailingPercent  *decimal.Decimal
	LimitOffset      *decimal.Decimal
	TriggerStatus    TriggerStatus
	Currency         string
	OutsideRth       OutsideRTH
}

// AccountBalances has a AccountBalance list
type AccountBalances struct {
	List []*AccountBalance
}

// AccountBalance is user account balance
type AccountBalance struct {
	TotalCash              *decimal.Decimal
	MaxFinanceAmount       *decimal.Decimal
	RemainingFinanceAmount *decimal.Decimal
	RiskLevel              string
	MarginCall             *decimal.Decimal
	NetAssets              *decimal.Decimal // net asset
	InitMargin             *decimal.Decimal // initial margin
	MaintenanceMargin      *decimal.Decimal // maintenance margin
	Currency               string
	CashInfos              []*CashInfo
}

// FundPositions has a FundPosition list
type FundPositions struct {
	List []*FundPositionChannel
}

// FundPositionChannel is a account channel's fund position details
type FundPositionChannel struct {
	AccountChannel string
	Positions      []*FundPosition
}

// FundPosition is fund position details
type FundPosition struct {
	Symbol               string
	CurrentNetAssetValue *decimal.Decimal
	NetAssetValueDay     int64
	SymbolName           string
	Currency             string
	CostNetAssetValue    *decimal.Decimal
	HoldingUnits         *decimal.Decimal
}

// StockPositions has a StockPosition list
type StockPositions struct {
	List []*StockPositionChannel
}

// StockPositionChannel is a account channel's stock positions details
type StockPositionChannel struct {
	AccountChannel string
	Positions      []*StockPosition
}

// StockPosition is user stock position details
type StockPosition struct {
	Symbol            string
	SymbolName        string
	Quantity          string
	AvailableQuantity string
	Currency          string
	CostPrice         *decimal.Decimal
	Market            openapi.Market
}

// CashFlows has a CashFlow list
type CashFlows struct {
	List []*CashFlow
}

// CashFlow is cash flow details
type CashFlow struct {
	TransactionFlowName string
	Direction           OfDirection
	BusinessType        BalanceType
	Balance             *decimal.Decimal
	Currency            string
	BusinessTime        string
	Symbol              string
	Description         string
}

// CashInfo
type CashInfo struct {
	WithdrawCash  *decimal.Decimal
	AvailableCash *decimal.Decimal
	FrozenCash    *decimal.Decimal
	SettlingCash  *decimal.Decimal
	Currency      string
}

// PushEvent is quote context callback event
type PushEvent struct {
	Event string
	Data  *PushOrderChanged
}

// PushOrderChanged is order change event details
type PushOrderChanged struct {
	Side             OrderSide
	StockName        string
	Quantity         string
	Symbol           string
	OrderType        OrderType
	Price            *decimal.Decimal
	ExecutedQuantity string
	ExecutedPrice    *decimal.Decimal
	OrderId          string
	Currency         string
	Status           OrderStatus
	SubmittedAt      string
	UpdatedAt        string
	TriggerPrice     *decimal.Decimal
	Msg              string
	Tag              OrderTag
	TriggerStatus    TriggerStatus
	TriggerAt        string
	TrailingAmount   *decimal.Decimal
	TrailingPercent  string
	LimitOffset      string
	AccountNo        string
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

// MarginRatio contains some ratio
type MarginRatio struct {
	ImFactor *decimal.Decimal // Initial margin ratio
	MmFactor *decimal.Decimal // Maintain the initial margin ratio
	FmFactor *decimal.Decimal // Forced close-out margin ratio
}
