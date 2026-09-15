package options

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MarginAccountServiceTestSuite struct {
	baseTestSuite
}

func TestMarginAccountService(t *testing.T) {
	suite.Run(t, new(MarginAccountServiceTestSuite))
}

func (s *MarginAccountServiceTestSuite) TestMarginAccount() {
	data := []byte(`{
		"asset": [
		 {
		  "asset": "USDT",
		  "marginBalance": "10099.448",
		  "equity": "10094.44662",
		  "available": "8725.92524",
		  "initialMargin": "1084.52138",
		  "maintMargin": "151.00138",
		  "unrealizedPNL": "-5.00138",
		  "adjustedEquity": "34.13282285"
		 }
		],
		"greek":[{
			"underlying": "BTCUSDT",
			"delta": "-0.05",
			"gamma": "-0.002",
			"theta": "-0.05",
			"vega": "-0.002"
		}],
		"time": 1762843368098,
		"canTrade": true,
		"canDeposit": true,
		"canWithdraw": true,
		"reduceOnly": false,
		"tradeGroupId": -1
	   }`)
	s.mockDo(data, nil)
	defer s.assertDo()

	marginAccount, err := s.client.NewMarginAccountService().Do(newContext())
	targetMarginAccount := &MarginAccount{
		Asset: []*MarginAccountAsset{
			{
				Asset:          "USDT",
				MarginBalance:  "10099.448",
				Equity:         "10094.44662",
				Available:      "8725.92524",
				InitialMargin:  "1084.52138",
				MaintMargin:    "151.00138",
				UnrealizedPNL:  "-5.00138",
				AdjustedEquity: "34.13282285",
			},
		},
		Greek: []*Greek{
			{
				Underlying: "BTCUSDT",
				Delta:      "-0.05",
				Gamma:      "-0.002",
				Theta:      "-0.05",
				Vega:       "-0.002",
			},
		},
		Time:         1762843368098,
		CanTrade:     true,
		CanDeposit:   true,
		CanWithdraw:  true,
		ReduceOnly:   false,
		TradeGroupID: -1,
	}

	s.r().Equal(nil, err, "err != nil")

	r := s.r()
	for i := range marginAccount.Asset {
		r.Equal(marginAccount.Asset[i].Asset, targetMarginAccount.Asset[i].Asset, "MarginAccountAsset.Asset")
		r.Equal(marginAccount.Asset[i].MarginBalance, targetMarginAccount.Asset[i].MarginBalance, "MarginAccountAsset.MarginBalance")
		r.Equal(marginAccount.Asset[i].Equity, targetMarginAccount.Asset[i].Equity, "MarginAccountAsset.Equity")
		r.Equal(marginAccount.Asset[i].Available, targetMarginAccount.Asset[i].Available, "MarginAccountAsset.Available")
		r.Equal(marginAccount.Asset[i].InitialMargin, targetMarginAccount.Asset[i].InitialMargin, "MarginAccountAsset.InitialMargin")
		r.Equal(marginAccount.Asset[i].MaintMargin, targetMarginAccount.Asset[i].MaintMargin, "MarginAccountAsset.MaintMargin")
		r.Equal(marginAccount.Asset[i].UnrealizedPNL, targetMarginAccount.Asset[i].UnrealizedPNL, "MarginAccountAsset.UnrealizedPNL")
		r.Equal(marginAccount.Asset[i].AdjustedEquity, targetMarginAccount.Asset[i].AdjustedEquity, "MarginAccountAsset.AdjustedEquity")
	}
	for i := range marginAccount.Greek {
		r.Equal(marginAccount.Greek[i].Underlying, targetMarginAccount.Greek[i].Underlying, "Greek.Underlying")
		r.Equal(marginAccount.Greek[i].Delta, targetMarginAccount.Greek[i].Delta, "Greek.Delta")
		r.Equal(marginAccount.Greek[i].Gamma, targetMarginAccount.Greek[i].Gamma, "Greek.Gamma")
		r.Equal(marginAccount.Greek[i].Theta, targetMarginAccount.Greek[i].Theta, "Greek.Theta")
		r.Equal(marginAccount.Greek[i].Vega, targetMarginAccount.Greek[i].Vega, "Greek.Vega")
	}
	r.Equal(marginAccount.Time, targetMarginAccount.Time, "Time")
	r.Equal(marginAccount.CanTrade, targetMarginAccount.CanTrade, "CanTrade")
	r.Equal(marginAccount.CanDeposit, targetMarginAccount.CanDeposit, "CanDeposit")
	r.Equal(marginAccount.CanWithdraw, targetMarginAccount.CanWithdraw, "CanWithdraw")
	r.Equal(marginAccount.ReduceOnly, targetMarginAccount.ReduceOnly, "ReduceOnly")
	r.Equal(marginAccount.TradeGroupID, targetMarginAccount.TradeGroupID, "TradeGroupID")
}
