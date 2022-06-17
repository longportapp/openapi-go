package quote

import (
	"time"

	"github.com/longbridgeapp/openapi-go"
	"github.com/longbridgeapp/openapi-protobufs/gen/go/quote"
)

type TradeStatus int32
type TradeSession int32
type TradeSessionType int32
type EventType int8
type SubType uint8
type Period int32
type AdjustType int32

const (
	// SubType
	SubTypeUnknown SubType = SubType(quotev1.SubType_UNKNOWN_TYPE)
	SubTytpeQuote  SubType = SubType(quotev1.SubType_QUOTE)
	SubTypeDepth   SubType = SubType(quotev1.SubType_DEPTH)
	SubTypeBrokers SubType = SubType(quotev1.SubType_BROKERS)
	SubTypeTrade   SubType = SubType(quotev1.SubType_TRADE)

	// SubEvent
	EventQuote EventType = iota
	EventBroker
	EventTrade
	EventDepth

	// Period
	PeriodOneMinute     = Period(quotev1.Period_ONE_MINUTE)
	PeriodFiveMinute    = Period(quotev1.Period_FIVE_MINUTE)
	PeriodFifteenMinute = Period(quotev1.Period_FIFTEEN_MINUTE)
	PeriodThirtyMinute  = Period(quotev1.Period_THIRTY_MINUTE)
	PeriodSixtyMinute   = Period(quotev1.Period_SIXTY_MINUTE)
	PeriodDay           = Period(quotev1.Period_DAY)
	PeriodWeek          = Period(quotev1.Period_WEEK)
	PeriodMonth         = Period(quotev1.Period_MONTH)
	PeriodYear          = Period(quotev1.Period_YEAR)

	// AdjustType
	AdjustTypeNo      = AdjustType(quotev1.AdjustType_NO_ADJUST)
	AdjustTypeForward = AdjustType(quotev1.AdjustType_FORWARD_ADJUST)
)

type QotSubscription struct {
	Symbol        string
	Subscriptions []SubType
}

type PushEvent struct {
	Type     EventType
	Symbol   string
	Sequence int64
	Quote    *PushQuote
	Depth    *PushDepth
	Brokers  *PushBrokers
	Trade    *PushTrade
}

type PushQuote struct {
	Symbol       string
	Sequence     int64
	LastDone     string
	Open         string
	High         string
	Low          string
	Timestamp    int64
	Volume       int64
	Turnover     string
	TradeStatus  TradeStatus
	TradeSession TradeSessionType
}

func toPushQuote(origin *quotev1.PushQuote) *PushQuote {
	return &PushQuote{
		Symbol:       origin.GetSymbol(),
		Sequence:     origin.GetSequence(),
		LastDone:     origin.GetLastDone(),
		Open:         origin.GetOpen(),
		High:         origin.GetHigh(),
		Low:          origin.GetLow(),
		Timestamp:    origin.GetTimestamp(),
		Volume:       origin.GetVolume(),
		Turnover:     origin.GetTurnover(),
		TradeStatus:  TradeStatus(origin.GetTradeStatus()),
		TradeSession: TradeSessionType(origin.GetTradeSession()),
	}
}

type PushDepth struct {
	Symbol   string
	Sequence int64
	Ask      []*Depth
	Bid      []*Depth
}

func toPushDepth(origin *quotev1.PushDepth) *PushDepth {
	return &PushDepth{
		Symbol:   origin.GetSymbol(),
		Sequence: origin.GetSequence(),
		Ask:      toDepths(origin.GetAsk()),
		Bid:      toDepths(origin.GetBid()),
	}
}

type PushBrokers struct {
	Symbol     string
	Sequence   int64
	AskBrokers []*Brokers
	BidBrokers []*Brokers
}

func toPushBrokers(origin *quotev1.PushBrokers) *PushBrokers {
	return &PushBrokers{
		Symbol:     origin.GetSymbol(),
		Sequence:   origin.GetSequence(),
		AskBrokers: toBrokers(origin.GetAskBrokers()),
		BidBrokers: toBrokers(origin.GetBidBrokers()),
	}
}

type PushTrade struct {
	Symbol   string
	Sequence int64
	Trade    []*Trade
}

func toPushTrades(origin *quotev1.PushTrade) *PushTrade {
	return &PushTrade{
		Symbol:   origin.GetSymbol(),
		Sequence: origin.GetSequence(),
		Trade:    toTrades(origin.GetTrade()),
	}
}

type Depth struct {
	Position int32
	Price    string
	Volume   int64
	OrderNum int64
}

func toDepth(origin *quotev1.Depth) *Depth {
	return &Depth{
		Position: origin.GetPosition(),
		Price:    origin.GetPrice(),
		Volume:   origin.GetVolume(),
		OrderNum: origin.GetOrderNum(),
	}
}

func toDepths(origin []*quotev1.Depth) (depths []*Depth) {
	depths = make([]*Depth, 0, len(origin))
	for _, item := range origin {
		depths = append(depths, toDepth(item))
	}
	return
}

type Brokers struct {
	Position  int32
	BrokerIds []int32
}

func toBrokers(origin []*quotev1.Brokers) (brokers []*Brokers) {
	brokers = make([]*Brokers, 0, len(origin))
	for _, item := range origin {
		brokers = append(brokers, &Brokers{
			Position:  item.GetPosition(),
			BrokerIds: item.GetBrokerIds(),
		})
	}
	return
}

type Trade struct {
	Price  string
	Volume int64
	// Timestamp of trading
	Timestamp int64
	// TradeType
	// HK
	//
	// - `*` - Overseas trade
	// - `D` - Odd-lot trade
	// - `M` - Non-direct off-exchange trade
	// - `P` - Late trade (Off-exchange previous day)
	// - `U` - Auction trade
	// - `X` - Direct off-exchange trade
	// - `Y` - Automatch internalized
	// - `<empty string>` -  Automatch normal
	//
	// US
	//
	// - `<empty string>` - Regular sale
	// - `A` - Acquisition
	// - `B` - Bunched trade
	// - `D` - Distribution
	// - `F` - Intermarket sweep
	// - `G` - Bunched sold trades
	// - `H` - Price variation trade
	// - `I` - Odd lot trade
	// - `K` - Rule 155 trde(NYSE MKT)
	// - `M` - Market center close price
	// - `P` - Prior reference price
	// - `Q` - Market center open price
	// - `S` - Split trade
	// - `V` - Contingent trade
	// - `W` - Average price trade
	// - `X` - Cross trade
	// - `1` - Stopped stock(Regular trade)
	TradeType    string
	Direction    int32
	TradeSession TradeSession
}

func toTrades(origin []*quotev1.Trade) (trades []*Trade) {
	trades = make([]*Trade, 0, len(origin))
	for _, item := range origin {
		trades = append(trades, &Trade{
			Price:        item.GetPrice(),
			Volume:       item.GetVolume(),
			Timestamp:    item.GetTimestamp(),
			TradeType:    item.GetTradeType(),
			Direction:    item.GetDirection(),
			TradeSession: TradeSession(item.GetTradeSession()),
		})
	}
	return
}

type StaticInfo struct {
	Symbol            string
	NameCn            string
	NameEn            string
	NameHk            string
	Exchange          string
	Currency          string
	LotSize           int32
	TotalShares       int64
	CirculatingShares int64
	HkShares          int64
	Eps               string
	EpsTtm            string
	Bps               string
	DividendYield     string
	StockDerivatives  []int32
}

func toStaticInfos(origin []*quotev1.StaticInfo) (staticInfos []*StaticInfo) {
	staticInfos = make([]*StaticInfo, 0, len(origin))
	for _, item := range origin {
		staticInfos = append(staticInfos, &StaticInfo{
			Symbol:            item.GetSymbol(),
			NameCn:            item.GetNameCn(),
			NameEn:            item.GetNameEn(),
			NameHk:            item.GetNameHk(),
			Exchange:          item.GetExchange(),
			Currency:          item.GetCurrency(),
			LotSize:           item.GetLotSize(),
			TotalShares:       item.GetTotalShares(),
			CirculatingShares: item.GetCirculatingShares(),
			HkShares:          item.GetHkShares(),
			Eps:               item.GetEps(),
			EpsTtm:            item.GetEpsTtm(),
			Bps:               item.GetBps(),
			DividendYield:     item.GetDividendYield(),
			StockDerivatives:  item.GetStockDerivatives(),
		})
	}
	return
}

type Issuer struct {
	ID     int32
	NameCn string
	NameEn string
	NameHk string
}

type RealtimeOptionQuote struct {
	Symbol       string
	LastDone     float64
	PrevClose    float64
	Open         float64
	High         float64
	Low          float64
	Timestamp    int64
	Volume       int64
	Turnover     float64
	OptionExtend *OptionExtend
}

type OptionQuote struct {
	Symbol       string
	LastDone     string
	PrevClose    string
	Open         string
	High         string
	Low          string
	Timestamp    int64
	Volume       int64
	Turnover     string
	TradeStatus  TradeStatus
	OptionExtend *OptionExtend
}

func toOptionQuotes(origin []*quotev1.OptionQuote) (quotes []*OptionQuote) {
	quotes = make([]*OptionQuote, 0, len(origin))
	for _, item := range origin {
		quotes = append(quotes, &OptionQuote{
			Symbol:       item.GetSymbol(),
			LastDone:     item.GetLastDone(),
			Open:         item.GetOpen(),
			High:         item.GetHigh(),
			Low:          item.GetLow(),
			Timestamp:    item.GetTimestamp(),
			Volume:       item.GetVolume(),
			Turnover:     item.GetTurnover(),
			TradeStatus:  TradeStatus(item.GetTradeStatus()),
			OptionExtend: toOptionExtend(item.GetOptionExtend()),
		})
	}
	return
}

type OptionExtend struct {
	ImpliedVolatility    string
	OpenInterest         int64
	ExpiryDate           string // YYMMDD
	StrikePrice          string
	ContractMultiplier   string
	ContractType         string
	ContractSize         string
	Direction            string
	HistoricalVolatility string
	UnderlyingSymbol     string
}

func toOptionExtend(origin *quotev1.OptionExtend) *OptionExtend {
	return &OptionExtend{
		ImpliedVolatility:    origin.GetImpliedVolatility(),
		OpenInterest:         origin.GetOpenInterest(),
		ExpiryDate:           origin.GetExpiryDate(),
		StrikePrice:          origin.GetStrikePrice(),
		ContractMultiplier:   origin.GetContractMultiplier(),
		ContractType:         origin.GetContractType(),
		ContractSize:         origin.GetContractSize(),
		Direction:            origin.GetDirection(),
		HistoricalVolatility: origin.GetHistoricalVolatility(),
		UnderlyingSymbol:     origin.GetUnderlyingSymbol(),
	}
}

type StrikePriceInfo struct {
	Price      string
	CallSymbol string
	PutSymbol  string
	Standard   bool
}

func toStrikePriceInfos(origin []*quotev1.StrikePriceInfo) (priceInfos []*StrikePriceInfo) {
	priceInfos = make([]*StrikePriceInfo, 0, len(origin))
	for _, item := range origin {
		priceInfos = append(priceInfos, &StrikePriceInfo{
			Price:      item.GetPrice(),
			CallSymbol: item.GetCallSymbol(),
			PutSymbol:  item.GetPutSymbol(),
			Standard:   item.GetStandard(),
		})
	}
	return priceInfos
}

type WarrantExtend struct {
	ImpliedVolatility string
	ExpiryDate        string
	LastTradeDate     string
	OutstandingRatio  string
	OutstandingQty    int64
	ConversionRatio   string
	Category          string
	StrikePrice       string
	UpperStrikePrice  string
	LowerStrikePrice  string
	CallPrice         string
	UnderlyingSymbol  string
}

type RealtimeWarrantQuote struct {
	Symbol        string
	LastDone      float64
	PrevClose     float64
	Open          float64
	High          float64
	Low           float64
	Timestamp     int64
	Volume        int64
	Turnover      float64
	WarrantExtend *WarrantExtend
}

func toWarrantExtend(origin *quotev1.WarrantExtend) *WarrantExtend {
	return &WarrantExtend{
		ImpliedVolatility: origin.GetImpliedVolatility(),
		ExpiryDate:        origin.GetExpiryDate(),
		LastTradeDate:     origin.GetLastTradeDate(),
		OutstandingRatio:  origin.GetOutstandingRatio(),
		OutstandingQty:    origin.GetOutstandingQty(),
		ConversionRatio:   origin.GetConversionRatio(),
		Category:          origin.GetCategory(),
		StrikePrice:       origin.GetStrikePrice(),
		UpperStrikePrice:  origin.GetUpperStrikePrice(),
		LowerStrikePrice:  origin.GetLowerStrikePrice(),
		CallPrice:         origin.GetCallPrice(),
		UnderlyingSymbol:  origin.GetUnderlyingSymbol(),
	}
}

type WarrantQuote struct {
	Symbol        string
	LastDone      string
	PrevClose     string
	Open          string
	High          string
	Low           string
	Timestamp     int64
	Volume        int64
	Turnover      string
	TradeStatus   TradeStatus
	WarrantExtend *WarrantExtend
}

func toWarrantQuotes(origin []*quotev1.WarrantQuote) (warrantQuotes []*WarrantQuote) {
	warrantQuotes = make([]*WarrantQuote, 0, len(origin))
	for _, item := range origin {
		warrantQuotes = append(warrantQuotes, &WarrantQuote{
			Symbol:        item.GetSymbol(),
			LastDone:      item.GetLastDone(),
			PrevClose:     item.GetPrevClose(),
			Open:          item.GetOpen(),
			High:          item.GetHigh(),
			Low:           item.GetLow(),
			Timestamp:     item.GetTimestamp(),
			Volume:        item.GetVolume(),
			Turnover:      item.GetTurnover(),
			TradeStatus:   TradeStatus(item.GetTradeStatus()),
			WarrantExtend: toWarrantExtend(item.GetWarrantExtend()),
		})
	}
	return
}

type WarrantFilter struct {
	Symbol   string
	Language string

	SortBy     int32
	SortOrder  int32 // 0 Ascending 1 Desending
	SortOffset int32
	PageSize   int32 // Up to 500

	// The following are optional

	Type      []int32 // optional values: 0 - Call	1 - Put 2 - Bull 3 - Bear 4 - Inline
	IssuerIDs []int32

	// ExpiryDateType can have the following values.
	// 1 - Less than 3 months
	// 2 - 3 - 6 months
	// 3 - 6 - 12 months
	// 4 - greater than 12 months
	ExpiryDateType []int32

	// Optional values for PriceType
	// 1 - In bounds
	// 2 - Out bounds
	PriceType []int32

	// Optional values for Status:
	// 2 - Suspend trading
	// 3 - Papare List
	// 4 - Normal
	Status []int32
}

type Warrant struct {
	Symbol            string
	Name              string
	LastDone          float64
	ChangeRate        float64
	ChangeVal         float64
	Turnover          float64
	ExpiryDate        string // YYYYMMDD
	StrikePrice       float64
	UpperStrikePrice  float64
	LowerStrikePrice  float64
	OutstandingQty    float64
	OutstandingRatio  float64
	Premium           float64
	ItmOtm            float64
	ImpliedVolatility float64
	Delta             float64
	CallPrice         float64
	EffectiveLeverage float64
	LeverageRatio     float64
	ConversionRatio   float64
	BalancePoint      float64
	State             string
}

type TradeDate struct {
	Date          string
	TradeDateType int32 // 0 full day, 1 morning only, 2 afternoon only(not happened before)
}

type Candlestick struct {
	Close     string
	Open      string
	Low       string
	High      string
	Volume    int64
	Turnover  string
	Timestamp int64
}

func toCandlesticks(origin []*quotev1.Candlestick) (sticks []*Candlestick) {
	sticks = make([]*Candlestick, 0, len(origin))
	for _, item := range origin {
		sticks = append(sticks, &Candlestick{
			Close:     item.GetClose(),
			Open:      item.GetOpen(),
			Low:       item.GetLow(),
			High:      item.GetHigh(),
			Volume:    item.GetVolume(),
			Turnover:  item.GetTurnover(),
			Timestamp: item.GetTimestamp(),
		})
	}
	return
}

type Quote struct {
	Symbol       string
	LastDone     string
	Open         string
	High         string
	Low          string
	Timestamp    int64
	Volume       int64
	Turnover     string
	TradeStatus  TradeStatus
	TradeSession TradeSessionType
}

type SecurityQuote struct {
	Symbol          string
	LastDone        string
	PrevClose       string
	Open            string
	High            string
	Low             string
	Timestamp       int64
	Volume          int64
	Turnover        string
	TradeStatus     TradeStatus
	PreMarketQuote  *PrePostQuote
	PostMarketQuote *PrePostQuote
}

func toSecurityQuotes(origin []*quotev1.SecurityQuote) (quotes []*SecurityQuote) {
	quotes = make([]*SecurityQuote, 0, len(origin))
	for _, item := range origin {
		quotes = append(quotes, &SecurityQuote{
			Symbol:      item.GetSymbol(),
			LastDone:    item.GetLastDone(),
			PrevClose:   item.GetPrevClose(),
			Open:        item.GetOpen(),
			High:        item.GetHigh(),
			Low:         item.GetLow(),
			Timestamp:   item.GetTimestamp(),
			Volume:      item.GetVolume(),
			Turnover:    item.GetTurnover(),
			TradeStatus: TradeStatus(item.GetTradeStatus()),
			PreMarketQuote: &PrePostQuote{
				LastDone:  item.GetPreMarketQuote().GetLastDone(),
				Timestamp: item.GetPreMarketQuote().GetTimestamp(),
				Volume:    item.GetPreMarketQuote().GetVolume(),
				Turnover:  item.GetPreMarketQuote().GetTurnover(),
				High:      item.GetPreMarketQuote().GetHigh(),
				Low:       item.GetPreMarketQuote().GetHigh(),
				PrevClose: item.GetPreMarketQuote().GetPrevClose(),
			},
			PostMarketQuote: &PrePostQuote{
				LastDone:  item.GetPostMarketQuote().GetLastDone(),
				Timestamp: item.GetPostMarketQuote().GetTimestamp(),
				Volume:    item.GetPostMarketQuote().GetVolume(),
				Turnover:  item.GetPostMarketQuote().GetTurnover(),
				High:      item.GetPostMarketQuote().GetHigh(),
				Low:       item.GetPostMarketQuote().GetHigh(),
				PrevClose: item.GetPostMarketQuote().GetPrevClose(),
			},
		})
	}
	return
}

type PrePostQuote struct {
	LastDone  string
	Timestamp int64
	Volume    int64
	Turnover  string
	High      string
	Low       string
	PrevClose string
}

type SecurityDepth struct {
	Symbol string
	Ask    []*Depth
	Bid    []*Depth
}

func toSecurityDepth(origin *quotev1.SecurityDepthResponse) *SecurityDepth {
	return &SecurityDepth{
		Symbol: origin.GetSymbol(),
		Ask:    toDepths(origin.GetAsk()),
		Bid:    toDepths(origin.GetBid()),
	}
}

type SecurityBrokers struct {
	Symbol     string
	AskBrokers []*Brokers
	BidBrokers []*Brokers
}

func toSecurityBrokers(origin *quotev1.SecurityBrokersResponse) *SecurityBrokers {
	return &SecurityBrokers{
		Symbol:     origin.GetSymbol(),
		AskBrokers: toBrokers(origin.GetAskBrokers()),
		BidBrokers: toBrokers(origin.GetBidBrokers()),
	}
}

type ParticipantInfo struct {
	BrokerIds         []int32
	ParticipantNameCn string
	ParticipantNameEn string
	ParticipantNameHk string
}

func toParticipantInfos(origin []*quotev1.ParticipantInfo) (participantInfos []*ParticipantInfo) {
	participantInfos = make([]*ParticipantInfo, 0, len(origin))
	for _, item := range origin {
		participantInfos = append(participantInfos, &ParticipantInfo{
			BrokerIds:         item.GetBrokerIds(),
			ParticipantNameCn: item.GetParticipantNameCn(),
			ParticipantNameEn: item.GetParticipantNameEn(),
			ParticipantNameHk: item.GetParticipantNameHk(),
		})
	}
	return
}

type IntradayLine struct {
	Price     string
	Timestamp int64
	Volume    int64
	Turnover  string
	AvgPrice  string
}

func toIntradayLines(origin []*quotev1.Line) (lines []*IntradayLine) {
	lines = make([]*IntradayLine, 0, len(origin))
	for _, item := range origin {
		lines = append(lines, &IntradayLine{
			Price:     item.GetPrice(),
			Timestamp: item.GetTimestamp(),
			Turnover:  item.GetTurnover(),
			Volume:    item.GetVolume(),
			AvgPrice:  item.GetAvgPrice(),
		})
	}
	return
}

type IssuerInfo struct {
	Id     int32
	NameCn string
	NameEn string
	NameHk string
}

func toIssueInfos(origin []*quotev1.IssuerInfo) (infos []*IssuerInfo) {
	infos = make([]*IssuerInfo, 0, len(origin))
	for _, item := range origin {
		infos = append(infos, &IssuerInfo{
			Id:     item.GetId(),
			NameCn: item.GetNameCn(),
			NameEn: item.GetNameEn(),
			NameHk: item.GetNameHk(),
		})
	}
	return
}

type MarketTradingSession struct {
	Market       openapi.Market
	TradeSession []*TradePeriod
}

func toMarketTradingSessions(origin []*quotev1.MarketTradePeriod) (sessions []*MarketTradingSession) {
	sessions = make([]*MarketTradingSession, 0, len(origin))
	for _, item := range origin {
		sessions = append(sessions, &MarketTradingSession{
			Market:       openapi.Market(item.GetMarket()),
			TradeSession: toTradePeriods(item.GetTradeSession()),
		})
	}
	return
}

type TradePeriod struct {
	BegTime      int32
	EndTime      int32
	TradeSession TradeSession
}

func toTradePeriods(origin []*quotev1.TradePeriod) (periods []*TradePeriod) {
	periods = make([]*TradePeriod, 0, len(origin))
	for _, item := range origin {
		periods = append(periods, &TradePeriod{
			BegTime:      item.GetBegTime(),
			EndTime:      item.GetEndTime(),
			TradeSession: TradeSession(item.GetTradeSession()),
		})
	}
	return
}

type MarketTradingDay struct {
	TradeDay     []time.Time
	HalfTradeDay []time.Time
}

func toQuoteSubTypes(origin []SubType) []quotev1.SubType {
	subTypes := make([]quotev1.SubType, 0, len(origin))
	for _, item := range origin {
		subTypes = append(subTypes, quotev1.SubType(item))
	}
	return subTypes
}

func toSubTypes(origin []quotev1.SubType) []SubType {
	subTypes := make([]SubType, 0, len(origin))
	for _, item := range origin {
		subTypes = append(subTypes, SubType(item))
	}
	return subTypes
}
