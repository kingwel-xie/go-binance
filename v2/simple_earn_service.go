package binance

import (
	"context"
	"net/http"
)

// SimpleEarnAccountService gets simple-earn account info.
type SimpleEarnAccountService struct {
	c *Client
}

// Do sends the request.
func (s *SimpleEarnAccountService) Do(ctx context.Context) (res *SimpleEarnAccountResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/account",
		secType:  secTypeSigned,
	}

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SimpleEarnAccountResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnAccountResponse represents a simple-earn account response.
type SimpleEarnAccountResponse struct {
	TotalAmountInBTC          string `json:"totalAmountInBTC"`
	TotalAmountInUSDT         string `json:"totalAmountInUSDT"`
	TotalFlexibleAmountInBTC  string `json:"totalFlexibleAmountInBTC"`
	TotalFlexibleAmountInUSDT string `json:"totalFlexibleAmountInUSDT"`
	TotalLockedInBTC          string `json:"totalLockedInBTC"`
	TotalLockedInUSDT         string `json:"totalLockedInUSDT"`
}

// ListSimpleEarnFlexibleService list simple-earn flexible products.
type ListSimpleEarnFlexibleService struct {
	c       *Client
	asset   *string
	current *int32
	size    *int32
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnFlexibleService) Asset(asset string) *ListSimpleEarnFlexibleService {
	s.asset = &asset
	return s
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnFlexibleService) Current(current int32) *ListSimpleEarnFlexibleService {
	s.current = &current
	return s
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnFlexibleService) Size(size int32) *ListSimpleEarnFlexibleService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnFlexibleService) Do(ctx context.Context) (res *SimpleEarnFlexibleList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/list",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SimpleEarnFlexibleList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnFlexibleList represents a flexible simple-earn list.
type SimpleEarnFlexibleList struct {
	Rows []struct {
		Asset                      string `json:"asset"`
		LatestAnnualPercentageRate string `json:"latestAnnualPercentageRate"`
		TierAnnualPercentageRate   struct {
			//Zero5BTC  float64 `json:"0-5BTC"`
			//Five10BTC float64 `json:"5-10BTC"`
		} `json:"tierAnnualPercentageRate"`
		AirDropPercentageRate string `json:"airDropPercentageRate"`
		CanPurchase           bool   `json:"canPurchase"`
		CanRedeem             bool   `json:"canRedeem"`
		IsSoldOut             bool   `json:"isSoldOut"`
		Hot                   bool   `json:"hot"`
		MinPurchaseAmount     string `json:"minPurchaseAmount"`
		ProductID             string `json:"productId"`
		SubscriptionStartTime int64  `json:"subscriptionStartTime"`
		Status                string `json:"status"`
	} `json:"rows"`
	Total int `json:"total"`
}

// ListSimpleEarnLockedService list simple-earn locked products.
type ListSimpleEarnLockedService struct {
	c       *Client
	asset   *string
	current *int32
	size    *int32
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnLockedService) Asset(asset string) *ListSimpleEarnLockedService {
	s.asset = &asset
	return s
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnLockedService) Current(current int32) *ListSimpleEarnLockedService {
	s.current = &current
	return s
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnLockedService) Size(size int32) *ListSimpleEarnLockedService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnLockedService) Do(ctx context.Context) (res *SimpleEarnLockedList, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/locked/list",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SimpleEarnLockedList)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnLockedList represents a locked simple-earn list.
type SimpleEarnLockedList struct {
	Rows []struct {
		ProjectID string `json:"projectId"`
		Detail    struct {
			Asset                 string `json:"asset"`
			RewardAsset           string `json:"rewardAsset"`
			Duration              int    `json:"duration"`
			Renewable             bool   `json:"renewable"`
			IsSoldOut             bool   `json:"isSoldOut"`
			Apr                   string `json:"apr"`
			Status                string `json:"status"`
			SubscriptionStartTime int64  `json:"subscriptionStartTime"`
			ExtraRewardAsset      string `json:"extraRewardAsset"`
			ExtraRewardAPR        string `json:"extraRewardAPR"`
		} `json:"detail"`
		Quota struct {
			TotalPersonalQuota string `json:"totalPersonalQuota"`
			Minimum            string `json:"minimum"`
		} `json:"quota"`
	} `json:"rows"`
	Total int `json:"total"`
}

// SubscribeSimpleEarnFlexibleService subscribe to a simple-earn flexible product.
type SubscribeSimpleEarnFlexibleService struct {
	c         *Client
	productId string
	amount    float64
}

// Asset sets the asset parameter.
func (s *SubscribeSimpleEarnFlexibleService) ProductId(productId string) *SubscribeSimpleEarnFlexibleService {
	s.productId = productId
	return s
}

// Amount sets the Amount parameter (MANDATORY).
func (s *SubscribeSimpleEarnFlexibleService) Amount(v float64) *SubscribeSimpleEarnFlexibleService {
	s.amount = v
	return s
}

// Do sends the request.
func (s *SubscribeSimpleEarnFlexibleService) Do(ctx context.Context) (res *SubscribeSimpleEarnFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/flexible/subscribe",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	r.setParam("amount", s.amount)
	r.setParam("sourceAccount", "ALL")

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SubscribeSimpleEarnFlexibleResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// SubscribeSimpleEarnFlexibleResponse represents a response from subscribing flexible simple-earn.
type SubscribeSimpleEarnFlexibleResponse struct {
	PurchaseID int  `json:"purchaseId"`
	Success    bool `json:"success"`
}

// SubscribeSimpleEarnLockedService subscribe to a simple-earn locked product.
type SubscribeSimpleEarnLockedService struct {
	c         *Client
	projectId string
	amount    float64
}

// Asset sets the asset parameter.
func (s *SubscribeSimpleEarnLockedService) ProjectId(projectId string) *SubscribeSimpleEarnLockedService {
	s.projectId = projectId
	return s
}

// Amount sets the Amount parameter (MANDATORY).
func (s *SubscribeSimpleEarnLockedService) Amount(v float64) *SubscribeSimpleEarnLockedService {
	s.amount = v
	return s
}

// Do sends the request.
func (s *SubscribeSimpleEarnLockedService) Do(ctx context.Context) (res *SubscribeSimpleEarnLockedResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/locked/subscribe",
		secType:  secTypeSigned,
	}
	r.setParam("projectId", s.projectId)
	r.setParam("amount", s.amount)
	r.setParam("sourceAccount", "ALL")

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(SubscribeSimpleEarnLockedResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// SubscribeSimpleEarnLockedResponse represents a response from subscribing locked simple-earn.
type SubscribeSimpleEarnLockedResponse struct {
	PurchaseID int  `json:"purchaseId"`
	PositionID int  `json:"positionId"`
	Success    bool `json:"success"`
}

// RedeemSimpleEarnFlexibleService subscribe to a simple-earn flexible product.
type RedeemSimpleEarnFlexibleService struct {
	c         *Client
	productId string
	amount    float64
}

// ProductId sets the productId parameter.
func (s *RedeemSimpleEarnFlexibleService) ProductId(productId string) *RedeemSimpleEarnFlexibleService {
	s.productId = productId
	return s
}

// Amount sets the Amount parameter (MANDATORY when redeemAll is false).
func (s *RedeemSimpleEarnFlexibleService) Amount(v float64) *RedeemSimpleEarnFlexibleService {
	s.amount = v
	return s
}

// Do sends the request.
func (s *RedeemSimpleEarnFlexibleService) Do(ctx context.Context) (res *RedeemSimpleEarnFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/flexible/redeem",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	if s.amount == 0 {
		r.setParam("redeemAll", true)
	} else {
		r.setParam("amount", s.amount)
	}
	r.setParam("destAccount", "SPOT")

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(RedeemSimpleEarnFlexibleResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// RedeemSimpleEarnFlexibleResponse represents a response from redeeming flexible simple-earn.
type RedeemSimpleEarnFlexibleResponse struct {
	RedeemID int  `json:"redeemId"`
	Success  bool `json:"success"`
}

// RedeemSimpleEarnLockedService redeem a simple-earn locked product.
type RedeemSimpleEarnLockedService struct {
	c          *Client
	positionId string
}

// PositionId sets the positionId parameter.
func (s *RedeemSimpleEarnLockedService) PositionId(positionId string) *RedeemSimpleEarnLockedService {
	s.positionId = positionId
	return s
}

// Do sends the request.
func (s *RedeemSimpleEarnLockedService) Do(ctx context.Context) (res *RedeemSimpleEarnLockedResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/locked/redeem",
		secType:  secTypeSigned,
	}
	r.setParam("positionId", s.positionId)
	//r.setParam("destAccount", "SPOT")

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(RedeemSimpleEarnLockedResponse)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// RedeemSimpleEarnLockedResponse represents a response from redeeming flexible simple-earn.
type RedeemSimpleEarnLockedResponse struct {
	RedeemID int  `json:"redeemId"`
	Success  bool `json:"success"`
}

// Position related

// GetSimpleEarnFlexiblePositionService get simple-earn flexisible position..
type GetSimpleEarnFlexiblePositionService struct {
	c       *Client
	asset   *string
	current *int32
	size    *int32
}

// Asset sets the asset parameter.
func (s *GetSimpleEarnFlexiblePositionService) Asset(asset string) *GetSimpleEarnFlexiblePositionService {
	s.asset = &asset
	return s
}

// Current sets the current parameter.
func (s *GetSimpleEarnFlexiblePositionService) Current(current int32) *GetSimpleEarnFlexiblePositionService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *GetSimpleEarnFlexiblePositionService) Size(size int32) *GetSimpleEarnFlexiblePositionService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *GetSimpleEarnFlexiblePositionService) Do(ctx context.Context) (res *GetSimpleEarnFlexiblePositionResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/position",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(GetSimpleEarnFlexiblePositionResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// GetSimpleEarnFlexiblePositionResponse represents the simple-earn flexible postion.
type GetSimpleEarnFlexiblePositionResponse struct {
	Rows []struct {
		TotalAmount              string `json:"totalAmount"`
		TierAnnualPercentageRate struct {
		} `json:"tierAnnualPercentageRate"`
		LatestAnnualPercentageRate     string `json:"latestAnnualPercentageRate"`
		YesterdayAirdropPercentageRate string `json:"yesterdayAirdropPercentageRate"`
		Asset                          string `json:"asset"`
		AirDropAsset                   string `json:"airDropAsset"`
		CanRedeem                      bool   `json:"canRedeem"`
		CollateralAmount               string `json:"collateralAmount"`
		ProductID                      string `json:"productId"`
		YesterdayRealTimeRewards       string `json:"yesterdayRealTimeRewards"`
		CumulativeBonusRewards         string `json:"cumulativeBonusRewards"`
		CumulativeRealTimeRewards      string `json:"cumulativeRealTimeRewards"`
		CumulativeTotalRewards         string `json:"cumulativeTotalRewards"`
		AutoSubscribe                  bool   `json:"autoSubscribe"`
	} `json:"rows"`
	Total int `json:"total"`
}

// GetSimpleEarnLockedPositionService get simple-earn locked position..
type GetSimpleEarnLockedPositionService struct {
	c       *Client
	asset   *string
	current *int32
	size    *int32
}

// Asset sets the asset parameter.
func (s *GetSimpleEarnLockedPositionService) Asset(asset string) *GetSimpleEarnLockedPositionService {
	s.asset = &asset
	return s
}

// Asset sets the asset parameter.
func (s *GetSimpleEarnLockedPositionService) Current(current int32) *GetSimpleEarnLockedPositionService {
	s.current = &current
	return s
}

// Asset sets the asset parameter.
func (s *GetSimpleEarnLockedPositionService) Size(size int32) *GetSimpleEarnLockedPositionService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *GetSimpleEarnLockedPositionService) Do(ctx context.Context) (res *GetSimpleEarnLockedPositionResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/locked/position",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(GetSimpleEarnLockedPositionResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// GetSimpleEarnLockedPositionResponse represents the simple-earn flexible postion.
type GetSimpleEarnLockedPositionResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		PositionID        int    `json:"positionId"`
		ProjectID         string `json:"projectId"`
		Asset             string `json:"asset"`
		Amount            string `json:"amount"`
		PurchaseTime      int64  `json:"purchaseTime"`
		Duration          int    `json:"duration"`
		AccrualDays       int    `json:"accrualDays"`
		RewardAsset       string `json:"rewardAsset"`
		RewardAmt         string `json:"rewardAmt"`
		NextPay           string `json:"nextPay"`
		NextPayDate       int64  `json:"nextPayDate"`
		PayPeriod         int    `json:"payPeriod"`
		RedeemAmountEarly string `json:"redeemAmountEarly"`
		RewardsEndDate    int64  `json:"rewardsEndDate"`
		DeliverDate       int64  `json:"deliverDate"`
		RedeemPeriod      int    `json:"redeemPeriod"`
		CanRedeemEarly    bool   `json:"canRedeemEarly"`
		AutoSubscribe     bool   `json:"autoSubscribe"`
		Type              string `json:"type"`
		Status            string `json:"status"`
		CanReStake        bool   `json:"canReStake"`
		Apy               string `json:"apy"`
	} `json:"rows"`
}

// ListSimpleEarnFlexibleRateHistoryService list simple-earn locked products.
type ListSimpleEarnFlexibleRateHistoryService struct {
	c         *Client
	productId string
	current   *int32
	size      *int32
}

// ProductId sets the productId parameter.
func (s *ListSimpleEarnFlexibleRateHistoryService) ProductId(productId string) *ListSimpleEarnFlexibleRateHistoryService {
	s.productId = productId
	return s
}

// Current sets the current parameter.
func (s *ListSimpleEarnFlexibleRateHistoryService) Current(current int32) *ListSimpleEarnFlexibleRateHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *ListSimpleEarnFlexibleRateHistoryService) Size(size int32) *ListSimpleEarnFlexibleRateHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnFlexibleRateHistoryService) Do(ctx context.Context) (res *ListSimpleEarnFlexibleRateHistoryResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/history/rateHistory",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return
	}
	res = new(ListSimpleEarnFlexibleRateHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

type ListSimpleEarnFlexibleRateHistoryResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		Asset                string `json:"asset"`
		AnnualPercentageRate string `json:"annualPercentageRate"`
		ProductID            string `json:"productId"`
		Time                 int64  `json:"time"`
	} `json:"rows"`
}
