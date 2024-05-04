package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetBalanceService get account balance
type GetBalanceService struct {
	c *Client
}

// Do send request
func (s *GetBalanceService) Do(ctx context.Context, opts ...RequestOption) (res []*Balance, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/balance",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Balance{}, err
	}
	res = make([]*Balance, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Balance{}, err
	}
	return res, nil
}

// Balance define user balance of your account
type Balance struct {
	Asset               string `json:"asset"`
	TotalWalletBalance  string `json:"totalWalletBalance"`  // 钱包余额 =  全仓杠杆未锁定 + 全仓杠杆锁定 + u本位合约钱包余额 + 币本位合约钱包余额
	CrossMarginAsset    string `json:"crossMarginAsset"`    // 全仓资产 = 全仓杠杆未锁定 + 全仓杠杆锁定
	CrossMarginBorrowed string `json:"crossMarginBorrowed"` // 全仓杠杆借贷
	CrossMarginFree     string `json:"crossMarginFree"`     // 全仓杠杆未锁定
	CrossMarginInterest string `json:"crossMarginInterest"` // 全仓杠杆利息
	CrossMarginLocked   string `json:"crossMarginLocked"`   //全仓杠杆锁定
	UmWalletBalance     string `json:"umWalletBalance"`     // u本位合约钱包余额
	UmUnrealizedPNL     string `json:"umUnrealizedPNL"`     // u本位未实现盈亏
	CmWalletBalance     string `json:"cmWalletBalance"`     // 币本位合约钱包余额
	CmUnrealizedPNL     string `json:"cmUnrealizedPNL"`     // 币本位未实现盈亏
	UpdateTime          int64  `json:"updateTime"`
	NegativeBalance     string `json:"negativeBalance"`
}

// GetAccountService get account info
type GetAccountService struct {
	c *Client
}

// Do send request
func (s *GetAccountService) Do(ctx context.Context, opts ...RequestOption) (res *Account, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/account",
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Account)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Account struct {
	UniMMR                   string `json:"uniMMR"`        // 统一账户维持保证金率
	AccountEquity            string `json:"accountEquity"` // 以USD计价的账户权益
	ActualEquity             string `json:"actualEquity"`  // 不考虑质押率后的以USD计价账户权益
	AccountInitialMargin     string `json:"accountInitialMargin"`
	AccountMaintMargin       string `json:"accountMaintMargin"`
	AccountStatus            string `json:"accountStatus"`            // 统一账户账户状态："NORMAL", "MARGIN_CALL", "SUPPLY_MARGIN", "REDUCE_ONLY", "ACTIVE_LIQUIDATION", "FORCE_LIQUIDATION", "BANKRUPTED"
	VirtualMaxWithdrawAmount string `json:"virtualMaxWithdrawAmount"` // 以USD计价的最大可转出
	TotalAvailableBalance    string `json:"totalAvailableBalance"`
	TotalMarginOpenLoss      string `json:"totalMarginOpenLoss"`
	UpdateTime               int64  `json:"updateTime"`
}

// GetAccountExtService get account extend info
type GetAccountExtService struct {
	c     *Client
	which string // 'um' or 'cm'
}

// Which set which product
func (s *GetAccountExtService) Which(which string) *GetAccountExtService {
	s.which = which
	return s
}

// Do send request
func (s *GetAccountExtService) Do(ctx context.Context, opts ...RequestOption) (res *AccountExt, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}

	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/account", s.which),
		secType:  secTypeSigned,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(AccountExt)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type AccountExt struct {
	Assets []struct {
		Asset                  string `json:"asset"`
		CrossWalletBalance     string `json:"crossWalletBalance"`
		CrossUnPnl             string `json:"crossUnPnl"`
		MaintMargin            string `json:"maintMargin"`
		InitialMargin          string `json:"initialMargin"`
		PositionInitialMargin  string `json:"positionInitialMargin"`
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		UpdateTime             int64  `json:"updateTime"`
	} `json:"assets"`
	Positions []struct {
		Symbol                 string `json:"symbol"`
		InitialMargin          string `json:"initialMargin"`
		MaintMargin            string `json:"maintMargin"`
		UnrealizedProfit       string `json:"unrealizedProfit"`
		PositionInitialMargin  string `json:"positionInitialMargin"`
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		Leverage               string `json:"leverage"`
		EntryPrice             string `json:"entryPrice"`
		MaxNotional            string `json:"maxNotional"`
		MaxQty                 string `json:"maxQty"` // for CM
		BidNotional            string `json:"bidNotional"`
		AskNotional            string `json:"askNotional"`
		PositionSide           string `json:"positionSide"`
		PositionAmt            string `json:"positionAmt"`
		UpdateTime             int64  `json:"updateTime"`
		Notional               string `json:"notional"`
		BreakEvenPrice         string `json:"breakEvenPrice"`
	} `json:"positions"`
}

//type T struct {
//	Assets []struct {
//		Asset                  string `json:"asset"`
//		CrossWalletBalance     string `json:"crossWalletBalance"`
//		CrossUnPnl             string `json:"crossUnPnl"`
//		MaintMargin            string `json:"maintMargin"`
//		InitialMargin          string `json:"initialMargin"`
//		PositionInitialMargin  string `json:"positionInitialMargin"`
//		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
//		UpdateTime             int64  `json:"updateTime"`
//	} `json:"assets"`
//	Positions []struct {
//		Symbol                 string `json:"symbol"`
//		InitialMargin          string `json:"initialMargin"`
//		MaintMargin            string `json:"maintMargin"`
//		UnrealizedProfit       string `json:"unrealizedProfit"`
//		PositionInitialMargin  string `json:"positionInitialMargin"`
//		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
//		Leverage               string `json:"leverage"`
//		EntryPrice             string `json:"entryPrice"`
//		PositionSide           string `json:"positionSide"`
//		PositionAmt            string `json:"positionAmt"`
//		MaxQty                 string `json:"maxQty"`
//		UpdateTime             int64  `json:"updateTime"`
//		NotionalValue          string `json:"notionalValue"`
//		BreakEvenPrice         string `json:"breakEvenPrice"`
//	} `json:"positions"`
//}
