package binance

import (
	"context"
	stdjson "encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// SimpleEarnAccountService gets simple-earn account info.
// See https://binance-docs.github.io/apidocs/spot/en/#get-account-user_data
type SimpleEarnAccountService struct {
	c *Client
}

// Do sends the request.
func (s *SimpleEarnAccountService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnAccountResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/account",
		secType:  secTypeSigned,
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// TierAnnualPercentageRate represents the tiered APR returned by simple-earn
// flexible product list/position endpoints. The values can arrive either as
// JSON numbers or strings depending on the endpoint, so both are accepted and
// normalized to string.
type TierAnnualPercentageRate map[string]string

// UnmarshalJSON accepts both numeric and string tier APR values.
func (m *TierAnnualPercentageRate) UnmarshalJSON(data []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	res := make(TierAnnualPercentageRate, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			res[k] = val
		case float64:
			res[k] = strconv.FormatFloat(val, 'f', -1, 64)
		default:
			return fmt.Errorf("unsupported tierAnnualPercentageRate value type %T for key %q", v, k)
		}
	}
	*m = res
	return nil
}

// ListSimpleEarnFlexibleService list simple-earn flexible products.
// See https://binance-docs.github.io/apidocs/spot/en/#get-simple-earn-flexible-product-list-user_data
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

// Current sets the current parameter.
func (s *ListSimpleEarnFlexibleService) Current(current int32) *ListSimpleEarnFlexibleService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *ListSimpleEarnFlexibleService) Size(size int32) *ListSimpleEarnFlexibleService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnFlexibleService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnFlexibleList, err error) {
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

	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// SimpleEarnFlexibleProduct represents a flexible simple-earn product.
type SimpleEarnFlexibleProduct struct {
	Asset                      string                   `json:"asset"`
	LatestAnnualPercentageRate string                   `json:"latestAnnualPercentageRate"`
	TierAnnualPercentageRate   TierAnnualPercentageRate `json:"tierAnnualPercentageRate"`
	AirDropPercentageRate      string                   `json:"airDropPercentageRate"`
	CanPurchase                bool                     `json:"canPurchase"`
	CanRedeem                  bool                     `json:"canRedeem"`
	IsSoldOut                  bool                     `json:"isSoldOut"`
	Hot                        bool                     `json:"hot"`
	MinPurchaseAmount          string                   `json:"minPurchaseAmount"`
	ProductID                  string                   `json:"productId"`
	SubscriptionStartTime      int64                    `json:"subscriptionStartTime"`
	Status                     string                   `json:"status"`
}

// SimpleEarnFlexibleList represents a flexible simple-earn list.
type SimpleEarnFlexibleList struct {
	Rows  []SimpleEarnFlexibleProduct `json:"rows"`
	Total int                         `json:"total"`
}

// ListSimpleEarnLockedService list simple-earn locked products.
// See https://binance-docs.github.io/apidocs/spot/en/#get-simple-earn-locked-product-list-user_data
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

// Current sets the current parameter.
func (s *ListSimpleEarnLockedService) Current(current int32) *ListSimpleEarnLockedService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *ListSimpleEarnLockedService) Size(size int32) *ListSimpleEarnLockedService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnLockedService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnLockedList, err error) {
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
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// SimpleEarnLockedProduct represents a locked simple-earn product.
type SimpleEarnLockedProduct struct {
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
}

// SimpleEarnLockedList represents a locked simple-earn list.
type SimpleEarnLockedList struct {
	Rows  []SimpleEarnLockedProduct `json:"rows"`
	Total int                       `json:"total"`
}

// SubscribeSimpleEarnFlexibleService subscribe to a simple-earn flexible product.
// See https://binance-docs.github.io/apidocs/spot/en/#subscribe-flexible-product-trade
type SubscribeSimpleEarnFlexibleService struct {
	c             *Client
	productId     string
	amount        string
	autoSubscribe *bool
	sourceAccount SimpleEarnSubscribeSourceAccount
}

// ProductId sets the productId parameter.
func (s *SubscribeSimpleEarnFlexibleService) ProductId(productId string) *SubscribeSimpleEarnFlexibleService {
	s.productId = productId
	return s
}

// Amount sets the Amount parameter (MANDATORY).
func (s *SubscribeSimpleEarnFlexibleService) Amount(v string) *SubscribeSimpleEarnFlexibleService {
	s.amount = v
	return s
}

// AutoSubscribe sets the autoSubscribe parameter (true or false, default true).
func (s *SubscribeSimpleEarnFlexibleService) AutoSubscribe(autoSubscribe bool) *SubscribeSimpleEarnFlexibleService {
	s.autoSubscribe = &autoSubscribe
	return s
}

// SourceAccount sets the sourceAccount parameter (SPOT, FUND, ALL; default SPOT).
func (s *SubscribeSimpleEarnFlexibleService) SourceAccount(sourceAccount SimpleEarnSubscribeSourceAccount) *SubscribeSimpleEarnFlexibleService {
	s.sourceAccount = sourceAccount
	return s
}

// Do sends the request.
func (s *SubscribeSimpleEarnFlexibleService) Do(ctx context.Context, opts ...RequestOption) (res *SubscribeSimpleEarnFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/flexible/subscribe",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	r.setParam("amount", s.amount)
	if s.autoSubscribe != nil {
		r.setParam("autoSubscribe", *s.autoSubscribe)
	}
	if s.sourceAccount != "" {
		r.setParam("sourceAccount", string(s.sourceAccount))
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
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
// See https://binance-docs.github.io/apidocs/spot/en/#subscribe-locked-product-trade
type SubscribeSimpleEarnLockedService struct {
	c             *Client
	projectId     string
	amount        string
	autoSubscribe *bool
	sourceAccount SimpleEarnSubscribeSourceAccount
	redeemTo      RedeemTo
}

// ProjectId sets the projectId parameter.
func (s *SubscribeSimpleEarnLockedService) ProjectId(projectId string) *SubscribeSimpleEarnLockedService {
	s.projectId = projectId
	return s
}

// Amount sets the Amount parameter (MANDATORY).
func (s *SubscribeSimpleEarnLockedService) Amount(v string) *SubscribeSimpleEarnLockedService {
	s.amount = v
	return s
}

// AutoSubscribe sets the autoSubscribe parameter (true or false, default true).
func (s *SubscribeSimpleEarnLockedService) AutoSubscribe(autoSubscribe bool) *SubscribeSimpleEarnLockedService {
	s.autoSubscribe = &autoSubscribe
	return s
}

// SourceAccount sets the sourceAccount parameter (SPOT, FUND, ALL; default SPOT).
func (s *SubscribeSimpleEarnLockedService) SourceAccount(sourceAccount SimpleEarnSubscribeSourceAccount) *SubscribeSimpleEarnLockedService {
	s.sourceAccount = sourceAccount
	return s
}

// RedeemTo sets the redeemTo parameter (SPOT, FLEXIBLE; default SPOT).
func (s *SubscribeSimpleEarnLockedService) RedeemTo(redeemTo RedeemTo) *SubscribeSimpleEarnLockedService {
	s.redeemTo = redeemTo
	return s
}

// Do sends the request.
func (s *SubscribeSimpleEarnLockedService) Do(ctx context.Context, opts ...RequestOption) (res *SubscribeSimpleEarnLockedResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/locked/subscribe",
		secType:  secTypeSigned,
	}
	r.setParam("projectId", s.projectId)
	r.setParam("amount", s.amount)
	if s.autoSubscribe != nil {
		r.setParam("autoSubscribe", *s.autoSubscribe)
	}
	if s.sourceAccount != "" {
		r.setParam("sourceAccount", string(s.sourceAccount))
	}
	if s.redeemTo != "" {
		r.setParam("redeemTo", string(s.redeemTo))
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// RedeemSimpleEarnFlexibleService redeem a simple-earn flexible product.
// See https://binance-docs.github.io/apidocs/spot/en/#redeem-flexible-product-trade
type RedeemSimpleEarnFlexibleService struct {
	c           *Client
	productId   string
	redeemAll   *bool
	amount      *string
	destAccount string
}

// ProductId sets the productId parameter.
func (s *RedeemSimpleEarnFlexibleService) ProductId(productId string) *RedeemSimpleEarnFlexibleService {
	s.productId = productId
	return s
}

// RedeemAll sets the redeemAll parameter (true or false, default false).
func (s *RedeemSimpleEarnFlexibleService) RedeemAll(redeemAll bool) *RedeemSimpleEarnFlexibleService {
	s.redeemAll = &redeemAll
	return s
}

// Amount sets the Amount parameter (MANDATORY when redeemAll is false).
func (s *RedeemSimpleEarnFlexibleService) Amount(v *string) *RedeemSimpleEarnFlexibleService {
	s.amount = v
	return s
}

// DestAccount sets the destAccount parameter (SPOT, FUND; default SPOT).
func (s *RedeemSimpleEarnFlexibleService) DestAccount(destAccount string) *RedeemSimpleEarnFlexibleService {
	s.destAccount = destAccount
	return s
}

// Do sends the request.
func (s *RedeemSimpleEarnFlexibleService) Do(ctx context.Context, opts ...RequestOption) (res *RedeemSimpleEarnFlexibleResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/flexible/redeem",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	switch {
	case s.redeemAll != nil:
		r.setParam("redeemAll", *s.redeemAll)
		if !*s.redeemAll && s.amount != nil && *s.amount != "" {
			r.setParam("amount", *s.amount)
		}
	case s.amount != nil && *s.amount != "":
		r.setParam("amount", *s.amount)
	default:
		r.setParam("redeemAll", true)
	}
	if s.destAccount != "" {
		r.setParam("destAccount", s.destAccount)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
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
// See https://binance-docs.github.io/apidocs/spot/en/#redeem-locked-product-trade
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
func (s *RedeemSimpleEarnLockedService) Do(ctx context.Context, opts ...RequestOption) (res *RedeemSimpleEarnLockedResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/locked/redeem",
		secType:  secTypeSigned,
	}
	r.setParam("positionId", s.positionId)

	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// RedeemSimpleEarnLockedResponse represents a response from redeeming locked simple-earn.
type RedeemSimpleEarnLockedResponse struct {
	RedeemID int  `json:"redeemId"`
	Success  bool `json:"success"`
}

// GetSimpleEarnFlexiblePositionService get simple-earn flexible position.
// See https://binance-docs.github.io/apidocs/spot/en/#get-flexible-product-position-user_data
type GetSimpleEarnFlexiblePositionService struct {
	c         *Client
	asset     *string
	productId *string
	current   *int32
	size      *int32
}

// Asset sets the asset parameter.
func (s *GetSimpleEarnFlexiblePositionService) Asset(asset string) *GetSimpleEarnFlexiblePositionService {
	s.asset = &asset
	return s
}

// ProductId sets the productId parameter.
func (s *GetSimpleEarnFlexiblePositionService) ProductId(productId string) *GetSimpleEarnFlexiblePositionService {
	s.productId = &productId
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
func (s *GetSimpleEarnFlexiblePositionService) Do(ctx context.Context, opts ...RequestOption) (res *GetSimpleEarnFlexiblePositionResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/position",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.productId != nil {
		r.setParam("productId", *s.productId)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// SimpleEarnFlexiblePosition represents one simple-earn flexible position.
type SimpleEarnFlexiblePosition struct {
	TotalAmount                    string                   `json:"totalAmount"`
	TierAnnualPercentageRate       TierAnnualPercentageRate `json:"tierAnnualPercentageRate"`
	LatestAnnualPercentageRate     string                   `json:"latestAnnualPercentageRate"`
	YesterdayAirdropPercentageRate string                   `json:"yesterdayAirdropPercentageRate"`
	Asset                          string                   `json:"asset"`
	AirDropAsset                   string                   `json:"airDropAsset"`
	CanRedeem                      bool                     `json:"canRedeem"`
	CollateralAmount               string                   `json:"collateralAmount"`
	ProductID                      string                   `json:"productId"`
	YesterdayRealTimeRewards       string                   `json:"yesterdayRealTimeRewards"`
	CumulativeBonusRewards         string                   `json:"cumulativeBonusRewards"`
	CumulativeRealTimeRewards      string                   `json:"cumulativeRealTimeRewards"`
	CumulativeTotalRewards         string                   `json:"cumulativeTotalRewards"`
	AutoSubscribe                  bool                     `json:"autoSubscribe"`
}

// GetSimpleEarnFlexiblePositionResponse represents the simple-earn flexible position.
type GetSimpleEarnFlexiblePositionResponse struct {
	Rows  []SimpleEarnFlexiblePosition `json:"rows"`
	Total int                          `json:"total"`
}

// GetSimpleEarnLockedPositionService get simple-earn locked position.
// See https://binance-docs.github.io/apidocs/spot/en/#get-locked-product-position-user_data
type GetSimpleEarnLockedPositionService struct {
	c          *Client
	asset      *string
	positionId *int64
	projectId  *string
	current    *int32
	size       *int32
}

// Asset sets the asset parameter.
func (s *GetSimpleEarnLockedPositionService) Asset(asset string) *GetSimpleEarnLockedPositionService {
	s.asset = &asset
	return s
}

// PositionId sets the positionId parameter.
func (s *GetSimpleEarnLockedPositionService) PositionId(positionId int64) *GetSimpleEarnLockedPositionService {
	s.positionId = &positionId
	return s
}

// ProjectId sets the projectId parameter.
func (s *GetSimpleEarnLockedPositionService) ProjectId(projectId string) *GetSimpleEarnLockedPositionService {
	s.projectId = &projectId
	return s
}

// Current sets the current parameter.
func (s *GetSimpleEarnLockedPositionService) Current(current int32) *GetSimpleEarnLockedPositionService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *GetSimpleEarnLockedPositionService) Size(size int32) *GetSimpleEarnLockedPositionService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *GetSimpleEarnLockedPositionService) Do(ctx context.Context, opts ...RequestOption) (res *GetSimpleEarnLockedPositionResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/locked/position",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.positionId != nil {
		r.setParam("positionId", *s.positionId)
	}
	if s.projectId != nil {
		r.setParam("projectId", *s.projectId)
	}
	if s.current != nil {
		r.setParam("current", *s.current)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// SimpleEarnLockedPosition represents one simple-earn locked position.
type SimpleEarnLockedPosition struct {
	PositionID            int64  `json:"positionId"`
	ParentPositionID      int64  `json:"parentPositionId"`
	ProjectID             string `json:"projectId"`
	Asset                 string `json:"asset"`
	Amount                string `json:"amount"`
	PurchaseTime          int64  `json:"purchaseTime"`
	Duration              int    `json:"duration"`
	AccrualDays           int    `json:"accrualDays"`
	RewardAsset           string `json:"rewardAsset"`
	Apy                   string `json:"APY"`
	RewardAmt             string `json:"rewardAmt"`
	ExtraRewardAsset      string `json:"extraRewardAsset"`
	ExtraRewardAPR        string `json:"extraRewardAPR"`
	EstExtraRewardAmt     string `json:"estExtraRewardAmt"`
	NextPay               string `json:"nextPay"`
	NextPayDate           int64  `json:"nextPayDate"`
	PayPeriod             int    `json:"payPeriod"`
	RedeemAmountEarly     string `json:"redeemAmountEarly"`
	RewardsEndDate        int64  `json:"rewardsEndDate"`
	DeliverDate           int64  `json:"deliverDate"`
	RedeemPeriod          int    `json:"redeemPeriod"`
	RedeemingAmt          string `json:"redeemingAmt"`
	RedeemTo              string `json:"redeemTo"`
	PartialAmtDeliverDate int64  `json:"partialAmtDeliverDate"`
	CanRedeemEarly        bool   `json:"canRedeemEarly"`
	CanFastRedemption     bool   `json:"canFastRedemption"`
	AutoSubscribe         bool   `json:"autoSubscribe"`
	Type                  string `json:"type"`
	Status                string `json:"status"`
	CanReStake            bool   `json:"canReStake"`
}

// GetSimpleEarnLockedPositionResponse represents the simple-earn locked position.
type GetSimpleEarnLockedPositionResponse struct {
	Total int                        `json:"total"`
	Rows  []SimpleEarnLockedPosition `json:"rows"`
}

// ListSimpleEarnFlexibleRateHistoryService list simple-earn flexible rate history.
// See https://binance-docs.github.io/apidocs/spot/en/#get-rate-history-user_data
type ListSimpleEarnFlexibleRateHistoryService struct {
	c         *Client
	productId string
	aprPeriod string
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// ProductId sets the productId parameter.
func (s *ListSimpleEarnFlexibleRateHistoryService) ProductId(productId string) *ListSimpleEarnFlexibleRateHistoryService {
	s.productId = productId
	return s
}

// AprPeriod sets the aprPeriod parameter ("DAY", "YEAR"; default "DAY").
func (s *ListSimpleEarnFlexibleRateHistoryService) AprPeriod(aprPeriod string) *ListSimpleEarnFlexibleRateHistoryService {
	s.aprPeriod = aprPeriod
	return s
}

// StartTime sets the startTime parameter.
func (s *ListSimpleEarnFlexibleRateHistoryService) StartTime(startTime int64) *ListSimpleEarnFlexibleRateHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *ListSimpleEarnFlexibleRateHistoryService) EndTime(endTime int64) *ListSimpleEarnFlexibleRateHistoryService {
	s.endTime = &endTime
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
func (s *ListSimpleEarnFlexibleRateHistoryService) Do(ctx context.Context, opts ...RequestOption) (res *ListSimpleEarnFlexibleRateHistoryResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/history/rateHistory",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	if s.aprPeriod != "" {
		r.setParam("aprPeriod", s.aprPeriod)
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
	data, _, err := s.c.callAPI(ctx, r, opts...)
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

// ListSimpleEarnFlexibleRateHistoryResponse represents the flexible rate history.
// The official response example reports "total" as a quoted string, so
// json.Number is used to accept both string and number forms.
type ListSimpleEarnFlexibleRateHistoryResponse struct {
	Total stdjson.Number `json:"total"`
	Rows  []struct {
		Asset                string `json:"asset"`
		AnnualPercentageRate string `json:"annualPercentageRate"`
		ProductID            string `json:"productId"`
		Time                 int64  `json:"time"`
	} `json:"rows"`
}

// GetSimpleEarnFlexibleQuotaService get flexible personal left quota.
// See https://binance-docs.github.io/apidocs/spot/en/#get-flexible-personal-left-quota-user_data
type GetSimpleEarnFlexibleQuotaService struct {
	c         *Client
	productId string
}

// ProductId sets the productId parameter.
func (s *GetSimpleEarnFlexibleQuotaService) ProductId(productId string) *GetSimpleEarnFlexibleQuotaService {
	s.productId = productId
	return s
}

// Do sends the request.
func (s *GetSimpleEarnFlexibleQuotaService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnPersonalLeftQuotaResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/personalLeftQuota",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnPersonalLeftQuotaResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// GetSimpleEarnLockedQuotaService get locked personal left quota.
// See https://binance-docs.github.io/apidocs/spot/en/#get-locked-personal-left-quota-user_data
type GetSimpleEarnLockedQuotaService struct {
	c         *Client
	projectId string
}

// ProjectId sets the projectId parameter.
func (s *GetSimpleEarnLockedQuotaService) ProjectId(projectId string) *GetSimpleEarnLockedQuotaService {
	s.projectId = projectId
	return s
}

// Do sends the request.
func (s *GetSimpleEarnLockedQuotaService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnPersonalLeftQuotaResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/locked/personalLeftQuota",
		secType:  secTypeSigned,
	}
	r.setParam("projectId", s.projectId)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnPersonalLeftQuotaResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnPersonalLeftQuotaResponse represents a personal left quota response.
type SimpleEarnPersonalLeftQuotaResponse struct {
	LeftPersonalQuota string `json:"leftPersonalQuota"`
}

// SetSimpleEarnFlexibleAutoSubscribeService set flexible auto subscribe.
// See https://binance-docs.github.io/apidocs/spot/en/#set-flexible-auto-subscribe-user_data
type SetSimpleEarnFlexibleAutoSubscribeService struct {
	c             *Client
	productId     string
	autoSubscribe *bool
}

// ProductId sets the productId parameter.
func (s *SetSimpleEarnFlexibleAutoSubscribeService) ProductId(productId string) *SetSimpleEarnFlexibleAutoSubscribeService {
	s.productId = productId
	return s
}

// AutoSubscribe sets the autoSubscribe parameter (MANDATORY, true or false).
func (s *SetSimpleEarnFlexibleAutoSubscribeService) AutoSubscribe(autoSubscribe bool) *SetSimpleEarnFlexibleAutoSubscribeService {
	s.autoSubscribe = &autoSubscribe
	return s
}

// Do sends the request.
func (s *SetSimpleEarnFlexibleAutoSubscribeService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnSetAutoSubscribeResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/flexible/setAutoSubscribe",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	if s.autoSubscribe != nil {
		r.setParam("autoSubscribe", *s.autoSubscribe)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnSetAutoSubscribeResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SetSimpleEarnLockedAutoSubscribeService set locked auto subscribe.
// See https://binance-docs.github.io/apidocs/spot/en/#set-locked-auto-subscribe-user_data
type SetSimpleEarnLockedAutoSubscribeService struct {
	c             *Client
	positionId    string
	autoSubscribe *bool
}

// PositionId sets the positionId parameter.
func (s *SetSimpleEarnLockedAutoSubscribeService) PositionId(positionId string) *SetSimpleEarnLockedAutoSubscribeService {
	s.positionId = positionId
	return s
}

// AutoSubscribe sets the autoSubscribe parameter (MANDATORY, true or false).
func (s *SetSimpleEarnLockedAutoSubscribeService) AutoSubscribe(autoSubscribe bool) *SetSimpleEarnLockedAutoSubscribeService {
	s.autoSubscribe = &autoSubscribe
	return s
}

// Do sends the request.
func (s *SetSimpleEarnLockedAutoSubscribeService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnSetAutoSubscribeResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/locked/setAutoSubscribe",
		secType:  secTypeSigned,
	}
	r.setParam("positionId", s.positionId)
	if s.autoSubscribe != nil {
		r.setParam("autoSubscribe", *s.autoSubscribe)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnSetAutoSubscribeResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnSetAutoSubscribeResponse represents a set auto subscribe response.
type SimpleEarnSetAutoSubscribeResponse struct {
	Success bool `json:"success"`
}

// GetSimpleEarnFlexibleSubscriptionPreviewService get flexible subscription preview.
// See https://binance-docs.github.io/apidocs/spot/en/#get-flexible-subscription-preview-user_data
type GetSimpleEarnFlexibleSubscriptionPreviewService struct {
	c         *Client
	productId string
	amount    string
}

// ProductId sets the productId parameter.
func (s *GetSimpleEarnFlexibleSubscriptionPreviewService) ProductId(productId string) *GetSimpleEarnFlexibleSubscriptionPreviewService {
	s.productId = productId
	return s
}

// Amount sets the Amount parameter (MANDATORY).
func (s *GetSimpleEarnFlexibleSubscriptionPreviewService) Amount(amount string) *GetSimpleEarnFlexibleSubscriptionPreviewService {
	s.amount = amount
	return s
}

// Do sends the request.
func (s *GetSimpleEarnFlexibleSubscriptionPreviewService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnFlexibleSubscriptionPreviewResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/subscriptionPreview",
		secType:  secTypeSigned,
	}
	r.setParam("productId", s.productId)
	r.setParam("amount", s.amount)

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnFlexibleSubscriptionPreviewResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnFlexibleSubscriptionPreviewResponse represents a flexible subscription preview.
type SimpleEarnFlexibleSubscriptionPreviewResponse struct {
	TotalAmount             string `json:"totalAmount"`
	RewardAsset             string `json:"rewardAsset"`
	AirDropAsset            string `json:"airDropAsset"`
	EstDailyBonusRewards    string `json:"estDailyBonusRewards"`
	EstDailyRealTimeRewards string `json:"estDailyRealTimeRewards"`
	EstDailyAirdropRewards  string `json:"estDailyAirdropRewards"`
}

// GetSimpleEarnLockedSubscriptionPreviewService get locked subscription preview.
// See https://binance-docs.github.io/apidocs/spot/en/#get-locked-subscription-preview-user_data
type GetSimpleEarnLockedSubscriptionPreviewService struct {
	c             *Client
	projectId     string
	amount        string
	autoSubscribe *bool
}

// ProjectId sets the projectId parameter.
func (s *GetSimpleEarnLockedSubscriptionPreviewService) ProjectId(projectId string) *GetSimpleEarnLockedSubscriptionPreviewService {
	s.projectId = projectId
	return s
}

// Amount sets the Amount parameter (MANDATORY).
func (s *GetSimpleEarnLockedSubscriptionPreviewService) Amount(amount string) *GetSimpleEarnLockedSubscriptionPreviewService {
	s.amount = amount
	return s
}

// AutoSubscribe sets the autoSubscribe parameter (true or false, default true).
func (s *GetSimpleEarnLockedSubscriptionPreviewService) AutoSubscribe(autoSubscribe bool) *GetSimpleEarnLockedSubscriptionPreviewService {
	s.autoSubscribe = &autoSubscribe
	return s
}

// Do sends the request. The response is an array of previews.
func (s *GetSimpleEarnLockedSubscriptionPreviewService) Do(ctx context.Context, opts ...RequestOption) (res []*SimpleEarnLockedSubscriptionPreview, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/locked/subscriptionPreview",
		secType:  secTypeSigned,
	}
	r.setParam("projectId", s.projectId)
	r.setParam("amount", s.amount)
	if s.autoSubscribe != nil {
		r.setParam("autoSubscribe", *s.autoSubscribe)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = make([]*SimpleEarnLockedSubscriptionPreview, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnLockedSubscriptionPreview represents one locked subscription preview.
type SimpleEarnLockedSubscriptionPreview struct {
	RewardAsset            string `json:"rewardAsset"`
	TotalRewardAmt         string `json:"totalRewardAmt"`
	ExtraRewardAsset       string `json:"extraRewardAsset"`
	EstTotalExtraRewardAmt string `json:"estTotalExtraRewardAmt"`
	NextPay                string `json:"nextPay"`
	NextPayDate            int64  `json:"nextPayDate"`
	ValueDate              int64  `json:"valueDate"`
	RewardsEndDate         int64  `json:"rewardsEndDate"`
	DeliverDate            int64  `json:"deliverDate"`
	NextSubscriptionDate   int64  `json:"nextSubscriptionDate"`
}

// SetSimpleEarnLockedRedeemOptionService set locked redeem option.
// See https://binance-docs.github.io/apidocs/spot/en/#set-locked-product-redeem-option-user_data
type SetSimpleEarnLockedRedeemOptionService struct {
	c          *Client
	positionId string
	redeemTo   RedeemTo
}

// PositionId sets the positionId parameter.
func (s *SetSimpleEarnLockedRedeemOptionService) PositionId(positionId string) *SetSimpleEarnLockedRedeemOptionService {
	s.positionId = positionId
	return s
}

// RedeemTo sets the redeemTo parameter (SPOT, FLEXIBLE).
func (s *SetSimpleEarnLockedRedeemOptionService) RedeemTo(redeemTo RedeemTo) *SetSimpleEarnLockedRedeemOptionService {
	s.redeemTo = redeemTo
	return s
}

// Do sends the request.
func (s *SetSimpleEarnLockedRedeemOptionService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnSetRedeemOptionResponse, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/simple-earn/locked/setRedeemOption",
		secType:  secTypeSigned,
	}
	r.setParam("positionId", s.positionId)
	if s.redeemTo != "" {
		r.setParam("redeemTo", string(s.redeemTo))
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnSetRedeemOptionResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnSetRedeemOptionResponse represents a set redeem option response.
type SimpleEarnSetRedeemOptionResponse struct {
	Success bool `json:"success"`
}

// ListSimpleEarnFlexibleRewardsHistoryService list flexible rewards history (interest records).
// See https://binance-docs.github.io/apidocs/spot/en/#get-flexible-rewards-history-user_data
type ListSimpleEarnFlexibleRewardsHistoryService struct {
	c         *Client
	productId *string
	asset     *string
	rtype     SimpleEarnRewardsType
	startTime *int64
	endTime   *int64
	current   *int32
	size      *int32
}

// ProductId sets the productId parameter.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) ProductId(productId string) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.productId = &productId
	return s
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) Asset(asset string) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.asset = &asset
	return s
}

// Type sets the type parameter (BONUS, REALTIME, REWARDS, ALL; default ALL).
func (s *ListSimpleEarnFlexibleRewardsHistoryService) Type(rtype SimpleEarnRewardsType) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.rtype = rtype
	return s
}

// StartTime sets the startTime parameter.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) StartTime(startTime int64) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) EndTime(endTime int64) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) Current(current int32) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) Size(size int32) *ListSimpleEarnFlexibleRewardsHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnFlexibleRewardsHistoryService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnFlexibleRewardsHistoryResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/flexible/history/rewardsRecord",
		secType:  secTypeSigned,
	}
	if s.productId != nil {
		r.setParam("productId", *s.productId)
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.rtype != "" {
		r.setParam("type", string(s.rtype))
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
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnFlexibleRewardsHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnFlexibleRewardsHistory represents one flexible rewards (interest) record.
type SimpleEarnFlexibleRewardsHistory struct {
	Asset     string `json:"asset"`
	Rewards   string `json:"rewards"`
	ProjectID string `json:"projectId"`
	Type      string `json:"type"`
	Time      int64  `json:"time"`
}

// SimpleEarnFlexibleRewardsHistoryResponse represents the flexible rewards history.
type SimpleEarnFlexibleRewardsHistoryResponse struct {
	Rows  []SimpleEarnFlexibleRewardsHistory `json:"rows"`
	Total int                                `json:"total"`
}

// ListSimpleEarnLockedRewardsHistoryService list locked rewards history (interest records).
// See https://binance-docs.github.io/apidocs/spot/en/#get-locked-rewards-history-user_data
type ListSimpleEarnLockedRewardsHistoryService struct {
	c          *Client
	positionId *string
	asset      *string
	startTime  *int64
	endTime    *int64
	current    *int32
	size       *int32
}

// PositionId sets the positionId parameter.
func (s *ListSimpleEarnLockedRewardsHistoryService) PositionId(positionId string) *ListSimpleEarnLockedRewardsHistoryService {
	s.positionId = &positionId
	return s
}

// Asset sets the asset parameter.
func (s *ListSimpleEarnLockedRewardsHistoryService) Asset(asset string) *ListSimpleEarnLockedRewardsHistoryService {
	s.asset = &asset
	return s
}

// StartTime sets the startTime parameter.
func (s *ListSimpleEarnLockedRewardsHistoryService) StartTime(startTime int64) *ListSimpleEarnLockedRewardsHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime sets the endTime parameter.
func (s *ListSimpleEarnLockedRewardsHistoryService) EndTime(endTime int64) *ListSimpleEarnLockedRewardsHistoryService {
	s.endTime = &endTime
	return s
}

// Current sets the current parameter.
func (s *ListSimpleEarnLockedRewardsHistoryService) Current(current int32) *ListSimpleEarnLockedRewardsHistoryService {
	s.current = &current
	return s
}

// Size sets the size parameter.
func (s *ListSimpleEarnLockedRewardsHistoryService) Size(size int32) *ListSimpleEarnLockedRewardsHistoryService {
	s.size = &size
	return s
}

// Do sends the request.
func (s *ListSimpleEarnLockedRewardsHistoryService) Do(ctx context.Context, opts ...RequestOption) (res *SimpleEarnLockedRewardsHistoryResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/sapi/v1/simple-earn/locked/history/rewardsRecord",
		secType:  secTypeSigned,
	}
	if s.positionId != nil {
		r.setParam("positionId", *s.positionId)
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
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
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	res = new(SimpleEarnLockedRewardsHistoryResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	return res, nil
}

// SimpleEarnLockedRewardsHistory represents one locked rewards (interest) record.
type SimpleEarnLockedRewardsHistory struct {
	PositionID int64  `json:"positionId"`
	Time       int64  `json:"time"`
	Asset      string `json:"asset"`
	LockPeriod string `json:"lockPeriod"`
	Amount     string `json:"amount"`
}

// SimpleEarnLockedRewardsHistoryResponse represents the locked rewards history.
type SimpleEarnLockedRewardsHistoryResponse struct {
	Rows  []SimpleEarnLockedRewardsHistory `json:"rows"`
	Total int                              `json:"total"`
}
