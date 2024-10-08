package binance

import (
	"context"
	"net/http"
)

// ListLoanableCoinFlexibleService list flexible loanable coin data.
type ListLoanableCoinFlexibleService struct {
	c        *Client
	loanCoin *string
}

// LoanCoin sets the loanCoin parameter.
func (s *ListLoanableCoinFlexibleService) LoanCoin(coin string) *ListLoanableCoinFlexibleService {
	if len(coin) > 0 {
		s.loanCoin = &coin
	}
	return s
}

// Do sends the request.
func (s *ListLoanableCoinFlexibleService) Do(ctx context.Context) (res *LoanableCoinFlexibleList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/loan/flexible/loanable/data",
		secType:  secTypeSigned,
	}
	if s.loanCoin != nil {
		r.setParam("loanCoin", *s.loanCoin)
	}

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanableCoinFlexibleList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanableCoinFlexibleList represents a flexible loanable coin list.
type LoanableCoinFlexibleList struct {
	Rows []struct {
		LoanCoin             string `json:"loanCoin"`
		FlexibleInterestRate string `json:"flexibleInterestRate"`
		FlexibleMinLimit     string `json:"flexibleMinLimit"`
		FlexibleMaxLimit     string `json:"flexibleMaxLimit"`
	} `json:"rows"`
	Total int `json:"total"`
}

// ListCollateralCoinFlexibleService list flexible collateral coin data.
type ListCollateralCoinFlexibleService struct {
	c              *Client
	collateralCoin *string
}

// CollateralCoin sets the collateral coin parameter.
func (s *ListCollateralCoinFlexibleService) CollateralCoin(coin string) *ListCollateralCoinFlexibleService {
	if len(coin) > 0 {
		s.collateralCoin = &coin
	}
	return s
}

// Do sends the request.
func (s *ListCollateralCoinFlexibleService) Do(ctx context.Context) (res *CollateralCoinFlexibleList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/loan/flexible/collateral/data",
		secType:  secTypeSigned,
	}
	if s.collateralCoin != nil {
		r.setParam("collateralCoin", *s.collateralCoin)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(CollateralCoinFlexibleList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// CollateralCoinFlexibleList represents a list of flexible collateral coins.
type CollateralCoinFlexibleList struct {
	Rows []struct {
		CollateralCoin string `json:"collateralCoin"`
		InitialLTV     string `json:"initialLTV"`
		MarginCallLTV  string `json:"marginCallLTV"`
		LiquidationLTV string `json:"liquidationLTV"`
		MaxLimit       string `json:"maxLimit"`
	} `json:"rows"`
	Total int `json:"total"`
}

// LoanBorrowFlexibleService borrow flexible product.
type LoanBorrowFlexibleService struct {
	c                *Client
	loanCoin         string
	collateralCoin   string
	collateralAmount string
}

// LoanCoin sets the loan coin parameter (MANDATORY).
func (s *LoanBorrowFlexibleService) LoanCoin(lonaCoin string) *LoanBorrowFlexibleService {
	s.loanCoin = lonaCoin
	return s
}

// CollateralCoin sets the collateral coin parameter (MANDATORY).
func (s *LoanBorrowFlexibleService) CollateralCoin(collateralCoin string) *LoanBorrowFlexibleService {
	s.collateralCoin = collateralCoin
	return s
}

// CollateralAmount sets the CollateralAmount parameter (MANDATORY).
func (s *LoanBorrowFlexibleService) CollateralAmount(v string) *LoanBorrowFlexibleService {
	s.collateralAmount = v
	return s
}

// Do sends the request.
func (s *LoanBorrowFlexibleService) Do(ctx context.Context) (res *LoanBorrowFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/borrow",
		secType:  secTypeSigned,
	}
	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("collateralAmount", s.collateralAmount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanBorrowFlexibleResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanBorrowFlexibleResponse represents a response from borrow flexible loan.
type LoanBorrowFlexibleResponse struct {
	LoanCoin         string `json:"loanCoin"`
	LoanAmount       string `json:"loanAmount"`
	CollateralCoin   string `json:"collateralCoin"`
	CollateralAmount string `json:"collateralAmount"`
	Status           string `json:"status"` //Succeeds, Failed, Processing
}

// LoanRepayFlexibleService repay flexible product.
type LoanRepayFlexibleService struct {
	c              *Client
	loanCoin       string
	collateralCoin string
	repayAmount    string
}

// LoanCoin sets the loanCoin parameter (MANDATORY).
func (s *LoanRepayFlexibleService) LoanCoin(loanCoin string) *LoanRepayFlexibleService {
	s.loanCoin = loanCoin
	return s
}

// CollateralCoin sets the collateralCoin parameter (MANDATORY).
func (s *LoanRepayFlexibleService) CollateralCoin(collateralCoin string) *LoanRepayFlexibleService {
	s.collateralCoin = collateralCoin
	return s
}

// RepayAmount sets the repayAmount parameter (MANDATORY).
func (s *LoanRepayFlexibleService) RepayAmount(repayAmount string) *LoanRepayFlexibleService {
	s.repayAmount = repayAmount
	return s
}

// Do sends the request.
func (s *LoanRepayFlexibleService) Do(ctx context.Context) (res *LoanRepayFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/repay",
		secType:  secTypeSigned,
	}
	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("repayAmount", s.repayAmount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanRepayFlexibleResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanRepayFlexibleResponse represents a response from repay flexible loan.
type LoanRepayFlexibleResponse struct {
	LoanCoin            string `json:"loanCoin"`
	CollateralCoin      string `json:"collateralCoin"`
	RemainingDebt       string `json:"remainingDebt"`
	RemainingCollateral string `json:"remainingCollateral"`
	FullRepayment       bool   `json:"fullRepayment"`
	CurrentLTV          string `json:"currentLTV"`
	RepayStatus         string `json:"repayStatus"`
}

// ListLoanFlexibleService list flexible loan debt data.
type ListLoanFlexibleService struct {
	c        *Client
	loanCoin *string
	limit    *int
}

// LoanCoin sets the loan coin parameter.
func (s *ListLoanFlexibleService) LoanCoin(coin string) *ListLoanFlexibleService {
	if len(coin) > 0 {
		s.loanCoin = &coin
	}
	return s
}

// Limit set limit
func (s *ListLoanFlexibleService) Limit(limit int) *ListLoanFlexibleService {
	s.limit = &limit
	return s
}

// Do sends the request.
func (s *ListLoanFlexibleService) Do(ctx context.Context) (res *LoanOrderFlexibleList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/loan/flexible/ongoing/orders",
		secType:  secTypeSigned,
	}
	if s.loanCoin != nil {
		r.setParam("loanCoin", *s.loanCoin)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanOrderFlexibleList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanOrderFlexibleList represents a list of flexible loan orders.
type LoanOrderFlexibleList struct {
	Rows []struct {
		LoanCoin         string `json:"loanCoin"`
		TotalDebt        string `json:"totalDebt"`
		CollateralCoin   string `json:"collateralCoin"`
		CollateralAmount string `json:"collateralAmount"`
		CurrentLTV       string `json:"currentLTV"`
	} `json:"rows"`
	Total int `json:"total"`
}

// AdjustLtvLoanFlexibleService adjust flexible loan LTV.
type AdjustLtvLoanFlexibleService struct {
	c                *Client
	loanCoin         string
	collateralCoin   string
	adjustmentAmount string
	direction        string
}

// LoanCoin sets the loan coin parameter.
func (s *AdjustLtvLoanFlexibleService) LoanCoin(coin string) *AdjustLtvLoanFlexibleService {
	s.loanCoin = coin
	return s
}

// CollateralCoin set collateral coin parameter.
func (s *AdjustLtvLoanFlexibleService) CollateralCoin(collateralCoin string) *AdjustLtvLoanFlexibleService {
	s.collateralCoin = collateralCoin
	return s
}

// AdjustmentAmount set collateral adjustment amount parameter.
func (s *AdjustLtvLoanFlexibleService) AdjustmentAmount(adjustmentAmount string) *AdjustLtvLoanFlexibleService {
	s.adjustmentAmount = adjustmentAmount
	return s
}

// Direction set direction parameter, "ADDITIONAL", "REDUCED".
func (s *AdjustLtvLoanFlexibleService) Direction(direction string) *AdjustLtvLoanFlexibleService {
	s.direction = direction
	return s
}

// Do sends the request.
func (s *AdjustLtvLoanFlexibleService) Do(ctx context.Context) (res *AdjustLtvLoanFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/loan/flexible/adjust/ltv",
		secType:  secTypeSigned,
	}
	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("adjustmentAmount", s.adjustmentAmount)
	r.setParam("direction", s.direction)
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(AdjustLtvLoanFlexibleResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// AdjustLtvLoanFlexibleResponse represents a response of adjust LTV of flexible loan.
type AdjustLtvLoanFlexibleResponse struct {
	LoanCoin         string `json:"loanCoin"`
	CollateralCoin   string `json:"collateralCoin"`
	Direction        string `json:"direction"`
	AdjustmentAmount string `json:"adjustmentAmount"`
	CurrentLTV       string `json:"currentLTV"`
}
