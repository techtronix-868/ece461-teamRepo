/*
 * ECE 461 - Fall 2021 - Project 2
 *
 * API for ECE 461/Fall 2021/Project 2: A Trustworthy Module Registry
 *
 * API version: 2.0.0
 * Contact: davisjam@purdue.edu
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// UserAuthenticationInfo - Authentication info for a user
type UserAuthenticationInfo struct {

	// Password for a user. Per the spec, this should be a \"strong\" password.
	Password string `json:"password"`
}
