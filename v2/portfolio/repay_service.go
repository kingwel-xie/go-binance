package portfolio

import (
	"context"
	"net/http"
)

// RepayService repay negative balance assets
type RepayService struct {
	c *Client
}

// Do send request
func (s *RepayService) Do(ctx context.Context, opts ...RequestOption) (ret bool, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/papi/v1/repay-futures-negative-balance",
		secType:  secTypeSigned,
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
