package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChangeLeverageService change user's initial leverage of specific symbol market
type ChangeLeverageService struct {
	c        *Client
	which    string // 'um' or 'cm'
	symbol   string
	leverage int
}

// Which set which product
func (s *ChangeLeverageService) Which(which string) *ChangeLeverageService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *ChangeLeverageService) Symbol(symbol string) *ChangeLeverageService {
	s.symbol = symbol
	return s
}

// Leverage set leverage
func (s *ChangeLeverageService) Leverage(leverage int) *ChangeLeverageService {
	s.leverage = leverage
	return s
}

// Do send request
func (s *ChangeLeverageService) Do(ctx context.Context, opts ...RequestOption) (res *SymbolLeverage, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodPost,
		endpoint: fmt.Sprintf("/papi/v1/%s/leverage", s.which),
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"symbol":   s.symbol,
		"leverage": s.leverage,
	})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(SymbolLeverage)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SymbolLeverage define leverage info of symbol
type SymbolLeverage struct {
	Leverage         int    `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"` // UM
	MaxQty           string `json:"maxQty"`           // CM
	Symbol           string `json:"symbol"`
}

// ChangePositionModeService change user's position mode
type ChangePositionModeService struct {
	c        *Client
	which    string // 'um' or 'cm'
	dualSide bool
}

// Which set which product
func (s *ChangePositionModeService) Which(which string) *ChangePositionModeService {
	s.which = which
	return s
}

// Change user's position mode: true - Hedge Mode, false - One-way Mode
func (s *ChangePositionModeService) DualSide(dualSide bool) *ChangePositionModeService {
	s.dualSide = dualSide
	return s
}

// Do send request
func (s *ChangePositionModeService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	if s.which == "" {
		return errWhichMissing
	}
	r := &request{
		method:   http.MethodPost,
		endpoint: fmt.Sprintf("/papi/v1/%s/positionSide/dual", s.which),
		secType:  secTypeSigned,
	}
	r.setFormParams(params{
		"dualSidePosition": s.dualSide,
	})
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// GetPositionModeService get user's position mode
type GetPositionModeService struct {
	c     *Client
	which string // 'um' or 'cm'
}

// Response of user's position mode
type PositionMode struct {
	DualSidePosition bool `json:"dualSidePosition"`
}

// Which set which product
func (s *GetPositionModeService) Which(which string) *GetPositionModeService {
	s.which = which
	return s
}

// Do send request
func (s *GetPositionModeService) Do(ctx context.Context, opts ...RequestOption) (res *PositionMode, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/positionSide/dual", s.which),
		secType:  secTypeSigned,
	}
	r.setFormParams(params{})
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = &PositionMode{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ChangeMultiAssetModeService change user's multi-asset mode
type ChangeMultiAssetModeService struct {
	c                 *Client
	multiAssetsMargin bool
}
