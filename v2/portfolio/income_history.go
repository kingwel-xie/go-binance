package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetIncomeHistoryService get position margin history service
type GetIncomeHistoryService struct {
	c          *Client
	which      string // 'um' or 'cm'
	symbol     string
	incomeType string
	startTime  *int64
	endTime    *int64
	limit      *int64
}

// Which set which product
func (s *GetIncomeHistoryService) Which(which string) *GetIncomeHistoryService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *GetIncomeHistoryService) Symbol(symbol string) *GetIncomeHistoryService {
	s.symbol = symbol
	return s
}

// IncomeType set income type
func (s *GetIncomeHistoryService) IncomeType(incomeType string) *GetIncomeHistoryService {
	s.incomeType = incomeType
	return s
}

// StartTime set startTime
func (s *GetIncomeHistoryService) StartTime(startTime int64) *GetIncomeHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *GetIncomeHistoryService) EndTime(endTime int64) *GetIncomeHistoryService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *GetIncomeHistoryService) Limit(limit int64) *GetIncomeHistoryService {
	s.limit = &limit
	return s
}

// Do send request
func (s *GetIncomeHistoryService) Do(ctx context.Context, opts ...RequestOption) (res []*IncomeHistory, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/income", s.which),
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.incomeType != "" {
		r.setParam("incomeType", s.incomeType)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*IncomeHistory, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// IncomeHistory define position margin history info
type IncomeHistory struct {
	Asset      string `json:"asset"`
	Income     string `json:"income"`
	IncomeType string `json:"incomeType"`
	Info       string `json:"info"`
	Symbol     string `json:"symbol"`
	Time       int64  `json:"time"`
	TranID     int64  `json:"tranId"`
	TradeID    string `json:"tradeId"`
}
