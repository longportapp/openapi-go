package quote

import (
	"context"
	"sync"
	"time"

	"github.com/longbridgeapp/openapi-go"
	"github.com/longbridgeapp/openapi-go/config"
	"github.com/longbridgeapp/openapi-go/http"
	"github.com/longbridgeapp/openapi-go/internal/util"
	"github.com/longbridgeapp/openapi-go/log"

	"github.com/longbridgeapp/openapi-protobufs/gen/go/quote"
	protocol "github.com/longbridgeapp/openapi-protocol/go"
	"github.com/longbridgeapp/openapi-protocol/go/client"
	"github.com/pkg/errors"
)

type core struct {
	client        *client.Client
	url           string
	mu            sync.Mutex
	subscriptions map[string][]SubType
	store         *store
}

func newCore(url string, httpClient *http.Client) (*core, error) {
	getOTP := func() (string, error) {
		otp, err := httpClient.GetOTP(context.Background())
		if err != nil {
			return "", errors.Wrap(err, "failed to get otp")
		}
		return otp, nil
	}
	cl := client.New()
	err := cl.Dial(context.Background(), url, &protocol.Handshake{
		Version:  1,
		Codec:    protocol.CodecProtobuf,
		Platform: protocol.PlatformOpenapi,
	}, client.WithAuthTokenGetter(getOTP))
	if err != nil {
		return nil, err
	}
	cl.Logger.SetLevel(config.GetLogLevelFromEnv())
	core := &core{
		client:        cl,
		url:           url,
		subscriptions: make(map[string][]SubType),
		store:         newStore(),
	}
	core.client.AfterReconnected(func() {
		if err := core.resubscribe(context.Background()); err != nil {
			log.Error(err)
		}
	})
	return core, nil
}

func (c *core) SetQuoteHandler(f func(*PushQuote)) {
	c.client.Subscribe(uint32(quotev1.Command_PushQuoteData), parsePushQuoteFunc(f, c))
}

func (c *core) SetTradeHandler(f func(*PushTrade)) {
	c.client.Subscribe(uint32(quotev1.Command_PushTradeData), parsePushTradeFunc(f, c))
}

func (c *core) SetDepthHandler(f func(*PushDepth)) {
	c.client.Subscribe(uint32(quotev1.Command_PushDepthData), parsePushDepthFunc(f, c))
}

func (c *core) SetBrokersHandler(f func(*PushBrokers)) {
	c.client.Subscribe(uint32(quotev1.Command_PushBrokersData), parsePushBrokersFunc(f, c))
}

func (c *core) Subscribe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.doSubscirbe(ctx, symbols, subTypes, isFirstPush)
	return
}

func (c *core) doSubscirbe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	req := &quotev1.SubscribeRequest{
		IsFirstPush: isFirstPush,
		Symbol:      symbols,
	}
	err = util.Copy(&req.SubType, subTypes)
	if err != nil {
		return
	}
	_, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_Subscribe), Body: req})
	if err != nil {
		return
	}
	for _, symbol := range symbols {
		c.subscriptions[symbol] = subTypes
	}
	return
}

func (c *core) Unsubscribe(ctx context.Context, unSubAll bool, symbols []string, subTypes []SubType) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	req := &quotev1.UnsubscribeRequest{
		Symbol:   symbols,
		UnsubAll: unSubAll,
	}
	err = util.Copy(&req.SubType, subTypes)
	if err != nil {
		return
	}
	_, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_Unsubscribe), Body: req})
	if err != nil {
		return
	}
	if unSubAll {
		c.subscriptions = make(map[string][]SubType)
	}
	for _, symbol := range symbols {
		delete(c.subscriptions, symbol)
	}
	return
}

func (c *core) resubscribe(ctx context.Context) error {
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

func (c *core) Subscriptions(ctx context.Context) (subscriptions map[string][]SubType, err error) {
	req := &quotev1.SubscriptionRequest{}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_Subscription), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SubscriptionResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	subscriptions = make(map[string][]SubType, len(ret.GetSubList()))
	for _, item := range ret.GetSubList() {
		sublist := make([]SubType, 0, len(item.GetSubType()))
		err = util.Copy(&sublist, item.GetSubType())
		if err != nil {
			return
		}
		subscriptions[item.GetSymbol()] = sublist
	}
	return
}

func (c *core) StaticInfo(ctx context.Context, symbols []string) (staticInfos []*StaticInfo, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QuerySecurityStaticInfo), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityStaticInfoResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&staticInfos, ret.GetSecuStaticInfo())
	return
}

func (c *core) Quote(ctx context.Context, symbols []string) (quotes []*SecurityQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QuerySecurityQuote), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityQuoteResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&quotes, ret.GetSecuQuote())
	return
}

func (c *core) OptionQuote(ctx context.Context, symbols []string) (optionQuotes []*OptionQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QuerySecurityQuote), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.OptionQuoteResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&optionQuotes, ret.GetSecuQuote())
	return

}

func (c *core) WarrantQuote(ctx context.Context, symbols []string) (warrantQuotes []*WarrantQuote, err error) {
	req := &quotev1.MultiSecurityRequest{
		Symbol: symbols,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryWarrantQuote), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.WarrantQuoteResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&warrantQuotes, ret.GetSecuQuote())
	return
}

func (c *core) Depth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryDepth), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityDepthResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	securityDepth = &SecurityDepth{}
	err = util.Copy(&securityDepth, ret)
	return
}

func (c *core) Brokers(ctx context.Context, symbol string) (securityBorkers *SecurityBrokers, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryBrokers), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityBrokersResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	securityBorkers = &SecurityBrokers{}
	err = util.Copy(securityBorkers, ret)
	return
}

func (c *core) Participants(ctx context.Context) (infos []*ParticipantInfo, err error) {
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryParticipantBrokerIds)})
	if err != nil {
		return
	}
	var ret quotev1.ParticipantBrokerIdsResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&infos, ret.GetParticipantBrokerNumbers())
	return
}

func (c *core) Trades(ctx context.Context, symbol string, count int32) (trades []*Trade, err error) {
	req := &quotev1.SecurityTradeRequest{
		Symbol: symbol,
		Count:  count,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryTrade), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityTradeResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&trades, ret.GetTrades())
	return
}

func (c *core) Intraday(ctx context.Context, symbol string) (lines []*IntradayLine, err error) {
	req := &quotev1.SecurityIntradayRequest{
		Symbol: symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryIntraday), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityIntradayResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&lines, ret.GetLines())
	return
}

func (c *core) Candlesticks(ctx context.Context, symbol string, period Period, count int32, adjustType AdjustType) (sticks []*Candlestick, err error) {
	req := &quotev1.SecurityCandlestickRequest{
		Symbol:     symbol,
		Period:     quotev1.Period(period),
		Count:      count,
		AdjustType: quotev1.AdjustType(adjustType),
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryCandlestick), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.SecurityCandlestickResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&sticks, ret.GetCandlesticks())
	return
}

func (c *core) OptionChainExpiryDateList(ctx context.Context, symbol string) (times []time.Time, err error) {
	req := &quotev1.SecurityRequest{
		Symbol: symbol,
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryOptionChainDate), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.OptionChainDateListResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	times, err = toTimes(ret.GetExpiryDate())
	if err != nil {
		return
	}
	return
}

func (c *core) OptionChainInfoByDate(ctx context.Context, symbol string, expiryDate *time.Time) (times []time.Time, err error) {
	req := &quotev1.OptionChainDateStrikeInfoRequest{
		Symbol:     symbol,
		ExpiryDate: util.FormatDateSimple(expiryDate),
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryOptionChainDate), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.OptionChainDateListResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	times, err = toTimes(ret.GetExpiryDate())
	if err != nil {
		return
	}
	return
}

func (c *core) WarrantIssuers(ctx context.Context) (infos []*IssuerInfo, err error) {
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryWarrantIssuerInfo)})
	if err != nil {
		return
	}
	var ret quotev1.IssuerInfoResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&infos, ret.GetIssuerInfo())
	return
}

func (c *core) TradingSession(ctx context.Context) (sessions []*MarketTradingSession, err error) {
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryMarketTradePeriod)})
	if err != nil {
		return
	}
	var ret quotev1.MarketTradePeriodResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	err = util.Copy(&sessions, ret.GetMarketTradeSession())
	return
}

func (c *core) TradingDays(ctx context.Context, market openapi.Market, begin *time.Time, end *time.Time) (days *MarketTradingDay, err error) {
	var (
		tradingDays   []time.Time
		halfTradeDays []time.Time
	)

	req := &quotev1.MarketTradeDayRequest{
		Market: string(market),
		BegDay: util.FormatDateSimple(begin),
		EndDay: util.FormatDateSimple(end),
	}
	var res *protocol.Packet
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(quotev1.Command_QueryMarketTradeDay), Body: req})
	if err != nil {
		return
	}
	var ret quotev1.MarketTradeDayResponse
	err = res.Unmarshal(&ret)
	if err != nil {
		return
	}
	tradingDays, err = toTimes(ret.GetTradeDay())
	if err != nil {
		return
	}
	halfTradeDays, err = toTimes(ret.GetHalfTradeDay())
	if err != nil {
		return
	}
	days = &MarketTradingDay{
		TradeDay:     tradingDays,
		HalfTradeDay: halfTradeDays,
	}
	return
}

func (c *core) RealtimeQuote(ctx context.Context, symbols []string) (quotes []*Quote, err error) {
	quotes = make([]*Quote, 0, len(symbols))
	for _, symbol := range symbols {
		quotes = append(quotes, c.store.GetQuote(symbol))
	}
	return
}

func (c *core) RealtimeDepth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	askDepths, bidDepths := c.store.GetDepth(symbol)
	return &SecurityDepth{
		Symbol: symbol,
		Ask:    askDepths,
		Bid:    bidDepths,
	}, nil
}

func (c *core) RealtimeTrades(ctx context.Context, symbol string) (trades []*Trade, err error) {
	return c.store.GetTrades(symbol), nil
}

func (c *core) RealtimeBrokers(ctx context.Context, symbol string) (securityBorkers *SecurityBrokers, err error) {
	askBrokers, bidBorkers := c.store.GetBrokers(symbol)
	return &SecurityBrokers{
		Symbol:     symbol,
		AskBrokers: askBrokers,
		BidBrokers: bidBorkers,
	}, nil
}

func (c *core) Close() error {
	return c.client.Close(nil)
}

func parsePushQuoteFunc(f func(*PushQuote), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushQuote
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote push event, unmarshal error:%v", err)
			return
		}
		var pb PushQuote
		if err = util.Copy(&pb, data); err != nil {
			log.Errorf("quote push event, copy data error:%v", err)
			return
		}
		core.store.MergeQuote(&pb)
		f(&pb)
	}
}

func parsePushDepthFunc(f func(*PushDepth), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushDepth
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote depth push event, unmarshal error:%v", err)
			return
		}
		var pd PushDepth
		if err = util.Copy(&pd, data); err != nil {
			log.Errorf("quote depth push event, copy data error:%v", err)
			return
		}
		core.store.MergeDepth(&pd)
		f(&pd)
	}
}

func parsePushBrokersFunc(f func(*PushBrokers), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushDepth
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote brokers push event, unmarshal error:%v", err)
			return
		}
		var pb PushBrokers
		if err = util.Copy(&pb, data); err != nil {
			log.Errorf("quote brokers push event, copy data error:%v", err)
			return
		}
		core.store.MergeBroker(&pb)
		f(&pb)
	}
}

func parsePushTradeFunc(f func(*PushTrade), core *core) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var err error
		var data quotev1.PushDepth
		if err = packet.Unmarshal(&data); err != nil {
			log.Errorf("quote trade push event, unmarshal error:%v", err)
			return
		}
		var pt PushTrade
		if err = util.Copy(&pt, data); err != nil {
			log.Errorf("quote trade push event, copy data error:%v", err)
			return
		}
		core.store.MergeTrade(&pt)
		f(&pt)
	}
}

func toTimes(origin []string) (times []time.Time, err error) {
	times = make([]time.Time, 0, len(origin))
	for _, dateStr := range origin {
		dt, err := util.ParseDateSimple(dateStr)
		if err != nil {
			return nil, err
		}
		times = append(times, dt)
	}
	return
}
