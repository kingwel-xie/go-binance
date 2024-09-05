package futures

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type WsClientState int

const (
	WsInit         WsClientState = 0
	WsConnecting   WsClientState = 1
	WsConnected    WsClientState = 2
	WsAdminClosing WsClientState = 3
)

// Endpoints
var (
	WsAPIMainURL    = "wss://ws-fapi.binance.com/ws-fapi/v1"
	WsAPITestnetURL = "wss://testnet.binancefuture.com/ws-fapi/v1"
)

type _ResponseMap struct {
	lock sync.Mutex
	d    map[string]chan *WsApiResponse
}

func (m _ResponseMap) LoadAndDelete(id string) chan *WsApiResponse {
	m.lock.Lock()
	defer m.lock.Unlock()
	if a := m.d[id]; a != nil {
		delete(m.d, id)
		return a
	}
	return nil
}

func (m _ResponseMap) Set(id string, ch chan *WsApiResponse) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.d[id] = ch
}

var apiResponses = _ResponseMap{d: make(map[string]chan *WsApiResponse)}

type WsApiResponse struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Error  struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
	Result     ObjectType `json:"result"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
		Count         int    `json:"count"`
	} `json:"rateLimits"`
}

// ObjectType 自定义类型 for unmarshall
type ObjectType string

// UnmarshalJSON 实现 UnmarshalJSON 方法
func (o *ObjectType) UnmarshalJSON(data []byte) error {
	// 将 JSON 对象解析为字符串
	*o = ObjectType(data) // 直接将原始数据赋值为字符串
	return nil
}

// getWsAPIEndpoint return the base endpoint of the WebSocket API according the UseTestnet flag
func getWsAPIEndpoint() string {
	if UseTestnet {
		return WsAPITestnetURL
	}
	return WsAPIMainURL
}

func makeConn() (*websocket.Conn, chan struct{}, chan struct{}) {
	Dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true, // important for huge size message
	}

	c, _, err := Dialer.Dial(getWsAPIEndpoint(), nil)
	if err != nil {
		return nil, nil, nil
	}
	c.SetReadLimit(655350)
	doneC := make(chan struct{})
	stopC := make(chan struct{})
	disconnectedC := make(chan struct{})
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		if WebsocketKeepalive {
			keepAlive(c, WebsocketTimeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		adminForced := false
		go func() {
			select {
			case <-stopC:
				adminForced = true
			case <-doneC:
				close(disconnectedC)
			}
			_ = c.Close()

		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !adminForced {
					fmt.Println("ws error:", err)
				}
				return
			}
			res := new(WsApiResponse)
			err = json.Unmarshal(message, res)
			if err != nil {
				//fmt.Println("unmarshal error:", err)
				return
			}
			if a := apiResponses.LoadAndDelete(res.Id); a != nil {
				a <- res
				close(a)
			}
		}
	}()

	return c, stopC, disconnectedC
}

func (c *Client) handleDisconnected(ch chan struct{}) {
	go func() {
		select {
		case <-ch:
		}
		// if it is triggered by AdminClose, just ignore
		if c.wsState == WsAdminClosing {
			return
		}
		c.wsState = WsConnecting
		c.debug("disconnected, try reconnecting later...")

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			// 使用 select 语句同时等待超时和上下文完成
			select {
			//case <-sm.ctx.Done():
			//	// 上下文已完成，可能是超时或取消
			//	log.Info("%s stream, context terminated...", sm.name)
			//	return
			case <-ticker.C:
				conn, stopC, disconnectedC := makeConn()
				if conn != nil {
					c.Conn = conn
					c.StopC = stopC
					c.wsState = WsConnected
					// well done, break the loop
					c.handleDisconnected(disconnectedC)

					c.debug("reconnected with %s", c.BaseURL)
					return
				}
				c.debug("failed to connect to %s, retrying later...", c.BaseURL)
			}
		}
	}()
}

func (c *Client) Close() {
	if c.wsState == WsConnected {
		c.wsState = WsAdminClosing
		close(c.StopC)
	}
}

// Encode encodes the values into “URL encoded” form
// ("bar=baz&foo=quux") sorted by key.
func (v params) Encode() string {
	if len(v) == 0 {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(fmt.Sprintf("%v", vs))
	}
	return buf.String()
}

func (c *Client) parseWsRequest(r *request) (err error) {
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.wsParams == nil {
		r.wsParams = params{}
	}
	// collect params from query & form, construct wsParams
	for k, v := range r.query {
		r.wsParams[k] = v[0]
	}
	for k, v := range r.form {
		r.wsParams[k] = v[0]
	}

	if r.secType == secTypeSigned {
		r.wsParams[timestampKey] = currentTimestamp() - c.TimeOffset
		r.wsParams[apiKey] = c.APIKey
		raw := r.wsParams.Encode()

		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		r.wsParams[signatureKey] = fmt.Sprintf("%x", (mac.Sum(nil)))
	}

	c.debug("ws-method: %s, params: %v", r.wsMethod, r.wsParams)

	return nil
}

func (c *Client) callWsAPI(ctx context.Context, r *request) ([]byte, error) {
	err := c.parseWsRequest(r)
	if err != nil {
		return []byte{}, err
	}

	// allocate channel, size 1
	id, ch := uuid.NewString(), make(chan *WsApiResponse, 1)

	req := map[string]interface{}{
		"id":     id,
		"method": r.wsMethod,
	}
	if len(r.wsParams) > 0 {
		req["params"] = r.wsParams
	}

	c.debug("request: %#v", req)

	apiResponses.Set(id, ch)
	err = c.Conn.WriteJSON(req)
	//f := c.do
	//if f == nil {
	//	f = c.HTTPClient.Do
	//}
	//res, err := f(req)
	if err != nil {
		return []byte{}, err
	}

	// timeout context
	ctx2, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return []byte{}, ctx.Err()

	case <-ctx2.Done():
		return []byte{}, ctx2.Err()

	case res := <-ch:
		c.debug("response status code: %d", res.Status)
		c.debug("response raw: %s", string(res.Result))
		c.debug("response: %#v", res.Error)

		if res.Status >= http.StatusBadRequest {
			apiErr := new(common.APIError)
			apiErr.Code = res.Error.Code
			apiErr.Message = res.Error.Msg
			return nil, apiErr
		}
		return []byte(res.Result), nil
	}
}
