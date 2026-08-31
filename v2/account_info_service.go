package binance

import (
	"context"
	"net/http"
)

// GetAccountInfoService fetches account info detail, including the current VIP level.
//
// See https://developers.binance.com/docs/wallet/account/api-rest-api/account-info
type GetAccountInfoService struct {
	c *Client
}

// Do sends the request.
func (s *GetAccountInfoService) Do(ctx context.Context) (res *AccountInfo, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/account/info",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(AccountInfo)
	err = json.Unmarshal(data, res)
	return
}

// AccountInfo represents account info detail.
type AccountInfo struct {
	VipLevel        int  `json:"vipLevel"`
	IsMarginEnabled bool `json:"isMarginEnabled"`
	IsFutureEnabled bool `json:"isFutureEnabled"`
}