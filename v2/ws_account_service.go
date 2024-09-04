package binance

import (
	"context"
)

// WsGetAccountService get account info
type WsGetAccountService struct {
	c                *Client
	omitZeroBalances *bool
}

// OmitZeroBalances ignores the zero balance.
func (s *WsGetAccountService) OmitZeroBalances(v bool) *WsGetAccountService {
	s.omitZeroBalances = &v
	return s
}

// Do send wsRequest
func (s *WsGetAccountService) Do(ctx context.Context) (res *Account, err error) {
	r := &wsRequest{
		method:  "account.status",
		secType: secTypeSigned,
	}
	if s.omitZeroBalances != nil {
		r.setParam("omitZeroBalances", *s.omitZeroBalances)
	}
	data, err := s.c.callWsAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res = new(Account)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
