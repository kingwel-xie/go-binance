package binance

import (
	"context"
)

// WsPingService ping server
type WsPingService struct {
	c *Client
}

// Do send wsRequest
func (s *WsPingService) Do(ctx context.Context) (err error) {
	r := &wsRequest{
		method: "ping",
	}
	_, err = s.c.callWsAPI(ctx, r)
	return err
}

// WsServerTimeService get server time
type WsServerTimeService struct {
	c *Client
}

// Do send request
func (s *WsServerTimeService) Do(ctx context.Context) (serverTime int64, err error) {
	r := &wsRequest{
		method: "time",
	}
	data, err := s.c.callWsAPI(ctx, r)
	if err != nil {
		return 0, err
	}
	j, err := newJSON(data)
	if err != nil {
		return 0, err
	}
	serverTime = j.Get("serverTime").MustInt64()
	return serverTime, nil
}

//
//// SetServerTimeService set server time
//type SetServerTimeService struct {
//	c *Client
//}
//
//// Do send request
//func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
//	serverTime, err := s.c.NewServerTimeService().Do(ctx)
//	if err != nil {
//		return 0, err
//	}
//	timeOffset = currentTimestamp() - serverTime
//	s.c.TimeOffset = timeOffset
//	return timeOffset, nil
//}
