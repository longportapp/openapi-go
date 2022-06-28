package trade

import (
	"net/url"
)

func (req *GetHistoryExecutions) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", req.Symbol)
	p.AddDate("start_at", req.StartAt)
	p.AddDate("end_at", req.EndAt)
	return p.Values()
}

func (req *GetTodayExecutions) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", req.Symbol)
	p.Add("order_id", req.OrderId)
	return p.Values()
}

func (r *GetHistoryOrders) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", string(r.Symbol))
	p.Add("side", string(r.Side))
	p.Add("market", string(r.Market))
	p.AddInt("start_at", r.StartAt)
	p.AddInt("end_at", r.EndAt)
	vals := p.Values()
	for _, s := range r.Status {
		vals.Add("status", string(s))
	}
	return vals
}

func (r *GetStockPositions) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	vals := url.Values{}
	for _, s := range r.Symbols {
		vals.Add("symbols", string(s))
	}
	return vals
}

func (r *GetCashFlow) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", r.Symbol)
	p.AddInt("start_at", r.StartAt)
	p.AddInt("end_at", r.EndAt)
	p.AddOptInt("page", r.Page)
	p.AddOptInt("size", r.Size)
	p.AddInt("business_type", int64(r.BusinessType))
	return p.Values()
}

func (r *GetFundPositions) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	vals := url.Values{}
	for _, s := range r.Symbols {
		vals.Add("symbols", string(s))
	}
	return vals
}

func (r *GetTodayOrders) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("symbol", string(r.Symbol))
	p.Add("side", string(r.Side))
	p.Add("market", string(r.Market))
	vals := p.Values()
	for _, s := range r.Status {
		vals.Add("status", string(s))
	}
	return vals
}
