package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CommissionRateService struct {
	c      *Client
	which  string // 'um' or 'cm'
	symbol string
}

// Symbol set symbol
func (s *CommissionRateService) Symbol(symbol string) *CommissionRateService {
	s.symbol = symbol
	return s
}

// Which set which product
func (s *CommissionRateService) Which(which string) *CommissionRateService {
	s.which = which
	return s
}

// Do send request
func (s *CommissionRateService) Do(ctx context.Context, opts ...RequestOption) (res *CommissionRate, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/commissionRate", s.which),
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CommissionRate)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Commission Rate
type CommissionRate struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"`
	TakerCommissionRate string `json:"takerCommissionRate"`
}
