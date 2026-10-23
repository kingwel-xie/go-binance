package binance

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type convertTradeTestSuite struct {
	baseTestSuite
}

func TestConvertTradeService(t *testing.T) {
	suite.Run(t, new(convertTradeTestSuite))
}

func (s *convertTradeTestSuite) TestConvertTradeHistory() {
	data := []byte(`{
	   "list": [
			{
				"quoteId": "f3b91c525b2644c7bc1e1cd31b6e1aa6",
				"orderId": 940708407462087195,
				"orderStatus": "SUCCESS",
				"fromAsset": "USDT",
				"fromAmount": "20",
				"toAsset": "BNB",
				"toAmount": "0.06154036",
				"ratio": "0.00307702",
				"inverseRatio": "324.99",
				"createTime": 1624248872184
			}
	   ],
		"startTime": 1623824139000,
		"endTime": 1626416139000,
		"limit": 100,
		"moreData": false
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	startTime := time.Now().AddDate(0, 0, -7).Unix() * 1000
	endTime := time.Now().Unix() * 1000
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"startTime": startTime,
			"endTime":   endTime,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewConvertTradeHistoryService().
		StartTime(startTime).
		EndTime(endTime).
		Do(newContext())
	s.r().NoError(err)
	e := &ConvertTradeHistory{
		List: []ConvertTradeHistoryItem{
			{
				QuoteId:      "f3b91c525b2644c7bc1e1cd31b6e1aa6",
				OrderId:      940708407462087195,
				OrderStatus:  "SUCCESS",
				FromAsset:    "USDT",
				FromAmount:   "20",
				ToAsset:      "BNB",
				ToAmount:     "0.06154036",
				Ratio:        "0.00307702",
				InverseRatio: "324.99",
				CreateTime:   1624248872184,
			},
		},
		StartTime: 1623824139000,
		EndTime:   1626416139000,
		Limit:     100,
		MoreData:  false,
	}
	s.assertConvertTradeHistoryEqual(e, res)
}

func (s *convertTradeTestSuite) TestConvertExchangeInfo() {
	data := []byte(`[
		{
			"fromAsset": "BTC",
			"toAsset": "USDT",
			"fromAssetMinAmount": "0.0004",
			"fromAssetMaxAmount": "50",
			"toAssetMinAmount": "20",
			"toAssetMaxAmount": "9E+24"
		}
	]`)
	s.mockDo(data, nil)
	defer s.assertDo()

	fromAsset := "BTC"
	toAsset := "USDT"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"fromAsset": fromAsset,
			"toAsset":   toAsset,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewConvertExchangeInfoService().
		FromAsset(fromAsset).
		ToAsset(toAsset).
		Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 1)
	s.r().Equal(fromAsset, res[0].FromAsset)
	s.r().Equal(toAsset, res[0].ToAsset)
	s.r().Equal("0.0004", res[0].FromAssetMinAmount)
	s.r().Equal("50", res[0].FromAssetMaxAmount)
	s.r().Equal("20", res[0].ToAssetMinAmount)
	s.r().Equal("9E+24", res[0].ToAssetMaxAmount)
}

func (s *convertTradeTestSuite) TestConvertAssetInfo() {
	data := []byte(`[
		{
			"asset": "BTC",
			"fraction": 8
		}
	]`)
	s.mockDo(data, nil)
	defer s.assertDo()

	s.assertReq(func(r *request) {
		s.assertRequestEqual(newSignedRequest(), r)
	})

	res, err := s.client.NewConvertAssetInfoService().Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 1)
	s.r().Equal("BTC", res[0].Asset)
	s.r().Equal(8, res[0].Fraction)
}

func (s *convertTradeTestSuite) TestConvertGetQuote() {
	data := []byte(`{
		"quoteId": "12415572564",
		"ratio": "38163.7",
		"inverseRatio": "0.0000262",
		"validTimestamp": 1623319461670,
		"toAmount": "3816.37",
		"fromAmount": "0.1"
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	fromAsset := "BTC"
	toAsset := "USDT"
	fromAmount := "0.1"
	toAmount := "3816.37"
	walletType := "SPOT"
	validTime := "10s"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"fromAsset":  fromAsset,
			"toAsset":    toAsset,
			"fromAmount": fromAmount,
			"toAmount":   toAmount,
			"walletType": walletType,
			"validTime":  validTime,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewConvertQuoteService().
		FromAsset(fromAsset).
		ToAsset(toAsset).
		FromAmount(fromAmount).
		ToAmount(toAmount).
		WalletType(walletType).
		ValidTime(validTime).
		Do(newContext())
	s.r().NoError(err)
	s.r().Equal("12415572564", res.QuoteId)
	s.r().Equal("38163.7", res.Ratio)
	s.r().Equal("0.0000262", res.InverseRatio)
	s.r().Equal(int64(1623319461670), res.ValidTimestamp)
	s.r().Equal("3816.37", res.ToAmount)
	s.r().Equal("0.1", res.FromAmount)
}

func (s *convertTradeTestSuite) TestConvertAcceptQuote() {
	data := []byte(`{
		"orderId": "933256278426274426",
		"createTime": 1623381330472,
		"orderStatus": "PROCESS"
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	quoteId := "12415572564"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"quoteId": quoteId,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewConvertAcceptQuoteService().QuoteId(quoteId).Do(newContext())
	s.r().NoError(err)
	s.r().Equal("933256278426274426", res.OrderId)
	s.r().Equal(int64(1623381330472), res.CreateTime)
	s.r().Equal("PROCESS", res.OrderStatus)
}

func (s *convertTradeTestSuite) TestConvertOrderStatus() {
	data := []byte(`{
		"orderId": 933256278426274426,
		"orderStatus": "SUCCESS",
		"fromAsset": "BTC",
		"fromAmount": "0.00054414",
		"toAsset": "USDT",
		"toAmount": "20",
		"ratio": "36755",
		"inverseRatio": "0.00002721",
		"createTime": 1623381330472
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	orderId := "933256278426274426"
	quoteId := "12415572564"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setParams(params{
			"orderId": orderId,
			"quoteId": quoteId,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewConvertOrderStatusService().
		OrderId(orderId).
		QuoteId(quoteId).
		Do(newContext())
	s.r().NoError(err)
	s.r().Equal(int64(933256278426274426), res.OrderId)
	s.r().Equal("SUCCESS", res.OrderStatus)
	s.r().Equal("BTC", res.FromAsset)
	s.r().Equal("0.00054414", res.FromAmount)
	s.r().Equal("USDT", res.ToAsset)
	s.r().Equal("20", res.ToAmount)
	s.r().Equal("36755", res.Ratio)
	s.r().Equal("0.00002721", res.InverseRatio)
	s.r().Equal(int64(1623381330472), res.CreateTime)
}

func (s *convertTradeTestSuite) assertConvertTradeHistoryEqual(e, a *ConvertTradeHistory) {
	r := s.r()

	r.Len(a.List, len(e.List))
	for i := 0; i < len(a.List); i++ {
		s.assertConvertTradeHistoryItemEqual(&e.List[i], &a.List[i])
	}

	r.Equal(e.StartTime, a.StartTime, "StartTime")
	r.Equal(e.EndTime, a.EndTime, "EndTime")
	r.Equal(e.Limit, a.Limit, "Limit")
	r.Equal(e.MoreData, a.MoreData, "MoreData")
}

func (s *convertTradeTestSuite) assertConvertTradeHistoryItemEqual(e, a *ConvertTradeHistoryItem) {
	r := s.r()
	r.Equal(e.QuoteId, a.QuoteId, "QuoteId")
	r.Equal(e.OrderId, a.OrderId, "OrderId")
	r.Equal(e.OrderStatus, a.OrderStatus, "OrderStatus")
	r.Equal(e.FromAsset, a.FromAsset, "FromAsset")
	r.Equal(e.FromAmount, a.FromAmount, "FromAmount")
	r.Equal(e.ToAsset, a.ToAsset, "ToAsset")
	r.Equal(e.ToAmount, a.ToAmount, "ToAmount")
	r.Equal(e.Ratio, a.Ratio, "Ratio")
	r.Equal(e.InverseRatio, a.InverseRatio, "InverseRatio")
	r.Equal(e.CreateTime, a.CreateTime, "CreateTime")
}
