package portfolio

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type accountServiceTestSuite struct {
	baseTestSuite
}

func TestAccountService(t *testing.T) {
	suite.Run(t, new(accountServiceTestSuite))
}

func (s *accountServiceTestSuite) TestGetBalance() {
	data := []byte(`[
		{
			"asset": "USDT",
			"totalWalletBalance": "8.74947592",
			"crossMarginAsset": "23.72469206",
			"crossMarginBorrowed": "0.00000000",
			"crossMarginFree": "10.00000000",
			"crossMarginInterest": "0.00000000",
			"crossMarginLocked": "5.00000000",
			"umWalletBalance": "8.74947592",
			"umUnrealizedPNL": "0.10000000",
			"cmWalletBalance": "0.00000000",
			"cmUnrealizedPNL": "0.00000000",
			"updateTime": 1625474304765,
			"negativeBalance": "0.00000000"
		}
	]`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetBalanceService().Do(newContext())
	s.r().NoError(err)
	s.r().Len(res, 1)
	e := &Balance{
		Asset:               "USDT",
		TotalWalletBalance:  "8.74947592",
		CrossMarginAsset:    "23.72469206",
		CrossMarginBorrowed: "0.00000000",
		CrossMarginFree:     "10.00000000",
		CrossMarginInterest: "0.00000000",
		CrossMarginLocked:   "5.00000000",
		UmWalletBalance:     "8.74947592",
		UmUnrealizedPNL:     "0.10000000",
		CmWalletBalance:     "0.00000000",
		CmUnrealizedPNL:     "0.00000000",
		UpdateTime:          1625474304765,
		NegativeBalance:     "0.00000000",
	}
	s.assertBalanceEqual(e, res[0])
}

func (s *accountServiceTestSuite) assertBalanceEqual(e, a *Balance) {
	r := s.r()
	r.Equal(e.Asset, a.Asset, "Asset")
	r.Equal(e.TotalWalletBalance, a.TotalWalletBalance, "TotalWalletBalance")
	r.Equal(e.CrossMarginAsset, a.CrossMarginAsset, "CrossMarginAsset")
	r.Equal(e.CrossMarginBorrowed, a.CrossMarginBorrowed, "CrossMarginBorrowed")
	r.Equal(e.CrossMarginFree, a.CrossMarginFree, "CrossMarginFree")
	r.Equal(e.CrossMarginInterest, a.CrossMarginInterest, "CrossMarginInterest")
	r.Equal(e.CrossMarginLocked, a.CrossMarginLocked, "CrossMarginLocked")
	r.Equal(e.UmWalletBalance, a.UmWalletBalance, "UmWalletBalance")
	r.Equal(e.UmUnrealizedPNL, a.UmUnrealizedPNL, "UmUnrealizedPNL")
	r.Equal(e.CmWalletBalance, a.CmWalletBalance, "CmWalletBalance")
	r.Equal(e.CmUnrealizedPNL, a.CmUnrealizedPNL, "CmUnrealizedPNL")
	r.Equal(e.UpdateTime, a.UpdateTime, "UpdateTime")
	r.Equal(e.NegativeBalance, a.NegativeBalance, "NegativeBalance")
}

func (s *accountServiceTestSuite) TestGetAccount() {
	data := []byte(`{
		"uniMMR": "1.5",
		"accountEquity": "122607.35137903",
		"actualEquity": "122607.35137903",
		"accountInitialMargin": "0.00000000",
		"accountMaintMargin": "0.00000000",
		"accountStatus": "NORMAL",
		"virtualMaxWithdrawAmount": "100000.00000000",
		"totalAvailableBalance": "100000.00000000",
		"totalMarginOpenLoss": "0.00000000",
		"updateTime": 1625474304765
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()
	s.assertReq(func(r *request) {
		e := newSignedRequest()
		s.assertRequestEqual(e, r)
	})

	res, err := s.client.NewGetAccountService().Do(newContext())
	s.r().NoError(err)
	e := &Account{
		UniMMR:                   "1.5",
		AccountEquity:            "122607.35137903",
		ActualEquity:             "122607.35137903",
		AccountInitialMargin:     "0.00000000",
		AccountMaintMargin:       "0.00000000",
		AccountStatus:            "NORMAL",
		VirtualMaxWithdrawAmount: "100000.00000000",
		TotalAvailableBalance:    "100000.00000000",
		TotalMarginOpenLoss:      "0.00000000",
		UpdateTime:               1625474304765,
	}
	s.assertAccountEqual(e, res)
}

func (s *accountServiceTestSuite) assertAccountEqual(e, a *Account) {
	r := s.r()
	r.Equal(e.UniMMR, a.UniMMR, "UniMMR")
	r.Equal(e.AccountEquity, a.AccountEquity, "AccountEquity")
	r.Equal(e.ActualEquity, a.ActualEquity, "ActualEquity")
	r.Equal(e.AccountInitialMargin, a.AccountInitialMargin, "AccountInitialMargin")
	r.Equal(e.AccountMaintMargin, a.AccountMaintMargin, "AccountMaintMargin")
	r.Equal(e.AccountStatus, a.AccountStatus, "AccountStatus")
	r.Equal(e.VirtualMaxWithdrawAmount, a.VirtualMaxWithdrawAmount, "VirtualMaxWithdrawAmount")
	r.Equal(e.TotalAvailableBalance, a.TotalAvailableBalance, "TotalAvailableBalance")
	r.Equal(e.TotalMarginOpenLoss, a.TotalMarginOpenLoss, "TotalMarginOpenLoss")
	r.Equal(e.UpdateTime, a.UpdateTime, "UpdateTime")
}