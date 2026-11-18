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
	r.Equal("BNB", res[0].Asset)
	r.Len(res[0].CollateralInfo, 2)
	tier := res[0].CollateralInfo[0]
	r.Equal("0.0000", tier.TierFloor)
	r.Equal("50.0000", tier.TierCap)
	r.Equal("1.0000", tier.CollateralRate)
	r.Equal("0.0000", tier.Cum)
	r.Equal("0.9000", res[0].CollateralInfo[1].CollateralRate)
	r.Equal("BTC", res[1].Asset)
	r.Len(res[1].CollateralInfo, 1)
}

func (s *portfolioMarginServiceTestSuite) TestGetCollateralRate() {
	s.mockDo(s.collateralRateData(), nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetPortfolioMarginCollateralRateService().Do(newContext())
	s.r().NoError(err)
	s.assertCollateralRate(res)
}

func (s *portfolioMarginServiceTestSuite) TestGetTieredCollateralRate() {
	s.mockDo(s.collateralRateData(), nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetPortfolioMarginTieredCollateralRateService().Do(newContext())
	s.r().NoError(err)
	s.assertCollateralRate(res)
}
