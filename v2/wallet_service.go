package binance

import (
	"context"
	"net/http"
)

// GetFundingAssetService get funding asset
type GetFundingAssetService struct {
	c                *Client
	asset            *string
	needBtcValuation bool
}

func (s *GetFundingAssetService) Asset(asset string) *GetFundingAssetService {
	s.asset = &asset
	return s
}

func (s *GetFundingAssetService) NeedBtcValuation(val bool) *GetFundingAssetService {
	s.needBtcValuation = val
	return s
}

// Do send request
func (s *GetFundingAssetService) Do(ctx context.Context, opts ...RequestOption) (res []FundingAssetRecord, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/sapi/v1/asset/get-funding-asset",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.needBtcValuation {
		r.setParam("needBtcValuation", s.needBtcValuation)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FundingAssetRecord represents funding asset
type FundingAssetRecord struct {
	Asset        string `json:"asset"`
	Free         string `json:"free"`
	Locked       string `json:"locked"`
	Freeze       string `json:"freeze"`
	Withdrawing  string `json:"withdrawing"`
	BtcValuation string `json:"btcValuation"`
}
