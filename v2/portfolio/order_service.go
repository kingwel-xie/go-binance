package portfolio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// CreateOrderService create order
type CreateOrderService struct {
	c                *Client
	which            string // 'um' or 'cm'
	symbol           string
	side             SideType
	positionSide     *PositionSideType
	orderType        OrderType
	timeInForce      *TimeInForceType
	quantity         string
	reduceOnly       *bool
	price            *string
	newClientOrderID *string
	//stopPrice        *string
	//workingType      *WorkingType
	//activationPrice  *string
	//callbackRate     *string
	//priceProtect     *bool
	newOrderRespType NewOrderRespType
	//closePosition    *bool
}

// Which set which product
func (s *CreateOrderService) Which(which string) *CreateOrderService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateOrderService) Side(side SideType) *CreateOrderService {
	s.side = side
	return s
}

// PositionSide set side
func (s *CreateOrderService) PositionSide(positionSide PositionSideType) *CreateOrderService {
	s.positionSide = &positionSide
	return s
}

// Type set type
func (s *CreateOrderService) Type(orderType OrderType) *CreateOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CreateOrderService) TimeInForce(timeInForce TimeInForceType) *CreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateOrderService) Quantity(quantity string) *CreateOrderService {
	s.quantity = quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *CreateOrderService) ReduceOnly(reduceOnly bool) *CreateOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// Price set price
func (s *CreateOrderService) Price(price string) *CreateOrderService {
	s.price = &price
	return s
}

// NewClientOrderID set newClientOrderID
func (s *CreateOrderService) NewClientOrderID(newClientOrderID string) *CreateOrderService {
	s.newClientOrderID = &newClientOrderID
	return s
}

//
//// StopPrice set stopPrice
//func (s *CreateOrderService) StopPrice(stopPrice string) *CreateOrderService {
//	s.stopPrice = &stopPrice
//	return s
//}
//
//// WorkingType set workingType
//func (s *CreateOrderService) WorkingType(workingType WorkingType) *CreateOrderService {
//	s.workingType = &workingType
//	return s
//}
//
//// ActivationPrice set activationPrice
//func (s *CreateOrderService) ActivationPrice(activationPrice string) *CreateOrderService {
//	s.activationPrice = &activationPrice
//	return s
//}
//
//// CallbackRate set callbackRate
//func (s *CreateOrderService) CallbackRate(callbackRate string) *CreateOrderService {
//	s.callbackRate = &callbackRate
//	return s
//}
//
//// PriceProtect set priceProtect
//func (s *CreateOrderService) PriceProtect(priceProtect bool) *CreateOrderService {
//	s.priceProtect = &priceProtect
//	return s
//}

// NewOrderResponseType set newOrderResponseType
func (s *CreateOrderService) NewOrderResponseType(newOrderResponseType NewOrderRespType) *CreateOrderService {
	s.newOrderRespType = newOrderResponseType
	return s
}

//// ClosePosition set closePosition
//func (s *CreateOrderService) ClosePosition(closePosition bool) *CreateOrderService {
//	s.closePosition = &closePosition
//	return s
//}

func (s *CreateOrderService) createOrder(ctx context.Context, endpoint string, opts ...RequestOption) (data []byte, header *http.Header, err error) {

	r := &request{
		method:   http.MethodPost,
		endpoint: endpoint,
		secType:  secTypeSigned,
	}
	m := params{
		"symbol":           s.symbol,
		"side":             s.side,
		"type":             s.orderType,
		"newOrderRespType": s.newOrderRespType,
	}
	if s.quantity != "" {
		m["quantity"] = s.quantity
	}
	if s.positionSide != nil {
		m["positionSide"] = *s.positionSide
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.reduceOnly != nil {
		m["reduceOnly"] = *s.reduceOnly
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	//if s.stopPrice != nil {
	//	m["stopPrice"] = *s.stopPrice
	//}
	//if s.workingType != nil {
	//	m["workingType"] = *s.workingType
	//}
	//if s.priceProtect != nil {
	//	m["priceProtect"] = *s.priceProtect
	//}
	//if s.activationPrice != nil {
	//	m["activationPrice"] = *s.activationPrice
	//}
	//if s.callbackRate != nil {
	//	m["callbackRate"] = *s.callbackRate
	//}
	//if s.closePosition != nil {
	//	m["closePosition"] = *s.closePosition
	//}
	r.setFormParams(m)
	data, header, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []byte{}, &http.Header{}, err
	}
	return data, header, nil
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CreateOrderResponse, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	endpoint := fmt.Sprintf("/papi/v1/%s/order", s.which)
	data, header, err := s.createOrder(ctx, endpoint, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CreateOrderResponse)
	err = json.Unmarshal(data, res)
	res.RateLimitOrder10s = header.Get("X-Mbx-Order-Count-10s")
	res.RateLimitOrder1m = header.Get("X-Mbx-Order-Count-1m")

	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateOrderResponse define create order response
type CreateOrderResponse struct {
	ClientOrderId           string           `json:"clientOrderId"`
	CumQty                  string           `json:"cumQty"`
	CumQuote                string           `json:"cumQuote"` // UM only
	CumBase                 string           `json:"cumBase"`  // CM only
	ExecutedQty             string           `json:"executedQty"`
	OrderId                 int              `json:"orderId"`
	AvgPrice                string           `json:"avgPrice"`
	OrigQty                 string           `json:"origQty"`
	Price                   string           `json:"price"`
	ReduceOnly              bool             `json:"reduceOnly"`
	Side                    SideType         `json:"side"`
	PositionSide            PositionSideType `json:"positionSide"`
	Status                  OrderStatusType  `json:"status"`
	Symbol                  string           `json:"symbol"`
	Pair                    string           `json:"pair"` // CM
	TimeInForce             TimeInForceType  `json:"timeInForce"`
	OrderType               OrderType        `json:"type"`
	SelfTradePreventionMode string           `json:"selfTradePreventionMode"` // UM
	GoodTillDate            int64            `json:"goodTillDate"`            // UM
	UpdateTime              int64            `json:"updateTime"`
	// Conditional Order
	BookTime       int64       `json:"bookTime"` // Conditional Order book time
	StrategyId     int         `json:"strategyId"`
	StrategyStatus string      `json:"strategyStatus"`
	StrategyType   string      `json:"strategyType"`
	StopPrice      string      `json:"stopPrice"`     // please ignore when order type is TRAILING_STOP_MARKET
	WorkingType    WorkingType `json:"workingType"`   //
	ActivatePrice  string      `json:"activatePrice"` // activation price, only return with TRAILING_STOP_MARKET order
	PriceProtect   bool        `json:"priceProtect"`  // if conditional order trigger is protected
	// extra info
	RateLimitOrder10s string `json:"rateLimitOrder10s,omitempty"` //
	RateLimitOrder1m  string `json:"rateLimitOrder1m,omitempty"`  //
}

// ListOpenOrdersService list opened orders
type ListOpenOrdersService struct {
	c      *Client
	which  string // 'um' or 'cm'
	symbol string
}

// Which set which product
func (s *ListOpenOrdersService) Which(which string) *ListOpenOrdersService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *ListOpenOrdersService) Symbol(symbol string) *ListOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *ListOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/openOrders", s.which),
		secType:  secTypeSigned,
	}
	if s.symbol != "" {
		r.setParam("symbol", s.symbol)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Order{}, err
	}
	res = make([]*Order, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Order{}, err
	}
	return res, nil
}

// GetOpenOrderService query current open order
type GetOpenOrderService struct {
	c                 *Client
	which             string // 'um' or 'cm'
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Which set which product
func (s *GetOpenOrderService) Which(which string) *GetOpenOrderService {
	s.which = which
	return s
}

func (s *GetOpenOrderService) Symbol(symbol string) *GetOpenOrderService {
	s.symbol = symbol
	return s
}

func (s *GetOpenOrderService) OrderID(orderID int64) *GetOpenOrderService {
	s.orderID = &orderID
	return s
}

func (s *GetOpenOrderService) OrigClientOrderID(origClientOrderID string) *GetOpenOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

func (s *GetOpenOrderService) Do(ctx context.Context, opts ...RequestOption) (res *Order, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/openOrder", s.which),
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID == nil && s.origClientOrderID == nil {
		return nil, errors.New("either orderId or origClientOrderId must be sent")
	}
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setParam("origClientOrderId", *s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetOrderService get an order
type GetOrderService struct {
	c                 *Client
	which             string // 'um' or 'cm'
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Which set which product
func (s *GetOrderService) Which(which string) *GetOrderService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *GetOrderService) Symbol(symbol string) *GetOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *GetOrderService) OrderID(orderID int64) *GetOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *GetOrderService) OrigClientOrderID(origClientOrderID string) *GetOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *GetOrderService) Do(ctx context.Context, opts ...RequestOption) (res *Order, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/order", s.which),
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setParam("origClientOrderId", *s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(Order)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Order define order info
type Order struct {
	Symbol                  string           `json:"symbol"`
	OrderID                 int64            `json:"orderId"`
	ClientOrderID           string           `json:"clientOrderId"`
	Price                   string           `json:"price"`
	ReduceOnly              bool             `json:"reduceOnly"`
	OrigQuantity            string           `json:"origQty"`
	ExecutedQuantity        string           `json:"executedQty"`
	CumQuantity             string           `json:"cumQty"`
	CumQuote                string           `json:"cumQuote"`
	Status                  OrderStatusType  `json:"status"`
	TimeInForce             TimeInForceType  `json:"timeInForce"`
	Type                    OrderType        `json:"type"`
	Side                    SideType         `json:"side"`
	StopPrice               string           `json:"stopPrice"`
	Time                    int64            `json:"time"`
	UpdateTime              int64            `json:"updateTime"`
	WorkingType             WorkingType      `json:"workingType"`
	ActivatePrice           string           `json:"activatePrice"`
	PriceRate               string           `json:"priceRate"`
	AvgPrice                string           `json:"avgPrice"`
	OrigType                OrderType        `json:"origType"`
	PositionSide            PositionSideType `json:"positionSide"`
	PriceProtect            bool             `json:"priceProtect"`
	ClosePosition           bool             `json:"closePosition"`
	PriceMatch              string           `json:"priceMatch"`
	SelfTradePreventionMode string           `json:"selfTradePreventionMode"`
	GoodTillDate            int64            `json:"goodTillDate"`
}

// ListOrdersService all account orders; active, canceled, or filled
type ListOrdersService struct {
	c         *Client
	which     string // 'um' or 'cm'
	symbol    string
	orderID   *int64
	startTime *int64
	endTime   *int64
	limit     *int
}

// Which set which product
func (s *ListOrdersService) Which(which string) *ListOrdersService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *ListOrdersService) Symbol(symbol string) *ListOrdersService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *ListOrdersService) OrderID(orderID int64) *ListOrdersService {
	s.orderID = &orderID
	return s
}

// StartTime set starttime
func (s *ListOrdersService) StartTime(startTime int64) *ListOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endtime
func (s *ListOrdersService) EndTime(endTime int64) *ListOrdersService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListOrdersService) Limit(limit int) *ListOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*Order, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/allOrders", s.which),
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setParam("orderId", *s.orderID)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*Order{}, err
	}
	res = make([]*Order, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*Order{}, err
	}
	return res, nil
}

// CancelOrderService cancel an order
type CancelOrderService struct {
	c                 *Client
	which             string // 'um' or 'cm'
	symbol            string
	orderID           *int64
	origClientOrderID *string
}

// Which set which product
func (s *CancelOrderService) Which(which string) *CancelOrderService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

// OrderID set orderID
func (s *CancelOrderService) OrderID(orderID int64) *CancelOrderService {
	s.orderID = &orderID
	return s
}

// OrigClientOrderID set origClientOrderID
func (s *CancelOrderService) OrigClientOrderID(origClientOrderID string) *CancelOrderService {
	s.origClientOrderID = &origClientOrderID
	return s
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) (res *CancelOrderResponse, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodDelete,
		endpoint: fmt.Sprintf("/papi/v1/%s/order", s.which),
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	if s.orderID != nil {
		r.setFormParam("orderId", *s.orderID)
	}
	if s.origClientOrderID != nil {
		r.setFormParam("origClientOrderId", *s.origClientOrderID)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(CancelOrderResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelOrderResponse define response of canceling order
type CancelOrderResponse struct {
	ClientOrderID    string           `json:"clientOrderId"`
	CumQuantity      string           `json:"cumQty"`
	CumQuote         string           `json:"cumQuote"`
	ExecutedQuantity string           `json:"executedQty"`
	OrderID          int64            `json:"orderId"`
	OrigQuantity     string           `json:"origQty"`
	Price            string           `json:"price"`
	ReduceOnly       bool             `json:"reduceOnly"`
	Side             SideType         `json:"side"`
	Status           OrderStatusType  `json:"status"`
	StopPrice        string           `json:"stopPrice"`
	Symbol           string           `json:"symbol"`
	TimeInForce      TimeInForceType  `json:"timeInForce"`
	Type             OrderType        `json:"type"`
	UpdateTime       int64            `json:"updateTime"`
	WorkingType      WorkingType      `json:"workingType"`
	ActivatePrice    string           `json:"activatePrice"`
	PriceRate        string           `json:"priceRate"`
	OrigType         string           `json:"origType"`
	PositionSide     PositionSideType `json:"positionSide"`
	PriceProtect     bool             `json:"priceProtect"`
}

// CancelAllOpenOrdersService cancel all open orders
type CancelAllOpenOrdersService struct {
	c      *Client
	which  string // 'um' or 'cm'
	symbol string
}

// Which set which product
func (s *CancelAllOpenOrdersService) Which(which string) *CancelAllOpenOrdersService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *CancelAllOpenOrdersService) Symbol(symbol string) *CancelAllOpenOrdersService {
	s.symbol = symbol
	return s
}

// Do send request
func (s *CancelAllOpenOrdersService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	if s.which == "" {
		return errWhichMissing
	}
	r := &request{
		method:   http.MethodDelete,
		endpoint: "/papi/v1/allOpenOrders",
		secType:  secTypeSigned,
	}
	r.setFormParam("symbol", s.symbol)
	_, _, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}

// ListUserLiquidationOrdersService lists user's liquidation orders
type ListUserLiquidationOrdersService struct {
	c             *Client
	which         string // 'um' or 'cm'
	symbol        *string
	autoCloseType ForceOrderCloseType
	startTime     *int64
	endTime       *int64
	limit         *int
}

// Which set which product
func (s *ListUserLiquidationOrdersService) Which(which string) *ListUserLiquidationOrdersService {
	s.which = which
	return s
}

// Symbol set symbol
func (s *ListUserLiquidationOrdersService) Symbol(symbol string) *ListUserLiquidationOrdersService {
	s.symbol = &symbol
	return s
}

// AutoCloseType set symbol
func (s *ListUserLiquidationOrdersService) AutoCloseType(autoCloseType ForceOrderCloseType) *ListUserLiquidationOrdersService {
	s.autoCloseType = autoCloseType
	return s
}

// StartTime set startTime
func (s *ListUserLiquidationOrdersService) StartTime(startTime int64) *ListUserLiquidationOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *ListUserLiquidationOrdersService) EndTime(endTime int64) *ListUserLiquidationOrdersService {
	s.endTime = &endTime
	return s
}

// Limit set limit
func (s *ListUserLiquidationOrdersService) Limit(limit int) *ListUserLiquidationOrdersService {
	s.limit = &limit
	return s
}

// Do send request
func (s *ListUserLiquidationOrdersService) Do(ctx context.Context, opts ...RequestOption) (res []*UserLiquidationOrder, err error) {
	if s.which == "" {
		return nil, errWhichMissing
	}
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("/papi/v1/%s/forceOrders", s.which),
		secType:  secTypeSigned,
	}

	r.setParam("autoCloseType", s.autoCloseType)
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*UserLiquidationOrder{}, err
	}
	res = make([]*UserLiquidationOrder, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*UserLiquidationOrder{}, err
	}
	return res, nil
}

// UserLiquidationOrder defines user's liquidation order
type UserLiquidationOrder struct {
	OrderId       int64  `json:"orderId"`
	Symbol        string `json:"symbol"`
	Pair          string `json:"pair"` // CM only
	Status        string `json:"status"`
	ClientOrderId string `json:"clientOrderId"`
	Price         string `json:"price"`
	AvgPrice      string `json:"avgPrice"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	CumQuote      string `json:"cumQuote"` // UM
	CumBase       string `json:"cumBase"`  // CM
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	OrigType      string `json:"origType"`
	Time          int64  `json:"time"`
	UpdateTime    int64  `json:"updateTime"`
}

// ListMarginForceOrdersService lists margin liquidation orders
type ListMarginForceOrdersService struct {
	c         *Client
	startTime *int64
	endTime   *int64
}

// StartTime set startTime
func (s *ListMarginForceOrdersService) StartTime(startTime int64) *ListMarginForceOrdersService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *ListMarginForceOrdersService) EndTime(endTime int64) *ListMarginForceOrdersService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *ListMarginForceOrdersService) Do(ctx context.Context, opts ...RequestOption) (res *MarginForceOrders, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/papi/v1/margin/forceOrders",
		secType:  secTypeSigned,
	}

	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	r.setParam("size", 100)
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(MarginForceOrders)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type MarginForceOrders struct {
	Rows []struct {
		AvgPrice    string `json:"avgPrice"`
		ExecutedQty string `json:"executedQty"`
		OrderId     int    `json:"orderId"`
		Price       string `json:"price"`
		Qty         string `json:"qty"`
		Side        string `json:"side"`
		Symbol      string `json:"symbol"`
		TimeInForce string `json:"timeInForce"`
		UpdatedTime int64  `json:"updatedTime"`
	} `json:"rows"`
	Total int `json:"total"`
}
