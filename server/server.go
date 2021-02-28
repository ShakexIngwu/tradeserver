package server

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	model "github.com/ShakexIngwu/tradeserver/webullmodel"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	defaultPort    = 8000
	currentVersion = "api/v1"
)

type TradeServer struct {
	Port   int
	Router *gin.Engine
}

func NewTradeServer() (*TradeServer, error) {
	err := NewAccounts()
	if err != nil {
		Log(Error, "Failed to load account information: %s", err.Error())
		return nil, err
	}

	return &TradeServer{
		Port:   defaultPort,
		Router: gin.Default(),
	}, nil
}

func (t *TradeServer) Work(wg *sync.WaitGroup) {
	defer wg.Done()

	t.configureRoutes()
	if err := t.Router.Run(fmt.Sprintf(":%d", t.Port)); err != nil {
		panic(fmt.Sprintf("Unexpected error happens when running API router: %s", err.Error()))
		return
	}
}

func (t *TradeServer) configureRoutes() {
	v1 := t.Router.Group(currentVersion)
	t.buildOrderRoutes(v1)
}

func (t *TradeServer) buildOrderRoutes(parent *gin.RouterGroup) {
	order := parent.Group("order")
	order.GET("/", t.getOrders)
	order.GET("/open", t.getOpenOrders)
	order.POST("/", t.placeOrder)
	order.PUT("/", t.modifyOrder)
	order.DELETE("/", t.deleteOrder)
}

func (t *TradeServer) getOpenOrders(c *gin.Context) {
	var accKeys []string
	var res []GetOpenOrdersResponse
	if err := c.BindQuery(&accKeys); err != nil {
		Log(Info, "Invalid Account keys %v, error: %s", accKeys, err.Error())
		c.JSON(http.StatusBadRequest, "Account keys are invalid.")
		return
	}
	// If no account keys is passed in, then by default we will get open orders for all available accounts
	if accKeys == nil || len(accKeys) == 0 {
		for k := range accounts {
			accKeys = append(accKeys, k)
		}
	}
	Log(Debug, "Getting open orders for accounts %v", accKeys)
	for _, accKey := range accKeys {
		if account, ok := accounts[accKey]; ok {
			openOrders, err := GetOpenOrders(account.accountID, account.client)
			if err != nil {
				Log(Error, "Failed to get open orders for account %s", accKey)
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get open orders for account %s, error: %s", accKey, err.Error()))
				return
			}
			// Update account cached open orders each time
			account.openOrders = openOrders
			openOrdersRes := GetOpenOrdersResponse{
				Orders:    account.openOrders,
				AccountID: account.accountID,
				Username:  account.accountInfo.Username,
			}
			res = append(res, openOrdersRes)
		} else {
			Log(Warn, "Cannot find account %s", accKey)
		}
	}
	c.JSON(http.StatusOK, res)
}

func (t *TradeServer) getOrders(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, nil)
}

func (t *TradeServer) placeOrder(c *gin.Context) {
	var reqBody PostPlaceOrderRequest
	var succeed []string
	var failed []string
	if err := c.BindJSON(&reqBody); err != nil {
		Log(Info, "Invalid request body, caught error: %s", err.Error())
		c.JSON(http.StatusBadRequest, "Cannot parse request body.")
		return
	}
	for _, accKey := range reqBody.AccKeys {
		if account, ok := accounts[accKey]; ok {
			// Prepare stock order request params
			tickerIDStr, err := account.client.GetTickerID(reqBody.Symbol)
			if err != nil {
				Log(Error, "Cannot find ticker by symbol %s, caught error: %s", reqBody.Symbol, err.Error())
				c.JSON(http.StatusBadRequest, fmt.Sprintf("Cannot find ticker ID by symbol %s, error: %s", reqBody.Symbol, err.Error()))
				return
			}
			tickerID, err := strconv.Atoi(tickerIDStr)
			if err != nil {
				Log(Error, "TickerID string %s cannot be converted into int32: %s", tickerIDStr, err.Error())
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("TickerID string %s cannot be converted into int32: %s", tickerIDStr, err.Error()))
				return
			}
			if reqBody.OrderType == "" {
				reqBody.OrderType = model.LMT
			}
			if reqBody.TimeInForce == "" {
				reqBody.TimeInForce = model.GTC
			}
			orderRequest := model.PostStockOrderRequest{
				Action:                    reqBody.Action,
				ComboType:                 "stock",
				LmtPrice:                  reqBody.LmtPrice,
				OrderType:                 reqBody.OrderType,
				OutsideRegularTradingHour: reqBody.OutsideRegularTradingHour,
				Quantity:                  reqBody.Quantity,
				SerialId:                  uuid.NewString(),
				TickerId:                  int32(tickerID),
				TimeInForce:               reqBody.TimeInForce,
			}
			_, err = account.client.PlaceOrder(account.accountID, orderRequest)
			if err != nil {
				Log(Error, "Failed to place order for account %s, caught error: %s", accKey, err.Error())
				failed = append(failed, accKey)
			} else {
				succeed = append(succeed, accKey)
			}
		} else {
			Log(Error, "Cannot find account %s when placing order", accKey)
			failed = append(failed, accKey)
		}
	}
	Log(Debug, "Order placed successfully for users %v, failed for users %v", succeed, failed)
	c.JSON(http.StatusOK, fmt.Sprintf("Order succeeded for users %v;\nOrder failed for users %v", succeed, failed))
}

func (t *TradeServer) modifyOrder(c *gin.Context) {
	var reqBody PutModifyOrderRequest
	var succeed []string
	var missed []string
	var failed []string
	if err := c.BindJSON(&reqBody); err != nil {
		Log(Info, "Invalid request body, caught error: %s", err.Error())
		c.JSON(http.StatusBadRequest, "Cannot parse request body.")
		return
	}
	// Update open orders for users that need order modification
	for _, accKey := range reqBody.AccKeys {
		if account, ok := accounts[accKey]; ok {
			openOrders, err := GetOpenOrders(account.accountID, account.client)
			if err != nil {
				Log(Error, "Failed to get open orders for account %s, caught error: %s", accKey, err.Error())
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get open orders for account %s, error: %s", accKey, err.Error()))
				return
			}
			// update account open orders in case the order has been filled or cancelled.
			account.openOrders = openOrders

			var matchedOrder Order
			for _, order := range account.openOrders {
				if order.Action == reqBody.Action && order.Symbol == reqBody.Symbol && order.LmtPrice == reqBody.OldLmtPrice {
					matchedOrder = order
					break
				}
			}
			if &matchedOrder == nil {
				Log(Warn, "Cannot find order for account %s", accKey)
				missed = append(missed, accKey)
			} else {
				orderRequest := model.PostStockOrderRequest{
					Action:                    reqBody.Action,
					ComboType:                 matchedOrder.ComboTickerType,
					LmtPrice:                  reqBody.NewLmtPrice,
					OrderType:                 reqBody.OrderType,
					OutsideRegularTradingHour: reqBody.OutsideRegularTradingHour,
					Quantity:                  reqBody.Quantity,
					SerialId:                  uuid.NewString(),
					TickerId:                  matchedOrder.TickerId,
					TimeInForce:               reqBody.TimeInForce,
				}
				_, err = account.client.ModifyOrder(account.accountID, fmt.Sprint(matchedOrder.OrderID), orderRequest)
				if err != nil {
					Log(Error, "Failed to modify order for account %s, caught error: %s", accKey, err.Error())
					failed = append(failed, accKey)
				} else {
					succeed = append(succeed, accKey)
				}
			}
		} else {
			Log(Error, "Cannot find account %s when modifying order", accKey)
			failed = append(failed, accKey)
		}
	}
	Log(Debug, "Order modification succeeded for users %v, failed for users %v, missed for users %v", succeed, failed, missed)
	c.JSON(http.StatusOK, fmt.Sprintf("Order modification succeeded for users %v;\nOrder modification failed for users %v;\nOrder modification missed for users %v.", succeed, failed, missed))
}

func (t *TradeServer) deleteOrder(c *gin.Context) {
	var reqBody DeleteOrderRequest
	var succeed []string
	var missed []string
	var failed []string
	if err := c.BindJSON(&reqBody); err != nil {
		Log(Info, "Invalid request body, caught error: %s", err.Error())
		c.JSON(http.StatusBadRequest, "Cannot parse request body.")
		return
	}
	// Update open orders for users that need order modification
	for _, accKey := range reqBody.AccKeys {
		if account, ok := accounts[accKey]; ok {
			openOrders, err := GetOpenOrders(account.accountID, account.client)
			if err != nil {
				Log(Error, "Failed to get open orders for account %s, caught error: %s", accKey, err.Error())
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get open orders for account %s, error: %s", accKey, err.Error()))
				return
			}
			// update account open orders in case the order has been filled or cancelled.
			account.openOrders = openOrders

			orderID := ""
			for _, order := range account.openOrders {
				if order.Action == reqBody.Action && order.Symbol == reqBody.Symbol && order.LmtPrice == reqBody.LmtPrice {
					orderID = fmt.Sprint(order.OrderID)
					break
				}
			}
			if orderID == "" {
				Log(Warn, "Cannot find order for account %s", accKey)
				missed = append(missed, accKey)
			} else {
				_, err = account.client.CancelOrder(account.accountID, orderID)
				if err != nil {
					Log(Error, "Failed to cancel order for account %s, caught error: %s", accKey, err.Error())
					failed = append(failed, accKey)
				} else {
					succeed = append(succeed, accKey)
				}
			}
		} else {
			Log(Error, "Cannot find account %s when cancelling order", accKey)
			failed = append(failed, accKey)
		}
	}
	Log(Debug, "Order cancellation succeeded for users %v, failed for users %v, missed for users %v", succeed, failed, missed)
	c.JSON(http.StatusOK, fmt.Sprintf("Order cancellation succeeded for users %v;\nOrder cancellation failed for users %v;\nOrder cancellation missed for users %v.", succeed, failed, missed))
}
