package portfolio

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetInterestHistoryService get position interest history service
type GetInterestHistoryService struct {
	c         *Client
	asset     *string
	startTime *int64
	endTime   *int64
	size      *int32
}

// Asset set asset
func (s *GetInterestHistoryService) Asset(asset string) *GetInterestHistoryService {
	s.asset = &asset
	return s
}

// StartTime set startTime
func (s *GetInterestHistoryService) StartTime(startTime int64) *GetInterestHistoryService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *GetInterestHistoryService) EndTime(endTime int64) *GetInterestHistoryService {
	s.endTime = &endTime
	return s
}

// Size sets the size parameter.
func (s *GetInterestHistoryService) Size(size int32) *GetInterestHistoryService {
	s.size = &size
	return s
}

// Do send request
func (s *GetInterestHistoryService) Do(ctx context.Context, opts ...RequestOption) (res []*InterestHistory, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/portfolio/interest-history",
		secType:  secTypeSigned,
	}
	if s.asset != nil {
		r.setParam("asset", *s.asset)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.size != nil {
		r.setParam("size", *s.size)
	}

	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = make([]*InterestHistory, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// InterestHistory define position margin history info
type InterestHistory struct {
	Asset               string `json:"asset"`
	Interest            string `json:"interest"`
	InterestAccuredTime int64  `json:"interestAccruedTime"`
	InterestRate        string `json:"interestRate"`
	Principal           string `json:"principal"`
}
