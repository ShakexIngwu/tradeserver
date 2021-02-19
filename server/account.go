package tradeserver

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ShakexIngwu/tradeserver/webull"
	model "github.com/ShakexIngwu/tradeserver/webullmodel"
)

const AccInfoJsonFile = "/opt/kerish/acc_info.json"

type AccountInfo struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"device_name"`
	TradePIN   string `json:"trade_pin"`
	MFA        string `json:"mfa"`
}

type order struct {
	action                    model.OrderSide
	ComboType                 string
	filledQuantity            int32
	lmtPrice                  float32
	orderID                   int32
	orderType                 model.OrderType
	outsideRegularTradingHour bool
	remainQuantity            int32
	status                    model.OrderStatus
	symbol                    string
	tickerId                  int32
	timeInForce               model.TimeInForce
	totalQuantity             int32
}

type account struct {
	accountID  string
	client     *webull.Client
	// openOrders will not always be updated, be sure to run GetOpenOrders first to update before using it
	openOrders []*order
}

var accounts = make(map[string]account)

func NewAccounts() error {
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

		openOrders, err := GetOpenOrders(accountID, client)
		if err != nil {
			return err
		}

		accounts[accKey] = account{
			accountID:  accountID,
			client:     client,
			openOrders: openOrders,
		}
	}
	return nil
}

func GetOpenOrders(accountID string, client *webull.Client) ([]*order, error) {
	orderItems, err := client.GetOrders(accountID, model.WORKING, 10)
	if err != nil {
		return nil, err
	}

	var openOrders []*order
	for _, orderItem := range *orderItems {
		if orderItem.Orders != nil && len(orderItem.Orders) == 1 {
			filledQuantity, err := strconv.Atoi(orderItem.Orders[0].FilledQuantity)
			if err != nil {
				continue
			}
			remainQuantity, err := strconv.Atoi(orderItem.Orders[0].RemainQuantity)
			if err != nil {
				continue
			}
			totalQuantity, err := strconv.Atoi(orderItem.Orders[0].TotalQuantity)
			if err != nil {
				continue
			}
			lmtPrice, err := strconv.ParseFloat(orderItem.LmtPrice, 32)
			if err != nil {
				continue
			}
			order := order{
				action:                    orderItem.Action,
				ComboType:                 orderItem.ComboType,
				filledQuantity:            int32(filledQuantity),
				lmtPrice:                  float32(lmtPrice),
				orderID:                   orderItem.Orders[0].OrderId,
				orderType:                 orderItem.Orders[0].OrderType,
				outsideRegularTradingHour: orderItem.OutsideRegularTradingHour,
				remainQuantity:            int32(remainQuantity),
				status:                    orderItem.Status,
				symbol:                    orderItem.Orders[0].Symbol,
				tickerId:                  orderItem.Orders[0].TickerId,
				timeInForce:               orderItem.Orders[0].TimeInForce,
				totalQuantity:             int32(totalQuantity),
			}
			openOrders = append(openOrders, &order)
		}
	}
	return openOrders, nil
}
