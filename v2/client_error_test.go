package binance

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adshao/go-binance/v2/common"
)

type clientErrorHandlingTestSuite struct {
	baseTestSuite
}

func TestClientErrorHandling(t *testing.T) {
	suite.Run(t, new(clientErrorHandlingTestSuite))
}

func (s *clientErrorHandlingTestSuite) TestGatewayNotFound() {
	s.mockDo([]byte(`{"timestamp":1794992372163,"error":"Not Found","path":"/sapi/v2/portfolio/retail/collateralRate","status":404,"message":"No message available"}`), nil, 404)
	defer s.assertDo()

	err := s.client.NewPingService().Do(newContext())
	r := s.r()
	r.Error(err)
	r.True(common.IsAPIError(err))
	apiErr, ok := err.(*common.APIError)
	r.True(ok)
	r.Equal(404, apiErr.StatusCode)
	r.Equal("Not Found: No message available", apiErr.Message)
	require.False(s.T(), common.IsRateLimitError(err))
}

func (s *clientErrorHandlingTestSuite) TestRateLimitError() {
	s.mockDo([]byte(`{"code":-1003,"msg":"Too many requests."}`), nil, 429)
	defer s.assertDo()

	err := s.client.NewPingService().Do(newContext())
	r := s.r()
	r.Error(err)
	apiErr, ok := err.(*common.APIError)
	r.True(ok)
	r.Equal(429, apiErr.StatusCode)
	r.Equal(int64(-1003), apiErr.Code)
	r.Equal("Too many requests.", apiErr.Message)
	r.True(common.IsRateLimitError(err))
	r.False(common.IsIPBannedError(err))
}

func (s *clientErrorHandlingTestSuite) TestIPBanError() {
	s.mockDo([]byte("banned"), nil, 418)
	defer s.assertDo()

	err := s.client.NewPingService().Do(newContext())
	r := s.r()
	r.Error(err)
	apiErr, ok := err.(*common.APIError)
	r.True(ok)
	r.Equal(418, apiErr.StatusCode)
	r.False(apiErr.IsValid())
	r.Equal("<APIError> http=418, rsp=banned", apiErr.Error())
	r.True(common.IsIPBannedError(err))
}
