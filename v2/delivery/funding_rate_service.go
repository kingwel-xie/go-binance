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
	startTime *int64
	endTime *int64
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

// StartTime set startTime
func (s *ListFundingRateService) StartTime(startTime int64) *ListFundingRateService {
	s.startTime = &startTime
	return s
}

// EndTime set startTime
func (s *ListFundingRateService) EndTime(endTime int64) *ListFundingRateService {
	s.endTime = &endTime
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
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
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

// ListFundingInfoService list funding rate info
type ListFundingInfoService struct {
	c *Client
}

// Do send request
func (s *ListFundingInfoService) Do(ctx context.Context, opts ...RequestOption) (res []*FundingInfo, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/fundingInfo",
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res = make([]*FundingInfo, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FundingInfo define funding rate entry info
type FundingInfo struct {
	Symbol                   string `json:"symbol"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int    `json:"fundingIntervalHours"`
	Disclaimer               bool   `json:"disclaimer"`
}
