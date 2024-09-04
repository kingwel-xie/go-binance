package binance

import (
	"context"
)

// WsDepthService show depth info
type WsDepthService struct {
	c      *Client
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *WsDepthService) Symbol(symbol string) *WsDepthService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *WsDepthService) Limit(limit int) *WsDepthService {
	s.limit = &limit
	return s
}

// Do send wsRequest
func (s *WsDepthService) Do(ctx context.Context) (res *DepthResponse, err error) {
	r := &wsRequest{
		method: "depth",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callWsAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	j, err := newJSON(data)
	if err != nil {
		return nil, err
	}
	res = new(DepthResponse)
	res.LastUpdateID = j.Get("lastUpdateId").MustInt64()
	bidsLen := len(j.Get("bids").MustArray())
	res.Bids = make([]Bid, bidsLen)
	for i := 0; i < bidsLen; i++ {
		item := j.Get("bids").GetIndex(i)
		res.Bids[i] = Bid{
			Price:    item.GetIndex(0).MustString(),
			Quantity: item.GetIndex(1).MustString(),
		}
	}
	asksLen := len(j.Get("asks").MustArray())
	res.Asks = make([]Ask, asksLen)
	for i := 0; i < asksLen; i++ {
		item := j.Get("asks").GetIndex(i)
		res.Asks[i] = Ask{
			Price:    item.GetIndex(0).MustString(),
			Quantity: item.GetIndex(1).MustString(),
		}
	}
	return res, nil
}
