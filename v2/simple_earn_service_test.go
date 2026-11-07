package binance

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type simpleEarnServiceTestSuite struct {
	baseTestSuite
}

func TestSimpleEarnService(t *testing.T) {
	suite.Run(t, new(simpleEarnServiceTestSuite))
}

func (s *simpleEarnServiceTestSuite) TestGetAccount() {
	data := []byte(`{
    "totalAmountInBTC": "0.01067982",
    "totalAmountInUSDT": "77.13289230",
    "totalFlexibleAmountInBTC": "0.00000000",
    "totalFlexibleAmountInUSDT": "0.00000000",
    "totalLockedInBTC": "0.01067982",
    "totalLockedInUSDT": "77.13289230"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewSimpleEarnAccountService().Do(newContext())
	s.r().NoError(err)
	e := &SimpleEarnAccountResponse{
		TotalAmountInBTC:          "0.01067982",
		TotalAmountInUSDT:         "77.13289230",
		TotalFlexibleAmountInBTC:  "0.00000000",
		TotalFlexibleAmountInUSDT: "0.00000000",
		TotalLockedInBTC:          "0.01067982",
		TotalLockedInUSDT:         "77.13289230",
	}
	s.Equal(e, res)
}

func (s *simpleEarnServiceTestSuite) TestListFlexibleProducts() {
	data := []byte(`{
    "rows": [
        {
            "asset": "BTC",
            "latestAnnualPercentageRate": "0.05000000",
            "tierAnnualPercentageRate": {
                "0-5BTC": 0.05,
                "5-10BTC": 0.03
            },
            "airDropPercentageRate": "0.05000000",
            "canPurchase": true,
            "canRedeem": true,
            "isSoldOut": true,
            "hot": true,
            "minPurchaseAmount": "0.01000000",
            "productId": "BTC001",
            "subscriptionStartTime": 1646182276000,
            "status": "PURCHASING"
        }
    ],
    "total": 1
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("asset", "BTC").setParam("current", 1).setParam("size", 10)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListSimpleEarnFlexibleService().Asset("BTC").Current(1).Size(10).Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 1)
	s.r().Equal(1, res.Total)
	row := res.Rows[0]
	s.r().Equal("BTC", row.Asset)
	s.r().Equal("0.05000000", row.LatestAnnualPercentageRate)
	s.r().Equal(TierAnnualPercentageRate{"0-5BTC": "0.05", "5-10BTC": "0.03"}, row.TierAnnualPercentageRate)
	s.r().Equal("0.05000000", row.AirDropPercentageRate)
	s.r().True(row.CanPurchase)
	s.r().True(row.CanRedeem)
	s.r().True(row.IsSoldOut)
	s.r().True(row.Hot)
	s.r().Equal("0.01000000", row.MinPurchaseAmount)
	s.r().Equal("BTC001", row.ProductID)
	s.r().Equal(int64(1646182276000), row.SubscriptionStartTime)
	s.r().Equal("PURCHASING", row.Status)
}

func (s *simpleEarnServiceTestSuite) TestListLockedProducts() {
	data := []byte(`{
    "rows": [
        {
            "projectId": "Bnb*120",
            "detail": {
                "asset": "BNB",
                "rewardAsset": "BNB",
                "duration": 120,
                "renewable": false,
                "isSoldOut": false,
                "apr": "0.0235",
                "status": "PURCHASING",
                "subscriptionStartTime": 1646182276000,
                "extraRewardAsset": "",
                "extraRewardAPR": ""
            },
            "quota": {
                "totalPersonalQuota": "10",
                "minimum": "0.00001"
            }
        }
    ],
    "total": 1
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("asset", "BNB").setParam("current", 1).setParam("size", 10)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListSimpleEarnLockedService().Asset("BNB").Current(1).Size(10).Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 1)
	s.r().Equal(1, res.Total)
	row := res.Rows[0]
	s.r().Equal("Bnb*120", row.ProjectID)
	s.r().Equal("BNB", row.Detail.Asset)
	s.r().Equal(120, row.Detail.Duration)
	s.r().False(row.Detail.Renewable)
	s.r().Equal("0.0235", row.Detail.Apr)
	s.r().Equal("PURCHASING", row.Detail.Status)
	s.r().Equal("10", row.Quota.TotalPersonalQuota)
	s.r().Equal("0.00001", row.Quota.Minimum)
}

func (s *simpleEarnServiceTestSuite) TestSubscribeFlexibleProduct() {
	data := []byte(`{
    "purchaseId": 40607,
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("productId", "BTC001").
			setParam("amount", "0.1").
			setParam("autoSubscribe", true).
			setParam("sourceAccount", "SPOT")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewSubscribeSimpleEarnFlexibleService().
		ProductId("BTC001").
		Amount("0.1").
		AutoSubscribe(true).
		SourceAccount(SimpleEarnSourceAccountSpot).
		Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SubscribeSimpleEarnFlexibleResponse{PurchaseID: 40607, Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestSubscribeLockedProduct() {
	data := []byte(`{
    "purchaseId": 40607,
    "positionId": 1234567,
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("projectId", "Bnb*120").
			setParam("amount", "0.1").
			setParam("autoSubscribe", false).
			setParam("redeemTo", "FLEXIBLE")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewSubscribeSimpleEarnLockedService().
		ProjectId("Bnb*120").
		Amount("0.1").
		AutoSubscribe(false).
		RedeemTo(RedeemToFlexible).
		Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SubscribeSimpleEarnLockedResponse{PurchaseID: 40607, PositionID: 1234567, Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestRedeemFlexibleProduct() {
	data := []byte(`{
    "redeemId": 40607,
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("productId", "BTC001").
			setParam("redeemAll", false).
			setParam("amount", "0.1").
			setParam("destAccount", "SPOT")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewRedeemSimpleEarnFlexibleService().
		ProductId("BTC001").
		RedeemAll(false).
		Amount(newStringPointer("0.1")).
		DestAccount("SPOT").
		Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&RedeemSimpleEarnFlexibleResponse{RedeemID: 40607, Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestRedeemFlexibleProductRedeemAll() {
	data := []byte(`{
    "redeemId": 40607,
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("productId", "BTC001").setParam("redeemAll", true)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewRedeemSimpleEarnFlexibleService().ProductId("BTC001").Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&RedeemSimpleEarnFlexibleResponse{RedeemID: 40607, Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestRedeemLockedProduct() {
	data := []byte(`{
    "redeemId": 40607,
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("positionId", "1234567")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewRedeemSimpleEarnLockedService().PositionId("1234567").Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&RedeemSimpleEarnLockedResponse{RedeemID: 40607, Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestGetFlexiblePosition() {
	data := []byte(`{
    "rows": [
        {
            "totalAmount": "75.46000000",
            "tierAnnualPercentageRate": {
                "0-5BTC": 0.05,
                "5-10BTC": 0.03
            },
            "latestAnnualPercentageRate": "0.02599895",
            "yesterdayAirdropPercentageRate": "0.02599895",
            "asset": "USDT",
            "airDropAsset": "BETH",
            "canRedeem": true,
            "collateralAmount": "232.23123213",
            "productId": "USDT001",
            "yesterdayRealTimeRewards": "0.10293829",
            "cumulativeBonusRewards": "0.22759183",
            "cumulativeRealTimeRewards": "0.22759183",
            "cumulativeTotalRewards": "0.45459183",
            "autoSubscribe": true
        }
    ],
    "total": 1
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("asset", "USDT").setParam("productId", "USDT001")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetSimpleEarnFlexiblePositionService().Asset("USDT").ProductId("USDT001").Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 1)
	s.r().Equal(1, res.Total)
	row := res.Rows[0]
	s.r().Equal("75.46000000", row.TotalAmount)
	s.r().Equal(TierAnnualPercentageRate{"0-5BTC": "0.05", "5-10BTC": "0.03"}, row.TierAnnualPercentageRate)
	s.r().Equal("USDT", row.Asset)
	s.r().Equal("BETH", row.AirDropAsset)
	s.r().True(row.CanRedeem)
	s.r().True(row.AutoSubscribe)
	s.r().Equal("0.45459183", row.CumulativeTotalRewards)
}

func (s *simpleEarnServiceTestSuite) TestGetLockedPosition() {
	data := []byte(`{
    "rows": [
        {
            "positionId": 347608038,
            "parentPositionId": 0,
            "projectId": "Bnb*120",
            "asset": "BNB",
            "amount": "0.05",
            "purchaseTime": 1783439066000,
            "duration": 120,
            "accrualDays": 0,
            "rewardAsset": "BNB",
            "APY": "0.004",
            "rewardAmt": "0",
            "extraRewardAsset": "",
            "extraRewardAPR": "",
            "estExtraRewardAmt": "",
            "nextPay": "0.00000054",
            "nextPayDate": 1783555200000,
            "payPeriod": 1,
            "redeemAmountEarly": "0.05",
            "rewardsEndDate": 1793836800000,
            "deliverDate": 1793959200000,
            "redeemPeriod": 1,
            "redeemingAmt": "",
            "redeemTo": "SPOT",
            "partialAmtDeliverDate": 0,
            "canRedeemEarly": true,
            "canFastRedemption": true,
            "autoSubscribe": false,
            "type": "AUTO",
            "status": "HOLDING",
            "canReStake": true
        }
    ],
    "total": 1
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("asset", "BNB").setParam("positionId", 347608038).setParam("projectId", "Bnb*120")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetSimpleEarnLockedPositionService().
		Asset("BNB").
		PositionId(347608038).
		ProjectId("Bnb*120").
		Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 1)
	s.r().Equal(1, res.Total)
	row := res.Rows[0]
	s.r().Equal(int64(347608038), row.PositionID)
	s.r().Equal("Bnb*120", row.ProjectID)
	s.r().Equal("BNB", row.Asset)
	s.r().Equal("0.05", row.Amount)
	s.r().Equal(120, row.Duration)
	s.r().Equal("0.004", row.Apy)
	s.r().True(row.CanRedeemEarly)
	s.r().True(row.CanFastRedemption)
	s.r().False(row.AutoSubscribe)
	s.r().Equal("HOLDING", row.Status)
	s.r().True(row.CanReStake)
	s.r().Equal("SPOT", row.RedeemTo)
}

func (s *simpleEarnServiceTestSuite) TestListFlexibleRateHistory() {
	data := []byte(`{
    "rows": [
        {
            "productId": "BUSD001",
            "asset": "BUSD",
            "annualPercentageRate": "0.00006408",
            "time": 1577233578000
        }
    ],
    "total": "1"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("productId", "BUSD001").
			setParam("aprPeriod", "DAY").
			setParam("startTime", 1575018814000).
			setParam("endTime", 1577257222000)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListSimpleEarnFlexibleRateHistoryService().
		ProductId("BUSD001").
		AprPeriod("DAY").
		StartTime(1575018814000).
		EndTime(1577257222000).
		Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 1)
	s.r().Equal("1", res.Total.String())
	s.r().Equal("BUSD001", res.Rows[0].ProductID)
	s.r().Equal("0.00006408", res.Rows[0].AnnualPercentageRate)
}

func (s *simpleEarnServiceTestSuite) TestGetFlexibleQuota() {
	data := []byte(`{
    "leftPersonalQuota": "0.0224"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("productId", "BTC001")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetSimpleEarnFlexibleQuotaService().ProductId("BTC001").Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SimpleEarnPersonalLeftQuotaResponse{LeftPersonalQuota: "0.0224"}, res)
}

func (s *simpleEarnServiceTestSuite) TestGetLockedQuota() {
	data := []byte(`{
    "leftPersonalQuota": "0.0224"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("projectId", "Bnb*120")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetSimpleEarnLockedQuotaService().ProjectId("Bnb*120").Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SimpleEarnPersonalLeftQuotaResponse{LeftPersonalQuota: "0.0224"}, res)
}

func (s *simpleEarnServiceTestSuite) TestSetFlexibleAutoSubscribe() {
	data := []byte(`{
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("productId", "BTC001").setParam("autoSubscribe", true)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewSetSimpleEarnFlexibleAutoSubscribeService().ProductId("BTC001").AutoSubscribe(true).Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SimpleEarnSetAutoSubscribeResponse{Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestSetLockedAutoSubscribe() {
	data := []byte(`{
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("positionId", "1234567").setParam("autoSubscribe", false)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewSetSimpleEarnLockedAutoSubscribeService().PositionId("1234567").AutoSubscribe(false).Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SimpleEarnSetAutoSubscribeResponse{Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestGetFlexibleSubscriptionPreview() {
	data := []byte(`{
    "totalAmount": "75.46000000",
    "rewardAsset": "USDT",
    "airDropAsset": "BETH",
    "estDailyBonusRewards": "0.10293829",
    "estDailyRealTimeRewards": "0.22759183",
    "estDailyAirdropRewards": "0.22759183"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("productId", "USDT001").setParam("amount", "1")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetSimpleEarnFlexibleSubscriptionPreviewService().ProductId("USDT001").Amount("1").Do(newContext())
	s.r().NoError(err)
	e := &SimpleEarnFlexibleSubscriptionPreviewResponse{
		TotalAmount:             "75.46000000",
		RewardAsset:             "USDT",
		AirDropAsset:            "BETH",
		EstDailyBonusRewards:    "0.10293829",
		EstDailyRealTimeRewards: "0.22759183",
		EstDailyAirdropRewards:  "0.22759183",
	}
	s.Equal(e, res)
}

func (s *simpleEarnServiceTestSuite) TestGetLockedSubscriptionPreview() {
	data := []byte(`[
    {
        "rewardAsset": "BTC",
        "totalRewardAmt": "0.125",
        "extraRewardAsset": "BNB",
        "estTotalExtraRewardAmt": "0.0012",
        "nextPay": "0.000625",
        "nextPayDate": 1642565111000,
        "valueDate": 1642565811000,
        "rewardsEndDate": 1645157511000,
        "deliverDate": 1645157611000,
        "nextSubscriptionDate": 1645157711000
    }
]`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("projectId", "Bnb*120").
			setParam("amount", "1").
			setParam("autoSubscribe", true)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetSimpleEarnLockedSubscriptionPreviewService().
		ProjectId("Bnb*120").
		Amount("1").
		AutoSubscribe(true).
		Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 1)
	row := res[0]
	s.r().Equal("BTC", row.RewardAsset)
	s.r().Equal("0.125", row.TotalRewardAmt)
	s.r().Equal("BNB", row.ExtraRewardAsset)
	s.r().Equal("0.000625", row.NextPay)
	s.r().Equal(int64(1642565111000), row.NextPayDate)
	s.r().Equal(int64(1645157611000), row.DeliverDate)
}

func (s *simpleEarnServiceTestSuite) TestSetLockedRedeemOption() {
	data := []byte(`{
    "success": true
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParam("positionId", "1234567").setParam("redeemTo", "FLEXIBLE")
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewSetSimpleEarnLockedRedeemOptionService().PositionId("1234567").RedeemTo(RedeemToFlexible).Do(newContext())
	s.r().NoError(err)
	s.r().Equal(&SimpleEarnSetRedeemOptionResponse{Success: true}, res)
}

func (s *simpleEarnServiceTestSuite) TestListFlexibleRewardsHistory() {
	data := []byte(`{
    "rows": [
        {
            "asset": "BUSD",
            "rewards": "0.00006408",
            "projectId": "USDT001",
            "type": "BONUS",
            "time": 1577233578000
        },
        {
            "asset": "USDT",
            "rewards": "0.00687654",
            "projectId": "USDT001",
            "type": "REALTIME",
            "time": 1577233562000
        }
    ],
    "total": 2
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("productId", "USDT001").
			setParam("asset", "USDT").
			setParam("type", "REALTIME").
			setParam("startTime", 1575018814000).
			setParam("endTime", 1577257222000).
			setParam("current", 1).
			setParam("size", 10)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListSimpleEarnFlexibleRewardsHistoryService().
		ProductId("USDT001").
		Asset("USDT").
		Type(SimpleEarnRewardsTypeRealtime).
		StartTime(1575018814000).
		EndTime(1577257222000).
		Current(1).
		Size(10).
		Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 2)
	s.r().Equal(2, res.Total)
	s.r().Equal("0.00006408", res.Rows[0].Rewards)
	s.r().Equal("BONUS", res.Rows[0].Type)
	s.r().Equal(int64(1577233562000), res.Rows[1].Time)
}

func (s *simpleEarnServiceTestSuite) TestListLockedRewardsHistory() {
	data := []byte(`{
    "rows": [
        {
            "positionId": 123123,
            "time": 1577257222000,
            "asset": "BNB",
            "lockPeriod": "30",
            "amount": "21312.23223"
        }
    ],
    "total": 1
}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest().
			setParam("positionId", "123123").
			setParam("asset", "BNB").
			setParam("startTime", 1575018814000).
			setParam("endTime", 1577257222000).
			setParam("current", 1).
			setParam("size", 10)
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListSimpleEarnLockedRewardsHistoryService().
		PositionId("123123").
		Asset("BNB").
		StartTime(1575018814000).
		EndTime(1577257222000).
		Current(1).
		Size(10).
		Do(newContext())
	s.r().NoError(err)
	s.r().Len(res.Rows, 1)
	s.r().Equal(1, res.Total)
	row := res.Rows[0]
	s.r().Equal(int64(123123), row.PositionID)
	s.r().Equal("BNB", row.Asset)
	s.r().Equal("30", row.LockPeriod)
	s.r().Equal("21312.23223", row.Amount)
}

func newStringPointer(v string) *string {
	return &v
}
