package binance

import (
	"context"
	"net/http"
)

// GetPortfolioMarginCollateralRateService gets the collateral rate of each asset
// for portfolio margin accounts (MARKET_DATA, API key only).
// See https://developers.binance.com/en/docs/derivatives-trading-portfolio-margin-pro/market-data/rest-api
type GetPortfolioMarginCollateralRateService struct {
	c *Client
}

// Do sends the request.
func (s *GetPortfolioMarginCollateralRateService) Do(ctx context.Context, opts ...RequestOption) (res []*PortfolioMarginCollateralRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/portfolio/collateralRate",
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

// GetPortfolioMarginTieredCollateralRateService gets the full tiered collateral
// rate configuration for portfolio margin pro accounts (USER_DATA).
// Note: unified-account and unified-account-pro tiered collateral rates are
// configured independently.
// See https://developers.binance.com/en/docs/derivatives-trading-portfolio-margin-pro/market-data/rest-api
type GetPortfolioMarginTieredCollateralRateService struct {
	c *Client
}

// Do sends the request.
func (s *GetPortfolioMarginTieredCollateralRateService) Do(ctx context.Context, opts ...RequestOption) (res []*PortfolioMarginTieredCollateralRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/portfolio/collateralRate",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = make([]*PortfolioMarginTieredCollateralRate, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// GetPortfolioMarginRetailTieredCollateralRateService gets the tiered collateral
// rate configuration applied to portfolio margin (retail) accounts
// (MARKET_DATA, API key only).
// See https://developers.binance.com/en/docs/derivatives-trading-portfolio-margin-pro/market-data/rest-api
type GetPortfolioMarginRetailTieredCollateralRateService struct {
	c *Client
}

// Do sends the request.
func (s *GetPortfolioMarginRetailTieredCollateralRateService) Do(ctx context.Context, opts ...RequestOption) (res []*PortfolioMarginTieredCollateralRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/portfolio/retail/collateralRate",
		secType:  secTypeAPIKey,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = make([]*PortfolioMarginTieredCollateralRate, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// PortfolioMarginCollateralRate represents the collateral rate of one asset.
type PortfolioMarginCollateralRate struct {
	Asset          string `json:"asset"`
	CollateralRate string `json:"collateralRate"`
}

// PortfolioMarginTieredCollateralRate represents the tiered collateral rate of one asset.
type PortfolioMarginTieredCollateralRate struct {
	Asset          string                                          `json:"asset"`
	CollateralInfo []PortfolioMarginTieredCollateralRateTierDetail `json:"collateralInfo"`
}

// PortfolioMarginTieredCollateralRateTierDetail represents one collateral rate tier.
type PortfolioMarginTieredCollateralRateTierDetail struct {
	TierFloor      string `json:"tierFloor"`
	TierCap        string `json:"tierCap"`
	CollateralRate string `json:"collateralRate"`
	Cum            string `json:"cum"`
}
