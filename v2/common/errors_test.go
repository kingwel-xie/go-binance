package common

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIErrorBinanceStyle(t *testing.T) {
	e := NewAPIError(400, []byte(`{"code":-1121,"msg":"Invalid symbol."}`))
	assert.Equal(t, int64(-1121), e.Code)
	assert.Equal(t, "Invalid symbol.", e.Message)
	assert.Equal(t, 400, e.StatusCode)
	assert.True(t, e.IsValid())
	assert.Equal(t, "<APIError> http=400, code=-1121, msg=Invalid symbol.", e.Error())
}

func TestNewAPIErrorGatewayStyle(t *testing.T) {
	body := []byte(`{"timestamp":1794992372163,"error":"Not Found","path":"/sapi/v2/portfolio/retail/collateralRate","status":404,"message":"No message available"}`)
	e := NewAPIError(http.StatusNotFound, body)
	assert.Equal(t, 404, e.StatusCode)
	assert.Equal(t, int64(0), e.Code)
	assert.Equal(t, "Not Found: No message available", e.Message)
	assert.True(t, e.IsValid())
	assert.Equal(t, "<APIError> http=404, msg=Not Found: No message available", e.Error())
}

func TestNewAPIErrorHTMLBody(t *testing.T) {
	body := []byte("<html><body>502 Bad Gateway</body></html>")
	e := NewAPIError(http.StatusBadGateway, body)
	assert.Equal(t, 502, e.StatusCode)
	assert.Equal(t, int64(0), e.Code)
	assert.Equal(t, "", e.Message)
	assert.False(t, e.IsValid())
	assert.Equal(t, "<APIError> http=502, rsp=<html><body>502 Bad Gateway</body></html>", e.Error())
}

func TestNewAPIErrorEmptyBody(t *testing.T) {
	e := NewAPIError(http.StatusInternalServerError, nil)
	assert.Equal(t, 500, e.StatusCode)
	assert.False(t, e.IsValid())
	assert.Equal(t, "<APIError> http=500, rsp=", e.Error())
}

func TestNewAPIErrorMessageOnlyBody(t *testing.T) {
	e := NewAPIError(403, []byte(`{"message":"Invalid API-key, IP, or permissions for action."}`))
	assert.Equal(t, 403, e.StatusCode)
	assert.Equal(t, "Invalid API-key, IP, or permissions for action.", e.Message)
	assert.True(t, e.IsValid())
}

func TestIsRateLimitError(t *testing.T) {
	assert.True(t, IsRateLimitError(NewAPIError(http.StatusTooManyRequests, []byte(`{"code":-1003,"msg":"Too many requests."}`))))
	assert.True(t, IsRateLimitError(NewAPIError(200, []byte(`{"code":-1003,"msg":"Too many requests."}`))))
	assert.True(t, IsRateLimitError(NewAPIError(200, []byte(`{"code":-1015,"msg":"Too many requests."}`))))
	assert.False(t, IsRateLimitError(NewAPIError(400, []byte(`{"code":-1121,"msg":"Invalid symbol."}`))))
	assert.False(t, IsRateLimitError(nil))
}

func TestIsIPBannedError(t *testing.T) {
	assert.True(t, IsIPBannedError(NewAPIError(http.StatusTeapot, nil)))
	assert.False(t, IsIPBannedError(NewAPIError(http.StatusTooManyRequests, nil)))
	assert.False(t, IsIPBannedError(nil))
}

func TestIsAPIError(t *testing.T) {
	assert.True(t, IsAPIError(NewAPIError(404, []byte(`{"error":"Not Found"}`))))
	assert.False(t, IsAPIError(assert.AnError))
}

func TestAPIErrorWithoutStatusCode(t *testing.T) {
	// Errors built by service level parsers carry no HTTP status, so the
	// http= prefix must be omitted rather than printed as http=0.
	e := APIError{Code: -2014, Message: "API-key format invalid."}
	assert.Equal(t, "<APIError> code=-2014, msg=API-key format invalid.", e.Error())
	assert.Equal(t, "<APIError> msg=oops", APIError{Message: "oops"}.Error())
}

func TestMaskSignature(t *testing.T) {
	assert.Equal(t,
		"https://api.binance.com/api/v3/order?symbol=BTCUSDT&signature=***",
		MaskSignature("https://api.binance.com/api/v3/order?symbol=BTCUSDT&signature=abcdef123456"))
	assert.Equal(t,
		"https://api.binance.com/api/v3/order?signature=***&recvWindow=5000",
		MaskSignature("https://api.binance.com/api/v3/order?signature=abcdef123456&recvWindow=5000"))
	assert.Equal(t,
		"https://api.binance.com/api/v3/time",
		MaskSignature("https://api.binance.com/api/v3/time"))
}
