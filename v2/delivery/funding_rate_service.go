package delivery

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListFundingRateService list depth
type ListFundingRateService struct {
	c      *Client
	symbol string
	limit  *int
}

// Symbol set symbol
func (s *ListFundingRateService) Symbol(symbol string) *ListFundingRateService {
	s.symbol = symbol
	return s
}

// Limit set limit
func (s *ListFundingRateService) Limit(limit int) *ListFundingRateService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListFundingRateService) Do(ctx context.Context, opts ...RequestOption) (res []*FundingRate, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/fundingRate",
	}
	r.setParam("symbol", s.symbol)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res = make([]*FundingRate, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FundingRate define funding rate history entry info
type FundingRate struct {
	Symbol      string `json:"symbol"`
	FundingTime int64  `json:"fundingTime"`
	FundingRate string `json:"fundingRate"`
	MarkPrice   string `json:"markPrice"`
}
