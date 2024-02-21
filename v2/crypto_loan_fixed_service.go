package binance

import (
	"context"
	"net/http"
)

// ListLoanableCoinService list loanable coin data.
type ListLoanableCoinService struct {
	c        *Client
	loanCoin *string
	vipLevel *int32
}

// LoanCoin sets the loanCoin parameter.
func (s *ListLoanableCoinService) LoanCoin(coin string) *ListLoanableCoinService {
	if len(coin) > 0 {
		s.loanCoin = &coin
	}
	return s
}

// VipLevel sets the vip Level.
func (s *ListLoanableCoinService) VipLevel(vipLevel int32) *ListLoanableCoinService {
	s.vipLevel = &vipLevel
	return s
}

// Do sends the request.
func (s *ListLoanableCoinService) Do(ctx context.Context) (res *LoanableCoinList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/loan/loanable/data",
		secType:  secTypeSigned,
	}
	if s.loanCoin != nil {
		r.setParam("loanCoin", *s.loanCoin)
	}
	if s.vipLevel != nil {
		r.setParam("vipLevel", *s.vipLevel)
	}

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanableCoinList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanableCoinList represents a loanable coin list.
type LoanableCoinList struct {
	Rows []struct {
		LoanCoin                   string `json:"loanCoin"`
		Seven_DHourlyInterestRate  string `json:"_7dHourlyInterestRate"`
		Seven_DDailyInterestRate   string `json:"_7dDailyInterestRate"`
		One_4DHourlyInterestRate   string `json:"_14dHourlyInterestRate"`
		One_4DDailyInterestRate    string `json:"_14dDailyInterestRate"`
		Three_0DHourlyInterestRate string `json:"_30dHourlyInterestRate"`
		Three_0DDailyInterestRate  string `json:"_30dDailyInterestRate"`
		Nine_0DHourlyInterestRate  string `json:"_90dHourlyInterestRate"`
		Nine_0DDailyInterestRate   string `json:"_90dDailyInterestRate"`
		One_80DHourlyInterestRate  string `json:"_180dHourlyInterestRate"`
		One_80DDailyInterestRate   string `json:"_180dDailyInterestRate"`
		MinLimit                   string `json:"minLimit"`
		MaxLimit                   string `json:"maxLimit"`
		VipLevel                   int    `json:"vipLevel"`
	} `json:"rows"`
	Total int `json:"total"`
}

// ListCollateralCoinService list collateral coin data.
type ListCollateralCoinService struct {
	c              *Client
	collateralCoin *string
	vipLevel       *int32
}

// CollateralCoin sets the collateral coin parameter.
func (s *ListCollateralCoinService) CollateralCoin(coin string) *ListCollateralCoinService {
	if len(coin) > 0 {
		s.collateralCoin = &coin
	}
	return s
}

// VipLevel sets the vip level.
func (s *ListCollateralCoinService) VipLevel(vipLevel int32) *ListCollateralCoinService {
	s.vipLevel = &vipLevel
	return s
}

// Do sends the request.
func (s *ListCollateralCoinService) Do(ctx context.Context) (res *CollateralCoinList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/loan/collateral/data",
		secType:  secTypeSigned,
	}
	if s.collateralCoin != nil {
		r.setParam("collateralCoin", *s.collateralCoin)
	}
	if s.vipLevel != nil {
		r.setParam("vipLevel", *s.vipLevel)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(CollateralCoinList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// CollateralCoinList represents a list of collateral coins.
type CollateralCoinList struct {
	Rows []struct {
		CollateralCoin string `json:"collateralCoin"`
		InitialLTV     string `json:"initialLTV"`
		MarginCallLTV  string `json:"marginCallLTV"`
		LiquidationLTV string `json:"liquidationLTV"`
		MaxLimit       string `json:"maxLimit"`
		VipLevel       int    `json:"vipLevel"`
	} `json:"rows"`
	Total int `json:"total"`
}

// LoanBorrowLockedService borrow locked product.
type LoanBorrowLockedService struct {
	c                *Client
	loanCoin         string
	collateralCoin   string
	collateralAmount string
	loanTerm         int
}

// LoanCoin sets the loan coin parameter (MANDATORY).
func (s *LoanBorrowLockedService) LoanCoin(lonaCoin string) *LoanBorrowLockedService {
	s.loanCoin = lonaCoin
	return s
}

// CollateralCoin sets the collateral coin parameter (MANDATORY).
func (s *LoanBorrowLockedService) CollateralCoin(collateralCoin string) *LoanBorrowLockedService {
	s.collateralCoin = collateralCoin
	return s
}

// CollateralAmount sets the CollateralAmount parameter (MANDATORY).
func (s *LoanBorrowLockedService) CollateralAmount(v string) *LoanBorrowLockedService {
	s.collateralAmount = v
	return s
}

// LoanTerm sets the LoanTerm parameter (MANDATORY).
func (s *LoanBorrowLockedService) LoanTerm(v int) *LoanBorrowLockedService {
	s.loanTerm = v
	return s
}

// Do sends the request.
func (s *LoanBorrowLockedService) Do(ctx context.Context) (res *LoanBorrowLockedResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/loan/borrow",
		secType:  secTypeSigned,
	}
	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("collateralAmount", s.collateralAmount)
	r.setParam("loanTerm", s.loanTerm)

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanBorrowLockedResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanBorrowLockedResponse represents a response from borrow locked loan.
type LoanBorrowLockedResponse struct {
	LoanCoin           string `json:"loanCoin"`
	LoanAmount         string `json:"loanAmount"`
	CollateralCoin     string `json:"collateralCoin"`
	CollateralAmount   string `json:"collateralAmount"`
	HourlyInterestRate string `json:"hourlyInterestRate"`
	OrderID            string `json:"orderId"`
}

// LoanRepayLockedService repay locked product.
type LoanRepayLockedService struct {
	c         *Client
	orderId   int64
	amount    string
	repayType int
}

// OrderId sets the orderId parameter (MANDATORY).
func (s *LoanRepayLockedService) OrderId(orderId int64) *LoanRepayLockedService {
	s.orderId = orderId
	return s
}

// Amount sets the amount parameter (MANDATORY).
func (s *LoanRepayLockedService) Amount(amount string) *LoanRepayLockedService {
	s.amount = amount
	return s
}

// Type sets the repayType parameter (MANDATORY).
func (s *LoanRepayLockedService) Type(repayType int) *LoanRepayLockedService {
	s.repayType = repayType
	return s
}

// Do sends the request.
func (s *LoanRepayLockedService) Do(ctx context.Context) (res *LoanRepayLockedResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/loan/repay",
		secType:  secTypeSigned,
	}
	r.setParam("orderId", s.orderId)
	r.setParam("amount", s.amount)
	r.setParam("type", s.repayType)

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanRepayLockedResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanRepayLockedResponse represents a response from repay locked loan.
type LoanRepayLockedResponse struct {
	LoanCoin            string `json:"loanCoin"`
	RemainingPrincipal  string `json:"remainingPrincipal"`
	RemainingInterest   string `json:"remainingInterest"`
	CollateralCoin      string `json:"collateralCoin"`
	RemainingCollateral string `json:"remainingCollateral"`
	CurrentLTV          string `json:"currentLTV"`
	RepayStatus         string `json:"repayStatus"`
}

// ListLoanLockedService list flexible loan debt data.
type ListLoanLockedService struct {
	c        *Client
	loanCoin *string
	limit    *int
}

// LoanCoin sets the loan coin parameter.
func (s *ListLoanLockedService) LoanCoin(coin string) *ListLoanLockedService {
	if len(coin) > 0 {
		s.loanCoin = &coin
	}
	return s
}

// Limit set limit
func (s *ListLoanLockedService) Limit(limit int) *ListLoanLockedService {
	s.limit = &limit
	return s
}

// Do sends the request.
func (s *ListLoanLockedService) Do(ctx context.Context) (res *LoanOrderLockedList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/loan/ongoing/orders",
		secType:  secTypeSigned,
	}
	if s.loanCoin != nil {
		r.setParam("loanCoin", *s.loanCoin)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(LoanOrderLockedList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// LoanOrderLockedList represents a list of flexible loan orders.
type LoanOrderLockedList struct {
	Rows []struct {
		OrderID          int    `json:"orderId"`
		LoanCoin         string `json:"loanCoin"`
		TotalDebt        string `json:"totalDebt"`
		ResidualInterest string `json:"residualInterest"`
		CollateralCoin   string `json:"collateralCoin"`
		CollateralAmount string `json:"collateralAmount"`
		CurrentLTV       string `json:"currentLTV"`
		ExpirationTime   int64  `json:"expirationTime"`
	} `json:"rows"`
	Total int `json:"total"`
}
