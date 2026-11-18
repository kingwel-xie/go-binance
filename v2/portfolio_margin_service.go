package binance

import (
	"context"
	"net/http"
)

// GetPortfolioMarginCollateralRateService gets portfolio margin collateral rate (USER_DATA).
// See https://developers.binance.com/en/docs/derivatives-trading-portfolio-margin-pro/market-data/rest-api
type GetPortfolioMarginCollateralRateService struct {
	c *Client
}

// Do sends the request.
func (s *GetPortfolioMarginCollateralRateService) Do(ctx context.Context, opts ...RequestOption) (res []*PortfolioMarginCollateralRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/portfolio/collateralRate",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = make([]*PortfolioMarginCollateralRate, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// GetPortfolioMarginTieredCollateralRateService gets the full tiered collateral rate
// configuration applied to portfolio margin accounts (MARKET_DATA).
// Note: the unified-account and unified-account-pro tiered collateral rates are
// configured independently; this endpoint returns the one for portfolio margin accounts.
// See https://developers.binance.com/en/docs/derivatives-trading-portfolio-margin-pro/market-data/rest-api
type GetPortfolioMarginTieredCollateralRateService struct {
	c *Client
}

// Do sends the request.
func (s *GetPortfolioMarginTieredCollateralRateService) Do(ctx context.Context, opts ...RequestOption) (res []*PortfolioMarginCollateralRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/portfolio/retail/collateralRate",
		secType:  secTypeAPIKey,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = make([]*PortfolioMarginCollateralRate, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// PortfolioMarginCollateralRate represents the collateral rate of one asset.
type PortfolioMarginCollateralRate struct {
	Asset          string                                    `json:"asset"`
	CollateralInfo []PortfolioMarginCollateralRateTierDetail `json:"collateralInfo"`
}

// PortfolioMarginCollateralRateTierDetail represents one collateral rate tier.
type PortfolioMarginCollateralRateTierDetail struct {
	TierFloor      string `json:"tierFloor"`
	TierCap        string `json:"tierCap"`
	CollateralRate string `json:"collateralRate"`
	Cum            string `json:"cum"`
}
