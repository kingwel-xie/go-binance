package binance

import (
	"context"
	"net/http"
)

// SolStakingAccountService fetches the staking product positions
type SolStakingAccountService struct {
	c *Client
}

// Do sends the request.
func (s *SolStakingAccountService) Do(ctx context.Context) (*SolStakingAccountResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/sol-staking/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(SolStakingAccountResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SolStakingAccountResponse represents a SOL staking account.
type SolStakingAccountResponse struct {
	BnsolAmount           string `json:"bnsolAmount"`
	HoldingInSOL          string `json:"holdingInSOL"`
	ThirtyDaysProfitInSOL string `json:"thirtyDaysProfitInSOL"`
}

// SolStakingHistoryService fetches the staking history
type SolStakingHistoryService struct {
	c         *Client
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// StartTime sets the startTime parameter.
func (s *SolStakingHistoryService) StartTime(startTime int64) *SolStakingHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *SolStakingHistoryService) EndTime(endTime int64) *SolStakingHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *SolStakingHistoryService) Current(current int32) *SolStakingHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *SolStakingHistoryService) Size(size int32) *SolStakingHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *SolStakingHistoryService) Do(ctx context.Context) (*SolStakingHistoryResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/sol-staking/sol/history/stakingHistory",
		secType:  secTypeSigned,
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(SolStakingHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SolStakingHistoryResponse represents the SOL staking history.
type SolStakingHistoryResponse struct {
	Rows []struct {
		Time             int64  `json:"time"`
		Asset            string `json:"asset"`
		Amount           string `json:"amount"`
		DistributeAsset  string `json:"distributeAsset"` // BNSOL
		DistributeAmount string `json:"distributeAmount"`
		ExchangeRate     string `json:"exchangeRate"`
		Status           string `json:"status"` //PENDING,SUCCESS,FAILED
	} `json:"rows"`
	Total int `json:"total"`
}

// SolStakingRewardsHistoryService fetches the staking history
type SolStakingRewardsHistoryService struct {
	c         *Client
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// StartTime sets the startTime parameter.
func (s *SolStakingRewardsHistoryService) StartTime(startTime int64) *SolStakingRewardsHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *SolStakingRewardsHistoryService) EndTime(endTime int64) *SolStakingRewardsHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *SolStakingRewardsHistoryService) Current(current int32) *SolStakingRewardsHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *SolStakingRewardsHistoryService) Size(size int32) *SolStakingRewardsHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *SolStakingRewardsHistoryService) Do(ctx context.Context) (*SolStakingRewardsHistoryResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/sol-staking/sol/history/bnsolRewardsHistory",
		secType:  secTypeSigned,
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(SolStakingRewardsHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SolStakingRewardsHistoryResponse represents the SOL staking BNSOL rewards history.
type SolStakingRewardsHistoryResponse struct {
	Total           int    `json:"total"`
	EstRewardsInSOL string `json:"estRewardsInSOL"`
	Rows            []struct {
		Time                 int64  `json:"time"`
		AmountInSOL          string `json:"amountInSOL"`
		Holding              string `json:"holding"`
		HoldingInSOL         string `json:"holdingInSOL"`
		AnnualPercentageRate string `json:"annualPercentageRate"`
	} `json:"rows"`
}

// SolStakingRedemptionHistoryService fetches the staking history
type SolStakingRedemptionHistoryService struct {
	c         *Client
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// StartTime sets the startTime parameter.
func (s *SolStakingRedemptionHistoryService) StartTime(startTime int64) *SolStakingRedemptionHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *SolStakingRedemptionHistoryService) EndTime(endTime int64) *SolStakingRedemptionHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *SolStakingRedemptionHistoryService) Current(current int32) *SolStakingRedemptionHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *SolStakingRedemptionHistoryService) Size(size int32) *SolStakingRedemptionHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *SolStakingRedemptionHistoryService) Do(ctx context.Context) (*SolStakingRedemptionHistoryResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/sol-staking/sol/history/redemptionHistory",
		secType:  secTypeSigned,
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(SolStakingRedemptionHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SolStakingRedemptionHistoryResponse represents the SOL redemption history.
type SolStakingRedemptionHistoryResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		Time             int64  `json:"time"`
		ArrivalTime      int64  `json:"arrivalTime"`
		Asset            string `json:"asset"`
		Amount           string `json:"amount"`
		DistributeAsset  string `json:"distributeAsset"`
		DistributeAmount string `json:"distributeAmount"`
		ExchangeRate     string `json:"exchangeRate"`
		Status           string `json:"status"` // PENDING,SUCCESS,FAILED
	} `json:"rows"`
}

// SolStakingService stake SOL.
type SolStakingService struct {
	c      *Client
	amount string
}

// Amount sets the amount parameter (MANDATORY).
func (s *SolStakingService) Amount(amount string) *SolStakingService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *SolStakingService) Do(ctx context.Context) (res *SolStakingResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/sol-staking/sol/stake",
		secType:  secTypeSigned,
	}
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SolStakingResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// SolStakingResponse represents a response from SOL staking.
type SolStakingResponse struct {
	Success      bool   `json:"success"`
	BnsolAmount  string `json:"bnsolAmount"`
	ExchangeRate string `json:"exchangeRate"`
}

// SolRedeemService redeem BNSOL.
type SolRedeemService struct {
	c      *Client
	amount string
}

// Amount sets the amount parameter (MANDATORY).
func (s *SolRedeemService) Amount(amount string) *SolRedeemService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *SolRedeemService) Do(ctx context.Context) (res *SolRedeemResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/sol-staking/sol/redeem",
		secType:  secTypeSigned,
	}
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SolRedeemResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// SolRedeemResponse represents a response from SOL wrapping.
type SolRedeemResponse struct {
	Success      bool   `json:"success"`
	SolAmount    string `json:"solAmount"`
	ExchangeRate string `json:"exchangeRate"`
	ArrivalTime  int64  `json:"arrivalTime"`
}
