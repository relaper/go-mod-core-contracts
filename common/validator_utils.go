package common

// 代码拷贝自 github.com/go-playground/validator/util.go

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	namespaceSeparator = "."
	leftBracket        = "["
	rightBracket       = "]"
)

var (
	timeType = reflect.TypeOf(time.Time{})
)

// getStructFieldOKInternal traverses a struct to retrieve a specific field denoted by the provided namespace and
// returns the field, field kind and whether is was successful in retrieving the field at all.
//
// NOTE: when not successful ok will be false, this can happen when a nested struct is nil and so the field
// could not be retrieved because it didn't exist.
func getStructFieldOKInternal(val reflect.Value, namespace string) (current reflect.Value, kind reflect.Kind, nullable bool, found bool) {

BEGIN:
	current, kind, nullable = extractTypeInternal(val, false)
	if kind == reflect.Invalid {
		return
	}

	if namespace == "" {
		found = true
		return
	}

	switch kind {

	case reflect.Ptr, reflect.Interface:
		return

	case reflect.Struct:

		typ := current.Type()
		fld := namespace
		var ns string

		if typ != timeType {

			idx := strings.Index(namespace, namespaceSeparator)

			if idx != -1 {
				fld = namespace[:idx]
				ns = namespace[idx+1:]
			} else {
				ns = ""
			}

			bracketIdx := strings.Index(fld, leftBracket)
			if bracketIdx != -1 {
				fld = fld[:bracketIdx]

				ns = namespace[bracketIdx:]
			}

			val = current.FieldByName(fld)
			namespace = ns
			goto BEGIN
		}

	case reflect.Array, reflect.Slice:
		idx := strings.Index(namespace, leftBracket)
		idx2 := strings.Index(namespace, rightBracket)

		arrIdx, _ := strconv.Atoi(namespace[idx+1 : idx2])

		if arrIdx >= current.Len() {
			return
		}

		startIdx := idx2 + 1

		if startIdx < len(namespace) {
			if namespace[startIdx:startIdx+1] == namespaceSeparator {
				startIdx++
			}
		}

		val = current.Index(arrIdx)
		namespace = namespace[startIdx:]
		goto BEGIN

	case reflect.Map:
		idx := strings.Index(namespace, leftBracket) + 1
		idx2 := strings.Index(namespace, rightBracket)

		endIdx := idx2

		if endIdx+1 < len(namespace) {
			if namespace[endIdx+1:endIdx+2] == namespaceSeparator {
				endIdx++
			}
		}

		key := namespace[idx:idx2]

		switch current.Type().Key().Kind() {
		case reflect.Int:
			i, _ := strconv.Atoi(key)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Int8:
			i, _ := strconv.ParseInt(key, 10, 8)
			val = current.MapIndex(reflect.ValueOf(int8(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int16:
			i, _ := strconv.ParseInt(key, 10, 16)
			val = current.MapIndex(reflect.ValueOf(int16(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int32:
			i, _ := strconv.ParseInt(key, 10, 32)
			val = current.MapIndex(reflect.ValueOf(int32(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int64:
			i, _ := strconv.ParseInt(key, 10, 64)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Uint:
			i, _ := strconv.ParseUint(key, 10, 0)
			val = current.MapIndex(reflect.ValueOf(uint(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint8:
			i, _ := strconv.ParseUint(key, 10, 8)
			val = current.MapIndex(reflect.ValueOf(uint8(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint16:
			i, _ := strconv.ParseUint(key, 10, 16)
			val = current.MapIndex(reflect.ValueOf(uint16(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint32:
			i, _ := strconv.ParseUint(key, 10, 32)
			val = current.MapIndex(reflect.ValueOf(uint32(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint64:
			i, _ := strconv.ParseUint(key, 10, 64)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Float32:
			f, _ := strconv.ParseFloat(key, 32)
			val = current.MapIndex(reflect.ValueOf(float32(f)))
			namespace = namespace[endIdx+1:]

		case reflect.Float64:
			f, _ := strconv.ParseFloat(key, 64)
			val = current.MapIndex(reflect.ValueOf(f))
			namespace = namespace[endIdx+1:]

		case reflect.Bool:
			b, _ := strconv.ParseBool(key)
			val = current.MapIndex(reflect.ValueOf(b))
			namespace = namespace[endIdx+1:]

		// reflect.Type = string
		default:
			val = current.MapIndex(reflect.ValueOf(key))
			namespace = namespace[endIdx+1:]
		}

		goto BEGIN
	}

	// if got here there was more namespace, cannot go any deeper
	panic("Invalid field namespace")
}

// extractTypeInternal gets the actual underlying type of field value.
// It will dive into pointers, customTypes and return you the
// underlying value and it's kind.
func extractTypeInternal(current reflect.Value, nullable bool) (reflect.Value, reflect.Kind, bool) {

BEGIN:
	switch current.Kind() {
	case reflect.Ptr:

		nullable = true

		if current.IsNil() {
			return current, reflect.Ptr, nullable
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:

		nullable = true

		if current.IsNil() {
			return current, reflect.Interface, nullable
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current, reflect.Invalid, nullable

	default:

		return current, current.Kind(), nullable
	}
}
