package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/adshao/go-binance/v2/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
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
	WsAPIMainURL    = "wss://ws-api.binance.com:443/ws-api/v3"
	WsAPITestnetURL = "wss://testnet.binance.vision/ws-api/v3"
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
	RawMessage []byte
}

// 自定义类型
type ObjectType string

// 实现 UnmarshalJSON 方法
func (o *ObjectType) UnmarshalJSON(data []byte) error {
	// 将 JSON 对象解析为字符串
	*o = ObjectType(data) // 直接将原始数据赋值为字符串
	return nil
}

// Global enums
//const (
//	SideTypeBuy  SideType = "BUY"
//	SideTypeSell SideType = "SELL"
//
//	OrderTypeLimit           OrderType = "LIMIT"
//	OrderTypeMarket          OrderType = "MARKET"
//	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
//	OrderTypeStopLoss        OrderType = "STOP_LOSS"
//	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
//	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
//	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
//
//	TimeInForceTypeGTC TimeInForceType = "GTC"
//	TimeInForceTypeIOC TimeInForceType = "IOC"
//	TimeInForceTypeFOK TimeInForceType = "FOK"
//
//	NewOrderRespTypeACK    NewOrderRespType = "ACK"
//	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
//	NewOrderRespTypeFULL   NewOrderRespType = "FULL"
//
//	OrderStatusTypeNew             OrderStatusType = "NEW"
//	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
//	OrderStatusTypeFilled          OrderStatusType = "FILLED"
//	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
//	OrderStatusTypePendingCancel   OrderStatusType = "PENDING_CANCEL"
//	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
//	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
//	OrderStatusExpiredInMatch      OrderStatusType = "EXPIRED_IN_MATCH" // STP Expired
//
//	SymbolTypeSpot SymbolType = "SPOT"
//
//	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
//	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
//	SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
//	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
//	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
//	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
//	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"
//
//	SymbolFilterTypeLotSize            SymbolFilterType = "LOT_SIZE"
//	SymbolFilterTypePriceFilter        SymbolFilterType = "PRICE_FILTER"
//	SymbolFilterTypePercentPriceBySide SymbolFilterType = "PERCENT_PRICE_BY_SIDE"
//	SymbolFilterTypeMinNotional        SymbolFilterType = "MIN_NOTIONAL"
//	SymbolFilterTypeNotional           SymbolFilterType = "NOTIONAL"
//	SymbolFilterTypeIcebergParts       SymbolFilterType = "ICEBERG_PARTS"
//	SymbolFilterTypeMarketLotSize      SymbolFilterType = "MARKET_LOT_SIZE"
//	SymbolFilterTypeMaxNumOrders       SymbolFilterType = "MAX_NUM_ORDERS"
//	SymbolFilterTypeMaxNumAlgoOrders   SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
//	SymbolFilterTypeTrailingDelta      SymbolFilterType = "TRAILING_DELTA"
//
//	UserDataEventTypeOutboundAccountPosition UserDataEventType = "outboundAccountPosition"
//	UserDataEventTypeBalanceUpdate           UserDataEventType = "balanceUpdate"
//	UserDataEventTypeExecutionReport         UserDataEventType = "executionReport"
//	UserDataEventTypeListStatus              UserDataEventType = "ListStatus"
//
//	MarginTransferTypeToMargin MarginTransferType = 1
//	MarginTransferTypeToMain   MarginTransferType = 2
//
//	FuturesTransferTypeToFutures FuturesTransferType = 1
//	FuturesTransferTypeToMain    FuturesTransferType = 2
//
//	MarginLoanStatusTypePending   MarginLoanStatusType = "PENDING"
//	MarginLoanStatusTypeConfirmed MarginLoanStatusType = "CONFIRMED"
//	MarginLoanStatusTypeFailed    MarginLoanStatusType = "FAILED"
//
//	MarginRepayStatusTypePending   MarginRepayStatusType = "PENDING"
//	MarginRepayStatusTypeConfirmed MarginRepayStatusType = "CONFIRMED"
//	MarginRepayStatusTypeFailed    MarginRepayStatusType = "FAILED"
//
//	FuturesTransferStatusTypePending   FuturesTransferStatusType = "PENDING"
//	FuturesTransferStatusTypeConfirmed FuturesTransferStatusType = "CONFIRMED"
//	FuturesTransferStatusTypeFailed    FuturesTransferStatusType = "FAILED"
//
//	SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
//	SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
//	SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"
//
//	TransactionTypeDeposit  TransactionType = "0"
//	TransactionTypeWithdraw TransactionType = "1"
//	TransactionTypeBuy      TransactionType = "0"
//	TransactionTypeSell     TransactionType = "1"
//
//	LendingTypeFlexible LendingType = "DAILY"
//	LendingTypeFixed    LendingType = "CUSTOMIZED_FIXED"
//	LendingTypeActivity LendingType = "ACTIVITY"
//
//	LiquidityOperationTypeCombination LiquidityOperationType = "COMBINATION"
//	LiquidityOperationTypeSingle      LiquidityOperationType = "SINGLE"
//
//	timestampKey  = "timestamp"
//	signatureKey  = "signature"
//	recvWindowKey = "recvWindow"
//
//	StakingProductLockedStaking       = "STAKING"
//	StakingProductFlexibleDeFiStaking = "F_DEFI"
//	StakingProductLockedDeFiStaking   = "L_DEFI"
//
//	StakingTransactionTypeSubscription = "SUBSCRIPTION"
//	StakingTransactionTypeRedemption   = "REDEMPTION"
//	StakingTransactionTypeInterest     = "INTEREST"
//
//	SwappingStatusPending SwappingStatus = 0
//	SwappingStatusDone    SwappingStatus = 1
//	SwappingStatusFailed  SwappingStatus = 2
//
//	RewardTypeTrading   LiquidityRewardType = 0
//	RewardTypeLiquidity LiquidityRewardType = 1
//
//	RewardClaimPending RewardClaimStatus = 0
//	RewardClaimDone    RewardClaimStatus = 1
//
//	RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
//	RateLimitTypeOrders        RateLimitType = "ORDERS"
//	RateLimitTypeRawRequests   RateLimitType = "RAW_REQUESTS"
//
//	RateLimitIntervalSecond RateLimitInterval = "SECOND"
//	RateLimitIntervalMinute RateLimitInterval = "MINUTE"
//	RateLimitIntervalDay    RateLimitInterval = "DAY"
//
//	AccountTypeSpot           AccountType = "SPOT"
//	AccountTypeMargin         AccountType = "MARGIN"
//	AccountTypeIsolatedMargin AccountType = "ISOLATED_MARGIN"
//	AccountTypeUSDTFuture     AccountType = "USDT_FUTURE"
//	AccountTypeCoinFuture     AccountType = "COIN_FUTURE"
//)
//
//func currentTimestamp() int64 {
//	return FormatTimestamp(time.Now())
//}
//
//// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
//func FormatTimestamp(t time.Time) int64 {
//	return t.UnixNano() / int64(time.Millisecond)
//}
//
//func newJSON(data []byte) (j *simplejson.Json, err error) {
//	j, err = simplejson.NewJson(data)
//	if err != nil {
//		return nil, err
//	}
//	return j, nil
//}

// getWsAPIEndpoint return the base endpoint of the WebSocket API according the UseTestnet flag
func getWsAPIEndpoint() string {
	if UseTestnet {
		return WsAPITestnetURL
	}
	return WsAPIMainURL
}

// NewWsClient initialize an WS-API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewWsClient(apiKey, secretKey string) *WsClient {
	c, stopC, disconnectedC := makeConn()
	if c == nil {
		return nil
	}

	client := &WsClient{
		APIKey:    apiKey,
		SecretKey: secretKey,
		BaseURL:   getWsAPIEndpoint(),
		UserAgent: "Binance/golang",
		Conn:      c,
		Logger:    log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
		StopC:     stopC,
		Debug:     true,
		state:     WsConnected,
	}

	client.handleDisconnected(disconnectedC)

	return client
}

func makeConn() (*websocket.Conn, chan struct{}, chan struct{}) {
	Dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
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
				res.RawMessage = message
				a <- res
				close(a)
			}
		}
	}()

	return c, stopC, disconnectedC
}

//type doFunc func(req *http.Request) (*http.Response, error)

// WsClient define API client
type WsClient struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	Conn       *websocket.Conn
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
	StopC      chan struct{}

	state WsClientState // init/connecting/connected
}

func (c *WsClient) handleDisconnected(ch chan struct{}) {
	go func() {
		select {
		case <-ch:
		}
		// if it is triggered by AdminClose, just ignore
		if c.state == WsAdminClosing {
			return
		}
		c.state = WsConnecting
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
					c.state = WsConnected
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

func (c *WsClient) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *WsClient) Close() {
	c.state = WsAdminClosing
	close(c.StopC)
}

// wsRequest define an API wsRequest
type wsRequest struct {
	method     string
	query      params
	recvWindow int64
	secType    secType
	ch         chan interface{}
}

// addParam add param with key/value to query string
func (r *wsRequest) addParam(key string, value interface{}) *wsRequest {
	if r.query == nil {
		r.query = params{}
	}
	r.query[key] = value
	return r
}

// setParam set param with key/value to query string
func (r *wsRequest) setParam(key string, value interface{}) *wsRequest {
	if r.query == nil {
		r.query = params{}
	}

	r.query[key] = value
	return r
}

// setParams set params with key/values to query string
func (r *wsRequest) setParams(m params) *wsRequest {
	for k, v := range m {
		r.setParam(k, v)
	}
	return r
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

func (c *WsClient) parseRequest(r *wsRequest) (err error) {
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
		r.setParam(apiKey, c.APIKey)
		raw := r.query.Encode()

		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		r.setParam(signatureKey, fmt.Sprintf("%x", (mac.Sum(nil))))
	}
	c.debug("method: %s, body: %v", r.method, r.query)

	return nil
}

func (c *WsClient) callAPI(ctx context.Context, r *wsRequest) ([]byte, error) {
	// check client state
	if c.state != WsConnected {
		return []byte{}, fmt.Errorf("not connected")
	}

	err := c.parseRequest(r)
	if err != nil {
		return []byte{}, err
	}

	// allocate channel, size 1
	id, ch := uuid.NewString(), make(chan *WsApiResponse, 1)

	req := map[string]interface{}{
		"id":     id,
		"method": r.method,
	}
	if len(r.query) > 0 {
		req["params"] = r.query
	}

	c.debug("wsRequest: %#v", req)

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
	ctx2, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
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

// SetApiEndpoint set api Endpoint
func (c *WsClient) SetApiEndpoint(url string) *WsClient {
	c.BaseURL = url
	return c
}

// NewPingService init ping service
func (c *WsClient) NewPingService() *WsPingService {
	return &WsPingService{c: c}
}

// NewServerTimeService init server time service
func (c *WsClient) NewServerTimeService() *WsServerTimeService {
	return &WsServerTimeService{c: c}
}

// NewDepthService init depth service
func (c *WsClient) NewDepthService() *WsDepthService {
	return &WsDepthService{c: c}
}

//// NewAggTradesService init aggregate trades service
//func (c *WsClient) NewAggTradesService() *AggTradesService {
//	return &AggTradesService{c: c}
//}
//
//// NewRecentTradesService init recent trades service
//func (c *WsClient) NewRecentTradesService() *RecentTradesService {
//	return &RecentTradesService{c: c}
//}
//
//// NewKlinesService init klines service
//func (c *WsClient) NewKlinesService() *KlinesService {
//	return &KlinesService{c: c}
//}
//
//// NewListPriceChangeStatsService init list prices change stats service
//func (c *WsClient) NewListPriceChangeStatsService() *ListPriceChangeStatsService {
//	return &ListPriceChangeStatsService{c: c}
//}
//
//// NewListPricesService init listing prices service
//func (c *WsClient) NewListPricesService() *ListPricesService {
//	return &ListPricesService{c: c}
//}
//
//// NewListBookTickersService init listing booking tickers service
//func (c *WsClient) NewListBookTickersService() *ListBookTickersService {
//	return &ListBookTickersService{c: c}
//}
//
//// NewListSymbolTickerService init listing symbols tickers
//func (c *WsClient) NewListSymbolTickerService() *ListSymbolTickerService {
//	return &ListSymbolTickerService{c: c}
//}
//
//// NewCreateOrderService init creating order service
//func (c *WsClient) NewCreateOrderService() *CreateOrderService {
//	return &CreateOrderService{c: c}
//}
//
//// NewCreateOCOService init creating OCO service
//func (c *WsClient) NewCreateOCOService() *CreateOCOService {
//	return &CreateOCOService{c: c}
//}
//
//// NewCancelOCOService init cancel OCO service
//func (c *WsClient) NewCancelOCOService() *CancelOCOService {
//	return &CancelOCOService{c: c}
//}
//
//// NewGetOrderService init get order service
//func (c *WsClient) NewGetOrderService() *GetOrderService {
//	return &GetOrderService{c: c}
//}
//
//// NewCancelOrderService init cancel order service
//func (c *WsClient) NewCancelOrderService() *CancelOrderService {
//	return &CancelOrderService{c: c}
//}
//
//// NewCancelOpenOrdersService init cancel open orders service
//func (c *WsClient) NewCancelOpenOrdersService() *CancelOpenOrdersService {
//	return &CancelOpenOrdersService{c: c}
//}
//
//// NewListOpenOrdersService init list open orders service
//func (c *WsClient) NewListOpenOrdersService() *ListOpenOrdersService {
//	return &ListOpenOrdersService{c: c}
//}
//
//// NewListOpenOrderService init list open order service
//func (c *WsClient) NewListOpenOrderService() *ListOpenOrderService {
//	return &ListOpenOrderService{c: c}
//}
//
//// NewListOrdersService init listing orders service
//func (c *WsClient) NewListOrdersService() *ListOrdersService {
//	return &ListOrdersService{c: c}
//}

// NewGetAccountService init getting account service
func (c *WsClient) NewGetAccountService() *WsGetAccountService {
	return &WsGetAccountService{c: c}
}

//
//// NewGetAPIKeyPermission init getting API key permission
//func (c *WsClient) NewGetAPIKeyPermission() *GetAPIKeyPermission {
//	return &GetAPIKeyPermission{c: c}
//}
//
//// NewSavingFlexibleProductPositionsService get flexible products positions (Savings)
//func (c *WsClient) NewSavingFlexibleProductPositionsService() *SavingFlexibleProductPositionsService {
//	return &SavingFlexibleProductPositionsService{c: c}
//}
//
//// NewSavingFixedProjectPositionsService get fixed project positions (Savings)
//func (c *WsClient) NewSavingFixedProjectPositionsService() *SavingFixedProjectPositionsService {
//	return &SavingFixedProjectPositionsService{c: c}
//}
//
//// NewListSavingsFlexibleProductsService get flexible products list (Savings)
//func (c *WsClient) NewListSavingsFlexibleProductsService() *ListSavingsFlexibleProductsService {
//	return &ListSavingsFlexibleProductsService{c: c}
//}
//
//// NewPurchaseSavingsFlexibleProductService purchase a flexible product (Savings)
//func (c *WsClient) NewPurchaseSavingsFlexibleProductService() *PurchaseSavingsFlexibleProductService {
//	return &PurchaseSavingsFlexibleProductService{c: c}
//}
//
//// NewRedeemSavingsFlexibleProductService redeem a flexible product (Savings)
//func (c *WsClient) NewRedeemSavingsFlexibleProductService() *RedeemSavingsFlexibleProductService {
//	return &RedeemSavingsFlexibleProductService{c: c}
//}
//
//// NewListSavingsFixedAndActivityProductsService get fixed and activity product list (Savings)
//func (c *WsClient) NewListSavingsFixedAndActivityProductsService() *ListSavingsFixedAndActivityProductsService {
//	return &ListSavingsFixedAndActivityProductsService{c: c}
//}
//
//// NewGetAccountSnapshotService init getting account snapshot service
//func (c *WsClient) NewGetAccountSnapshotService() *GetAccountSnapshotService {
//	return &GetAccountSnapshotService{c: c}
//}
//
//// NewListTradesService init listing trades service
//func (c *WsClient) NewListTradesService() *ListTradesService {
//	return &ListTradesService{c: c}
//}
//
//// NewHistoricalTradesService init listing trades service
//func (c *WsClient) NewHistoricalTradesService() *HistoricalTradesService {
//	return &HistoricalTradesService{c: c}
//}
//
//// NewListDepositsService init listing deposits service
//func (c *WsClient) NewListDepositsService() *ListDepositsService {
//	return &ListDepositsService{c: c}
//}
//
//// NewGetDepositAddressService init getting deposit address service
//func (c *WsClient) NewGetDepositAddressService() *GetDepositsAddressService {
//	return &GetDepositsAddressService{c: c}
//}
//
//// NewCreateWithdrawService init creating withdraw service
//func (c *WsClient) NewCreateWithdrawService() *CreateWithdrawService {
//	return &CreateWithdrawService{c: c}
//}
//
//// NewListWithdrawsService init listing withdraw service
//func (c *WsClient) NewListWithdrawsService() *ListWithdrawsService {
//	return &ListWithdrawsService{c: c}
//}
//
//// NewStartUserStreamService init starting user stream service
//func (c *WsClient) NewStartUserStreamService() *StartUserStreamService {
//	return &StartUserStreamService{c: c}
//}
//
//// NewKeepaliveUserStreamService init keep alive user stream service
//func (c *WsClient) NewKeepaliveUserStreamService() *KeepaliveUserStreamService {
//	return &KeepaliveUserStreamService{c: c}
//}
//
//// NewCloseUserStreamService init closing user stream service
//func (c *WsClient) NewCloseUserStreamService() *CloseUserStreamService {
//	return &CloseUserStreamService{c: c}
//}
//
//// NewExchangeInfoService init exchange info service
//func (c *WsClient) NewExchangeInfoService() *ExchangeInfoService {
//	return &ExchangeInfoService{c: c}
//}
//
//// NewRateLimitService init rate limit service
//func (c *WsClient) NewRateLimitService() *RateLimitService {
//	return &RateLimitService{c: c}
//}
//
//// NewGetAssetDetailService init get asset detail service
//func (c *WsClient) NewGetAssetDetailService() *GetAssetDetailService {
//	return &GetAssetDetailService{c: c}
//}
//
//// NewAveragePriceService init average price service
//func (c *WsClient) NewAveragePriceService() *AveragePriceService {
//	return &AveragePriceService{c: c}
//}
//
//// NewMarginTransferService init margin account transfer service
//func (c *WsClient) NewMarginTransferService() *MarginTransferService {
//	return &MarginTransferService{c: c}
//}
//
//// NewMarginLoanService init margin account loan service
//func (c *WsClient) NewMarginLoanService() *MarginLoanService {
//	return &MarginLoanService{c: c}
//}
//
//// NewMarginRepayService init margin account repay service
//func (c *WsClient) NewMarginRepayService() *MarginRepayService {
//	return &MarginRepayService{c: c}
//}
//
//// NewCreateMarginOrderService init creating margin order service
//func (c *WsClient) NewCreateMarginOrderService() *CreateMarginOrderService {
//	return &CreateMarginOrderService{c: c}
//}
//
//// NewCancelMarginOrderService init cancel order service
//func (c *WsClient) NewCancelMarginOrderService() *CancelMarginOrderService {
//	return &CancelMarginOrderService{c: c}
//}
//
//// NewCreateMarginOCOService init creating margin order service
//func (c *WsClient) NewCreateMarginOCOService() *CreateMarginOCOService {
//	return &CreateMarginOCOService{c: c}
//}
//
//// NewCancelMarginOCOService init cancel order service
//func (c *WsClient) NewCancelMarginOCOService() *CancelMarginOCOService {
//	return &CancelMarginOCOService{c: c}
//}
//
//// NewGetMarginOrderService init get order service
//func (c *WsClient) NewGetMarginOrderService() *GetMarginOrderService {
//	return &GetMarginOrderService{c: c}
//}
//
//// NewListMarginLoansService init list margin loan service
//func (c *WsClient) NewListMarginLoansService() *ListMarginLoansService {
//	return &ListMarginLoansService{c: c}
//}
//
//// NewListMarginRepaysService init list margin repay service
//func (c *WsClient) NewListMarginRepaysService() *ListMarginRepaysService {
//	return &ListMarginRepaysService{c: c}
//}
//
//// NewGetMarginAccountService init get margin account service
//func (c *WsClient) NewGetMarginAccountService() *GetMarginAccountService {
//	return &GetMarginAccountService{c: c}
//}
//
//// NewGetIsolatedMarginAccountService init get isolated margin asset service
//func (c *WsClient) NewGetIsolatedMarginAccountService() *GetIsolatedMarginAccountService {
//	return &GetIsolatedMarginAccountService{c: c}
//}
//
//func (c *WsClient) NewIsolatedMarginTransferService() *IsolatedMarginTransferService {
//	return &IsolatedMarginTransferService{c: c}
//}
//
//// NewGetMarginAssetService init get margin asset service
//func (c *WsClient) NewGetMarginAssetService() *GetMarginAssetService {
//	return &GetMarginAssetService{c: c}
//}
//
//// NewGetMarginPairService init get margin pair service
//func (c *WsClient) NewGetMarginPairService() *GetMarginPairService {
//	return &GetMarginPairService{c: c}
//}
//
//// NewGetMarginAllPairsService init get margin all pairs service
//func (c *WsClient) NewGetMarginAllPairsService() *GetMarginAllPairsService {
//	return &GetMarginAllPairsService{c: c}
//}
//
//// NewGetMarginPriceIndexService init get margin price index service
//func (c *WsClient) NewGetMarginPriceIndexService() *GetMarginPriceIndexService {
//	return &GetMarginPriceIndexService{c: c}
//}
//
//// NewListMarginOpenOrdersService init list margin open orders service
//func (c *WsClient) NewListMarginOpenOrdersService() *ListMarginOpenOrdersService {
//	return &ListMarginOpenOrdersService{c: c}
//}
//
//// NewListMarginOrdersService init list margin all orders service
//func (c *WsClient) NewListMarginOrdersService() *ListMarginOrdersService {
//	return &ListMarginOrdersService{c: c}
//}
//
//// NewListMarginTradesService init list margin trades service
//func (c *WsClient) NewListMarginTradesService() *ListMarginTradesService {
//	return &ListMarginTradesService{c: c}
//}
//
//// NewGetMaxBorrowableService init get max borrowable service
//func (c *WsClient) NewGetMaxBorrowableService() *GetMaxBorrowableService {
//	return &GetMaxBorrowableService{c: c}
//}
//
//// NewGetMaxTransferableService init get max transferable service
//func (c *WsClient) NewGetMaxTransferableService() *GetMaxTransferableService {
//	return &GetMaxTransferableService{c: c}
//}
//
//// NewStartMarginUserStreamService init starting margin user stream service
//func (c *WsClient) NewStartMarginUserStreamService() *StartMarginUserStreamService {
//	return &StartMarginUserStreamService{c: c}
//}
//
//// NewKeepaliveMarginUserStreamService init keep alive margin user stream service
//func (c *WsClient) NewKeepaliveMarginUserStreamService() *KeepaliveMarginUserStreamService {
//	return &KeepaliveMarginUserStreamService{c: c}
//}
//
//// NewCloseMarginUserStreamService init closing margin user stream service
//func (c *WsClient) NewCloseMarginUserStreamService() *CloseMarginUserStreamService {
//	return &CloseMarginUserStreamService{c: c}
//}
//
//// NewStartIsolatedMarginUserStreamService init starting margin user stream service
//func (c *WsClient) NewStartIsolatedMarginUserStreamService() *StartIsolatedMarginUserStreamService {
//	return &StartIsolatedMarginUserStreamService{c: c}
//}
//
//// NewKeepaliveIsolatedMarginUserStreamService init keep alive margin user stream service
//func (c *WsClient) NewKeepaliveIsolatedMarginUserStreamService() *KeepaliveIsolatedMarginUserStreamService {
//	return &KeepaliveIsolatedMarginUserStreamService{c: c}
//}
//
//// NewCloseIsolatedMarginUserStreamService init closing margin user stream service
//func (c *WsClient) NewCloseIsolatedMarginUserStreamService() *CloseIsolatedMarginUserStreamService {
//	return &CloseIsolatedMarginUserStreamService{c: c}
//}
//
//// NewFuturesTransferService init futures transfer service
//func (c *WsClient) NewFuturesTransferService() *FuturesTransferService {
//	return &FuturesTransferService{c: c}
//}
//
//// NewListFuturesTransferService init list futures transfer service
//func (c *WsClient) NewListFuturesTransferService() *ListFuturesTransferService {
//	return &ListFuturesTransferService{c: c}
//}
//
//// NewListDustLogService init list dust log service
//func (c *WsClient) NewListDustLogService() *ListDustLogService {
//	return &ListDustLogService{c: c}
//}
//
//// NewDustTransferService init dust transfer service
//func (c *WsClient) NewDustTransferService() *DustTransferService {
//	return &DustTransferService{c: c}
//}
//
//// NewListDustService init dust list service
//func (c *WsClient) NewListDustService() *ListDustService {
//	return &ListDustService{c: c}
//}
//
//// NewTransferToSubAccountService transfer to subaccount service
//func (c *WsClient) NewTransferToSubAccountService() *TransferToSubAccountService {
//	return &TransferToSubAccountService{c: c}
//}
//
//// NewSubaccountAssetsService init list subaccount assets
//func (c *WsClient) NewSubaccountAssetsService() *SubaccountAssetsService {
//	return &SubaccountAssetsService{c: c}
//}
//
//// NewSubaccountSpotSummaryService init subaccount spot summary
//func (c *WsClient) NewSubaccountSpotSummaryService() *SubaccountSpotSummaryService {
//	return &SubaccountSpotSummaryService{c: c}
//}
//
//// NewSubaccountDepositAddressService init subaccount deposit address service
//func (c *WsClient) NewSubaccountDepositAddressService() *SubaccountDepositAddressService {
//	return &SubaccountDepositAddressService{c: c}
//}
//
//// NewSubAccountFuturesPositionRiskService init subaccount futures position risk service
//func (c *WsClient) NewSubAccountFuturesPositionRiskService() *SubAccountFuturesPositionRiskService {
//	return &SubAccountFuturesPositionRiskService{c: c}
//}
//
//// NewAssetDividendService init the asset dividend list service
//func (c *WsClient) NewAssetDividendService() *AssetDividendService {
//	return &AssetDividendService{c: c}
//}
//
//// NewUserUniversalTransferService
//func (c *WsClient) NewUserUniversalTransferService() *CreateUserUniversalTransferService {
//	return &CreateUserUniversalTransferService{c: c}
//}
//
//// NewAllCoinsInformation
//func (c *WsClient) NewGetAllCoinsInfoService() *GetAllCoinsInfoService {
//	return &GetAllCoinsInfoService{c: c}
//}
//
//// NewDustTransferService init Get All Margin Assets service
//func (c *WsClient) NewGetAllMarginAssetsService() *GetAllMarginAssetsService {
//	return &GetAllMarginAssetsService{c: c}
//}
//
//// NewFiatDepositWithdrawHistoryService init the fiat deposit/withdraw history service
//func (c *WsClient) NewFiatDepositWithdrawHistoryService() *FiatDepositWithdrawHistoryService {
//	return &FiatDepositWithdrawHistoryService{c: c}
//}
//
//// NewFiatPaymentsHistoryService init the fiat payments history service
//func (c *WsClient) NewFiatPaymentsHistoryService() *FiatPaymentsHistoryService {
//	return &FiatPaymentsHistoryService{c: c}
//}
//
//// NewPayTransactionService init the pay transaction service
//func (c *WsClient) NewPayTradeHistoryService() *PayTradeHistoryService {
//	return &PayTradeHistoryService{c: c}
//}
//
//// NewFiatPaymentsHistoryService init the spot rebate history service
//func (c *WsClient) NewSpotRebateHistoryService() *SpotRebateHistoryService {
//	return &SpotRebateHistoryService{c: c}
//}
//
//// NewConvertTradeHistoryService init the convert trade history service
//func (c *WsClient) NewConvertTradeHistoryService() *ConvertTradeHistoryService {
//	return &ConvertTradeHistoryService{c: c}
//}
//
//// NewGetIsolatedMarginAllPairsService init get isolated margin all pairs service
//func (c *WsClient) NewGetIsolatedMarginAllPairsService() *GetIsolatedMarginAllPairsService {
//	return &GetIsolatedMarginAllPairsService{c: c}
//}
//
//// NewInterestHistoryService init the interest history service
//func (c *WsClient) NewInterestHistoryService() *InterestHistoryService {
//	return &InterestHistoryService{c: c}
//}
//
//// NewTradeFeeService init the trade fee service
//func (c *WsClient) NewTradeFeeService() *TradeFeeService {
//	return &TradeFeeService{c: c}
//}
//
//// NewC2CTradeHistoryService init the c2c trade history service
//func (c *WsClient) NewC2CTradeHistoryService() *C2CTradeHistoryService {
//	return &C2CTradeHistoryService{c: c}
//}
//
//// NewStakingProductPositionService init the staking product position service
//func (c *WsClient) NewStakingProductPositionService() *StakingProductPositionService {
//	return &StakingProductPositionService{c: c}
//}
//
//// NewStakingHistoryService init the staking history service
//func (c *WsClient) NewStakingHistoryService() *StakingHistoryService {
//	return &StakingHistoryService{c: c}
//}
//
//// NewGetAllLiquidityPoolService init the get all swap pool service
//func (c *WsClient) NewGetAllLiquidityPoolService() *GetAllLiquidityPoolService {
//	return &GetAllLiquidityPoolService{c: c}
//}
//
//// NewGetLiquidityPoolDetailService init the get liquidity pool detial service
//func (c *WsClient) NewGetLiquidityPoolDetailService() *GetLiquidityPoolDetailService {
//	return &GetLiquidityPoolDetailService{c: c}
//}
//
//// NewAddLiquidityPreviewService init the add liquidity preview service
//func (c *WsClient) NewAddLiquidityPreviewService() *AddLiquidityPreviewService {
//	return &AddLiquidityPreviewService{c: c}
//}
//
//// NewGetSwapQuoteService init the add liquidity preview service
//func (c *WsClient) NewGetSwapQuoteService() *GetSwapQuoteService {
//	return &GetSwapQuoteService{c: c}
//}
//
//// NewSwapService init the swap service
//func (c *WsClient) NewSwapService() *SwapService {
//	return &SwapService{c: c}
//}
//
//// NewAddLiquidityService init the add liquidity service
//func (c *WsClient) NewAddLiquidityService() *AddLiquidityService {
//	return &AddLiquidityService{c: c}
//}
//
//// NewGetUserSwapRecordsService init the service for listing the swap records
//func (c *WsClient) NewGetUserSwapRecordsService() *GetUserSwapRecordsService {
//	return &GetUserSwapRecordsService{c: c}
//}
//
//// NewClaimRewardService init the service for liquidity pool rewarding
//func (c *WsClient) NewClaimRewardService() *ClaimRewardService {
//	return &ClaimRewardService{c: c}
//}
//
//// NewRemoveLiquidityService init the service to remvoe liquidity
//func (c *WsClient) NewRemoveLiquidityService() *RemoveLiquidityService {
//	return &RemoveLiquidityService{c: c, assets: []string{}}
//}
//
//// NewQueryClaimedRewardHistoryService init the service to query reward claiming history
//func (c *WsClient) NewQueryClaimedRewardHistoryService() *QueryClaimedRewardHistoryService {
//	return &QueryClaimedRewardHistoryService{c: c}
//}
//
//// NewGetBNBBurnService init the service to get BNB Burn on spot trade and margin interest
//func (c *WsClient) NewGetBNBBurnService() *GetBNBBurnService {
//	return &GetBNBBurnService{c: c}
//}
//
//// NewToggleBNBBurnService init the service to toggle BNB Burn on spot trade and margin interest
//func (c *WsClient) NewToggleBNBBurnService() *ToggleBNBBurnService {
//	return &ToggleBNBBurnService{c: c}
//}
//
//// NewInternalUniversalTransferService Universal Transfer (For Master Account)
//func (c *WsClient) NewInternalUniversalTransferService() *InternalUniversalTransferService {
//	return &InternalUniversalTransferService{c: c}
//}
//
//// NewInternalUniversalTransferHistoryService Query Universal Transfer History (For Master Account)
//func (c *WsClient) NewInternalUniversalTransferHistoryService() *InternalUniversalTransferHistoryService {
//	return &InternalUniversalTransferHistoryService{c: c}
//}
//
//// NewSubAccountListService Query Sub-account List (For Master Account)
//func (c *WsClient) NewSubAccountListService() *SubAccountListService {
//	return &SubAccountListService{c: c}
//}
//
//// NewGetUserAsset Get user assets, just for positive data
//func (c *WsClient) NewGetUserAsset() *GetUserAssetService {
//	return &GetUserAssetService{c: c}
//}
//
//// NewManagedSubAccountDepositService Deposit Assets Into The Managed Sub-account（For Investor Master Account）
//func (c *WsClient) NewManagedSubAccountDepositService() *ManagedSubAccountDepositService {
//	return &ManagedSubAccountDepositService{c: c}
//}
//
//// NewManagedSubAccountWithdrawalService Withdrawal Assets From The Managed Sub-account（For Investor Master Account）
//func (c *WsClient) NewManagedSubAccountWithdrawalService() *ManagedSubAccountWithdrawalService {
//	return &ManagedSubAccountWithdrawalService{c: c}
//}
//
//// NewManagedSubAccountAssetsService Withdrawal Assets From The Managed Sub-account（For Investor Master Account）
//func (c *WsClient) NewManagedSubAccountAssetsService() *ManagedSubAccountAssetsService {
//	return &ManagedSubAccountAssetsService{c: c}
//}
//
//// NewSubAccountFuturesAccountService Get Detail on Sub-account's Futures Account (For Master Account)
//func (c *WsClient) NewSubAccountFuturesAccountService() *SubAccountFuturesAccountService {
//	return &SubAccountFuturesAccountService{c: c}
//}
//
//// NewSubAccountFuturesSummaryV1Service Get Summary of Sub-account's Futures Account (For Master Account)
//func (c *WsClient) NewSubAccountFuturesSummaryV1Service() *SubAccountFuturesSummaryV1Service {
//	return &SubAccountFuturesSummaryV1Service{c: c}
//}
//
//// NewSimpleEarnAccountService init simple-earn account service
//func (c *WsClient) NewSimpleEarnAccountService() *SimpleEarnAccountService {
//	return &SimpleEarnAccountService{c: c}
//}
//
//// NewListSimpleEarnFlexibleService init listing simple-earn flexible service
//func (c *WsClient) NewListSimpleEarnFlexibleService() *ListSimpleEarnFlexibleService {
//	return &ListSimpleEarnFlexibleService{c: c}
//}
//
//// NewListSimpleEarnLockedService init listing simple-earn locked service
//func (c *WsClient) NewListSimpleEarnLockedService() *ListSimpleEarnLockedService {
//	return &ListSimpleEarnLockedService{c: c}
//}
//
//// NewSubscribeSimpleEarnFlexibleService subscribe to simple-earn flexible service
//func (c *WsClient) NewSubscribeSimpleEarnFlexibleService() *SubscribeSimpleEarnFlexibleService {
//	return &SubscribeSimpleEarnFlexibleService{c: c}
//}
//
//// NewSubscribeSimpleEarnLockedService subscribe to simple-earn locked service
//func (c *WsClient) NewSubscribeSimpleEarnLockedService() *SubscribeSimpleEarnLockedService {
//	return &SubscribeSimpleEarnLockedService{c: c}
//}
//
//// NewRedeemSimpleEarnFlexibleService redeem simple-earn flexible service
//func (c *WsClient) NewRedeemSimpleEarnFlexibleService() *RedeemSimpleEarnFlexibleService {
//	return &RedeemSimpleEarnFlexibleService{c: c}
//}
//
//// NewRedeemSimpleEarnLockedService redeem simple-earn locked service
//func (c *WsClient) NewRedeemSimpleEarnLockedService() *RedeemSimpleEarnLockedService {
//	return &RedeemSimpleEarnLockedService{c: c}
//}
//
//// NewGetSimpleEarnFlexiblePositionService returns simple-earn flexible position service
//func (c *WsClient) NewGetSimpleEarnFlexiblePositionService() *GetSimpleEarnFlexiblePositionService {
//	return &GetSimpleEarnFlexiblePositionService{c: c}
//}
//
//// NewListSimpleEarnFlexibleRateHistoryService returns simple-earn listing flexible rate history
//func (c *WsClient) NewListSimpleEarnFlexibleRateHistoryService() *ListSimpleEarnFlexibleRateHistoryService {
//	return &ListSimpleEarnFlexibleRateHistoryService{c: c}
//}
//
//// NewGetSimpleEarnLockedPositionService returns simple-earn locked position service
//func (c *WsClient) NewGetSimpleEarnLockedPositionService() *GetSimpleEarnLockedPositionService {
//	return &GetSimpleEarnLockedPositionService{c: c}
//}
//
//// NewListLoanableCoinService returns crypto-loan list locked loanable data service
//func (c *WsClient) NewListLoanableCoinService() *ListLoanableCoinService {
//	return &ListLoanableCoinService{c: c}
//}
//
//// NewListCollateralCoinService returns crypto-loan list locked collateral data service
//func (c *WsClient) NewListCollateralCoinService() *ListCollateralCoinService {
//	return &ListCollateralCoinService{c: c}
//}
//
//// NewLoanBorrowLockedService returns crypto-loan locked borrow service
//func (c *WsClient) NewLoanBorrowLockedService() *LoanBorrowLockedService {
//	return &LoanBorrowLockedService{c: c}
//}
//
//// NewLoanRepayLockedService returns crypto-loan locked repay service
//func (c *WsClient) NewLoanRepayLockedService() *LoanRepayLockedService {
//	return &LoanRepayLockedService{c: c}
//}
//
//// NewListLoanableCoinFlexibleService returns crypto-loan list flexible loanable data service
//func (c *WsClient) NewListLoanableCoinFlexibleService() *ListLoanableCoinFlexibleService {
//	return &ListLoanableCoinFlexibleService{c: c}
//}
//
//// NewListCollateralCoinFlexibleService returns crypto-loan list flexible collateral data service
//func (c *WsClient) NewListCollateralCoinFlexibleService() *ListCollateralCoinFlexibleService {
//	return &ListCollateralCoinFlexibleService{c: c}
//}
//
//// NewLoanBorrowFlexibleService returns crypto-loan flexible borrow service
//func (c *WsClient) NewLoanBorrowFlexibleService() *LoanBorrowFlexibleService {
//	return &LoanBorrowFlexibleService{c: c}
//}
//
//// NewLoanRepayFlexibleService returns crypto-loan flexible repay service
//func (c *WsClient) NewLoanRepayFlexibleService() *LoanRepayFlexibleService {
//	return &LoanRepayFlexibleService{c: c}
//}
//
//// NewEthStakingAccountService returns eth-stake account service
//func (c *WsClient) NewEthStakingAccountService() *EthStakingAccountService {
//	return &EthStakingAccountService{c: c}
//}
//
//// NewEthStakingHistoryService returns eth-stake staking history service
//func (c *WsClient) NewEthStakingHistoryService() *EthStakingHistoryService {
//	return &EthStakingHistoryService{c: c}
//}
//
//// NewEthStakingRedemptionHistoryService returns eth-stake redemption history service
//func (c *WsClient) NewEthStakingRedemptionHistoryService() *EthStakingRedemptionHistoryService {
//	return &EthStakingRedemptionHistoryService{c: c}
//}
//
//// NewEthStakingRewardsHistoryService returns eth-stake rewards history service
//func (c *WsClient) NewEthStakingRewardsHistoryService() *EthStakingRewardsHistoryService {
//	return &EthStakingRewardsHistoryService{c: c}
//}
//
//// NewEthStakingService returns eth-stake staking service
//func (c *WsClient) NewEthStakingService() *EthStakingService {
//	return &EthStakingService{c: c}
//}
//
//// NewEthWrappingService returns eth-stake wrapping service
//func (c *WsClient) NewEthWrappingService() *EthWrappingService {
//	return &EthWrappingService{c: c}
//}
//
//// NewEthRedeemService returns eth-stake redeem service
//func (c *WsClient) NewEthRedeemService() *EthRedeemService {
//	return &EthRedeemService{c: c}
//}
//
//// NewGetFundingAssetService returns wallet get funding asset service
//func (c *WsClient) NewGetFundingAssetService() *GetFundingAssetService {
//	return &GetFundingAssetService{c: c}
//}
//
//// NewListLoanFlexibleService returns list crypto loan flexible order service
//func (c *WsClient) NewListLoanFlexibleService() *ListLoanFlexibleService {
//	return &ListLoanFlexibleService{c: c}
//}
//
//// NewAdjustLtvLoanFlexibleService returns adjust crypto loan flexible LTV service
//func (c *WsClient) NewAdjustLtvLoanFlexibleService() *AdjustLtvLoanFlexibleService {
//	return &AdjustLtvLoanFlexibleService{c: c}
//}
//
//// NewListLoanLockedService returns list crypto loan locked order service
//func (c *WsClient) NewListLoanLockedService() *ListLoanLockedService {
//	return &ListLoanLockedService{c: c}
//}
//
//// VIP Loan
//
//// NewListVipLoanableCoinService returns crypto-loan list vip loanable data service
//func (c *WsClient) NewListVipLoanableCoinService() *ListVipLoanableCoinService {
//	return &ListVipLoanableCoinService{c: c}
//}
//
//// NewListVipCollateralCoinService returns crypto-loan list vip collateral data service
//func (c *WsClient) NewListVipCollateralCoinService() *ListVipCollateralCoinService {
//	return &ListVipCollateralCoinService{c: c}
//}
//
//// NewVipLoanBorrowService returns crypto-loan vip borrow service
//func (c *WsClient) NewVipLoanBorrowService() *VipLoanBorrowService {
//	return &VipLoanBorrowService{c: c}
//}
//
//// NewVipLoanRepayService returns crypto-loan vip repay service
//func (c *WsClient) NewVipLoanRepayService() *VipLoanRepayService {
//	return &VipLoanRepayService{c: c}
//}
//
//// NewListVipLoanService returns list crypto-loan vip ongoing order service
//func (c *WsClient) NewListVipLoanService() *ListVipLoanService {
//	return &ListVipLoanService{c: c}
//}
