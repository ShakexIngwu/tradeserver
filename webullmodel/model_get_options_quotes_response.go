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
// GetOptionsQuotesResponse struct for GetOptionsQuotesResponse
type GetOptionsQuotesResponse struct {
	Data []map[string]interface{} `json:"data,omitempty"`
	DisExchangeCode string `json:"disExchangeCode,omitempty"`
	DisSymbol string `json:"disSymbol,omitempty"`
	Name string `json:"name,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	TickerId int32 `json:"tickerId,omitempty"`
}