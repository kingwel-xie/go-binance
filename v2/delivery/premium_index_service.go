package delivery

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetPremiumIndexService list premium index and mark price, last funding rate
type GetPremiumIndexService struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *GetPremiumIndexService) Symbol(symbol string) *GetPremiumIndexService {
	if len(symbol) > 0 {
		s.symbol = &symbol
	}
	return s
}

// Do send request
func (s *GetPremiumIndexService) Do(ctx context.Context, opts ...RequestOption) (res []*PremiumIndex, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/dapi/v1/premiumIndex",
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res = make([]*PremiumIndex, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// PremiumIndex define premium index entry info
type PremiumIndex struct {
	Symbol               string `json:"symbol"`
	Pair                 string `json:"pair"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	InterestRate         string `json:"interestRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	Time                 int64  `json:"time"`
}
