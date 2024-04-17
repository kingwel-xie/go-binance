package portfolio

import (
	"context"
	"fmt"
	"net/http"
)

// AssetCollectionService collect all/specified assets
type AssetCollectionService struct {
	c     *Client
	asset *string
}

// Asset sets the asset parameter.
func (s *AssetCollectionService) Asset(asset string) *AssetCollectionService {
	s.asset = &asset
	return s
}

// Do send request
func (s *AssetCollectionService) Do(ctx context.Context, opts ...RequestOption) (ret bool, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/auto-collection",
		secType:  secTypeAPIKey,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
		r.endpoint = "/papi/v1/asset-collection"
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	j, err := newJSON(data)
	if err != nil {
		return
	}
	msg := j.Get("msg").MustString()
	if msg == "success" {
		ret = true
	}
	return
}

// BnbTransferService transfer BNB between PM and UM wallet
type BnbTransferService struct {
	c            *Client
	amount       string
	transferSide string
}

// Amount sets the amount parameter.
func (s *BnbTransferService) Amount(amount string) *BnbTransferService {
	s.amount = amount
	return s
}

// TransferSide sets the transferSide parameter. "TO_UM","FROM_UM"
func (s *BnbTransferService) TransferSide(transferSide string) *BnbTransferService {
	s.transferSide = transferSide
	return s
}

// Do send request
func (s *BnbTransferService) Do(ctx context.Context, opts ...RequestOption) (tranId string, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/auto-collection",
		secType:  secTypeAPIKey,
	}
	r.setParam("amount", s.amount)
	r.setParam("transferSide", s.transferSide)
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return
	}
	j, err := newJSON(data)
	if err != nil {
		return
	}
	tranId = fmt.Sprintf("%d", j.Get("tranId").MustInt64())
	return
}
