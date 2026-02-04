package binance

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

// StartUserStreamService create listen key for user stream service
type StartUserStreamService struct {
	c *Client
}

// Do send request
func (s *StartUserStreamService) Do(ctx context.Context, opts ...RequestOption) (listenKey string, err error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/api/v3/userDataStream",
		secType:  secTypeAPIKey,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return "", err
	}
	j, err := newJSON(data)
	if err != nil {
		return "", err
	}
	listenKey = j.Get("listenKey").MustString()
	return listenKey, nil
}

// KeepaliveUserStreamService update listen key
type KeepaliveUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *KeepaliveUserStreamService) ListenKey(listenKey string) *KeepaliveUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *KeepaliveUserStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodPut,
		endpoint: "/api/v3/userDataStream",
		secType:  secTypeAPIKey,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// CloseUserStreamService delete listen key
type CloseUserStreamService struct {
	c         *Client
	listenKey string
}

// ListenKey set listen key
func (s *CloseUserStreamService) ListenKey(listenKey string) *CloseUserStreamService {
	s.listenKey = listenKey
	return s
}

// Do send request
func (s *CloseUserStreamService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/api/v3/userDataStream",
		secType:  secTypeAPIKey,
	}
	r.setFormParam("listenKey", s.listenKey)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}

func NewDataStreamClient(oc *Client, handler WsUserDataHandler, errHandler ErrHandler) (*Client, error) {
	c := makeConn(handler, errHandler)
	if c == nil {
		return nil, fmt.Errorf("error to establish websocket connnetion")
	}

	client := &Client{
		APIKey:     oc.APIKey,
		SecretKey:  oc.SecretKey,
		BaseURL:    getAPIEndpoint(),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
		WsURL:      getWsAPIEndpoint(),
		WsConn:     c,
		wsState:    WsConnected,
	}
	client.handleDisconnected(c.Done, handler, errHandler)

	r := &request{
		secType:  secTypeSigned,
		wsMethod: "userDataStream.subscribe.signature",
	}
	data, _, err := client.callWsAPI(context.TODO(), r)
	if err != nil {
		return nil, err
	}
	res := new(_SubscriptionResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	//client.Logger.Printf("userdata stream %d\n", res.SubscriptionId)

	return client, nil
}

type _SubscriptionResponse struct {
	SubscriptionId int `json:"subscriptionId"`
}
