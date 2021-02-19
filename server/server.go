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
		c.JSON(http.StatusBadRequest, "Account keys are invalid.")
		return
	}
	for _, accKey := range accKeys {
		if account, ok := accounts[accKey]; ok {
			openOrders, err := GetOpenOrders(account.accountID, account.client)
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get open orders for account %s, error: %s", accKey, err.Error()))
				return
			}
			// Update account cached open orders each time
			account.openOrders = openOrders
			openOrdersRes := GetOpenOrdersResponse{
				orders:    account.openOrders,
				accountID: account.accountID,
				Username:  account.client.Username,
			}
			res = append(res, openOrdersRes)
		} else {
			fmt.Printf("Cannot find account %s", accKey)
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
		c.JSON(http.StatusBadRequest, "Cannot parse request body.")
		return
	}
	for _, accKey := range reqBody.AccKeys {
		if account, ok := accounts[accKey]; ok {
			// Prepare stock order request params
			tickerIDStr, err := account.client.GetTickerID(reqBody.Symbol)
			if err != nil {
				c.JSON(http.StatusBadRequest, fmt.Sprintf("Cannot find ticker ID by symbol %s, error: %s", reqBody.Symbol, err.Error()))
				return
			}
			tickerID, err := strconv.Atoi(tickerIDStr)
			if err != nil {
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
				failed = append(failed, accKey)
			} else {
				succeed = append(succeed, accKey)
			}
		} else {
			failed = append(failed, accKey)
		}
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Order succeeded for users %v;\nOrder failed for users %v", succeed, failed))
}

func (t *TradeServer) modifyOrder(c *gin.Context) {
	var reqBody PutModifyOrderRequest
	var succeed []string
	var missed []string
	var failed []string
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, "Cannot parse request body.")
		return
	}
	// Update open orders for users that need order modification
	for _, accKey := range reqBody.AccKeys {
		if account, ok := accounts[accKey]; ok {
			openOrders, err := GetOpenOrders(account.accountID, account.client)
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get open orders for account %s, error: %s", accKey, err.Error()))
				return
			}
			// update account open orders in case the order has been filled or cancelled.
			account.openOrders = openOrders

			var matchedOrder *order
			for _, order := range account.openOrders {
				if order.action == reqBody.Action && order.symbol == reqBody.Symbol && order.lmtPrice == reqBody.OldLmtPrice{
					matchedOrder = order
					break
				}
			}
			if matchedOrder == nil {
				missed = append(missed, accKey)
			} else {
				orderRequest := model.PostStockOrderRequest{
					Action:                    reqBody.Action,
					ComboType:                 matchedOrder.ComboType,
					LmtPrice:                  reqBody.NewLmtPrice,
					OrderType:                 reqBody.OrderType,
					OutsideRegularTradingHour: reqBody.OutsideRegularTradingHour,
					Quantity:                  reqBody.Quantity,
					SerialId:                  uuid.NewString(),
					TickerId:                  matchedOrder.tickerId,
					TimeInForce:               reqBody.TimeInForce,
				}
				_, err = account.client.ModifyOrder(account.accountID, fmt.Sprint(matchedOrder.orderID), orderRequest)
				if err != nil {
					failed = append(failed, accKey)
				} else {
					succeed = append(succeed, accKey)
				}
			}
		} else {
			failed = append(failed, accKey)
		}
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Order modification succeeded for users %v;\nOrder modification failed for users %v;\nOrder modification missed for users %v.", succeed, failed, missed))
}

func (t *TradeServer) deleteOrder(c *gin.Context) {
	var reqBody DeleteOrderRequest
	var succeed []string
	var missed []string
	var failed []string
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, "Cannot parse request body.")
		return
	}
	// Update open orders for users that need order modification
	for _, accKey := range reqBody.AccKeys {
		if account, ok := accounts[accKey]; ok {
			openOrders, err := GetOpenOrders(account.accountID, account.client)
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("Cannot get open orders for account %s, error: %s", accKey, err.Error()))
				return
			}
			// update account open orders in case the order has been filled or cancelled.
			account.openOrders = openOrders

			orderID := ""
			for _, order := range account.openOrders {
				if order.action == reqBody.Action && order.symbol == reqBody.Symbol && order.lmtPrice == reqBody.LmtPrice {
					orderID = fmt.Sprint(order.orderID)
					break
				}
			}
			if orderID == "" {
				missed = append(missed, accKey)
			} else {
				_, err = account.client.CancelOrder(account.accountID, orderID)
				if err != nil {
					failed = append(failed, accKey)
				} else {
					succeed = append(succeed, accKey)
				}
			}
		} else {
			failed = append(failed, accKey)
		}
	}
	c.JSON(http.StatusOK, fmt.Sprintf("Order cancellation succeeded for users %v;\nOrder cancellation failed for users %v;\nOrder cancellation missed for users %v.", succeed, failed, missed))
}