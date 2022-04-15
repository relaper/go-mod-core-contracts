//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

// DeviceCoreCommand and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command/2.1.0#/DeviceCoreCommand
type DeviceCoreCommand struct {
	DeviceName   string        `json:"deviceName" validate:"required,edgex-dto-rfc3986-unreserved-chars" validate_name:"设备"`
	ProfileName  string        `json:"profileName" validate:"required,edgex-dto-rfc3986-unreserved-chars" validate_name:"模型"`
	CoreCommands []CoreCommand `json:"coreCommands,omitempty" validate:"dive" validate_name:"命令列表"`
}

// CoreCommand and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command/2.1.0#/CoreCommand
type CoreCommand struct {
	Name        string                 `json:"name" validate:"required,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars" validate_name:"命令标识"`
	Description string                 `json:"description" yaml:"description"`
	Get         bool                   `json:"get,omitempty" validate:"required_without=Set" validate_name:"是否读"`
	Set         bool                   `json:"set,omitempty" validate:"required_without=Get" validate_name:"是否写"`
	Path        string                 `json:"path,omitempty"`
	Url         string                 `json:"url,omitempty"`
	Parameters  []CoreCommandParameter `json:"parameters,omitempty"`
}

// CoreCommandParameter and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-command/2.1.0#/CoreCommandParameter
type CoreCommandParameter struct {
	ResourceName        string `json:"resourceName"`
	ResourceDescription string `json:"resourceDescription"`
	ValueType           string `json:"valueType"`
}
