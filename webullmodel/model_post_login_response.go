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
// PostLoginResponse struct for PostLoginResponse
type PostLoginResponse struct {
	AccessToken string `json:"accessToken,omitempty"`
	ExtInfo PostLoginResponseExtInfo `json:"extInfo,omitempty"`
	FirstTimeOfThird bool `json:"firstTimeOfThird,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	RegisterAddress int32 `json:"registerAddress,omitempty"`
	Settings PostLoginResponseSettings `json:"settings,omitempty"`
	TokenExpireTime string `json:"tokenExpireTime,omitempty"`
	Uuid string `json:"uuid,omitempty"`
}