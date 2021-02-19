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
// OrderStatus the model 'OrderStatus'
type OrderStatus string

// List of OrderStatus
const (
	QUEUED OrderStatus = "Queued"
	UNCONFIRMED OrderStatus = "Unconfirmed"
	CONFIRMED OrderStatus = "Confirmed"
	PARTIALLY_FILLED OrderStatus = "Partially Filled"
	FILLED OrderStatus = "Filled"
	REJECTED OrderStatus = "Rejected"
	CANCELLED OrderStatus = "Cancelled"
	FAILED OrderStatus = "Failed"
	WORKING OrderStatus = "Working"
)