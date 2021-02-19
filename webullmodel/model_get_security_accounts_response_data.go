/*
 * Webull API
 *
 * Webull API Documentation
 *
 * API version: 3.0.1
 * Contact: austin.millan@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package webullmodel
// GetSecurityAccountsResponseData struct for GetSecurityAccountsResponseData
type GetSecurityAccountsResponseData struct {
	AccountTypes []string `json:"accountTypes,omitempty"`
	AllowDeposit bool `json:"allowDeposit,omitempty"`
	BrokerAccountId string `json:"brokerAccountId,omitempty"`
	BrokerId int32 `json:"brokerId,omitempty"`
	BrokerName string `json:"brokerName,omitempty"`
	ComboTypes []string `json:"comboTypes,omitempty"`
	CustomerType string `json:"customerType,omitempty"`
	Deposit bool `json:"deposit,omitempty"`
	DepositStatus string `json:"depositStatus,omitempty"`
	GiftStockStatus int32 `json:"giftStockStatus,omitempty"`
	IsDefault bool `json:"isDefault,omitempty"`
	IsDefaultChecked bool `json:"isDefaultChecked,omitempty"`
	OpenAccountUrl string `json:"openAccountUrl,omitempty"`
	OptionOpenStatus string `json:"optionOpenStatus,omitempty"`
	RegisterRegionId int32 `json:"registerRegionId,omitempty"`
	RegisterTradeUrl string `json:"registerTradeUrl,omitempty"`
	SecAccountId int32 `json:"secAccountId,omitempty"`
	Status string `json:"status,omitempty"`
	SupportClickIPO bool `json:"supportClickIPO,omitempty"`
	SupportOpenOption bool `json:"supportOpenOption,omitempty"`
	SupportOutsideRth bool `json:"supportOutsideRth,omitempty"`
	TickerTypes []GetSecurityAccountsResponseTickerTypes `json:"tickerTypes,omitempty"`
	TimeInForces []GetSecurityAccountsResponseTimeInForces `json:"timeInForces,omitempty"`
	UserTradePermissionVOs []GetSecurityAccountsResponseUserTradePermissionVOs `json:"userTradePermissionVOs,omitempty"`
}