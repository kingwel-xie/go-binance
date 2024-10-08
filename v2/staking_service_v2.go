package binance

import (
	"context"
	"net/http"
)

// EthStakingAccountService fetches the staking product positions
type EthStakingAccountService struct {
	c *Client
}

// Do sends the request.
func (s *EthStakingAccountService) Do(ctx context.Context) (*EthStakingAccountResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v2/eth-staking/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}
	res := new(EthStakingAccountResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// EthStakingAccountResponse represents a ETH staking account.
type EthStakingAccountResponse struct {
	HoldingInETH string `json:"holdingInETH"`
	Holdings     struct {
		WbethAmount string `json:"wbethAmount"`
		BethAmount  string `json:"bethAmount"`
	} `json:"holdings"`
	ThirtyDaysProfitInETH string `json:"thirtyDaysProfitInETH"`
	Profit                struct {
		AmountFromWBETH string `json:"amountFromWBETH"`
		AmountFromBETH  string `json:"amountFromBETH"`
	} `json:"profit"`
}

// EthStakingHistoryService fetches the staking history
type EthStakingHistoryService struct {
	c         *Client
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// StartTime sets the startTime parameter.
func (s *EthStakingHistoryService) StartTime(startTime int64) *EthStakingHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *EthStakingHistoryService) EndTime(endTime int64) *EthStakingHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *EthStakingHistoryService) Current(current int32) *EthStakingHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *EthStakingHistoryService) Size(size int32) *EthStakingHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *EthStakingHistoryService) Do(ctx context.Context) (*EthStakingHistoryResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/eth-staking/eth/history/stakingHistory",
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
	res := new(EthStakingHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// EthStakingHistoryResponse represents the ETH staking history.
type EthStakingHistoryResponse struct {
	Rows []struct {
		Time             int64  `json:"time"`
		Asset            string `json:"asset"`
		Amount           string `json:"amount"`
		Status           string `json:"status"` //PENDING,SUCCESS,FAILED
		DistributeAmount string `json:"distributeAmount"`
		ConversionRatio  string `json:"conversionRatio"`
	} `json:"rows"`
	Total int `json:"total"`
}

// EthStakingRewardsHistoryService fetches the staking history
type EthStakingRewardsHistoryService struct {
	c         *Client
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// StartTime sets the startTime parameter.
func (s *EthStakingRewardsHistoryService) StartTime(startTime int64) *EthStakingRewardsHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *EthStakingRewardsHistoryService) EndTime(endTime int64) *EthStakingRewardsHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *EthStakingRewardsHistoryService) Current(current int32) *EthStakingRewardsHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *EthStakingRewardsHistoryService) Size(size int32) *EthStakingRewardsHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *EthStakingRewardsHistoryService) Do(ctx context.Context) (*EthStakingRewardsHistoryResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/eth-staking/eth/history/wbethRewardsHistory",
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
	res := new(EthStakingRewardsHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// EthStakingRewardsHistoryResponse represents the ETH staking WBETH rewards history.
type EthStakingRewardsHistoryResponse struct {
	Total           int    `json:"total"`
	EstRewardsInETH string `json:"estRewardsInETH"`
	Rows            []struct {
		Time                 int64  `json:"time"`
		AmountInETH          string `json:"amountInETH"`
		Holding              string `json:"holding"`
		HoldingInETH         string `json:"holdingInETH"`
		AnnualPercentageRate string `json:"annualPercentageRate"`
	} `json:"rows"`
}

// EthStakingRedemptionHistoryService fetches the staking history
type EthStakingRedemptionHistoryService struct {
	c         *Client
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// StartTime sets the startTime parameter.
func (s *EthStakingRedemptionHistoryService) StartTime(startTime int64) *EthStakingRedemptionHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *EthStakingRedemptionHistoryService) EndTime(endTime int64) *EthStakingRedemptionHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *EthStakingRedemptionHistoryService) Current(current int32) *EthStakingRedemptionHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *EthStakingRedemptionHistoryService) Size(size int32) *EthStakingRedemptionHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *EthStakingRedemptionHistoryService) Do(ctx context.Context) (*EthStakingRedemptionHistoryResponse, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/eth-staking/eth/history/redemptionHistory",
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
	res := new(EthStakingRedemptionHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// EthStakingRedemptionHistoryResponse represents the ETH redemption history.
type EthStakingRedemptionHistoryResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		Time             int64  `json:"time"`
		ArrivalTime      int64  `json:"arrivalTime"`
		Asset            string `json:"asset"`
		Amount           string `json:"amount"`
		Status           string `json:"status"`
		DistributeAsset  string `json:"distributeAsset"`
		DistributeAmount string `json:"distributeAmount"`
		ConversionRatio  string `json:"conversionRatio"`
	} `json:"rows"`
}

// EthStakingService stake ETH.
type EthStakingService struct {
	c      *Client
	amount string
}

// Amount sets the amount parameter (MANDATORY).
func (s *EthStakingService) Amount(amount string) *EthStakingService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *EthStakingService) Do(ctx context.Context) (res *EthStakingResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v2/eth-staking/eth/stake",
		secType:  secTypeSigned,
	}
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(EthStakingResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// EthStakingResponse represents a response from ETH staking.
type EthStakingResponse struct {
	Success         bool   `json:"success"`
	WbethAmount     string `json:"wbethAmount"`
	ConversionRatio string `json:"conversionRatio"`
}

// EthWrappingService stake ETH.
type EthWrappingService struct {
	c      *Client
	amount string
}

// Amount sets the amount parameter (MANDATORY).
func (s *EthWrappingService) Amount(amount string) *EthWrappingService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *EthWrappingService) Do(ctx context.Context) (res *EthWrappingResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/eth-staking/wbeth/wrap",
		secType:  secTypeSigned,
	}
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(EthWrappingResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// EthWrappingResponse represents a response from ETH wrapping.
type EthWrappingResponse struct {
	Success      bool   `json:"success"`
	WbethAmount  string `json:"wbethAmount"`
	ExchangeRate string `json:"exchangeRate"`
}

// EthRedeemService redeem BETH/WBETH.
type EthRedeemService struct {
	c      *Client
	asset  string
	amount string
}

// Asset sets the asset parameter (MANDATORY).
func (s *EthRedeemService) Asset(asset string) *EthRedeemService {
	s.asset = asset
	return s
}

// Amount sets the amount parameter (MANDATORY).
func (s *EthRedeemService) Amount(amount string) *EthRedeemService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *EthRedeemService) Do(ctx context.Context) (res *EthRedeemResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/eth-staking/eth/redeem",
		secType:  secTypeSigned,
	}
	r.setParam("asset", s.asset)
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(EthRedeemResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// EthRedeemResponse represents a response from ETH wrapping.
type EthRedeemResponse struct {
	Success         bool   `json:"success"`
	ArrivalTime     int64  `json:"arrivalTime"`
	EthAmount       string `json:"ethAmount"`
	ConversionRatio string `json:"conversionRatio"`
}
