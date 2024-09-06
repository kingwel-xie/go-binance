package binance

import (
	"context"
	"net/http"
)

// ListVipLoanableCoinService list VIP loanable coin data.
type ListVipLoanableCoinService struct {
	c        *Client
	loanCoin *string
}

// LoanCoin sets the loanCoin parameter.
func (s *ListVipLoanableCoinService) LoanCoin(coin string) *ListVipLoanableCoinService {
	if len(coin) > 0 {
		s.loanCoin = &coin
	}
	return s
}

// Do sends the request.
func (s *ListVipLoanableCoinService) Do(ctx context.Context) (res *VipLoanableCoinList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/loan/vip/loanable/data",
		secType:  secTypeSigned,
	}
	if s.loanCoin != nil {
		r.setParam("loanCoin", *s.loanCoin)
	}

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(VipLoanableCoinList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// VipLoanableCoinList represents a VIP loanable coin list.
type VipLoanableCoinList struct {
	Rows []struct {
		LoanCoin                   string `json:"loanCoin"`
		HourlyInterestRate         string `json:"_flexibleHourlyInterestRate"`
		FlexibleYearlyInterestRate string `json:"_flexibleYearlyInterestRate"`
		DDailyInterestRate         string `json:"_30dDailyInterestRate"`
		DYearlyInterestRate        string `json:"_30dYearlyInterestRate"`
		DDailyInterestRate1        string `json:"_60dDailyInterestRate"`
		DYearlyInterestRate1       string `json:"_60dYearlyInterestRate"`
		MinLimit                   string `json:"minLimit"`
		MaxLimit                   string `json:"maxLimit"`
		VipLevel                   int    `json:"vipLevel"`
	} `json:"rows"`
	Total int `json:"total"`
}

// ListVipCollateralCoinService list flexible collateral coin data.
type ListVipCollateralCoinService struct {
	c              *Client
	collateralCoin *string
}

// CollateralCoin sets the collateral coin parameter.
func (s *ListVipCollateralCoinService) CollateralCoin(coin string) *ListVipCollateralCoinService {
	if len(coin) > 0 {
		s.collateralCoin = &coin
	}
	return s
}

// Do sends the request.
func (s *ListVipCollateralCoinService) Do(ctx context.Context) (res *VipCollateralCoinList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/loan/vip/collateral/data",
		secType:  secTypeSigned,
	}
	if s.collateralCoin != nil {
		r.setParam("collateralCoin", *s.collateralCoin)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(VipCollateralCoinList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// VipCollateralCoinList represents a list of collateral coins.
type VipCollateralCoinList struct {
	Rows []struct {
		CollateralCoin    string `json:"collateralCoin"`
		StCollateralRatio string `json:"_1stCollateralRatio"`
		StCollateralRange string `json:"_1stCollateralRange"`
		NdCollateralRatio string `json:"_2ndCollateralRatio"`
		NdCollateralRange string `json:"_2ndCollateralRange"`
		RdCollateralRatio string `json:"_3rdCollateralRatio"`
		RdCollateralRange string `json:"_3rdCollateralRange"`
		ThCollateralRatio string `json:"_4thCollateralRatio"`
		ThCollateralRange string `json:"_4thCollateralRange"`
	} `json:"rows"`
	Total int `json:"total"`
}

// VipLoanBorrowService borrow flexible product.
type VipLoanBorrowService struct {
	c                   *Client
	loanCoin            string
	collateralCoin      string
	loanAccountId       string
	collateralAccountId string
	loanAmount          string
	isFlexibleRate      bool
}

// LoanAccountId sets the loan AccountId parameter (MANDATORY).
func (s *VipLoanBorrowService) LoanAccountId(v string) *VipLoanBorrowService {
	s.loanAccountId = v
	return s
}

// CollateralAccountId sets the collateral AccountId parameter (MANDATORY).
func (s *VipLoanBorrowService) CollateralAccountId(v string) *VipLoanBorrowService {
	s.collateralAccountId = v
	return s
}

// LoanCoin sets the loan coin parameter (MANDATORY).
func (s *VipLoanBorrowService) LoanCoin(lonaCoin string) *VipLoanBorrowService {
	s.loanCoin = lonaCoin
	return s
}

// CollateralCoin sets the collateral coin parameter (MANDATORY).
func (s *VipLoanBorrowService) CollateralCoin(collateralCoin string) *VipLoanBorrowService {
	s.collateralCoin = collateralCoin
	return s
}

// LoanAmount sets the loanAmount parameter (MANDATORY).
func (s *VipLoanBorrowService) LoanAmount(v string) *VipLoanBorrowService {
	s.loanAmount = v
	return s
}

// IsFlexibleRate sets the isFlexibleRate parameter (MANDATORY).
func (s *VipLoanBorrowService) IsFlexibleRate(v bool) *VipLoanBorrowService {
	s.isFlexibleRate = v
	return s
}

// Do sends the request.
func (s *VipLoanBorrowService) Do(ctx context.Context) (res *VipLoanBorrowResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/loan/vip/borrow",
		secType:  secTypeSigned,
	}
	r.setParam("loanAccountId", s.loanAccountId)
	r.setParam("collateralAccountId", s.collateralAccountId)
	r.setParam("loanCoin", s.loanCoin)
	r.setParam("collateralCoin", s.collateralCoin)
	r.setParam("loanAmount", s.loanAmount)
	r.setParam("isFlexibleRate", s.isFlexibleRate)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(VipLoanBorrowResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// VipLoanBorrowResponse represents a response from borrow flexible loan.
type VipLoanBorrowResponse struct {
	LoanAccountId       string `json:"loanAccountId"`
	RequestId           string `json:"requestId"`
	LoanCoin            string `json:"loanCoin"`
	IsFlexibleRate      string `json:"isFlexibleRate"`
	LoanAmount          string `json:"loanAmount"`
	CollateralAccountId string `json:"collateralAccountId"`
	CollateralCoin      string `json:"collateralCoin"`
	LoanTerm            string `json:"loanTerm"`
}

// VipLoanRepayService repay flexible product.
type VipLoanRepayService struct {
	c       *Client
	orderId string
	amount  string
}

// OrderId sets the loan orderId parameter (MANDATORY).
func (s *VipLoanRepayService) OrderId(orderId string) *VipLoanRepayService {
	s.orderId = orderId
	return s
}

// Amount sets the amount parameter (MANDATORY).
func (s *VipLoanRepayService) Amount(amount string) *VipLoanRepayService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *VipLoanRepayService) Do(ctx context.Context) (res *VipLoanRepayResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/loan/vip/repay",
		secType:  secTypeSigned,
	}
	r.setParam("orderId", s.orderId)
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(VipLoanRepayResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// VipLoanRepayResponse represents a response from repay flexible loan.
type VipLoanRepayResponse struct {
	LoanCoin           string `json:"loanCoin"`
	RepayAmount        string `json:"repayAmount"`
	RemainingPrincipal string `json:"remainingPrincipal"`
	RemainingInterest  string `json:"remainingInterest"`
	CollateralCoin     string `json:"collateralCoin"`
	CurrentLTV         string `json:"currentLTV"`
	RepayStatus        string `json:"repayStatus"` // Repaid, Repaying, Failed
}

// ListVipLoanService list flexible loan debt data.
type ListVipLoanService struct {
	c       *Client
	orderId *string
	limit   *int
}

// OrderId sets the orderId parameter.
func (s *ListVipLoanService) OrderId(orderId string) *ListVipLoanService {
	if len(orderId) > 0 {
		s.orderId = &orderId
	}
	return s
}

// Limit set limit
func (s *ListVipLoanService) Limit(limit int) *ListVipLoanService {
	s.limit = &limit
	return s
}

// Do sends the request.
func (s *ListVipLoanService) Do(ctx context.Context) (res *VipLoanOrderList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/loan/vip/ongoing/orders",
		secType:  secTypeSigned,
	}
	if s.orderId != nil {
		r.setParam("orderId", *s.orderId)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(VipLoanOrderList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// VipLoanOrderList represents a list of ongoing loan orders.
type VipLoanOrderList struct {
	Rows []struct {
		OrderId                          int    `json:"orderId"`
		LoanCoin                         string `json:"loanCoin"`
		TotalDebt                        string `json:"totalDebt"`
		LoanRate                         string `json:"loanRate"` // 浮动利率为"flexible rate"
		ResidualInterest                 string `json:"residualInterest"`
		CollateralAccountId              string `json:"collateralAccountId"`
		CollateralCoin                   string `json:"collateralCoin"`
		TotalCollateralValueAfterHaircut string `json:"totalCollateralValueAfterHaircut"`
		LockedCollateralValue            string `json:"lockedCollateralValue"`
		CurrentLTV                       string `json:"currentLTV"`
		ExpirationTime                   int64  `json:"expirationTime"` // 活期则为0
		LoanDate                         string `json:"loanDate"`
		LoanTerm                         string `json:"loanTerm"` // 活期则为"open term"
		InitialLtv                       string `json:"initialLtv"`
		MarginCallLtv                    string `json:"marginCallLtv"`
		LiquidationLtv                   string `json:"liquidationLtv"`
	} `json:"rows"`
	Total int `json:"total"`
}

//
//// AdjustLtvLoanService adjust flexible loan LTV.
//type AdjustLtvLoanService struct {
//	c                *Client
//	loanCoin         string
//	collateralCoin   string
//	adjustmentAmount string
//	direction        string
//}
//
//// LoanCoin sets the loan coin parameter.
//func (s *AdjustLtvLoanService) LoanCoin(coin string) *AdjustLtvLoanService {
//	s.loanCoin = coin
//	return s
//}
//
//// CollateralCoin set collateral coin parameter.
//func (s *AdjustLtvLoanService) CollateralCoin(collateralCoin string) *AdjustLtvLoanService {
//	s.collateralCoin = collateralCoin
//	return s
//}
//
//// AdjustmentAmount set collateral adjustment amount parameter.
//func (s *AdjustLtvLoanService) AdjustmentAmount(adjustmentAmount string) *AdjustLtvLoanService {
//	s.adjustmentAmount = adjustmentAmount
//	return s
//}
//
//// Direction set direction parameter, "ADDITIONAL", "REDUCED".
//func (s *AdjustLtvLoanService) Direction(direction string) *AdjustLtvLoanService {
//	s.direction = direction
//	return s
//}
//
//// Do sends the request.
//func (s *AdjustLtvLoanService) Do(ctx context.Context) (res *AdjustLtvLoanFlexibleResponse, err error) {
//	r := &request{
//		method:   http.MethodPost,
//		endpoint: "/sapi/v2/loan/flexible/adjust/ltv",
//		secType:  secTypeSigned,
//	}
//	r.setParam("loanCoin", s.loanCoin)
//	r.setParam("collateralCoin", s.collateralCoin)
//	r.setParam("adjustmentAmount", s.adjustmentAmount)
//	r.setParam("direction", s.direction)
//	data, _, err := s.c.callAPI(ctx, r)
//	if err != nil {
//		return
//	}
//	res = new(AdjustLtvLoanFlexibleResponse)
//	err = json.Unmarshal(data, res)
//	if err != nil {
//		return
//	}
//	return res, nil
//}
//
//// AdjustLtvLoanFlexibleResponse represents a response of adjust LTV of flexible loan.
//type AdjustLtvLoanFlexibleResponse struct {
//	LoanCoin         string `json:"loanCoin"`
//	CollateralCoin   string `json:"collateralCoin"`
//	Direction        string `json:"direction"`
//	AdjustmentAmount string `json:"adjustmentAmount"`
//	CurrentLTV       string `json:"currentLTV"`
//}
