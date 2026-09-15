package options

import (
	"context"
	"encoding/json"
	"net/http"
)

// MarginAccountService queries the option margin account information
type MarginAccountService struct {
	c *Client
}

// Do send request
func (s *MarginAccountService) Do(ctx context.Context, opts ...RequestOption) (res *MarginAccount, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/eapi/v1/marginAccount",
		secType:  secTypeSigned,
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res = new(MarginAccount)
	err = json.Unmarshal(data, res)

	if err != nil {
		return nil, err
	}
	return res, nil
}

// MarginAccountAsset defines the asset balance of the option margin account
type MarginAccountAsset struct {
	Asset          string `json:"asset"`
	MarginBalance  string `json:"marginBalance"`
	Equity         string `json:"equity"`
	Available      string `json:"available"`
	InitialMargin  string `json:"initialMargin"`
	MaintMargin    string `json:"maintMargin"`
	UnrealizedPNL  string `json:"unrealizedPNL"`
	AdjustedEquity string `json:"adjustedEquity"`
}

// Greek defines the greeks of the option margin account
type Greek struct {
	Underlying string `json:"underlying"`
	Delta      string `json:"delta"`
	Gamma      string `json:"gamma"`
	Theta      string `json:"theta"`
	Vega       string `json:"vega"`
}

// MarginAccount define the option margin account response
type MarginAccount struct {
	Asset        []*MarginAccountAsset `json:"asset"`
	Greek        []*Greek              `json:"greek"`
	Time         uint64                `json:"time"`
	CanTrade     bool                  `json:"canTrade"`
	CanDeposit   bool                  `json:"canDeposit"`
	CanWithdraw  bool                  `json:"canWithdraw"`
	ReduceOnly   bool                  `json:"reduceOnly"`
	TradeGroupID int64                 `json:"tradeGroupId"`
}
