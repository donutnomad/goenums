package enums

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// GenericScanner is a generic Scanner implementation
type GenericScanner[T any] struct {
	value *T
}

func NewScanner[T any](value *T) *GenericScanner[T] {
	return &GenericScanner[T]{value: value}
}

func (s *GenericScanner[T]) Scan(src any) error {
	if src == nil {
		return nil // or return specific error based on requirements
	}

	var target T
	switch any(target).(type) {
	case string:
		return s.scanString(src)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return s.scanInt(src)
	case float32, float64:
		return s.scanFloat(src)
	case bool:
		return s.scanBool(src)
	case time.Time:
		return s.scanTime(src)
	default:
		return fmt.Errorf("unsupported target type: %T", target)
	}
}

func (s *GenericScanner[T]) scanString(src any) error {
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot convert %T to string", src)
	}

	*s.value = any(str).(T)
	return nil
}

func (s *GenericScanner[T]) scanInt(src any) error {
	var i int64
	var err error

	switch v := src.(type) {
	case int64:
		i = v
	case int:
		i = int64(v)
	case int32:
		i = int64(v)
	case int16:
		i = int64(v)
	case int8:
		i = int64(v)
	case uint, uint8, uint16, uint32, uint64:
		// Use reflection to handle unsigned integers
		rv := reflect.ValueOf(v)
		i = int64(rv.Uint())
	case float64:
		i = int64(v)
	case float32:
		i = int64(v)
	case []byte:
		i, err = strconv.ParseInt(string(v), 10, 64)
	case string:
		i, err = strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			i = 1
		} else {
			i = 0
		}
	default:
		return fmt.Errorf("cannot convert %T to integer", src)
	}

	if err != nil {
		return fmt.Errorf("failed to parse integer: %v", err)
	}

	// Convert based on target type
	switch ptr := any(s.value).(type) {
	case *int:
		*ptr = int(i)
	case *int8:
		*ptr = int8(i)
	case *int16:
		*ptr = int16(i)
	case *int32:
		*ptr = int32(i)
	case *int64:
		*ptr = i
	case *uint:
		*ptr = uint(i)
	case *uint8:
		*ptr = uint8(i)
	case *uint16:
		*ptr = uint16(i)
	case *uint32:
		*ptr = uint32(i)
	case *uint64:
		*ptr = uint64(i)
	default:
		return fmt.Errorf("unsupported integer type: %T", s.value)
	}

	return nil
}

func (s *GenericScanner[T]) scanFloat(src any) error {
	var f float64
	var err error

	switch v := src.(type) {
	case float64:
		f = v
	case float32:
		f = float64(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// Use reflection to handle all integer types
		rv := reflect.ValueOf(v)
		f = float64(rv.Int())
	case []byte:
		f, err = strconv.ParseFloat(string(v), 64)
	case string:
		f, err = strconv.ParseFloat(v, 64)
	default:
		return fmt.Errorf("cannot convert %T to float", src)
	}

	if err != nil {
		return fmt.Errorf("failed to parse float: %v", err)
	}

	// Convert based on target type
	switch ptr := any(s.value).(type) {
	case *float32:
		*ptr = float32(f)
	case *float64:
		*ptr = f
	default:
		return fmt.Errorf("unsupported float type: %T", s.value)
	}

	return nil
}

func (s *GenericScanner[T]) scanBool(src any) error {
	var b bool
	var err error

	switch v := src.(type) {
	case bool:
		b = v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// All non-zero integers converted to true
		rv := reflect.ValueOf(v)
		b = rv.Int() != 0
	case float32, float64:
		// All non-zero floating point numbers converted to true
		rv := reflect.ValueOf(v)
		b = rv.Float() != 0
	case []byte:
		b, err = strconv.ParseBool(string(v))
	case string:
		b, err = strconv.ParseBool(v)
	default:
		return fmt.Errorf("cannot convert %T to bool", src)
	}

	if err != nil {
		return fmt.Errorf("failed to parse bool: %v", err)
	}

	*s.value = any(b).(T)
	return nil
}

func (s *GenericScanner[T]) scanTime(src any) error {
	var t time.Time
	var err error

	switch v := src.(type) {
	case time.Time:
		t = v
	case []byte:
		t, err = time.Parse(time.RFC3339Nano, string(v))
	case string:
		t, err = time.Parse(time.RFC3339Nano, v)
	default:
		return fmt.Errorf("cannot convert %T to time.Time", src)
	}

	if err != nil {
		return fmt.Errorf("failed to parse time: %v", err)
	}

	*s.value = any(t).(T)
	return nil
}
