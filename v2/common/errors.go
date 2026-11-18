package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIError define API error when response status is 4xx or 5xx
type APIError struct {
	StatusCode int    `json:"-"` // HTTP status code of the response
	Code       int64  `json:"code"`
	Message    string `json:"msg"`
	Response   []byte `json:"-"` // Assign the body value when the Code and Message fields are invalid.
}

// Error return error code and message
func (e APIError) Error() string {
	// The HTTP status code is the most reliable discriminator when the body
	// carries no Binance error code (gateway errors, HTML/empty bodies), so it
	// is printed first whenever it is known.
	status := ""
	if e.StatusCode != 0 {
		status = fmt.Sprintf("http=%d, ", e.StatusCode)
	}
	switch {
	case e.Code != 0:
		return fmt.Sprintf("<APIError> %scode=%d, msg=%s", status, e.Code, e.Message)
	case e.Message != "":
		return fmt.Sprintf("<APIError> %smsg=%s", status, e.Message)
	}
	return fmt.Sprintf("<APIError> %srsp=%s", status, string(e.Response))
}

// IsValid check if the APIError carries a concrete API code and message
func (e APIError) IsValid() bool {
	return e.Code != 0 || e.Message != ""
}

// IsAPIError check if e is an API error
func IsAPIError(e error) bool {
	_, ok := e.(*APIError)
	return ok
}

// IsRateLimitError check if e is a rate limit error (HTTP 429 or Binance error codes)
func IsRateLimitError(e error) bool {
	apiErr, ok := e.(*APIError)
	if !ok {
		return false
	}
	if apiErr.StatusCode == http.StatusTooManyRequests {
		return true
	}
	switch apiErr.Code {
	case -1003, // TOO_MANY_REQUESTS
		-1015: // TOO_MANY_REQUESTS
		return true
	}
	return false
}

// IsIPBannedError check if e is an IP ban error (HTTP 418)
func IsIPBannedError(e error) bool {
	apiErr, ok := e.(*APIError)
	if !ok {
		return false
	}
	return apiErr.StatusCode == http.StatusTeapot
}

// NewAPIError builds an APIError from a non-2xx HTTP response body.
//
// Parsing order:
//  1. Binance style {"code":..,"msg":..} -> Code/Message.
//  2. Gateway style {"status":..,"error":..,"message":..} (e.g. Spring 404 pages,
//     rate limiting proxies) -> Message.
//  3. Anything else (HTML pages, empty or unknown bodies) -> kept in Response.
//
// The HTTP status code is always recorded in StatusCode.
func NewAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{StatusCode: statusCode, Response: body}
	if len(body) == 0 {
		return apiErr
	}
	ce := new(struct {
		Code    *int64 `json:"code"`
		Message string `json:"msg"`
	})
	if err := json.Unmarshal(body, ce); err == nil && ce.Code != nil {
		apiErr.Code = *ce.Code
		apiErr.Message = ce.Message
		return apiErr
	}
	gw := new(struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	})
	if err := json.Unmarshal(body, gw); err == nil && (gw.Error != "" || gw.Message != "") {
		parts := make([]string, 0, 2)
		for _, p := range []string{gw.Error, gw.Message} {
			if p != "" {
				parts = append(parts, p)
			}
		}
		if msg := strings.Join(parts, ": "); msg != "" {
			apiErr.Message = msg
			apiErr.Response = nil
			return apiErr
		}
	}
	return apiErr
}

// MaskSignature masks the signature parameter of a signed URL, so that the
// signature value never leaks into debug logs.
func MaskSignature(rawURL string) string {
	const key = "signature="
	idx := strings.Index(rawURL, key)
	if idx < 0 {
		return rawURL
	}
	end := idx + len(key)
	for end < len(rawURL) && rawURL[end] != '&' {
		end++
	}
	return rawURL[:idx+len(key)] + "***" + rawURL[end:]
}
