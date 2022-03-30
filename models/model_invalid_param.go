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

type InvalidParam struct {
	// simple type

	Param string `json:"param" yaml:"param" bson:"param" mapstructure:"Param"`

	Reason *string `json:"reason,omitempty" yaml:"reason" bson:"reason" mapstructure:"Reason"`
}
