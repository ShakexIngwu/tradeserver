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
// ReplaceOptionOrderRequest struct for ReplaceOptionOrderRequest
type ReplaceOptionOrderRequest struct {
	AuxPrice float32 `json:"auxPrice,omitempty"`
	ComboId string `json:"comboId,omitempty"`
	LmtPrice float32 `json:"lmtPrice,omitempty"`
	OrderType OrderType `json:"orderType,omitempty"`
	Orders []OptionOrder `json:"orders,omitempty"`
	SerialId string `json:"serialId,omitempty"`
	TimeInForce TimeInForce `json:"timeInForce,omitempty"`
}
