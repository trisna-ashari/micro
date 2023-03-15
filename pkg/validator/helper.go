package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// JoinSlice is a function uses to join slice of interfaces.
// It's used by RuleOpt to construct RuleOpt.Value.
func JoinSlice(sliceOfInterface []interface{}) string {
	var sliceOfString []string
	for _, slice := range sliceOfInterface {
		v := reflect.ValueOf(slice)
		switch v.Kind() {
		case reflect.Bool:
			sliceOfString = append(sliceOfString, strconv.FormatBool(slice.(bool)))
		case reflect.String:
			sliceOfString = append(sliceOfString, slice.(string))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sliceOfString = append(sliceOfString, fmt.Sprintf("%d", slice))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			sliceOfString = append(sliceOfString, fmt.Sprintf("%d", slice))
		case reflect.Float32, reflect.Float64:
			sliceOfString = append(sliceOfString, fmt.Sprintf("%v", slice))
		default:

		}
	}

	return strings.Join(sliceOfString, "/")
}

// JoinSliceWithSeparator is a function uses to join slice of interfaces with separator.
// It's used by RuleOpt to construct RuleOpt.Value.
func JoinSliceWithSeparator(sliceOfInterface []interface{}, separator string) string {
	var sliceOfString []string
	for _, slice := range sliceOfInterface {
		v := reflect.ValueOf(slice)
		switch v.Kind() {
		case reflect.Bool:
			sliceOfString = append(sliceOfString, strconv.FormatBool(slice.(bool)))
		case reflect.String:
			sliceOfString = append(sliceOfString, slice.(string))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sliceOfString = append(sliceOfString, fmt.Sprintf("%d", slice))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			sliceOfString = append(sliceOfString, fmt.Sprintf("%d", slice))
		case reflect.Float32, reflect.Float64:
			sliceOfString = append(sliceOfString, fmt.Sprintf("%v", slice))
		default:

		}
	}

	return strings.Join(sliceOfString, separator)
}

// EqualValue is a closure uses by ValidationRules.EqualTo rule.
func EqualValue(targetValue string) validation.RuleFunc {
	return func(fieldValue interface{}) error {
		s, _ := fieldValue.(string)
		if s != targetValue {
			return errors.New("validation.error.must_be_equal_to")
		}
		return nil
	}
}

// Required is a closure used by ValidationRules.Required rule.
func Required() validation.RuleFunc {
	return func(fieldValue interface{}) error {
		value, isNil := validation.Indirect(fieldValue)
		if isNil && IsEmpty(value) || (isNil || IsEmpty(value)) {
			return errors.New("validation.error.is_required")
		}

		return nil
	}
}

// IsEmpty is a function uses to check if a value is empty or not.
// A value is considered empty if
// - integer, float: zero
// - string, array: len() == 0
// - slice, map: nil or len() == 0
// - interface, pointer: nil or the referenced value is empty
func IsEmpty(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Invalid:
		return true
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	case reflect.Struct:
		v, ok := value.(time.Time)
		if ok && v.IsZero() {
			return true
		}
	}

	return false
}
