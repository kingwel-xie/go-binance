package futures

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type convertServiceTestSuite struct {
	baseTestSuite
}

func TestConvertService(t *testing.T) {
	suite.Run(t, new(convertServiceTestSuite))
}

func (s *convertServiceTestSuite) TestListConvertExchangeInfo() {
	data := []byte(`[
  {
    "fromAsset":"BTC",
    "toAsset":"USDT",
    "fromAssetMinAmount":"0.0004",
    "fromAssetMaxAmount":"50",
    "toAssetMinAmount":"20",
    "toAssetMaxAmount":"2500000"
  }
]`)
	s.mockDo(data, nil)
	defer s.assertDo()

	fromAsset := "BTC"
	toAsset := "USDT"
	s.assertReq(func(r *request) {
		e := newRequest().setParams(params{
			"fromAsset": fromAsset,
			"toAsset":   toAsset,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewListConvertExchangeInfoService().FromAsset(fromAsset).ToAsset(toAsset).Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 1)
	e := ConvertExchangeInfo{
		FromAsset:          fromAsset,
		ToAsset:            toAsset,
		FromAssetMinAmount: "0.0004",
		FromAssetMaxAmount: "50",
		ToAssetMinAmount:   "20",
		ToAssetMaxAmount:   "2500000",
	}
	s.r().Equal(e.FromAsset, res[0].FromAsset, "FromAsset")
	s.r().Equal(e.ToAsset, res[0].ToAsset, "ToAsset")
	s.r().Equal(e.FromAssetMinAmount, res[0].FromAssetMinAmount, "FromAssetMinAmount")
	s.r().Equal(e.FromAssetMaxAmount, res[0].FromAssetMaxAmount, "FromAssetMaxAmount")
	s.r().Equal(e.ToAssetMinAmount, res[0].ToAssetMinAmount, "ToAssetMinAmount")
	s.r().Equal(e.ToAssetMaxAmount, res[0].ToAssetMaxAmount, "ToAssetMaxAmount")
}

func (s *convertServiceTestSuite) TestCreateConvertQuote() {
	data := []byte(`{
   "quoteId":"12415572564",
   "ratio":"38163.7",
   "inverseRatio":"0.0000262",
   "validTimestamp":1623319461670,
   "toAmount":"3816.37",
   "fromAmount":"0.1"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	fromAsset := "BTC"
	toAsset := "USDT"
	fromAmount := "0.0004"
	toAmount := "20"
	validTime := ConvertValidTime10S
	s.assertReq(func(r *request) {
		e := newSignedRequest().setFormParams(params{
			"fromAsset":  fromAsset,
			"toAsset":    toAsset,
			"fromAmount": fromAmount,
			"toAmount":   toAmount,
			"validTime":  validTime,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewCreateConvertQuoteService().
		FromAsset(fromAsset).
		ToAsset(toAsset).
		FromAmount(fromAmount).
		ToAmount(toAmount).
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

func (s *convertServiceTestSuite) TestAcceptQuote() {
	data := []byte(`{
  "orderId":"933256278426274426",
  "createTime":1623381330472,
  "orderStatus":"PROCESS"
}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	quoteId := "12415572564"
	s.assertReq(func(r *request) {
		e := newSignedRequest().setFormParams(params{
			"quoteId": quoteId,
		})
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewConvertAcceptService().QuoteId(quoteId).Do(newContext())
	s.r().NoError(err)
	s.r().Equal("933256278426274426", res.OrderId)
	s.r().Equal(int64(1623381330472), res.CreateTime)
	s.r().Equal(ConvertAcceptStatusProcess, res.OrderStatus)
	s.r().Equal(ConvertAcceptStatus("FAIL"), ConvertAcceptStatusFailed)
}

func (s *convertServiceTestSuite) TestGetConvertStatus() {
	data := []byte(`{
  "orderId":933256278426274426,
  "orderStatus":"SUCCESS",
  "fromAsset":"BTC",
  "fromAmount":"0.00054414",
  "toAsset":"USDT",
  "toAmount":"20",
  "ratio":"36755",
  "inverseRatio":"0.00002721",
  "createTime":1623381330472
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

	res, err := s.client.NewGetConvertStatusService().
		OrderId(orderId).
		QuoteId(quoteId).
		Do(newContext())
	s.r().NoError(err)
	s.r().Equal(int64(933256278426274426), res.OrderId)
	s.r().Equal(ConvertAcceptStatusSuccess, res.OrderStatus)
	s.r().Equal("BTC", res.FromAsset)
	s.r().Equal("0.00054414", res.FromAmount)
	s.r().Equal("USDT", res.ToAsset)
	s.r().Equal("20", res.ToAmount)
	s.r().Equal("36755", res.Ratio)
	s.r().Equal("0.00002721", res.InverseRatio)
	s.r().Equal(int64(1623381330472), res.CreateTime)
}
