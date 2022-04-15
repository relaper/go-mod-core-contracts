//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
)

// BaseReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.1.0#/BaseReading
type BaseReading struct {
	Id            string `json:"id,omitempty"`
	Origin        int64  `json:"origin" validate:"required" validate_name:"原始时间"`
	DeviceName    string `json:"deviceName" validate:"required,edgex-dto-rfc3986-unreserved-chars" validate_name:"设备标识"`
	ResourceName  string `json:"resourceName" validate:"required,edgex-dto-rfc3986-unreserved-chars" validate_name:"属性标识"`
	ProfileName   string `json:"profileName" validate:"required,edgex-dto-rfc3986-unreserved-chars" validate_name:"模型标识"`
	ValueType     string `json:"valueType" validate:"required,edgex-dto-value-type" validate_name:"数据类型"`
	BinaryReading `json:",inline" validate:"-"`
	SimpleReading `json:",inline" validate:"-"`
	ObjectReading `json:",inline" validate:"-"`
}

// SimpleReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.1.0#/SimpleReading
type SimpleReading struct {
	Value string `json:"value,omitempty" validate:"required" validate_name:"数值"`
}

// BinaryReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.1.0#/BinaryReading
type BinaryReading struct {
	BinaryValue []byte `json:"binaryValue,omitempty" validate:"gt=0,required" validate_name:"二进制数值"`
	MediaType   string `json:"mediaType,omitempty" validate:"required" validate_name:"媒体类型"`
}

// ObjectReading and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.1.0#/ObjectReading
type ObjectReading struct {
	ObjectValue interface{} `json:"objectValue,omitempty" validate:"required" validate_name:"对象数值"`
}

func newBaseReading(profileName string, deviceName string, resourceName string, valueType string) BaseReading {
	return BaseReading{
		Id:           uuid.NewString(),
		Origin:       time.Now().UnixNano(),
		DeviceName:   deviceName,
		ResourceName: resourceName,
		ProfileName:  profileName,
		ValueType:    valueType,
	}
}

// NewSimpleReading creates and returns a new initialized BaseReading with its SimpleReading initialized
func NewSimpleReading(profileName string, deviceName string, resourceName string, valueType string, value interface{}) (BaseReading, error) {
	stringValue, err := convertInterfaceValue(valueType, value)
	if err != nil {
		return BaseReading{}, err
	}

	reading := newBaseReading(profileName, deviceName, resourceName, valueType)
	reading.SimpleReading = SimpleReading{
		Value: stringValue,
	}
	return reading, nil
}

// NewBinaryReading creates and returns a new initialized BaseReading with its BinaryReading initialized
func NewBinaryReading(profileName string, deviceName string, resourceName string, binaryValue []byte, mediaType string) BaseReading {
	reading := newBaseReading(profileName, deviceName, resourceName, common.ValueTypeBinary)
	reading.BinaryReading = BinaryReading{
		BinaryValue: binaryValue,
		MediaType:   mediaType,
	}
	return reading
}

// NewObjectReading creates and returns a new initialized BaseReading with its ObjectReading initialized
func NewObjectReading(profileName string, deviceName string, resourceName string, objectValue interface{}) BaseReading {
	reading := newBaseReading(profileName, deviceName, resourceName, common.ValueTypeObject)
	reading.ObjectReading = ObjectReading{
		ObjectValue: objectValue,
	}
	return reading
}

func convertInterfaceValue(valueType string, value interface{}) (string, error) {
	switch valueType {
	case common.ValueTypeBool:
		return convertSimpleValue(valueType, reflect.Bool, value)
	case common.ValueTypeString:
		return convertSimpleValue(valueType, reflect.String, value)

	case common.ValueTypeUint8:
		return convertSimpleValue(valueType, reflect.Uint8, value)
	case common.ValueTypeUint16:
		return convertSimpleValue(valueType, reflect.Uint16, value)
	case common.ValueTypeUint32:
		return convertSimpleValue(valueType, reflect.Uint32, value)
	case common.ValueTypeUint64:
		return convertSimpleValue(valueType, reflect.Uint64, value)

	case common.ValueTypeInt8:
		return convertSimpleValue(valueType, reflect.Int8, value)
	case common.ValueTypeInt16:
		return convertSimpleValue(valueType, reflect.Int16, value)
	case common.ValueTypeInt32:
		return convertSimpleValue(valueType, reflect.Int32, value)
	case common.ValueTypeInt64:
		return convertSimpleValue(valueType, reflect.Int64, value)

	case common.ValueTypeFloat32:
		return convertFloatValue(valueType, reflect.Float32, value)
	case common.ValueTypeFloat64:
		return convertFloatValue(valueType, reflect.Float64, value)

	case common.ValueTypeBoolArray:
		return convertSimpleArrayValue(valueType, reflect.Bool, value)
	case common.ValueTypeStringArray:
		return convertSimpleArrayValue(valueType, reflect.String, value)

	case common.ValueTypeUint8Array:
		return convertSimpleArrayValue(valueType, reflect.Uint8, value)
	case common.ValueTypeUint16Array:
		return convertSimpleArrayValue(valueType, reflect.Uint16, value)
	case common.ValueTypeUint32Array:
		return convertSimpleArrayValue(valueType, reflect.Uint32, value)
	case common.ValueTypeUint64Array:
		return convertSimpleArrayValue(valueType, reflect.Uint64, value)

	case common.ValueTypeInt8Array:
		return convertSimpleArrayValue(valueType, reflect.Int8, value)
	case common.ValueTypeInt16Array:
		return convertSimpleArrayValue(valueType, reflect.Int16, value)
	case common.ValueTypeInt32Array:
		return convertSimpleArrayValue(valueType, reflect.Int32, value)
	case common.ValueTypeInt64Array:
		return convertSimpleArrayValue(valueType, reflect.Int64, value)

	case common.ValueTypeFloat32Array:
		arrayValue, ok := value.([]float32)
		if !ok {
			return "", fmt.Errorf("转换数据 %s 为 []float32 失败", valueType)
		}

		return convertFloat32ArrayValue(arrayValue)
	case common.ValueTypeFloat64Array:
		arrayValue, ok := value.([]float64)
		if !ok {
			return "", fmt.Errorf("转换数据 %s 为 []float64 失败", valueType)
		}

		return convertFloat64ArrayValue(arrayValue)

	default:
		return "", fmt.Errorf("数据类型 %s 不合法", valueType)
	}
}

func convertSimpleValue(valueType string, kind reflect.Kind, value interface{}) (string, error) {
	if err := validateType(valueType, kind, value); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", value), nil
}

func convertFloatValue(valueType string, kind reflect.Kind, value interface{}) (string, error) {
	if err := validateType(valueType, kind, value); err != nil {
		return "", err
	}

	return fmt.Sprintf("%g", value), nil
}

func convertSimpleArrayValue(valueType string, kind reflect.Kind, value interface{}) (string, error) {
	if err := validateType(valueType, kind, value); err != nil {
		return "", err
	}

	result := fmt.Sprintf("%v", value)
	result = strings.ReplaceAll(result, " ", ", ")
	return result, nil
}

func convertFloat32ArrayValue(values []float32) (string, error) {
	result := "["
	first := true
	for _, value := range values {
		if first {
			floatValue, err := convertFloatValue(common.ValueTypeFloat32, reflect.Float32, value)
			if err != nil {
				return "", err
			}
			result += floatValue
			first = false
			continue
		}

		floatValue, err := convertFloatValue(common.ValueTypeFloat32, reflect.Float32, value)
		if err != nil {
			return "", err
		}
		result += ", " + floatValue
	}

	result += "]"
	return result, nil
}

func convertFloat64ArrayValue(values []float64) (string, error) {
	result := "["
	first := true
	for _, value := range values {
		if first {
			floatValue, err := convertFloatValue(common.ValueTypeFloat64, reflect.Float64, value)
			if err != nil {
				return "", err
			}
			result += floatValue
			first = false
			continue
		}

		floatValue, err := convertFloatValue(common.ValueTypeFloat64, reflect.Float64, value)
		if err != nil {
			return "", err
		}
		result += ", " + floatValue
	}

	result += "]"
	return result, nil
}

func validateType(valueType string, kind reflect.Kind, value interface{}) error {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		if kind != reflect.TypeOf(value).Elem().Kind() {
			return fmt.Errorf("`%s`切片 不匹配数据类型 '%s", kind.String(), valueType)
		}
		return nil
	}

	if kind != reflect.TypeOf(value).Kind() {
		return fmt.Errorf("`%s` 不匹配数据类型 '%s", kind.String(), valueType)
	}

	return nil
}

// Validate satisfies the Validator interface
func (b BaseReading) Validate() error {
	if b.ValueType == common.ValueTypeBinary {
		// validate the inner BinaryReading struct
		binaryReading := b.BinaryReading
		if err := common.Validate(binaryReading); err != nil {
			return err
		}
	} else if b.ValueType == common.ValueTypeObject {
		// validate the inner ObjectReading struct
		objectReading := b.ObjectReading
		if err := common.Validate(objectReading); err != nil {
			return err
		}
	} else {
		// validate the inner SimpleReading struct
		simpleReading := b.SimpleReading
		if err := common.Validate(simpleReading); err != nil {
			return err
		}
	}

	return nil
}

// Convert Reading DTO to Reading model
func ToReadingModel(r BaseReading) models.Reading {
	var readingModel models.Reading
	br := models.BaseReading{
		Id:           r.Id,
		Origin:       r.Origin,
		DeviceName:   r.DeviceName,
		ResourceName: r.ResourceName,
		ProfileName:  r.ProfileName,
		ValueType:    r.ValueType,
	}
	if r.ValueType == common.ValueTypeBinary {
		readingModel = models.BinaryReading{
			BaseReading: br,
			BinaryValue: r.BinaryValue,
			MediaType:   r.MediaType,
		}
	} else if r.ValueType == common.ValueTypeObject {
		readingModel = models.ObjectReading{
			BaseReading: br,
			ObjectValue: r.ObjectValue,
		}
	} else {
		readingModel = models.SimpleReading{
			BaseReading: br,
			Value:       r.Value,
		}
	}
	return readingModel
}

func FromReadingModelToDTO(reading models.Reading) BaseReading {
	var baseReading BaseReading
	switch r := reading.(type) {
	case models.BinaryReading:
		baseReading = BaseReading{
			Id:            r.Id,
			Origin:        r.Origin,
			DeviceName:    r.DeviceName,
			ResourceName:  r.ResourceName,
			ProfileName:   r.ProfileName,
			ValueType:     r.ValueType,
			BinaryReading: BinaryReading{BinaryValue: r.BinaryValue, MediaType: r.MediaType},
		}
	case models.ObjectReading:
		baseReading = BaseReading{
			Id:            r.Id,
			Origin:        r.Origin,
			DeviceName:    r.DeviceName,
			ResourceName:  r.ResourceName,
			ProfileName:   r.ProfileName,
			ValueType:     r.ValueType,
			ObjectReading: ObjectReading{ObjectValue: r.ObjectValue},
		}
	case models.SimpleReading:
		baseReading = BaseReading{
			Id:            r.Id,
			Origin:        r.Origin,
			DeviceName:    r.DeviceName,
			ResourceName:  r.ResourceName,
			ProfileName:   r.ProfileName,
			ValueType:     r.ValueType,
			SimpleReading: SimpleReading{Value: r.Value},
		}
	}

	return baseReading
}
