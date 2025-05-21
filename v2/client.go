package binance

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/adshao/go-binance/v2/portfolio"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	jsoniter "github.com/json-iterator/go"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/adshao/go-binance/v2/options"
)

// SideType define side type of order
type SideType string

// OrderType define order type
type OrderType string

// TimeInForceType define time in force type of order
type TimeInForceType string

// NewOrderRespType define response JSON verbosity
type NewOrderRespType string

// OrderStatusType define order status type
type OrderStatusType string

// SymbolType define symbol type
type SymbolType string

// SymbolStatusType define symbol status type
type SymbolStatusType string

// SymbolFilterType define symbol filter type
type SymbolFilterType string

// UserDataEventType define spot user data event type
type UserDataEventType string

// MarginTransferType define margin transfer type
type MarginTransferType int

// MarginLoanStatusType define margin loan status type
type MarginLoanStatusType string

// MarginRepayStatusType define margin repay status type
type MarginRepayStatusType string

// FuturesTransferStatusType define futures transfer status type
type FuturesTransferStatusType string

// SideEffectType define side effect type for orders
type SideEffectType string

// FuturesTransferType define futures transfer type
type FuturesTransferType int

// TransactionType define transaction type
type TransactionType string

// LendingType define the type of lending (flexible saving, activity, ...)
type LendingType string

// StakingProduct define the staking product (locked staking, flexible defi staking, locked defi staking, ...)
type StakingProduct string

// StakingTransactionType define the staking transaction type (subscription, redemption, interest)
type StakingTransactionType string

// LiquidityOperationType define the type of adding/removing liquidity to a liquidity pool(COMBINATION, SINGLE)
type LiquidityOperationType string

// SwappingStatus define the status of swap when querying the swap history
type SwappingStatus int

// LiquidityRewardType define the type of reward we'd claim
type LiquidityRewardType int

// RewardClaimStatus define the status of claiming a reward
type RewardClaimStatus int

// RateLimitType define the rate limitation types
// see https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#enum-definitions
type RateLimitType string

// RateLimitInterval define the rate limitation intervals
type RateLimitInterval string

// AccountType define the account types
type AccountType string

// Endpoints
var (
	BaseAPIMainURL    = "https://api.binance.com"
	BaseAPITestnetURL = "https://testnet.binance.vision"
)

// UseTestnet switch all the API endpoints from production to the testnet
var UseTestnet = false

// Redefining the standard package
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Global enums
const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"

	OrderTypeLimit           OrderType = "LIMIT"
	OrderTypeMarket          OrderType = "MARKET"
	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
	OrderTypeStopLoss        OrderType = "STOP_LOSS"
	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"

	TimeInForceTypeGTC TimeInForceType = "GTC"
	TimeInForceTypeIOC TimeInForceType = "IOC"
	TimeInForceTypeFOK TimeInForceType = "FOK"

	NewOrderRespTypeACK    NewOrderRespType = "ACK"
	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
	NewOrderRespTypeFULL   NewOrderRespType = "FULL"

	OrderStatusTypeNew             OrderStatusType = "NEW"
	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
	OrderStatusTypeFilled          OrderStatusType = "FILLED"
	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
	OrderStatusTypePendingCancel   OrderStatusType = "PENDING_CANCEL"
	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
	OrderStatusExpiredInMatch      OrderStatusType = "EXPIRED_IN_MATCH" // STP Expired

	SymbolTypeSpot SymbolType = "SPOT"

	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
	SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"

	SymbolFilterTypeLotSize            SymbolFilterType = "LOT_SIZE"
	SymbolFilterTypePriceFilter        SymbolFilterType = "PRICE_FILTER"
	SymbolFilterTypePercentPriceBySide SymbolFilterType = "PERCENT_PRICE_BY_SIDE"
	SymbolFilterTypeMinNotional        SymbolFilterType = "MIN_NOTIONAL"
	SymbolFilterTypeNotional           SymbolFilterType = "NOTIONAL"
	SymbolFilterTypeIcebergParts       SymbolFilterType = "ICEBERG_PARTS"
	SymbolFilterTypeMarketLotSize      SymbolFilterType = "MARKET_LOT_SIZE"
	SymbolFilterTypeMaxNumOrders       SymbolFilterType = "MAX_NUM_ORDERS"
	SymbolFilterTypeMaxNumAlgoOrders   SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
	SymbolFilterTypeTrailingDelta      SymbolFilterType = "TRAILING_DELTA"

	UserDataEventTypeOutboundAccountPosition UserDataEventType = "outboundAccountPosition"
	UserDataEventTypeBalanceUpdate           UserDataEventType = "balanceUpdate"
	UserDataEventTypeExecutionReport         UserDataEventType = "executionReport"
	UserDataEventTypeListStatus              UserDataEventType = "ListStatus"

	MarginTransferTypeToMargin MarginTransferType = 1
	MarginTransferTypeToMain   MarginTransferType = 2

	FuturesTransferTypeToFutures FuturesTransferType = 1
	FuturesTransferTypeToMain    FuturesTransferType = 2

	MarginLoanStatusTypePending   MarginLoanStatusType = "PENDING"
	MarginLoanStatusTypeConfirmed MarginLoanStatusType = "CONFIRMED"
	MarginLoanStatusTypeFailed    MarginLoanStatusType = "FAILED"

	MarginRepayStatusTypePending   MarginRepayStatusType = "PENDING"
	MarginRepayStatusTypeConfirmed MarginRepayStatusType = "CONFIRMED"
	MarginRepayStatusTypeFailed    MarginRepayStatusType = "FAILED"

	FuturesTransferStatusTypePending   FuturesTransferStatusType = "PENDING"
	FuturesTransferStatusTypeConfirmed FuturesTransferStatusType = "CONFIRMED"
	FuturesTransferStatusTypeFailed    FuturesTransferStatusType = "FAILED"

	SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
	SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
	SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"

	TransactionTypeDeposit  TransactionType = "0"
	TransactionTypeWithdraw TransactionType = "1"
	TransactionTypeBuy      TransactionType = "0"
	TransactionTypeSell     TransactionType = "1"

	LendingTypeFlexible LendingType = "DAILY"
	LendingTypeFixed    LendingType = "CUSTOMIZED_FIXED"
	LendingTypeActivity LendingType = "ACTIVITY"

	LiquidityOperationTypeCombination LiquidityOperationType = "COMBINATION"
	LiquidityOperationTypeSingle      LiquidityOperationType = "SINGLE"

	apiKey        = "apiKey"
	timestampKey  = "timestamp"
	signatureKey  = "signature"
	recvWindowKey = "recvWindow"

	StakingProductLockedStaking       = "STAKING"
	StakingProductFlexibleDeFiStaking = "F_DEFI"
	StakingProductLockedDeFiStaking   = "L_DEFI"

	StakingTransactionTypeSubscription = "SUBSCRIPTION"
	StakingTransactionTypeRedemption   = "REDEMPTION"
	StakingTransactionTypeInterest     = "INTEREST"

	SwappingStatusPending SwappingStatus = 0
	SwappingStatusDone    SwappingStatus = 1
	SwappingStatusFailed  SwappingStatus = 2

	RewardTypeTrading   LiquidityRewardType = 0
	RewardTypeLiquidity LiquidityRewardType = 1

	RewardClaimPending RewardClaimStatus = 0
	RewardClaimDone    RewardClaimStatus = 1

	RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
	RateLimitTypeOrders        RateLimitType = "ORDERS"
	RateLimitTypeRawRequests   RateLimitType = "RAW_REQUESTS"

	RateLimitIntervalSecond RateLimitInterval = "SECOND"
	RateLimitIntervalMinute RateLimitInterval = "MINUTE"
	RateLimitIntervalDay    RateLimitInterval = "DAY"

	AccountTypeSpot           AccountType = "SPOT"
	AccountTypeMargin         AccountType = "MARGIN"
	AccountTypeIsolatedMargin AccountType = "ISOLATED_MARGIN"
	AccountTypeUSDTFuture     AccountType = "USDT_FUTURE"
	AccountTypeCoinFuture     AccountType = "COIN_FUTURE"
)

type RateLimits struct {
	RequestWeight1M int
	RawRequest5M    int
	Order10s        int
	Order1m         int
}

func RateLimitsFromHeader(header *http.Header) *RateLimits {
	var String2Int = func(intStr string, args ...int) int {
		v, err := strconv.ParseInt(intStr, 10, 64)
		if err != nil && len(args) > 0 {
			v = int64(args[0])
		}
		return int(v)
	}
	return &RateLimits{
		RequestWeight1M: String2Int(header.Get("X-Mbx-Used-Weight-1m")),
		Order1m:         String2Int(header.Get("X-Mbx-Order-Count-1m")),
		Order10s:        String2Int(header.Get("X-Mbx-Order-Count-10s")),
		RawRequest5M:    String2Int(header.Get("X-Mbx-Raw-Quests-5m")),
	}
}

func RateLimitsFromWsResponse(response *WsApiResponse) *RateLimits {
	var locate = func(typ, interval string, num int) int {
		for _, e := range response.RateLimits {
			if e.RateLimitType == typ && e.Interval == interval && e.IntervalNum == num {
				return e.Count
			}
		}
		return 0
	}
	return &RateLimits{
		RequestWeight1M: locate("REQUEST_WEIGHT", "MINUTE", 1),
		RawRequest5M:    locate("RAW_QUESTS", "MINUTE", 5),
		Order10s:        locate("ORDERS", "SECOND", 10),
		Order1m:         locate("ORDERS", "MINUTE", 1),
	}
}

func currentTimestamp() int64 {
	return FormatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func FormatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func newJSON(data []byte) (j *simplejson.Json, err error) {
	j, err = simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// getAPIEndpoint return the base endpoint of the Rest API according the UseTestnet flag
func getAPIEndpoint() string {
	if UseTestnet {
		return BaseAPITestnetURL
	}
	return BaseAPIMainURL
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey, secretKey string) *Client {
	wsState := WsInit
	c, stopC, disconnectedC := makeConn()
	if c != nil {
		wsState = WsConnected
	}

	client := &Client{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		BaseURL:    getAPIEndpoint(),
		UserAgent:  "Binance/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Binance-golang ", log.LstdFlags),
		WsURL:      getWsAPIEndpoint(),
		Conn:       c,
		StopC:      stopC,
		wsState:    wsState,
	}
	client.handleDisconnected(disconnectedC)

	return client
}

// NewFuturesClient initialize client for futures API
func NewFuturesClient(apiKey, secretKey string) *futures.Client {
	return futures.NewClient(apiKey, secretKey)
}

// NewDeliveryClient initialize client for coin-M futures API
func NewDeliveryClient(apiKey, secretKey string) *delivery.Client {
	return delivery.NewClient(apiKey, secretKey)
}

// NewOptionsClient initialize client for options API
func NewOptionsClient(apiKey, secretKey string) *options.Client {
	return options.NewClient(apiKey, secretKey)
}

// NewPortfolioClient initialize client for portfolio API
func NewPortfolioClient(apiKey, secretKey string) *portfolio.Client {
	return portfolio.NewClient(apiKey, secretKey)
}

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	SecretKey  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger
	TimeOffset int64
	do         doFunc
	WsURL      string
	Conn       *websocket.Conn
	StopC      chan struct{}
	wsState    WsClientState // init/connecting/connected
}

func (c *Client) WsConnected() bool {
	return c.wsState == WsConnected
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	if r.recvWindow > 0 {
		r.setParam(recvWindowKey, r.recvWindow)
	}
	if r.secType == secTypeSigned {
		r.setParam(timestampKey, currentTimestamp()-c.TimeOffset)
	}
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}
	if r.secType == secTypeAPIKey || r.secType == secTypeSigned {
		header.Set("X-MBX-APIKEY", c.APIKey)
	}

	if r.secType == secTypeSigned {
		raw := fmt.Sprintf("%s%s", queryString, bodyString)
		mac := hmac.New(sha256.New, []byte(c.SecretKey))
		_, err = mac.Write([]byte(raw))
		if err != nil {
			return err
		}
		v := url.Values{}
		v.Set(signatureKey, fmt.Sprintf("%x", (mac.Sum(nil))))
		if queryString == "" {
			queryString = v.Encode()
		} else {
			queryString = fmt.Sprintf("%s&%s", queryString, v.Encode())
		}
	}
	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.debug("full url: %s, body: %s", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) ([]byte, *RateLimits, error) {
	var err error
	// prefer to WS API
	if c.WsConnected() && r.wsMethod != "" {
		return c.callWsAPI(ctx, r, opts...)
	}

	err = c.parseRequest(r, opts...)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(common.APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			apiErr.Code = int64(res.StatusCode)
			apiErr.Message = string(data)
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, nil, apiErr
	}
	return data, RateLimitsFromHeader(&res.Header), nil
}

// SetApiEndpoint set api Endpoint
func (c *Client) SetApiEndpoint(url string) *Client {
	c.BaseURL = url
	return c
}

// NewPingService init ping service
func (c *Client) NewPingService() *PingService {
	return &PingService{c: c}
}

// NewServerTimeService init server time service
func (c *Client) NewServerTimeService() *ServerTimeService {
	return &ServerTimeService{c: c}
}

// NewSetServerTimeService init set server time service
func (c *Client) NewSetServerTimeService() *SetServerTimeService {
	return &SetServerTimeService{c: c}
}

// NewDepthService init depth service
func (c *Client) NewDepthService() *DepthService {
	return &DepthService{c: c}
}

// NewAggTradesService init aggregate trades service
func (c *Client) NewAggTradesService() *AggTradesService {
	return &AggTradesService{c: c}
}

// NewRecentTradesService init recent trades service
func (c *Client) NewRecentTradesService() *RecentTradesService {
	return &RecentTradesService{c: c}
}

// NewKlinesService init klines service
func (c *Client) NewKlinesService() *KlinesService {
	return &KlinesService{c: c}
}

// NewListPriceChangeStatsService init list prices change stats service
func (c *Client) NewListPriceChangeStatsService() *ListPriceChangeStatsService {
	return &ListPriceChangeStatsService{c: c}
}

// NewListPricesService init listing prices service
func (c *Client) NewListPricesService() *ListPricesService {
	return &ListPricesService{c: c}
}

// NewListBookTickersService init listing booking tickers service
func (c *Client) NewListBookTickersService() *ListBookTickersService {
	return &ListBookTickersService{c: c}
}

// NewListSymbolTickerService init listing symbols tickers
func (c *Client) NewListSymbolTickerService() *ListSymbolTickerService {
	return &ListSymbolTickerService{c: c}
}

// NewCreateOrderService init creating order service
func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

// NewCreateOCOService init creating OCO service
func (c *Client) NewCreateOCOService() *CreateOCOService {
	return &CreateOCOService{c: c}
}

// NewCancelOCOService init cancel OCO service
func (c *Client) NewCancelOCOService() *CancelOCOService {
	return &CancelOCOService{c: c}
}

// NewGetOrderService init get order service
func (c *Client) NewGetOrderService() *GetOrderService {
	return &GetOrderService{c: c}
}

// NewCancelOrderService init cancel order service
func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}

// NewCancelOpenOrdersService init cancel open orders service
func (c *Client) NewCancelOpenOrdersService() *CancelOpenOrdersService {
	return &CancelOpenOrdersService{c: c}
}

// NewListOpenOrdersService init list open orders service
func (c *Client) NewListOpenOrdersService() *ListOpenOrdersService {
	return &ListOpenOrdersService{c: c}
}

// NewListOpenOrderService init list open order service
func (c *Client) NewListOpenOrderService() *ListOpenOrderService {
	return &ListOpenOrderService{c: c}
}

// NewListOrdersService init listing orders service
func (c *Client) NewListOrdersService() *ListOrdersService {
	return &ListOrdersService{c: c}
}

// NewGetAccountService init getting account service
func (c *Client) NewGetAccountService() *GetAccountService {
	return &GetAccountService{c: c}
}

// NewGetAPIKeyPermission init getting API key permission
func (c *Client) NewGetAPIKeyPermission() *GetAPIKeyPermission {
	return &GetAPIKeyPermission{c: c}
}

// NewSavingFlexibleProductPositionsService get flexible products positions (Savings)
func (c *Client) NewSavingFlexibleProductPositionsService() *SavingFlexibleProductPositionsService {
	return &SavingFlexibleProductPositionsService{c: c}
}

// NewSavingFixedProjectPositionsService get fixed project positions (Savings)
func (c *Client) NewSavingFixedProjectPositionsService() *SavingFixedProjectPositionsService {
	return &SavingFixedProjectPositionsService{c: c}
}

// NewListSavingsFlexibleProductsService get flexible products list (Savings)
func (c *Client) NewListSavingsFlexibleProductsService() *ListSavingsFlexibleProductsService {
	return &ListSavingsFlexibleProductsService{c: c}
}

// NewPurchaseSavingsFlexibleProductService purchase a flexible product (Savings)
func (c *Client) NewPurchaseSavingsFlexibleProductService() *PurchaseSavingsFlexibleProductService {
	return &PurchaseSavingsFlexibleProductService{c: c}
}

// NewRedeemSavingsFlexibleProductService redeem a flexible product (Savings)
func (c *Client) NewRedeemSavingsFlexibleProductService() *RedeemSavingsFlexibleProductService {
	return &RedeemSavingsFlexibleProductService{c: c}
}

// NewListSavingsFixedAndActivityProductsService get fixed and activity product list (Savings)
func (c *Client) NewListSavingsFixedAndActivityProductsService() *ListSavingsFixedAndActivityProductsService {
	return &ListSavingsFixedAndActivityProductsService{c: c}
}

// NewGetAccountSnapshotService init getting account snapshot service
func (c *Client) NewGetAccountSnapshotService() *GetAccountSnapshotService {
	return &GetAccountSnapshotService{c: c}
}

// NewListTradesService init listing trades service
func (c *Client) NewListTradesService() *ListTradesService {
	return &ListTradesService{c: c}
}

// NewHistoricalTradesService init listing trades service
func (c *Client) NewHistoricalTradesService() *HistoricalTradesService {
	return &HistoricalTradesService{c: c}
}

// NewListDepositsService init listing deposits service
func (c *Client) NewListDepositsService() *ListDepositsService {
	return &ListDepositsService{c: c}
}

// NewGetDepositAddressService init getting deposit address service
func (c *Client) NewGetDepositAddressService() *GetDepositsAddressService {
	return &GetDepositsAddressService{c: c}
}

// NewCreateWithdrawService init creating withdraw service
func (c *Client) NewCreateWithdrawService() *CreateWithdrawService {
	return &CreateWithdrawService{c: c}
}

// NewListWithdrawsService init listing withdraw service
func (c *Client) NewListWithdrawsService() *ListWithdrawsService {
	return &ListWithdrawsService{c: c}
}

// NewStartUserStreamService init starting user stream service
func (c *Client) NewStartUserStreamService() *StartUserStreamService {
	return &StartUserStreamService{c: c}
}

// NewKeepaliveUserStreamService init keep alive user stream service
func (c *Client) NewKeepaliveUserStreamService() *KeepaliveUserStreamService {
	return &KeepaliveUserStreamService{c: c}
}

// NewCloseUserStreamService init closing user stream service
func (c *Client) NewCloseUserStreamService() *CloseUserStreamService {
	return &CloseUserStreamService{c: c}
}

// NewExchangeInfoService init exchange info service
func (c *Client) NewExchangeInfoService() *ExchangeInfoService {
	return &ExchangeInfoService{c: c}
}

// NewRateLimitService init rate limit service
func (c *Client) NewRateLimitService() *RateLimitService {
	return &RateLimitService{c: c}
}

// NewGetAssetDetailService init get asset detail service
func (c *Client) NewGetAssetDetailService() *GetAssetDetailService {
	return &GetAssetDetailService{c: c}
}

// NewAveragePriceService init average price service
func (c *Client) NewAveragePriceService() *AveragePriceService {
	return &AveragePriceService{c: c}
}

// NewMarginTransferService init margin account transfer service
func (c *Client) NewMarginTransferService() *MarginTransferService {
	return &MarginTransferService{c: c}
}

// NewMarginLoanService init margin account loan service
func (c *Client) NewMarginLoanService() *MarginLoanService {
	return &MarginLoanService{c: c}
}

// NewMarginRepayService init margin account repay service
func (c *Client) NewMarginRepayService() *MarginRepayService {
	return &MarginRepayService{c: c}
}

// NewCreateMarginOrderService init creating margin order service
func (c *Client) NewCreateMarginOrderService() *CreateMarginOrderService {
	return &CreateMarginOrderService{c: c}
}

// NewCancelMarginOrderService init cancel order service
func (c *Client) NewCancelMarginOrderService() *CancelMarginOrderService {
	return &CancelMarginOrderService{c: c}
}

// NewCreateMarginOCOService init creating margin order service
func (c *Client) NewCreateMarginOCOService() *CreateMarginOCOService {
	return &CreateMarginOCOService{c: c}
}

// NewCancelMarginOCOService init cancel order service
func (c *Client) NewCancelMarginOCOService() *CancelMarginOCOService {
	return &CancelMarginOCOService{c: c}
}

// NewGetMarginOrderService init get order service
func (c *Client) NewGetMarginOrderService() *GetMarginOrderService {
	return &GetMarginOrderService{c: c}
}

// NewListMarginLoansService init list margin loan service
func (c *Client) NewListMarginLoansService() *ListMarginLoansService {
	return &ListMarginLoansService{c: c}
}

// NewListMarginRepaysService init list margin repay service
func (c *Client) NewListMarginRepaysService() *ListMarginRepaysService {
	return &ListMarginRepaysService{c: c}
}

// NewGetMarginAccountService init get margin account service
func (c *Client) NewGetMarginAccountService() *GetMarginAccountService {
	return &GetMarginAccountService{c: c}
}

// NewGetIsolatedMarginAccountService init get isolated margin asset service
func (c *Client) NewGetIsolatedMarginAccountService() *GetIsolatedMarginAccountService {
	return &GetIsolatedMarginAccountService{c: c}
}

func (c *Client) NewIsolatedMarginTransferService() *IsolatedMarginTransferService {
	return &IsolatedMarginTransferService{c: c}
}

// NewGetMarginAssetService init get margin asset service
func (c *Client) NewGetMarginAssetService() *GetMarginAssetService {
	return &GetMarginAssetService{c: c}
}

// NewGetMarginPairService init get margin pair service
func (c *Client) NewGetMarginPairService() *GetMarginPairService {
	return &GetMarginPairService{c: c}
}

// NewGetMarginAllPairsService init get margin all pairs service
func (c *Client) NewGetMarginAllPairsService() *GetMarginAllPairsService {
	return &GetMarginAllPairsService{c: c}
}

// NewGetMarginPriceIndexService init get margin price index service
func (c *Client) NewGetMarginPriceIndexService() *GetMarginPriceIndexService {
	return &GetMarginPriceIndexService{c: c}
}

// NewListMarginOpenOrdersService init list margin open orders service
func (c *Client) NewListMarginOpenOrdersService() *ListMarginOpenOrdersService {
	return &ListMarginOpenOrdersService{c: c}
}

// NewListMarginOrdersService init list margin all orders service
func (c *Client) NewListMarginOrdersService() *ListMarginOrdersService {
	return &ListMarginOrdersService{c: c}
}

// NewListMarginTradesService init list margin trades service
func (c *Client) NewListMarginTradesService() *ListMarginTradesService {
	return &ListMarginTradesService{c: c}
}

// NewGetMaxBorrowableService init get max borrowable service
func (c *Client) NewGetMaxBorrowableService() *GetMaxBorrowableService {
	return &GetMaxBorrowableService{c: c}
}

// NewGetMaxTransferableService init get max transferable service
func (c *Client) NewGetMaxTransferableService() *GetMaxTransferableService {
	return &GetMaxTransferableService{c: c}
}

// NewStartMarginUserStreamService init starting margin user stream service
func (c *Client) NewStartMarginUserStreamService() *StartMarginUserStreamService {
	return &StartMarginUserStreamService{c: c}
}

// NewKeepaliveMarginUserStreamService init keep alive margin user stream service
func (c *Client) NewKeepaliveMarginUserStreamService() *KeepaliveMarginUserStreamService {
	return &KeepaliveMarginUserStreamService{c: c}
}

// NewCloseMarginUserStreamService init closing margin user stream service
func (c *Client) NewCloseMarginUserStreamService() *CloseMarginUserStreamService {
	return &CloseMarginUserStreamService{c: c}
}

// NewStartIsolatedMarginUserStreamService init starting margin user stream service
func (c *Client) NewStartIsolatedMarginUserStreamService() *StartIsolatedMarginUserStreamService {
	return &StartIsolatedMarginUserStreamService{c: c}
}

// NewKeepaliveIsolatedMarginUserStreamService init keep alive margin user stream service
func (c *Client) NewKeepaliveIsolatedMarginUserStreamService() *KeepaliveIsolatedMarginUserStreamService {
	return &KeepaliveIsolatedMarginUserStreamService{c: c}
}

// NewCloseIsolatedMarginUserStreamService init closing margin user stream service
func (c *Client) NewCloseIsolatedMarginUserStreamService() *CloseIsolatedMarginUserStreamService {
	return &CloseIsolatedMarginUserStreamService{c: c}
}

// NewFuturesTransferService init futures transfer service
func (c *Client) NewFuturesTransferService() *FuturesTransferService {
	return &FuturesTransferService{c: c}
}

// NewListFuturesTransferService init list futures transfer service
func (c *Client) NewListFuturesTransferService() *ListFuturesTransferService {
	return &ListFuturesTransferService{c: c}
}

// NewListDustLogService init list dust log service
func (c *Client) NewListDustLogService() *ListDustLogService {
	return &ListDustLogService{c: c}
}

// NewDustTransferService init dust transfer service
func (c *Client) NewDustTransferService() *DustTransferService {
	return &DustTransferService{c: c}
}

// NewListDustService init dust list service
func (c *Client) NewListDustService() *ListDustService {
	return &ListDustService{c: c}
}

// NewTransferToSubAccountService transfer to subaccount service
func (c *Client) NewTransferToSubAccountService() *TransferToSubAccountService {
	return &TransferToSubAccountService{c: c}
}

// NewSubaccountAssetsService init list subaccount assets
func (c *Client) NewSubaccountAssetsService() *SubaccountAssetsService {
	return &SubaccountAssetsService{c: c}
}

// NewSubaccountSpotSummaryService init subaccount spot summary
func (c *Client) NewSubaccountSpotSummaryService() *SubaccountSpotSummaryService {
	return &SubaccountSpotSummaryService{c: c}
}

// NewSubaccountDepositAddressService init subaccount deposit address service
func (c *Client) NewSubaccountDepositAddressService() *SubaccountDepositAddressService {
	return &SubaccountDepositAddressService{c: c}
}

// NewSubAccountFuturesPositionRiskService init subaccount futures position risk service
func (c *Client) NewSubAccountFuturesPositionRiskService() *SubAccountFuturesPositionRiskService {
	return &SubAccountFuturesPositionRiskService{c: c}
}

// NewAssetDividendService init the asset dividend list service
func (c *Client) NewAssetDividendService() *AssetDividendService {
	return &AssetDividendService{c: c}
}

// NewUserUniversalTransferService
func (c *Client) NewUserUniversalTransferService() *CreateUserUniversalTransferService {
	return &CreateUserUniversalTransferService{c: c}
}

// NewAllCoinsInformation
func (c *Client) NewGetAllCoinsInfoService() *GetAllCoinsInfoService {
	return &GetAllCoinsInfoService{c: c}
}

// NewDustTransferService init Get All Margin Assets service
func (c *Client) NewGetAllMarginAssetsService() *GetAllMarginAssetsService {
	return &GetAllMarginAssetsService{c: c}
}

// NewFiatDepositWithdrawHistoryService init the fiat deposit/withdraw history service
func (c *Client) NewFiatDepositWithdrawHistoryService() *FiatDepositWithdrawHistoryService {
	return &FiatDepositWithdrawHistoryService{c: c}
}

// NewFiatPaymentsHistoryService init the fiat payments history service
func (c *Client) NewFiatPaymentsHistoryService() *FiatPaymentsHistoryService {
	return &FiatPaymentsHistoryService{c: c}
}

// NewPayTransactionService init the pay transaction service
func (c *Client) NewPayTradeHistoryService() *PayTradeHistoryService {
	return &PayTradeHistoryService{c: c}
}

// NewFiatPaymentsHistoryService init the spot rebate history service
func (c *Client) NewSpotRebateHistoryService() *SpotRebateHistoryService {
	return &SpotRebateHistoryService{c: c}
}

// NewConvertTradeHistoryService init the convert trade history service
func (c *Client) NewConvertTradeHistoryService() *ConvertTradeHistoryService {
	return &ConvertTradeHistoryService{c: c}
}

// NewGetIsolatedMarginAllPairsService init get isolated margin all pairs service
func (c *Client) NewGetIsolatedMarginAllPairsService() *GetIsolatedMarginAllPairsService {
	return &GetIsolatedMarginAllPairsService{c: c}
}

// NewInterestHistoryService init the interest history service
func (c *Client) NewInterestHistoryService() *InterestHistoryService {
	return &InterestHistoryService{c: c}
}

// NewTradeFeeService init the trade fee service
func (c *Client) NewTradeFeeService() *TradeFeeService {
	return &TradeFeeService{c: c}
}

// NewC2CTradeHistoryService init the c2c trade history service
func (c *Client) NewC2CTradeHistoryService() *C2CTradeHistoryService {
	return &C2CTradeHistoryService{c: c}
}

// NewStakingProductPositionService init the staking product position service
func (c *Client) NewStakingProductPositionService() *StakingProductPositionService {
	return &StakingProductPositionService{c: c}
}

// NewStakingHistoryService init the staking history service
func (c *Client) NewStakingHistoryService() *StakingHistoryService {
	return &StakingHistoryService{c: c}
}

// NewGetAllLiquidityPoolService init the get all swap pool service
func (c *Client) NewGetAllLiquidityPoolService() *GetAllLiquidityPoolService {
	return &GetAllLiquidityPoolService{c: c}
}

// NewGetLiquidityPoolDetailService init the get liquidity pool detial service
func (c *Client) NewGetLiquidityPoolDetailService() *GetLiquidityPoolDetailService {
	return &GetLiquidityPoolDetailService{c: c}
}

// NewAddLiquidityPreviewService init the add liquidity preview service
func (c *Client) NewAddLiquidityPreviewService() *AddLiquidityPreviewService {
	return &AddLiquidityPreviewService{c: c}
}

// NewGetSwapQuoteService init the add liquidity preview service
func (c *Client) NewGetSwapQuoteService() *GetSwapQuoteService {
	return &GetSwapQuoteService{c: c}
}

// NewSwapService init the swap service
func (c *Client) NewSwapService() *SwapService {
	return &SwapService{c: c}
}

// NewAddLiquidityService init the add liquidity service
func (c *Client) NewAddLiquidityService() *AddLiquidityService {
	return &AddLiquidityService{c: c}
}

// NewGetUserSwapRecordsService init the service for listing the swap records
func (c *Client) NewGetUserSwapRecordsService() *GetUserSwapRecordsService {
	return &GetUserSwapRecordsService{c: c}
}

// NewClaimRewardService init the service for liquidity pool rewarding
func (c *Client) NewClaimRewardService() *ClaimRewardService {
	return &ClaimRewardService{c: c}
}

// NewRemoveLiquidityService init the service to remvoe liquidity
func (c *Client) NewRemoveLiquidityService() *RemoveLiquidityService {
	return &RemoveLiquidityService{c: c, assets: []string{}}
}

// NewQueryClaimedRewardHistoryService init the service to query reward claiming history
func (c *Client) NewQueryClaimedRewardHistoryService() *QueryClaimedRewardHistoryService {
	return &QueryClaimedRewardHistoryService{c: c}
}

// NewGetBNBBurnService init the service to get BNB Burn on spot trade and margin interest
func (c *Client) NewGetBNBBurnService() *GetBNBBurnService {
	return &GetBNBBurnService{c: c}
}

// NewToggleBNBBurnService init the service to toggle BNB Burn on spot trade and margin interest
func (c *Client) NewToggleBNBBurnService() *ToggleBNBBurnService {
	return &ToggleBNBBurnService{c: c}
}

// NewInternalUniversalTransferService Universal Transfer (For Master Account)
func (c *Client) NewInternalUniversalTransferService() *InternalUniversalTransferService {
	return &InternalUniversalTransferService{c: c}
}

// NewInternalUniversalTransferHistoryService Query Universal Transfer History (For Master Account)
func (c *Client) NewInternalUniversalTransferHistoryService() *InternalUniversalTransferHistoryService {
	return &InternalUniversalTransferHistoryService{c: c}
}

// NewSubAccountListService Query Sub-account List (For Master Account)
func (c *Client) NewSubAccountListService() *SubAccountListService {
	return &SubAccountListService{c: c}
}

// NewGetUserAsset Get user assets, just for positive data
func (c *Client) NewGetUserAsset() *GetUserAssetService {
	return &GetUserAssetService{c: c}
}

// NewManagedSubAccountDepositService Deposit Assets Into The Managed Sub-account（For Investor Master Account）
func (c *Client) NewManagedSubAccountDepositService() *ManagedSubAccountDepositService {
	return &ManagedSubAccountDepositService{c: c}
}

// NewManagedSubAccountWithdrawalService Withdrawal Assets From The Managed Sub-account（For Investor Master Account）
func (c *Client) NewManagedSubAccountWithdrawalService() *ManagedSubAccountWithdrawalService {
	return &ManagedSubAccountWithdrawalService{c: c}
}

// NewManagedSubAccountAssetsService Withdrawal Assets From The Managed Sub-account（For Investor Master Account）
func (c *Client) NewManagedSubAccountAssetsService() *ManagedSubAccountAssetsService {
	return &ManagedSubAccountAssetsService{c: c}
}

// NewSubAccountFuturesAccountService Get Detail on Sub-account's Futures Account (For Master Account)
func (c *Client) NewSubAccountFuturesAccountService() *SubAccountFuturesAccountService {
	return &SubAccountFuturesAccountService{c: c}
}

// NewSubAccountFuturesSummaryV1Service Get Summary of Sub-account's Futures Account (For Master Account)
func (c *Client) NewSubAccountFuturesSummaryV1Service() *SubAccountFuturesSummaryV1Service {
	return &SubAccountFuturesSummaryV1Service{c: c}
}

// NewSimpleEarnAccountService init simple-earn account service
func (c *Client) NewSimpleEarnAccountService() *SimpleEarnAccountService {
	return &SimpleEarnAccountService{c: c}
}

// NewListSimpleEarnFlexibleService init listing simple-earn flexible service
func (c *Client) NewListSimpleEarnFlexibleService() *ListSimpleEarnFlexibleService {
	return &ListSimpleEarnFlexibleService{c: c}
}

// NewListSimpleEarnLockedService init listing simple-earn locked service
func (c *Client) NewListSimpleEarnLockedService() *ListSimpleEarnLockedService {
	return &ListSimpleEarnLockedService{c: c}
}

// NewSubscribeSimpleEarnFlexibleService subscribe to simple-earn flexible service
func (c *Client) NewSubscribeSimpleEarnFlexibleService() *SubscribeSimpleEarnFlexibleService {
	return &SubscribeSimpleEarnFlexibleService{c: c}
}

// NewSubscribeSimpleEarnLockedService subscribe to simple-earn locked service
func (c *Client) NewSubscribeSimpleEarnLockedService() *SubscribeSimpleEarnLockedService {
	return &SubscribeSimpleEarnLockedService{c: c}
}

// NewRedeemSimpleEarnFlexibleService redeem simple-earn flexible service
func (c *Client) NewRedeemSimpleEarnFlexibleService() *RedeemSimpleEarnFlexibleService {
	return &RedeemSimpleEarnFlexibleService{c: c}
}

// NewRedeemSimpleEarnLockedService redeem simple-earn locked service
func (c *Client) NewRedeemSimpleEarnLockedService() *RedeemSimpleEarnLockedService {
	return &RedeemSimpleEarnLockedService{c: c}
}

// NewGetSimpleEarnFlexiblePositionService returns simple-earn flexible position service
func (c *Client) NewGetSimpleEarnFlexiblePositionService() *GetSimpleEarnFlexiblePositionService {
	return &GetSimpleEarnFlexiblePositionService{c: c}
}

// NewListSimpleEarnFlexibleRateHistoryService returns simple-earn listing flexible rate history
func (c *Client) NewListSimpleEarnFlexibleRateHistoryService() *ListSimpleEarnFlexibleRateHistoryService {
	return &ListSimpleEarnFlexibleRateHistoryService{c: c}
}

// NewGetSimpleEarnLockedPositionService returns simple-earn locked position service
func (c *Client) NewGetSimpleEarnLockedPositionService() *GetSimpleEarnLockedPositionService {
	return &GetSimpleEarnLockedPositionService{c: c}
}

// NewListLoanableCoinService returns crypto-loan list locked loanable data service
func (c *Client) NewListLoanableCoinService() *ListLoanableCoinService {
	return &ListLoanableCoinService{c: c}
}

// NewListCollateralCoinService returns crypto-loan list locked collateral data service
func (c *Client) NewListCollateralCoinService() *ListCollateralCoinService {
	return &ListCollateralCoinService{c: c}
}

// NewLoanBorrowLockedService returns crypto-loan locked borrow service
func (c *Client) NewLoanBorrowLockedService() *LoanBorrowLockedService {
	return &LoanBorrowLockedService{c: c}
}

// NewLoanRepayLockedService returns crypto-loan locked repay service
func (c *Client) NewLoanRepayLockedService() *LoanRepayLockedService {
	return &LoanRepayLockedService{c: c}
}

// NewListLoanableCoinFlexibleService returns crypto-loan list flexible loanable data service
func (c *Client) NewListLoanableCoinFlexibleService() *ListLoanableCoinFlexibleService {
	return &ListLoanableCoinFlexibleService{c: c}
}

// NewListCollateralCoinFlexibleService returns crypto-loan list flexible collateral data service
func (c *Client) NewListCollateralCoinFlexibleService() *ListCollateralCoinFlexibleService {
	return &ListCollateralCoinFlexibleService{c: c}
}

// NewLoanBorrowFlexibleService returns crypto-loan flexible borrow service
func (c *Client) NewLoanBorrowFlexibleService() *LoanBorrowFlexibleService {
	return &LoanBorrowFlexibleService{c: c}
}

// NewLoanRepayFlexibleService returns crypto-loan flexible repay service
func (c *Client) NewLoanRepayFlexibleService() *LoanRepayFlexibleService {
	return &LoanRepayFlexibleService{c: c}
}

// NewEthStakingAccountService returns eth-stake account service
func (c *Client) NewEthStakingAccountService() *EthStakingAccountService {
	return &EthStakingAccountService{c: c}
}

// NewEthStakingHistoryService returns eth-stake staking history service
func (c *Client) NewEthStakingHistoryService() *EthStakingHistoryService {
	return &EthStakingHistoryService{c: c}
}

// NewEthStakingRedemptionHistoryService returns eth-stake redemption history service
func (c *Client) NewEthStakingRedemptionHistoryService() *EthStakingRedemptionHistoryService {
	return &EthStakingRedemptionHistoryService{c: c}
}

// NewEthStakingRewardsHistoryService returns eth-stake rewards history service
func (c *Client) NewEthStakingRewardsHistoryService() *EthStakingRewardsHistoryService {
	return &EthStakingRewardsHistoryService{c: c}
}

// NewEthStakingService returns eth-stake staking service
func (c *Client) NewEthStakingService() *EthStakingService {
	return &EthStakingService{c: c}
}

// NewEthWrappingService returns eth-stake wrapping service
func (c *Client) NewEthWrappingService() *EthWrappingService {
	return &EthWrappingService{c: c}
}

// NewEthRedeemService returns eth-stake redeem service
func (c *Client) NewEthRedeemService() *EthRedeemService {
	return &EthRedeemService{c: c}
}

// NewGetFundingAssetService returns wallet get funding asset service
func (c *Client) NewGetFundingAssetService() *GetFundingAssetService {
	return &GetFundingAssetService{c: c}
}

// NewListLoanFlexibleService returns list crypto loan flexible order service
func (c *Client) NewListLoanFlexibleService() *ListLoanFlexibleService {
	return &ListLoanFlexibleService{c: c}
}

// NewAdjustLtvLoanFlexibleService returns adjust crypto loan flexible LTV service
func (c *Client) NewAdjustLtvLoanFlexibleService() *AdjustLtvLoanFlexibleService {
	return &AdjustLtvLoanFlexibleService{c: c}
}

// NewListLoanLockedService returns list crypto loan locked order service
func (c *Client) NewListLoanLockedService() *ListLoanLockedService {
	return &ListLoanLockedService{c: c}
}

// VIP Loan

// NewListVipLoanableCoinService returns crypto-loan list vip loanable data service
func (c *Client) NewListVipLoanableCoinService() *ListVipLoanableCoinService {
	return &ListVipLoanableCoinService{c: c}
}

// NewListVipCollateralCoinService returns crypto-loan list vip collateral data service
func (c *Client) NewListVipCollateralCoinService() *ListVipCollateralCoinService {
	return &ListVipCollateralCoinService{c: c}
}

// NewVipLoanBorrowService returns crypto-loan vip borrow service
func (c *Client) NewVipLoanBorrowService() *VipLoanBorrowService {
	return &VipLoanBorrowService{c: c}
}

// NewVipLoanRepayService returns crypto-loan vip repay service
func (c *Client) NewVipLoanRepayService() *VipLoanRepayService {
	return &VipLoanRepayService{c: c}
}

// NewListVipLoanService returns list crypto-loan vip ongoing order service
func (c *Client) NewListVipLoanService() *ListVipLoanService {
	return &ListVipLoanService{c: c}
}

// SOL staking

// NewSolStakingAccountService returns sol-stake account service
func (c *Client) NewSolStakingAccountService() *SolStakingAccountService {
	return &SolStakingAccountService{c: c}
}

// NewSolStakingHistoryService returns sol-stake staking history service
func (c *Client) NewSolStakingHistoryService() *SolStakingHistoryService {
	return &SolStakingHistoryService{c: c}
}

// NewSolStakingRedemptionHistoryService returns sol-stake redemption history service
func (c *Client) NewSolStakingRedemptionHistoryService() *SolStakingRedemptionHistoryService {
	return &SolStakingRedemptionHistoryService{c: c}
}

// NewSolStakingRewardsHistoryService returns sol-stake rewards history service
func (c *Client) NewSolStakingRewardsHistoryService() *SolStakingRewardsHistoryService {
	return &SolStakingRewardsHistoryService{c: c}
}

// NewSolStakingService returns sol-stake staking service
func (c *Client) NewSolStakingService() *SolStakingService {
	return &SolStakingService{c: c}
}

// NewSolRedeemService returns sol-stake redeem service
func (c *Client) NewSolRedeemService() *SolRedeemService {
	return &SolRedeemService{c: c}
}
