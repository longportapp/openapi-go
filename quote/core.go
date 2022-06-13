package quote

import (
	"context"
	"sync"
	"time"

	"github.com/longbridgeapp/openapi-go"

	"github.com/longbridgeapp/openapi-protobufs/gen/go/quote"
	protocol "github.com/longbridgeapp/openapi-protocol/go"
	"github.com/longbridgeapp/openapi-protocol/go/client"
)


type Core struct {
    client *client.Client
	url string
	mu sync.Mutex
	subscriptions map[string][]SubFlag
}

func NewCore(url string, otp string, f func(*PushEvent)) (*Core, error) {
	cl := client.New()
	err := cl.Dial(context.Background(), url, &protocol.Handshake{
		Version: 1,
		Codec:  protocol.CodecProtobuf,
		Platform: protocol.PlatformOpenapi,
	}, client.WithAuthToken(otp))
	if err != nil {
		return nil, err
	}

	core := &Core{client: cl, url: url, subscriptions: make(map[string][]SubFlag)}
	cl.AfterReconnected(func (){
		core.resubscribe(context.Background())
	})
	cl.Subscribe(uint32(quotev1.Command_PushBrokersData), parsePushFunc(f))
	cl.Subscribe(uint32(quotev1.Command_PushDepthData), parsePushFunc(f))
	cl.Subscribe(uint32(quotev1.Command_PushTradeData), parsePushFunc(f))
	cl.Subscribe(uint32(quotev1.Command_PushQuoteData), parsePushFunc(f))

	return core, nil
}

func (c *Core) Subscribe(ctx context.Context, symbols []string, subFlags []SubFlag, isFirstPush bool) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.doSubscirbe(ctx, symbols, subFlags, isFirstPush)
	return
}

func (c *Core) doSubscirbe(ctx context.Context, symbols []string, subFlags []SubFlag, isFirstPush bool) (err error) {
	req := &quotev1.SubscribeRequest{
	    IsFirstPush: isFirstPush,
	    SubType: toSubTypes(subFlags),
	    Symbol: symbols,
	}
	_, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_Subscribe), Body: req})
	if err != nil {
		return
	}
	for _, symbol := range symbols {
		c.subscriptions[symbol] = subFlags
	}
	return
}

func (c *Core) Unsubscribe(ctx context.Context, unSubAll bool, symbols []string) (err error){
	c.mu.Lock()
	defer c.mu.Unlock()
	req := &quotev1.UnsubscribeRequest{
		Symbol: symbols,
		UnsubAll: unSubAll,
	}
	_, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_Unsubscribe), Body: req })
	if err != nil {
		return
	}
	if unSubAll {
		c.subscriptions = make(map[string][]SubFlag)
	}
	for _, symbol := range symbols {
		delete(c.subscriptions, symbol)
	}
	return
}

func (c *Core) resubscribe(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for symbol, subflags := range c.subscriptions {
		err := c.doSubscirbe(ctx, []string{symbol}, subflags, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Core) Subscriptions(ctx context.Context) (subscriptions map[string][]SubFlag, err error) {
	req := &quotev1.SubscriptionRequest{}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_Subscription), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.SubscriptionResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	subscriptions = make(map[string][]SubFlag, len(ret.GetSubList()))
	for _, item := range ret.GetSubList() {
		subscriptions[item.GetSymbol()] = toSubFlags(item.GetSubType())
	}
	return
}

func (c *Core) StaticInfo(ctx context.Context, symbols []string) (staticInfos []*SecurityStaticInfo, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol:  symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QuerySecurityStaticInfo), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.SecurityStaticInfoResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	staticInfos = make([]*SecurityStaticInfo, 0, len(ret.GetSecuStaticInfo()))
	for _, item := range ret.GetSecuStaticInfo() {
		staticInfos = append(staticInfos, &SecurityStaticInfo{
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

func (c *Core) Quote(ctx context.Context, symbols []string) (quotes []*SecurityQuote, err error){
	req := &quotev1.MultiSecurityRequest{
		Symbol:  symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QuerySecurityQuote), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.SecurityQuoteResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	quotes = toSecurityQuotes(ret.GetSecuQuote())
	return
}


func (c *Core) OptionQuote(ctx context.Context, symbols []string) (optionQuotes []*OptionQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol:  symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QuerySecurityQuote), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.OptionQuoteResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	optionQuotes = toOptionQuotes(ret.GetSecuQuote())
	return

}

func (c *Core) WarrantQuote(ctx context.Context, symbols []string) (warrantQuotes []*WarrantQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol:  symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryWarrantQuote), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.WarrantQuoteResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	warrantQuotes = toWarrantQuotes(ret.GetSecuQuote())
	return
}

func (c *Core) Depth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error){
	req := &quotev1.SecurityRequest{
		Symbol:  symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryDepth), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.SecurityDepthResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	securityDepth = toSecurityDepth(&ret)
	return
}

func (c *Core) Brokers(ctx context.Context, symbol string) (securityBorkers *SecurityBrokers, err error){
	req := &quotev1.SecurityRequest{
		Symbol:  symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryBrokers), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.SecurityBrokersResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	securityBorkers = toSecurityBrokers(&ret)
	return
}


func (c *Core) Participants(ctx context.Context) (infos []*ParticipantInfo, err error){
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryBrokers)})
	if err != nil {
		return
	}
	var ret quotev1.ParticipantBrokerIdsResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	infos = toParticipantInfos(ret.GetParticipantBrokerNumbers())
	return
}


func (c *Core) Trade(ctx context.Context, symbol string, count int32) (trades []*Trade, err error){
	req := &quotev1.SecurityTradeRequest{
		Symbol: symbol,
		Count: count,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryTrade), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityTradeResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	trades = toTrades(ret.GetTrades())
	return
}

func (c *Core) Intraday(ctx context.Context, symbol string) (lines []*IntradayLine, err error){
	req := &quotev1.SecurityIntradayRequest{
		Symbol: symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryIntraday), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityIntradayResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	lines = toIntradayLines(ret.GetLines())
	return
}


func (c *Core) Candlesticks(ctx context.Context, symbol string, period Period, count int32, adjustType AdjustType) (sticks []*Candlestick, err error){
	req := &quotev1.SecurityCandlestickRequest{
		Symbol: symbol,
		Period: quotev1.Period(period),
		Count: count,
		AdjustType: quotev1.AdjustType(adjustType),
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryCandlestick), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityCandlestickResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	sticks = toCandlesticks(ret.GetCandlesticks())
	return
}


func (c *Core) OptionChainExpiryDateList(ctx context.Context, symbol string) (times []*time.Time, err error){
	req := &quotev1.SecurityRequest{
		Symbol:  symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryOptionChainDate), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.OptionChainDateListResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	times = make([]*time.Time, len(ret.GetExpiryDate()))
	for _, dateStr := range ret.GetExpiryDate(){
		var dt time.Time
		dt, err = time.Parse(DateLayout, dateStr)
		if err != nil {
			return nil, err
		}
		times = append(times, &dt)
	}
	return
}

func (c *Core) OptionChainInfoByDate(ctx context.Context, symbol string, expiryDate *time.Time) (priceInfos []*StrikePriceInfo, err error){
	req := &quotev1.OptionChainDateStrikeInfoRequest{
		Symbol:  symbol,
		ExpiryDate: expiryDate.Format(DateLayout),
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryOptionChainDate), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.OptionChainDateStrikeInfoResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	priceInfos = toStrikePriceInfos(ret.GetStrikePriceInfo())
	return
}

func (c *Core) WarrantIssuers(ctx context.Context) (infos []*IssuerInfo, err error){
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryWarrantIssuerInfo)})
	if err != nil {
		return
	}
	var ret quotev1.IssuerInfoResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	infos = toIssueInfos(ret.GetIssuerInfo())
	return
}

func (c *Core) TradingSession(ctx context.Context) (sessions []*MarketTradingSession, err error){
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryMarketTradePeriod)})
	if err != nil {
		return
	}
	var ret quotev1.MarketTradePeriodResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	sessions = toMarketTradingSessions(ret.GetMarketTradeSession())
	return
}

func (c *Core) TradingDays(ctx context.Context, market openapi.Market, begin *time.Time, end *time.Time) (days *MarketTradingDay, err error){
	req := &quotev1.MarketTradeDayRequest{
		Market:  string(market),
		BegDay: begin.Format(DateLayout),
		EndDay: end.Format(DateLayout),
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryMarketTradeDay), Body: req })
	if err != nil {
		return
	}
	var ret quotev1.MarketTradeDayResponse
	err = res.Unmarshal(ret)
	if err != nil {
		return
	}
	var day time.Time
	tradingDays := make([]*time.Time, len(ret.GetTradeDay()))
	halfTradeDays := make([]*time.Time, len(ret.GetTradeDay()))
	for _, dateStr := range ret.GetTradeDay() {
		day, err = time.Parse(DateLayout, dateStr)
		if err != nil {
			return
		}
		tradingDays = append(tradingDays, &day)
	}
	for _, dateStr := range ret.GetHalfTradeDay() {
		day, err = time.Parse(DateLayout, dateStr)
		if err != nil {
			return
		}
		halfTradeDays = append(halfTradeDays, &day)
	}
	days = &MarketTradingDay{
		TradeDay: tradingDays,
		HalfTradeDay: halfTradeDays,
	}
	return
}

func parsePushFunc(f func(*PushEvent)) func(*protocol.Packet) {
	return func (packet *protocol.Packet) {
		f(nil)
	}
}
