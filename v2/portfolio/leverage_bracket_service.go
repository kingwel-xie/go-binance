package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetLeverageBracketService get leverage bracket info
type GetLeverageBracketService struct {
	c      *Client
	which  string // 'um' or 'cm'
	symbol *string
}

// Which set which product
func (s *GetLeverageBracketService) Which(which string) *GetLeverageBracketService {
	s.which = which
	return s
}

// Symbol set symbol.
func (s *GetLeverageBracketService) Symbol(symbol string) *GetLeverageBracketService {
	s.symbol = &symbol
	return s
}

// Do send request
func (s *GetLeverageBracketService) Do(ctx context.Context, opts ...RequestOption) (res []*LeverageBracket, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/leverageBracket", s.which),
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*LeverageBracket{}, err
	}
	res = make([]*LeverageBracket, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*LeverageBracket{}, err
	}
	return res, nil
}

// LeverageBracket define user CM leverage bracket of the account
type LeverageBracket struct {
	Symbol   string `json:"symbol"`
	Brackets []struct {
		Bracket          int     `json:"bracket"`
		InitialLeverage  int     `json:"initialLeverage"`
		QtyCap           int     `json:"qtyCap"`
		QtyFloor         int     `json:"qtyFloor"`
		MaintMarginRatio float64 `json:"maintMarginRatio"`
		Cum              float64 `json:"cum"`
	} `json:"brackets"`
}
