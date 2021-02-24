package webull

import (
	"fmt"
	"net/http"
	"time"

	model "github.com/ShakexIngwu/tradeserver/webullmodel"
)

// ClientMock is a mock for client unit testing.
type ClientMock struct {
	Username       string
	HashedPassword string
	AccountType    model.AccountType
	MFA            string
	UUID           string
	DeviceName     string

	AccessToken           string
	AccessTokenExpiration time.Time
	RefreshToken          string

	TradeToken           string
	TradeTokenExpiration time.Time

	DeviceID string

	httpClient *http.Client
	ClientItf
}

// The following functions are mocking functions for client orders
func (m *ClientMock) GetOrders(accountID string, status model.OrderStatus, count int32) (*[]model.GetOrdersItem, error) {
	if accountID == "test_ok" {
		return &[]model.GetOrdersItem{
			{
				Action:          "BUY",
				ComboId:         "test_id",
				ComboTickerType: "stock",
				ComboType:       "stock",
				LmtPrice:        "100.00",
				Orders: []model.GetOrdersItemOrders{
					{
						OrderId:        12345,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "9",
						TotalQuantity:  "10",
						FilledQuantity: "1",
					},
				},
				OutsideRegularTradingHour: true,
				Quantity:                  "10",
				Status:                    "Working",
				TimeInForce:               "GTC",
			},
		}, nil
	} else if accountID == "test_invalid_orders" {
		return &[]model.GetOrdersItem{
			{
				Action:          "BUY",
				ComboId:         "test_id",
				ComboTickerType: "stock",
				ComboType:       "stock",
				FilledQuantity:  "1",
				LmtPrice:        "100.00",
				Orders: []model.GetOrdersItemOrders{
					{
						OrderId:        12345,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "9",
						TotalQuantity:  "10",
						FilledQuantity: "1",
					},
					{
						OrderId:        23456,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "9",
						TotalQuantity:  "10",
						FilledQuantity: "1",
					},
				},
				OutsideRegularTradingHour: true,
				Quantity:                  "10",
				Status:                    "Working",
				TimeInForce:               "GTC",
			},
		}, nil
	} else if accountID == "test_invalid_remain_quantity" {
		return &[]model.GetOrdersItem{
			{
				Action:          "BUY",
				ComboId:         "test_id",
				ComboTickerType: "stock",
				ComboType:       "stock",
				FilledQuantity:  "1",
				LmtPrice:        "100.00",
				Orders: []model.GetOrdersItemOrders{
					{
						OrderId:        12345,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "invalid number",
						TotalQuantity:  "10",
						FilledQuantity: "1",
					},
				},
				OutsideRegularTradingHour: true,
				Quantity:                  "10",
				Status:                    "Working",
				TimeInForce:               "GTC",
			},
		}, nil
	} else if accountID == "test_invalid_total_quantity" {
		return &[]model.GetOrdersItem{
			{
				Action:          "BUY",
				ComboId:         "test_id",
				ComboTickerType: "stock",
				ComboType:       "stock",
				FilledQuantity:  "1",
				LmtPrice:        "100.00",
				Orders: []model.GetOrdersItemOrders{
					{
						OrderId:        12345,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "9",
						TotalQuantity:  "invalid",
						FilledQuantity: "1",
					},
				},
				OutsideRegularTradingHour: true,
				Quantity:                  "10",
				Status:                    "Working",
				TimeInForce:               "GTC",
			},
		}, nil
	} else if accountID == "test_invalid_lmt_price" {
		return &[]model.GetOrdersItem{
			{
				Action:          "BUY",
				ComboId:         "test_id",
				ComboTickerType: "stock",
				ComboType:       "stock",
				FilledQuantity:  "1",
				LmtPrice:        "invalid",
				Orders: []model.GetOrdersItemOrders{
					{
						OrderId:        12345,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "9",
						TotalQuantity:  "10",
						FilledQuantity: "1",
					},
				},
				OutsideRegularTradingHour: true,
				Quantity:                  "10",
				Status:                    "Working",
				TimeInForce:               "GTC",
			},
		}, nil
	} else if accountID == "test_invalid_filled_quantity" {
		return &[]model.GetOrdersItem{
			{
				Action:          "BUY",
				ComboId:         "test_id",
				ComboTickerType: "stock",
				ComboType:       "stock",
				LmtPrice:        "100.00",
				Orders: []model.GetOrdersItemOrders{
					{
						OrderId:        12345,
						OrderType:      "LMT",
						Symbol:         "AMZN",
						TickerId:       54321,
						TimeInForce:    "GTC",
						RemainQuantity: "9",
						TotalQuantity:  "10",
						FilledQuantity: "invalid",
					},
				},
				OutsideRegularTradingHour: true,
				Quantity:                  "10",
				Status:                    "Working",
				TimeInForce:               "GTC",
			},
		}, nil
	} else if accountID == "test_server_error" {
		return nil, fmt.Errorf("internal server error")
	} else {
		return nil, nil
	}
}

func (m *ClientMock) IsTradeable(tickerID string) (*model.GetIsTradeableResponse, error) {
	return nil, nil
}

func (m *ClientMock) PlaceOrder(accountID string, input model.PostStockOrderRequest) (*model.PostOrderResponse, error) {
	return nil, nil
}

func (m *ClientMock) CheckOtocoOrder(accountID string, input model.PostOtocoOrderRequest) (*interface{}, error) {
	return nil, nil
}

func (m *ClientMock) PlaceOtocoOrder(accountID string, input model.PostOtocoOrderRequest) (*interface{}, error) {
	return nil, nil
}

func (m *ClientMock) CancelOrder(accountID, orderID string) (*interface{}, error) {
	return nil, nil
}

func (m *ClientMock) ModifyOrder(accountID string, orderID string, input model.PostStockOrderRequest) (*interface{}, error) {
	return nil, nil
}

func (m *ClientMock) GetTickerID(symbol string) (string, error) {
	return "", nil
}
