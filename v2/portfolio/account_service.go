package portfolio

import (
	"context"
	"encoding/json"
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

// Account define account info
type Account struct {
	UniMMR                   string `json:"uniMMR"`
	AccountEquity            string `json:"accountEquity"` // 以USD计价的账户权益
	ActualEquity             string `json:"actualEquity"`  // 不考虑质押率后的以USD计价账户权益
	AccountInitialMargin     string `json:"accountInitialMargin"`
	AccountMaintMargin       string `json:"accountMaintMargin"`       // 以USD计价统一账户维持保证金
	AccountStatus            string `json:"accountStatus"`            // 统一账户账户状态："NORMAL", "MARGIN_CALL", "SUPPLY_MARGIN", "REDUCE_ONLY", "ACTIVE_LIQUIDATION", "FORCE_LIQUIDATION", "BANKRUPTED"
	VirtualMaxWithdrawAmount string `json:"virtualMaxWithdrawAmount"` // 以USD计价的最大可转出
	TotalAvailableBalance    string `json:"totalAvailableBalance"`
	TotalMarginOpenLoss      string `json:"totalMarginOpenLoss"`
	UpdateTime               int64  `json:"updateTime"`
}
