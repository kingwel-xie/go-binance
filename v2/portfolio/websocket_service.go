package portfolio

import (
	"encoding/json"
	"fmt"
	"time"
)

// Endpoints
const (
	baseWsMainUrl       = "wss://fstream.binance.com/pm/ws"
	baseCombinedMainURL = "wss://fstream.binance.com/stream?streams="
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = false
)

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag
func getWsEndpoint() string {
	return baseWsMainUrl
}

// WsUserDataEvent define user data event
type WsUserDataEvent struct {
	Event               UserDataEventType     `json:"e"`
	Time                int64                 `json:"E"`
	CrossWalletBalance  string                `json:"cw"`
	MarginCallPositions []WsPosition          `json:"p"`
	TransactionTime     int64                 `json:"T"`
	AccountUpdate       WsAccountUpdate       `json:"a"`
	OrderTradeUpdate    WsOrderTradeUpdate    `json:"o"`
	AccountConfigUpdate WsAccountConfigUpdate `json:"ac"`
}

// WsAccountUpdate define account update
type WsAccountUpdate struct {
	Reason    UserDataEventReasonType `json:"m"`
	Balances  []WsBalance             `json:"B"`
	Positions []WsPosition            `json:"P"`
}

// WsBalance define balance
type WsBalance struct {
	Asset              string `json:"a"`
	Balance            string `json:"wb"`
	CrossWalletBalance string `json:"cw"`
	ChangeBalance      string `json:"bc"`
}

// WsPosition define position
type WsPosition struct {
	Symbol                    string           `json:"s"`
	Side                      PositionSideType `json:"ps"`
	Amount                    string           `json:"pa"`
	MarginType                MarginType       `json:"mt"`
	IsolatedWallet            string           `json:"iw"`
	EntryPrice                string           `json:"ep"`
	MarkPrice                 string           `json:"mp"`
	UnrealizedPnL             string           `json:"up"`
	AccumulatedRealized       string           `json:"cr"`
	MaintenanceMarginRequired string           `json:"mm"`
}

// WsOrderTradeUpdate define order trade update
type WsOrderTradeUpdate struct {
	Symbol               string             `json:"s"`
	ClientOrderID        string             `json:"c"`
	Side                 SideType           `json:"S"`
	Type                 OrderType          `json:"o"`
	TimeInForce          TimeInForceType    `json:"f"`
	OriginalQty          string             `json:"q"`
	OriginalPrice        string             `json:"p"`
	AveragePrice         string             `json:"ap"`
	StopPrice            string             `json:"sp"`
	ExecutionType        OrderExecutionType `json:"x"`
	Status               OrderStatusType    `json:"X"`
	ID                   int64              `json:"i"`
	LastFilledQty        string             `json:"l"`
	AccumulatedFilledQty string             `json:"z"`
	LastFilledPrice      string             `json:"L"`
	CommissionAsset      string             `json:"N"`
	Commission           string             `json:"n"`
	TradeTime            int64              `json:"T"`
	TradeID              int64              `json:"t"`
	BidsNotional         string             `json:"b"`
	AsksNotional         string             `json:"a"`
	IsMaker              bool               `json:"m"`
	IsReduceOnly         bool               `json:"R"`
	WorkingType          WorkingType        `json:"wt"`
	OriginalType         OrderType          `json:"ot"`
	PositionSide         PositionSideType   `json:"ps"`
	IsClosingPosition    bool               `json:"cp"`
	ActivationPrice      string             `json:"AP"`
	CallbackRate         string             `json:"cr"`
	RealizedPnL          string             `json:"rp"`
}

// WsAccountConfigUpdate define account config update
type WsAccountConfigUpdate struct {
	Symbol   string `json:"s"`
	Leverage int64  `json:"l"`
}

// WsUserDataHandler handle WsUserDataEvent
type WsUserDataHandler func(event *WsUserDataEvent)

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		event := new(WsUserDataEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
