//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

// Interval and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-scheduler/2.1.0#/Interval
type Interval struct {
	DBTimestamp `json:",inline"`
	Id          string `json:"id,omitempty" validate:"omitempty,uuid" validate_name:"间隔ID"`
	Name        string `json:"name" validate:"edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars" validate_name:"间隔名称"`
	Start       string `json:"start,omitempty" validate:"omitempty,edgex-dto-interval-datetime" validate_name:"开始时间"`
	End         string `json:"end,omitempty" validate:"omitempty,edgex-dto-interval-datetime" validate_name:"结束时间"`
	Interval    string `json:"interval" validate:"required,edgex-dto-duration" validate_name:"间隔"`
}

// NewInterval creates interval DTO with required fields
func NewInterval(name, interval string) Interval {
	return Interval{Name: name, Interval: interval}
}

// UpdateInterval and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-scheduler/2.1.0#/UpdateInterval
type UpdateInterval struct {
	Id       *string `json:"id" validate:"required_without=Name,edgex-dto-uuid" validate_name:"间隔ID"`
	Name     *string `json:"name" validate:"required_without=Id,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars" validate_name:"间隔名称"`
	Start    *string `json:"start" validate:"omitempty,edgex-dto-interval-datetime" validate_name:"开始时间"`
	End      *string `json:"end" validate:"omitempty,edgex-dto-interval-datetime" validate_name:"结束时间"`
	Interval *string `json:"interval" validate:"omitempty,edgex-dto-duration" validate_name:"间隔"`
}

// NewUpdateInterval creates updateInterval DTO with required field
func NewUpdateInterval(name string) UpdateInterval {
	return UpdateInterval{Name: &name}
}

// ToIntervalModel transforms the Interval DTO to the Interval Model
func ToIntervalModel(dto Interval) models.Interval {
	var model models.Interval
	model.Id = dto.Id
	model.Name = dto.Name
	model.Start = dto.Start
	model.End = dto.End
	model.Interval = dto.Interval
	return model
}

// FromIntervalModelToDTO transforms the Interval Model to the Interval DTO
func FromIntervalModelToDTO(model models.Interval) Interval {
	var dto Interval
	dto.DBTimestamp = DBTimestamp(model.DBTimestamp)
	dto.Id = model.Id
	dto.Name = model.Name
	dto.Start = model.Start
	dto.End = model.End
	dto.Interval = model.Interval
	return dto
}
