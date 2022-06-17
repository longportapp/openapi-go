package quote

import (
	"sync"
)

type QuoteData struct {
	Sequence int64
	Quote    *PushQuote
}

type DepthData struct {
	Sequence int64
	Ask      []*Depth
	Bid      []*Depth
}

type BrokersData struct {
	Sequence   int64
	AskBrokers []*Brokers
	BidBrokers []*Brokers
}

type TradesData struct {
	Sequence int64
	Trades   []*Trade
}

// store is an storeage to save quote, brokers, depth,
// trades information from server push event.
type store struct {
	quoteMut    sync.RWMutex
	quoteData   map[string]*QuoteData
	brokersMut  sync.RWMutex
	brokersData map[string]*BrokersData
	tradesMut   sync.RWMutex
	tradesData  map[string]*TradesData
	depthMut    sync.RWMutex
	depthData   map[string]*DepthData
}

func newStore() *store {
	return &store{
		quoteData:   make(map[string]*QuoteData),
		brokersData: make(map[string]*BrokersData),
		tradesData:  make(map[string]*TradesData),
		depthData:   make(map[string]*DepthData),
	}
}

func (s *store) HandlePushEvent(event *PushEvent) {
	switch event.Type {
	case EventBroker:
		s.MergeBroker(event.Brokers)
	case EventDepth:
		s.MergeDepth(event.Depth)
	case EventQuote:
		s.MergeQuote(event.Quote)
	case EventTrade:
		s.MergeTrade(event.Trade)
	}
}

func (s *store) MergeBroker(brokers *PushBrokers) {
	s.brokersMut.Lock()
	defer s.brokersMut.Unlock()
	data := s.brokersData[brokers.Symbol]
	if data == nil {
		data = &BrokersData{
			AskBrokers: make([]*Brokers, 0),
			BidBrokers: make([]*Brokers, 0),
			Sequence:   -1,
		}
		s.brokersData[brokers.Symbol] = data
	}
	if brokers.Sequence <= data.Sequence {
		return
	}
	data.Sequence = brokers.Sequence
	data.AskBrokers = replaceBrokers(data.AskBrokers, brokers.AskBrokers)
	data.BidBrokers = replaceBrokers(data.BidBrokers, brokers.BidBrokers)
}

func (s *store) MergeDepth(depth *PushDepth) {
	s.depthMut.Lock()
	defer s.depthMut.Unlock()
	data := s.depthData[depth.Symbol]
	if data == nil {
		data = &DepthData{
			Ask:      make([]*Depth, 0),
			Bid:      make([]*Depth, 0),
			Sequence: -1,
		}
		s.depthData[depth.Symbol] = data
	}
	if depth.Sequence <= data.Sequence {
		return
	}
	data.Sequence = depth.Sequence
	data.Ask = replaceDepth(data.Ask, depth.Ask)
	data.Bid = replaceDepth(data.Bid, depth.Bid)
}

func (s *store) MergeQuote(quote *PushQuote) {
	s.quoteMut.Lock()
	defer s.quoteMut.Unlock()
	data := s.quoteData[quote.Symbol]
	if data == nil {
		data = &QuoteData{
			Quote: &PushQuote{
				Symbol: quote.Symbol,
			},
			Sequence: -1,
		}
		s.quoteData[quote.Symbol] = data
	}
	if quote.Sequence <= data.Sequence {
		return
	}
	data.Sequence = quote.Sequence
	newQuote := new(PushQuote)
	*newQuote = *data.Quote
	if quote.LastDone != "" {
		newQuote.LastDone = quote.LastDone
	}
	if quote.High != "" {
		newQuote.High = quote.High
	}
	if quote.Open != "" {
		newQuote.Open = quote.Open
	}
	if quote.Low != "" {
		newQuote.Low = quote.Low
	}
	if quote.Turnover != "" {
		newQuote.Turnover = quote.Turnover
	}
	if quote.Volume != 0 {
		newQuote.Volume = quote.Volume
	}
	newQuote.Timestamp = quote.Timestamp
	newQuote.TradeSession = quote.TradeSession
	newQuote.TradeStatus = quote.TradeStatus
	newQuote.Sequence = quote.Sequence
	data.Quote = newQuote
}

func (s *store) MergeTrade(trade *PushTrade) {
	s.tradesMut.Lock()
	defer s.tradesMut.Unlock()
	data := s.tradesData[trade.Symbol]
	if data == nil {
		data = &TradesData{
			Trades:   make([]*Trade, 0),
			Sequence: -1,
		}
		s.tradesData[trade.Symbol] = data
	}
	if trade.Sequence <= data.Sequence {
		return
	}
	data.Sequence = trade.Sequence
	data.Trades = append(data.Trades, trade.Trade...)
}

func (s *store) GetTrades(symbol string) []*Trade {
	s.tradesMut.RLock()
	defer s.tradesMut.RUnlock()
	data := s.tradesData[symbol]
	if data == nil {
		return nil
	}
	// copy
	trades := make([]*Trade, 0, len(data.Trades))
	for _, trade := range data.Trades {
		n := new(Trade)
		*n = *trade
		trades = append(trades, n)
	}
	return trades
}

func (s *store) GetBrokers(symbol string) ([]*Brokers, []*Brokers) {
	s.brokersMut.RLock()
	defer s.brokersMut.RUnlock()
	data := s.brokersData[symbol]
	if data == nil {
		return nil, nil
	}
	return copyBrokers(data.AskBrokers), copyBrokers(data.BidBrokers)
}

func (s *store) GetDepth(symbol string) ([]*Depth, []*Depth) {
	s.depthMut.RLock()
	defer s.depthMut.RUnlock()
	data := s.depthData[symbol]
	if data == nil {
		return nil, nil
	}
	return copyDepth(data.Ask), copyDepth(data.Bid)
}

func (s *store) GetQuote(symbol string) *Quote {
	s.quoteMut.RLock()
	defer s.quoteMut.RUnlock()
	data := s.quoteData[symbol]
	if data == nil {
		return nil
	}
	return &Quote{
		Symbol:       data.Quote.Symbol,
		Open:         data.Quote.Open,
		High:         data.Quote.High,
		Low:          data.Quote.Low,
		LastDone:     data.Quote.LastDone,
		Timestamp:    data.Quote.Timestamp,
		Volume:       data.Quote.Volume,
		Turnover:     data.Quote.Turnover,
		TradeStatus:  data.Quote.TradeStatus,
		TradeSession: data.Quote.TradeSession,
	}
}

func copyDepth(depths []*Depth) []*Depth {
	newDepths := make([]*Depth, 0, len(depths))
	for _, depth := range depths {
		n := new(Depth)
		*n = *depth
		newDepths = append(newDepths, n)
	}
	return newDepths
}

func copyBrokers(brokers []*Brokers) []*Brokers {
	newBrokers := make([]*Brokers, 0, len(brokers))
	for _, broker := range brokers {
		n := new(Brokers)
		*n = *broker
		newBrokers = append(newBrokers, n)
	}
	return newBrokers
}

func replaceBrokers(prevBrokers, newBorkers []*Brokers) []*Brokers {
	brokers := append([]*Brokers{}, prevBrokers...)
	for _, new := range newBorkers {
		rflag := -1
		iflag := -1
		for i, prev := range brokers {
			if new.Position == prev.Position {
				rflag = i
				break
			}
			if new.Position < prev.Position {
				iflag = i
				break
			}
		}
		if rflag != -1 {
			brokers[rflag] = new
		} else if iflag != -1 {
			brokers = append(brokers[:iflag+1], brokers[iflag:]...)
			brokers[iflag] = new
		} else {
			brokers = append(brokers, new)
		}
	}
	return brokers
}

func replaceDepth(prevDepths, newDepths []*Depth) []*Depth {
	depths := append([]*Depth{}, prevDepths...)
	for _, new := range newDepths {
		rflag := -1
		iflag := -1
		for i, prev := range depths {
			if new.Position == prev.Position {
				rflag = i
				break
			}
			if new.Position < prev.Position {
				iflag = i
				break
			}
		}
		if rflag != -1 {
			depths[rflag] = new
		} else if iflag != -1 {
			depths = append(depths[:iflag+1], depths[iflag:]...)
			depths[iflag] = new
		} else {
			depths = append(depths, new)
		}
	}
	return depths
}
