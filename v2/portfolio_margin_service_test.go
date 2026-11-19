package binance

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type portfolioMarginServiceTestSuite struct {
	baseTestSuite
}

func TestPortfolioMarginService(t *testing.T) {
	suite.Run(t, new(portfolioMarginServiceTestSuite))
}

func (s *portfolioMarginServiceTestSuite) collateralRateData() []byte {
	return []byte(`[
    {
        "asset": "USDC",
        "collateralRate": "1.0000"
    },
    {
        "asset": "BTC",
        "collateralRate": "0.9000"
    }
]`)
}

func (s *portfolioMarginServiceTestSuite) tieredCollateralRateData() []byte {
	return []byte(`[
    {
        "asset": "BNB",
        "collateralInfo": [
            {
                "tierFloor": "0.0000",
                "tierCap": "50.0000",
                "collateralRate": "1.0000",
                "cum": "0.0000"
            },
            {
                "tierFloor": "50.0000",
                "tierCap": "500.0000",
                "collateralRate": "0.9000",
                "cum": "5.0000"
            }
        ]
    },
    {
        "asset": "BTC",
        "collateralInfo": [
            {
                "tierFloor": "0.0000",
                "tierCap": "1000.0000",
                "collateralRate": "1.0000",
                "cum": "0.0000"
            }
        ]
    }
]`)
}

func (s *portfolioMarginServiceTestSuite) assertCollateralRate(res []*PortfolioMarginCollateralRate) {
	r := s.r()
	r.Len(res, 2)
	r.Equal("USDC", res[0].Asset)
	r.Equal("1.0000", res[0].CollateralRate)
	r.Equal("BTC", res[1].Asset)
	r.Equal("0.9000", res[1].CollateralRate)
}

func (s *portfolioMarginServiceTestSuite) assertTieredCollateralRate(res []*PortfolioMarginTieredCollateralRate) {
	r := s.r()
	r.Len(res, 2)
	r.Equal("BNB", res[0].Asset)
	r.Len(res[0].CollateralInfo, 2)
	tier := res[0].CollateralInfo[0]
	r.Equal("0.0000", tier.TierFloor)
	r.Equal("50.0000", tier.TierCap)
	r.Equal("1.0000", tier.CollateralRate)
	r.Equal("0.0000", tier.Cum)
	r.Equal("0.9000", res[0].CollateralInfo[1].CollateralRate)
	r.Equal("5.0000", res[0].CollateralInfo[1].Cum)
	r.Equal("BTC", res[1].Asset)
	r.Len(res[1].CollateralInfo, 1)
}

func (s *portfolioMarginServiceTestSuite) TestGetCollateralRate() {
	s.mockDo(s.collateralRateData(), nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetPortfolioMarginCollateralRateService().Do(newContext())
	s.r().NoError(err)
	s.assertCollateralRate(res)
}

func (s *portfolioMarginServiceTestSuite) TestGetTieredCollateralRate() {
	s.mockDo(s.tieredCollateralRateData(), nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetPortfolioMarginTieredCollateralRateService().Do(newContext())
	s.r().NoError(err)
	s.assertTieredCollateralRate(res)
}

func (s *portfolioMarginServiceTestSuite) TestGetRetailTieredCollateralRate() {
	s.mockDo(s.tieredCollateralRateData(), nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetPortfolioMarginRetailTieredCollateralRateService().Do(newContext())
	s.r().NoError(err)
	s.assertTieredCollateralRate(res)
}

func (s *portfolioMarginServiceTestSuite) TestGetAssetLeverage() {
	data := []byte(`[
    {
        "asset": "USDC",
        "leverage": 10
    },
    {
        "asset": "BTC",
        "leverage": 10
    }
]`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetPortfolioMarginAssetLeverageService().Do(newContext())
	r := s.r()
	r.NoError(err)
	r.Len(res, 2)
	r.Equal("USDC", res[0].Asset)
	r.Equal(int64(10), res[0].Leverage)
	r.Equal("BTC", res[1].Asset)
	r.Equal(int64(10), res[1].Leverage)
}
