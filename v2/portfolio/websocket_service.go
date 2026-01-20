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
	Event                 UserDataEventType       `json:"e"`
	BusinessLine          string                  `json:"fs"` // UM or CM
	Time                  int64                   `json:"E"`
	TransactionTime       int64                   `json:"T"`
	AccountUpdate         WsAccountUpdate         `json:"a"`
	OrderTradeUpdate      WsOrderTradeUpdate      `json:"o"`
	AccountConfigUpdate   WsAccountConfigUpdate   `json:"ac"`
	MarginCondOrderUpdate WsMarginCondOrderUpdate `json:"so"`
	RiskLevelChange       WsRiskLevelChange
	MarginOpenLoss        WsMarginOpenLoss
	MarginLiabilityChange WsMarginLiabilityChange
	MarginAccountUpdate   WsMarginAccountUpdate
	MarginBalanceUpdate   WsMarginBalanceUpdate
	MarginOrderUpdate     WsMarginOrderUpdate
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

type WsRiskLevelChange struct {
	UniMMR       string `json:"u"`
	MarginEvent  string `json:"s"`
	MarginUsd    string `json:"eq"`
	ActualEquity string `json:"ae"`
	MaintMargin  string `json:"m"`
}

type WsMarginOpenLoss struct {
	O []struct {
		Asset  string `json:"a"`
		Amount string `json:"o"`
	} `json:"O"`
}

type WsMarginLiabilityChange struct {
	Asset          string `json:"a"`
	Type           string `json:"t"`
	TxId           int64  `json:"T"`
	Principal      string `json:"p"`
	Interest       string `json:"i"`
	TotalLiability string `json:"l"`
}

type WsMarginAccountUpdate struct {
	B []struct {
		Asset  string `json:"a"`
		Free   string `json:"f"`
		Locked string `json:"l"`
	} `json:"B"`
}

type WsMarginBalanceUpdate struct {
	Asset           string `json:"a"`
	Change          string `json:"d"`
	TransactionTime int64  `json:"T"`
}

type WsMarginOrderUpdate struct {
	Symbol                  string          `json:"s"`
	ClientOrderId           string          `json:"c"`
	Side                    string          `json:"S"`
	Type                    string          `json:"o"`
	TimeInForce             TimeInForceType `json:"f"`
	Volume                  string          `json:"q"`
	Price                   string          `json:"p"`
	StopPrice               string          `json:"P"`
	IceBergVolume           string          `json:"F"`
	OrderListId             int64           `json:"g"` // for OCO
	OrigCustomOrderId       string          `json:"C"` // customized order ID for the original order
	ExecutionType           string          `json:"x"` // execution type for this event NEW/TRADE...
	Status                  string          `json:"X"` // order status
	RejectReason            string          `json:"r"`
	Id                      int64           `json:"i"` // order id
	LatestVolume            string          `json:"l"` // quantity for the latest trade
	FilledVolume            string          `json:"z"`
	LatestPrice             string          `json:"L"` // price for the latest trade
	FeeAsset                string          `json:"N"`
	FeeCost                 string          `json:"n"`
	TransactionTime         int64           `json:"T"`
	TradeId                 int64           `json:"t"`
	IgnoreI                 int64           `json:"I"` // ignore
	IsInOrderBook           bool            `json:"w"` // is the order in the order book?
	IsMaker                 bool            `json:"m"` // is this order maker?
	IgnoreM                 bool            `json:"M"` // ignore
	CreateTime              int64           `json:"O"`
	FilledQuoteVolume       string          `json:"Z"` // the quote volume that already filled
	LatestQuoteVolume       string          `json:"Y"` // the quote volume for the latest trade
	QuoteVolume             string          `json:"Q"`
	SelfTradePreventionMode string          `json:"V"`

	//These are fields that appear in the payload only if certain conditions are met.
	TrailingDelta              int64  `json:"d"` // Appears only for trailing stop orders.
	TrailingTime               int64  `json:"D"`
	StrategyId                 int64  `json:"j"` // Appears only if the strategyId parameter was provided upon order placement.
	StrategyType               int64  `json:"J"` // Appears only if the strategyType parameter was provided upon order placement.
	PreventedMatchId           int64  `json:"v"` // Appears only for orders that expired due to STP.
	PreventedQuantity          string `json:"A"`
	LastPreventedQuantity      string `json:"B"`
	TradeGroupId               int64  `json:"u"`
	CounterOrderId             int64  `json:"U"`
	CounterSymbol              string `json:"Cs"`
	PreventedExecutionQuantity string `json:"pl"`
	PreventedExecutionPrice    string `json:"pL"`
	PreventedExecutionQuoteQty string `json:"pY"`
	WorkingTime                int64  `json:"W"` // Appears when the order is working on the book
	MatchType                  string `json:"b"`
	AllocationId               int64  `json:"a"`
	WorkingFloor               string `json:"k"`  // Appears for orders that could potentially have allocations
	UsedSor                    bool   `json:"uS"` // Appears for orders that used SOR
}

// WsMarginCondOrderUpdate define conditional order update
type WsMarginCondOrderUpdate struct {
	Symbol        string          `json:"s"`
	ClientOrderId string          `json:"c"`
	StrategyId    int             `json:"si"`
	Side          string          `json:"S"`
	StrategyType  string          `json:"st"`
	TimeInForce   TimeInForceType `json:"f"`
	Volume        string          `json:"q"`
	Price         string          `json:"p"`
	Sp            string          `json:"sp"`
	OrderStatus   string          `json:"os"`
	T             int64           `json:"T"`
	Ut            int64           `json:"ut"`
	R             bool            `json:"R"`
	Wt            string          `json:"wt"`
	Ps            string          `json:"ps"`
	Cp            bool            `json:"cp"`
	AP            string          `json:"AP"`
	Cr            string          `json:"cr"`
	I             int             `json:"i"`
	V             string          `json:"V"`
	Gtd           int             `json:"gtd"`
}

// WsUserDataHandler handle WsUserDataEvent
type WsUserDataHandler func(event *WsUserDataEvent)

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(listenKey string, handler WsUserDataHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	endpoint := fmt.Sprintf("%s/%s", getWsEndpoint(), listenKey)
	cfg := newWsConfig(endpoint)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			errHandler(err)
			return
		}
		event := new(WsUserDataEvent)
		event.Event = UserDataEventType(j.Get("e").MustString())
		event.BusinessLine = j.Get("fs").MustString()
		event.Time = j.Get("E").MustInt64()
		event.TransactionTime = j.Get("T").MustInt64()

		switch event.Event {
		case UserDataEventTypeAccountUpdate:
			fallthrough
		case UserDataEventTypeOrderTradeUpdate:
			fallthrough
		case UserDataEventTypeCondOrderUpdate:
			fallthrough
		case UserDataEventTypeAccountConfigUpdate:
			err = json.Unmarshal(message, &event)
			if err != nil {
				errHandler(err)
				return
			}

		case UserDataEventTypeRiskLevelChange:
			err = json.Unmarshal(message, &event.RiskLevelChange)
			if err != nil {
				errHandler(err)
				return
			}
		case UserDataEventTypeMarginLiabilityChange:
			err = json.Unmarshal(message, &event.MarginLiabilityChange)
			if err != nil {
				errHandler(err)
				return
			}

		case UserDataEventTypeMarginBalanceUpdate:
			err = json.Unmarshal(message, &event.MarginBalanceUpdate)
			if err != nil {
				errHandler(err)
				return
			}
		case UserDataEventTypeMarginAccountUpdate:
			err = json.Unmarshal(message, &event.MarginAccountUpdate)
			if err != nil {
				errHandler(err)
				return
			}
		case UserDataEventTypeMarginOpenLoss:
			err = json.Unmarshal(message, &event.MarginOpenLoss)
			if err != nil {
				errHandler(err)
				return
			}
		case UserDataEventTypeMarginOrderUpdate:
			err = json.Unmarshal(message, &event.MarginOrderUpdate)
			if err != nil {
				errHandler(err)
				return
			}
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
