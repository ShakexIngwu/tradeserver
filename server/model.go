package tradeserver

import (
	model "tradeserver/webullmodel"
)

type GetOpenOrdersResponse struct {
	orders    []*order
	accountID string
	Username  string
}

type PostPlaceOrderRequest struct {
	AccKeys                   []string          `json:"acc_keys" binding:"required"`
	Action                    model.OrderSide   `json:"action" binding:"required"`
	LmtPrice                  float32           `json:"lmt_price" binding:"required"`
	Symbol                    string            `json:"symbol" binding:"required"`
	OrderType                 model.OrderType   `json:"order_type"`
	OutsideRegularTradingHour bool              `json:"outside_regular_trading_hour"`
	Quantity                  int32             `json:"quantity" binding:"required"`
	TimeInForce               model.TimeInForce `json:"time_in_force"`
}

type PutModifyOrderRequest struct {
	AccKeys                   []string          `json:"acc_keys" binding:"required"`
	Action                    model.OrderSide   `json:"action" binding:"required"`
	OldLmtPrice               float32           `json:"old_lmt_price" binding:"required"`
	NewLmtPrice               float32           `json:"new_lmt_price" binding:"required"`
	Symbol                    string            `json:"symbol" binding:"required"`
	OrderType                 model.OrderType   `json:"order_type"`
	OutsideRegularTradingHour bool              `json:"outside_regular_trading_hour"`
	Quantity                  int32             `json:"quantity" binding:"required"`
	TimeInForce               model.TimeInForce `json:"time_in_force"`
}

type DeleteOrderRequest struct {
	AccKeys  []string        `json:"acc_keys" binding:"required"`
	Action   model.OrderSide `json:"action" binding:"required"`
	LmtPrice float32         `json:"lmt_price" binding:"required"`
	Symbol   string          `json:"symbol" binding:"required"`
}
