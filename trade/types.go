package trade

import "time"

type OrderType string
type OrderSide string
type TimeInForce string
type OutsideRTH string // Outside regular trading hours
type OrderStatus string
type Market string
type OrderTag string
type TriggerStatus string

type Execution struct  {
    // Order ID
    OrderId string
    // Execution ID
    TradeId string
    // Security code
    Symbol string
    // Trade done time
    TradeDoneAt time.Time
    // Executed quantity
    Quantity string
    // Executed price
    Price string
}

type SubmitOrderResponse struct {
    OrderId string
}


type Order struct {
    // Order ID
    OrderId string
    // Order status
    Status OrderStatus
    // Stock name
    StockName string
    // Submitted quantity
    Quantity string
    // Executed quantity
    ExecutedQuantity string
    // Submitted price
    Price string
    // Executed price
    ExecutedPrice string
    // Submitted time
    SubmittedAt time.Time
    // Order side
    Side OrderSide
    // Security code
    Symbol string
    // Order type
    OrderType OrderType
    // Last done
    LastDone string
    // `LIT` / `MIT` Order Trigger Price
    TriggerPrice string
    // Rejected Message or remark
    Msg string
    // Order ta
    Tag OrderTag
    // Time in force type
    TimeInForce TimeInForceType
    // Long term order expire date
    ExpireDate string
    // Last updated time
    UpdatedAt time.Time
    // Conditional order trigger time
    TriggerAt time.Time
    // `TSMAMT` / `TSLPAMT` order trailing amount
    TrailingAmount string
    // `TSMPCT` / `TSLPPCT` order trailing percent
    TrailingPercent string
    // `TSLPAMT` / `TSLPPCT` order limit offset amount
    LimitOffset string
    // Conditional order trigger status
    TriggerStatus TriggerStatus
    // Currency
    Currency string
    // Enable or disable outside regular trading hours
    OutsideRth OutsideRTH
}


type AccountBalance struct {

}

type FundPositions struct {

}

type StockPositions struct {

}

type CashFlow struct {

}

type PushEvent struct {
	orderChanged *PushOrderChanged
}

type PushOrderChanged struct {
    /// Order side
    Side OrderSide
    /// Stock name
    StockName string
    /// Submitted quantity
    Quantity string
    /// Order symbol
    Symbol string
    /// Order type
    OrderType OrderType
    /// Submitted price
    Price string
    /// Executed quantity
    ExecutedQuantity int64
    /// Executed price
    ExecutedPrice string
    /// Order ID
    OrderId string
    /// Currency
    Currency string
    /// Order status
    Status OrderStatus
    /// Submitted time
    SubmittedAt time.Time
    /// Last updated time
    UpdatedAt time.Time
    /// Order trigger price
    TriggerPrice string
    /// Rejected message or remark
    Msg string
    /// Order tag
    Tag OrderTag
    /// Conditional order trigger status
    TriggerStatus *TriggerStatus
    /// Conditional order trigger time
    TriggerAt time.Time
    /// Trailing amount
    TrailingAmount string
    /// Trailing percent
    TrailingPercent string
    /// Limit offset amount
    LimitOffset string
    // Account no
    AccountNo string
}

type SubResponse struct {
    Success []string
    Fail    []*SubResponseFail
    Current []string
}

type SubResponseFail struct {
    Topic string
    Reason string
}

type UnsubResponse struct {
    Current []string
}
