/*
 * Nsoraf_SOR
 *
 * Nsoraf Steering Of Roaming Service. © 2021, 3GPP Organizational Partners (ARIB, ATIS, CCSA, ETSI, TSDSI, TTA, TTC). All rights reserved.
 *
 * API version: 1.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package datarepository

import ()

type ProblemDetails struct {
	// simple type

	Type *string `json:"type,omitempty" yaml:"type" bson:"type" mapstructure:"Type"`

	Title *string `json:"title,omitempty" yaml:"title" bson:"title" mapstructure:"Title"`

	Status *int32 `json:"status,omitempty" yaml:"status" bson:"status" mapstructure:"Status"`

	Detail *string `json:"detail,omitempty" yaml:"detail" bson:"detail" mapstructure:"Detail"`

	Instance *string `json:"instance,omitempty" yaml:"instance" bson:"instance" mapstructure:"Instance"`

	Cause *string `json:"cause,omitempty" yaml:"cause" bson:"cause" mapstructure:"Cause"`

	InvalidParams *[]InvalidParam `json:"invalidParams,omitempty" yaml:"invalidParams" bson:"invalidParams" mapstructure:"InvalidParams"`

	SupportedFeatures *string `json:"supportedFeatures,omitempty" yaml:"supportedFeatures" bson:"supportedFeatures" mapstructure:"SupportedFeatures"`

	AccessTokenError *AccessTokenErr `json:"accessTokenError,omitempty" yaml:"accessTokenError" bson:"accessTokenError" mapstructure:"AccessTokenError"`

	AccessTokenRequest *AccessTokenReq `json:"accessTokenRequest,omitempty" yaml:"accessTokenRequest" bson:"accessTokenRequest" mapstructure:"AccessTokenRequest"`

	NrfId *string `json:"nrfId,omitempty" yaml:"nrfId" bson:"nrfId" mapstructure:"NrfId"`
}
