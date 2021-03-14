package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ShakexIngwu/tradeserver/webull"
	model "github.com/ShakexIngwu/tradeserver/webullmodel"
)

const AccInfoJsonFile = "/opt/tradeserver/config/acc_info.json"

type AccountInfo struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"device_name"`
	TradePIN   string `json:"trade_pin"`
	MFA        string `json:"mfa"`
}

type Order struct {
	Action                    model.OrderSide   `json:"action"`
	ComboTickerType           string            `json:"combo_ticker_type"`
	FilledQuantity            int32             `json:"filled_quantity"`
	LmtPrice                  float32           `json:"lmt_price"`
	OrderID                   int32             `json:"order_id"`
	OrderType                 model.OrderType   `json:"order_type"`
	OutsideRegularTradingHour bool              `json:"outside_regular_trading_hour"`
	RemainQuantity            int32             `json:"remain_quantity"`
	Status                    model.OrderStatus `json:"status"`
	Symbol                    string            `json:"symbol"`
	TickerId                  int32             `json:"ticker_id"`
	TimeInForce               model.TimeInForce `json:"time_in_force"`
	TotalQuantity             int32             `json:"total_quantity"`
}

type account struct {
	accountDetails	*model.GetAccountResponse
	accountID   	string
	accountInfo 	AccountInfo
	client      	webull.ClientItf
	// openOrders will not always be updated, be sure to run GetOpenOrders first to update before using it
	openOrders []Order
}

var accounts = make(map[string]account)

func NewAccounts() error {
	Log(Debug, "Initializing accounts...")
	accountInfoMap := make(map[string]AccountInfo)
	jsonFile, err := os.Open(AccInfoJsonFile)
	if err != nil {
		return err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &accountInfoMap)
	if err != nil {
		return err
	}

	for accKey, accInfo := range accountInfoMap {
		cred := &webull.Credentials{
			Username:    accInfo.Username,
			Password:    accInfo.Password,
			TradePIN:    accInfo.TradePIN,
			MFA:         accInfo.MFA,
			DeviceName:  accInfo.DeviceName,
			AccountType: model.AccountType(2), // 1: phone number, 2: email
		}
		client, err := webull.NewClient(cred)
		if err != nil {
			return err
		}

		err = client.TradeLogin(*cred)
		if err != nil {
			return err
		}

		accountID, err := client.GetAccountID()
		if err != nil {
			return err
		}

		accountIdInt, _ := strconv.Atoi(accountID)
		accountDetails, err := client.GetAccount(accountIdInt)
		if err != nil {
			return err
		}

		openOrders, err := GetOpenOrders(accountID, client)
		if err != nil {
			return err
		}

		accounts[accKey] = account{
			accountDetails:	accountDetails,
			accountID:   	accountID,
			accountInfo: 	accInfo,
			client:      	client,
			openOrders:  	openOrders,
		}
		Log(Debug, "Loaded account information for user %s.", accKey)
	}
	Log(Debug, "Successfully loaded all account information.")
	return nil
}

func GetOpenOrders(accountID string, client webull.ClientItf) ([]Order, error) {
	Log(Debug, "Start to get open orders for account ID %s.", accountID)
	orderItems, err := client.GetOrders(accountID, model.WORKING, 10)
	if err != nil {
		return nil, err
	}

	var openOrders []Order
	skipped := 0
	for _, orderItem := range *orderItems {
		if orderItem.Orders != nil && len(orderItem.Orders) == 1 {
			filledQuantity, err := strconv.Atoi(orderItem.Orders[0].FilledQuantity)
			if err != nil {
				Log(Warn, "Failed to parse FilledQuantity from string to integer, skipping order...")
				skipped += 1
				continue
			}
			remainQuantity, err := strconv.Atoi(orderItem.Orders[0].RemainQuantity)
			if err != nil {
				Log(Warn, "Failed to parse RemainQuantity from string to integer, skipping order...")
				skipped += 1
				continue
			}
			totalQuantity, err := strconv.Atoi(orderItem.Orders[0].TotalQuantity)
			if err != nil {
				Log(Warn, "Failed to parse TotalQuantity from string to integer, skipping order...")
				skipped += 1
				continue
			}
			lmtPrice, err := strconv.ParseFloat(orderItem.LmtPrice, 32)
			if err != nil {
				Log(Warn, "Failed to parse limit price from string to float, skipping order...")
				skipped += 1
				continue
			}
			order := Order{
				Action:                    orderItem.Action,
				ComboTickerType:           orderItem.ComboTickerType,
				FilledQuantity:            int32(filledQuantity),
				LmtPrice:                  float32(lmtPrice),
				OrderID:                   orderItem.Orders[0].OrderId,
				OrderType:                 orderItem.Orders[0].OrderType,
				OutsideRegularTradingHour: orderItem.OutsideRegularTradingHour,
				RemainQuantity:            int32(remainQuantity),
				Status:                    orderItem.Status,
				Symbol:                    orderItem.Orders[0].Symbol,
				TickerId:                  orderItem.Orders[0].TickerId,
				TimeInForce:               orderItem.Orders[0].TimeInForce,
				TotalQuantity:             int32(totalQuantity),
			}
			openOrders = append(openOrders, order)
		}
	}
	Log(Debug, "Successfully updated %d open orders for account ID %s, skipped %d orders", len(openOrders), accountID, skipped)
	return openOrders, nil
}
