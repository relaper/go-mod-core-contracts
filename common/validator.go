//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
)

var val *validator.Validate

const (
	dtoDurationTag              = "edgex-dto-duration"
	dtoUuidTag                  = "edgex-dto-uuid"
	dtoNoneEmptyStringTag       = "edgex-dto-none-empty-string"
	dtoValueType                = "edgex-dto-value-type"
	dtoRFC3986UnreservedCharTag = "edgex-dto-rfc3986-unreserved-chars"
	dtoInterDatetimeTag         = "edgex-dto-interval-datetime"
)

const (
	// Per https://tools.ietf.org/html/rfc3986#section-2.3, unreserved characters= ALPHA / DIGIT / "-" / "." / "_" / "~"
	// Also due to names used in topics for Redis Pub/Sub, "."are not allowed
	rFC3986UnreservedCharsRegexString = "^[a-zA-Z0-9-_~]+$"
	intervalDatetimeLayout            = "20060102T150405"
	name                              = "Name"
)

var (
	rFC3986UnreservedCharsRegex = regexp.MustCompile(rFC3986UnreservedCharsRegexString)
)

func init() {
	val = validator.New()
	val.RegisterValidation(dtoDurationTag, ValidateDuration)
	val.RegisterValidation(dtoUuidTag, ValidateDtoUuid)
	val.RegisterValidation(dtoNoneEmptyStringTag, ValidateDtoNoneEmptyString)
	val.RegisterValidation(dtoValueType, ValidateValueType)
	val.RegisterValidation(dtoRFC3986UnreservedCharTag, ValidateDtoRFC3986UnreservedChars)
	val.RegisterValidation(dtoInterDatetimeTag, ValidateIntervalDatetime)
}

// Validate function will use the validator package to validate the struct annotation
func Validate(a interface{}) error {
	aVal := reflect.ValueOf(a)
	err := val.Struct(a)

	// translate all error at once
	if err != nil {
		errs := err.(validator.ValidationErrors)
		var errMsg []string
		for _, e := range errs {
			errMsg = append(errMsg, getErrorMessage(aVal, e))
		}
		return errors.NewCommonEdgeX(errors.KindContractInvalid, strings.Join(errMsg, "; "), nil)
	}
	return nil
}

// Internal: generate representative validation error messages
func getErrorMessage(aVal reflect.Value, e validator.FieldError) string {
	ns := e.StructNamespace()
	f := strings.Index(ns, ".")
	l := strings.LastIndex(ns, ".")

	fieldName := e.StructNamespace()
	var field *reflect.StructField
	if f == l {
		if sd, ok := aVal.Type().FieldByName(e.StructField()); ok {
			field = &sd
		}
	} else {
		parentNs := ns[f+1 : l]

		if parent, _, _, ok := getStructFieldOKInternal(aVal, parentNs); ok {
			if parent.Kind() == reflect.Struct {
				if sd, ok := parent.Type().FieldByName(e.StructField()); ok {
					field = &sd
				}
			}
		}
	}
	if field != nil {
		if validateName := field.Tag.Get("validate_name"); validateName != "" {
			fieldName = validateName
		}
	}

	tag := e.Tag()
	// StructNamespace returns the namespace for the field error, with the field's actual name.
	fieldValue := e.Param()
	var msg string
	switch tag {
	case "uuid":
		msg = fmt.Sprintf("%s须为UUID", fieldName)
	case "required":
		msg = fmt.Sprintf("%s为必填项", fieldName)
	case "required_without":
		msg = fmt.Sprintf("%s为必填项（当%s为空时）", fieldName, fieldValue)
	case "len":
		msg = fmt.Sprintf("%s长度不为%s", fieldName, fieldValue)
	case "oneof":
		msg = fmt.Sprintf("%s需为%s", fieldName, fieldValue)
	case "gt":
		msg = fmt.Sprintf("%s需大于%s", fieldName, fieldValue)
	case dtoDurationTag:
		msg = fmt.Sprintf("%s不是合法时间间隔格式（ISO 8601 Durations：100ms, 24h）", fieldName)
	case dtoUuidTag:
		msg = fmt.Sprintf("%s须为UUID", fieldName)
	case dtoNoneEmptyStringTag:
		msg = fmt.Sprintf("%s不能为空", fieldName)
	case dtoRFC3986UnreservedCharTag:
		msg = fmt.Sprintf("%s需由字母、数字、_、-或~组成", fieldName)
	default:
		msg = fmt.Sprintf("%s参数验证失败：%s", fieldName, tag)
	}
	return msg
}

// ValidateDuration validate field which should follow the ISO 8601 Durations format
func ValidateDuration(fl validator.FieldLevel) bool {
	_, err := time.ParseDuration(fl.Field().String())
	return err == nil
}

// ValidateDtoUuid used to check the UpdateDTO uuid pointer value
// Currently, required_without can not correct work with other tag, so write custom tag instead.
// Issue can refer to https://github.com/go-playground/validator/issues/624
func ValidateDtoUuid(fl validator.FieldLevel) bool {
	idField := fl.Field()
	// Skip the validation if the pointer value is nil
	if isNilPointer(idField) {
		return true
	}

	// The Id field should accept the empty string if the Name field is provided
	nameField := fl.Parent().FieldByName(name)
	if len(strings.TrimSpace(idField.String())) == 0 && !isNilPointer(nameField) && len(nameField.Elem().String()) > 0 {
		return true
	}

	_, err := uuid.Parse(idField.String())
	return err == nil
}

// ValidateDtoNoneEmptyString used to check the UpdateDTO name pointer value
func ValidateDtoNoneEmptyString(fl validator.FieldLevel) bool {
	val := fl.Field()
	// Skip the validation if the pointer value is nil
	if isNilPointer(val) {
		return true
	}
	// The string value should not be empty
	if len(strings.TrimSpace(val.String())) > 0 {
		return true
	} else {
		return false
	}
}

// ValidateValueType checks whether the valueType is valid
func ValidateValueType(fl validator.FieldLevel) bool {
	valueType := fl.Field().String()
	for _, v := range valueTypes {
		if strings.ToLower(valueType) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

// ValidateDtoRFC3986UnreservedChars used to check if DTO's name pointer value only contains unreserved characters as
// defined in https://tools.ietf.org/html/rfc3986#section-2.3
func ValidateDtoRFC3986UnreservedChars(fl validator.FieldLevel) bool {
	val := fl.Field()
	// Skip the validation if the pointer value is nil
	if isNilPointer(val) {
		return true
	} else {
		return rFC3986UnreservedCharsRegex.MatchString(val.String())
	}
}

// ValidateIntervalDatetime validate Interval's datetime field which should follow the ISO 8601 format YYYYMMDD'T'HHmmss
func ValidateIntervalDatetime(fl validator.FieldLevel) bool {
	_, err := time.Parse(intervalDatetimeLayout, fl.Field().String())
	return err == nil
}

func isNilPointer(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr && value.IsNil()
}
