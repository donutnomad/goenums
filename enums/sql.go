package enums

import (
	"fmt"
	"math"
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
		return nil
	}

	// 获取源值的反射值
	v := reflect.ValueOf(src)

	// 如果是指针，获取其元素
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// 获取目标类型
	var target T
	targetType := reflect.TypeOf(target)

	// 根据目标类型的种类进行处理
	switch targetType.Kind() {
	case reflect.String:
		return s.scanString(src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return s.scanInt(src)
	case reflect.Float32, reflect.Float64:
		return s.scanFloat(src)
	case reflect.Bool:
		return s.scanBool(src)
	default:
		// 处理特殊类型
		if targetType == reflect.TypeOf(time.Time{}) {
			return s.scanTime(src)
		}
		return fmt.Errorf("unsupported target type: %v", targetType)
	}
}

func (s *GenericScanner[T]) scanString(src any) error {
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	var str string
	switch v.Kind() {
	case reflect.String:
		str = v.String()
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			str = string(v.Bytes())
		} else {
			return fmt.Errorf("cannot convert slice of %v to string", v.Type().Elem())
		}
	default:
		return fmt.Errorf("cannot convert %v to string", v.Type())
	}

	reflect.ValueOf(s.value).Elem().Set(reflect.ValueOf(str))
	return nil
}

func (s *GenericScanner[T]) scanInt(src any) error {
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	var i int64
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := v.Uint()
		if u > math.MaxInt64 {
			return fmt.Errorf("unsigned integer value %d overflows int64", u)
		}
		i = int64(u)
	case reflect.Float32, reflect.Float64:
		i = int64(v.Float())
	case reflect.Bool:
		if v.Bool() {
			i = 1
		}
	case reflect.String:
		var err error
		i, err = strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse integer from string: %v", err)
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			var err error
			i, err = strconv.ParseInt(string(v.Bytes()), 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse integer from bytes: %v", err)
			}
		} else {
			return fmt.Errorf("cannot convert slice of %v to integer", v.Type().Elem())
		}
	default:
		return fmt.Errorf("cannot convert %v to integer", v.Type())
	}

	// 获取目标值的反射值
	targetValue := reflect.ValueOf(s.value).Elem()

	// 根据目标类型设置值
	switch targetValue.Kind() {
	case reflect.Int:
		if i > math.MaxInt || i < math.MinInt {
			return fmt.Errorf("value %d overflows int", i)
		}
		targetValue.SetInt(i)
	case reflect.Int8:
		if i > math.MaxInt8 || i < math.MinInt8 {
			return fmt.Errorf("value %d overflows int8", i)
		}
		targetValue.SetInt(i)
	case reflect.Int16:
		if i > math.MaxInt16 || i < math.MinInt16 {
			return fmt.Errorf("value %d overflows int16", i)
		}
		targetValue.SetInt(i)
	case reflect.Int32:
		if i > math.MaxInt32 || i < math.MinInt32 {
			return fmt.Errorf("value %d overflows int32", i)
		}
		targetValue.SetInt(i)
	case reflect.Int64:
		targetValue.SetInt(i)
	case reflect.Uint:
		if i < 0 || uint64(i) > math.MaxUint {
			return fmt.Errorf("value %d cannot be converted to uint", i)
		}
		targetValue.SetUint(uint64(i))
	case reflect.Uint8:
		if i < 0 || i > math.MaxUint8 {
			return fmt.Errorf("value %d overflows uint8", i)
		}
		targetValue.SetUint(uint64(i))
	case reflect.Uint16:
		if i < 0 || i > math.MaxUint16 {
			return fmt.Errorf("value %d overflows uint16", i)
		}
		targetValue.SetUint(uint64(i))
	case reflect.Uint32:
		if i < 0 || i > math.MaxUint32 {
			return fmt.Errorf("value %d overflows uint32", i)
		}
		targetValue.SetUint(uint64(i))
	case reflect.Uint64:
		if i < 0 {
			return fmt.Errorf("negative value %d cannot be converted to uint64", i)
		}
		targetValue.SetUint(uint64(i))
	default:
		return fmt.Errorf("unsupported integer type: %v", targetValue.Type())
	}

	return nil
}

func (s *GenericScanner[T]) scanFloat(src any) error {
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	var f float64
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		f = v.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f = float64(v.Uint())
	case reflect.String:
		var err error
		f, err = strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return fmt.Errorf("failed to parse float from string: %v", err)
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			var err error
			f, err = strconv.ParseFloat(string(v.Bytes()), 64)
			if err != nil {
				return fmt.Errorf("failed to parse float from bytes: %v", err)
			}
		} else {
			return fmt.Errorf("cannot convert slice of %v to float", v.Type().Elem())
		}
	default:
		return fmt.Errorf("cannot convert %v to float", v.Type())
	}

	targetValue := reflect.ValueOf(s.value).Elem()
	switch targetValue.Kind() {
	case reflect.Float32:
		if f > math.MaxFloat32 || f < -math.MaxFloat32 {
			return fmt.Errorf("value %f overflows float32", f)
		}
		targetValue.SetFloat(f)
	case reflect.Float64:
		targetValue.SetFloat(f)
	default:
		return fmt.Errorf("unsupported float type: %v", targetValue.Type())
	}

	return nil
}

func (s *GenericScanner[T]) scanBool(src any) error {
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	var b bool
	switch v.Kind() {
	case reflect.Bool:
		b = v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b = v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b = v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		b = v.Float() != 0
	case reflect.String:
		var err error
		b, err = strconv.ParseBool(v.String())
		if err != nil {
			return fmt.Errorf("failed to parse bool from string: %v", err)
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			var err error
			b, err = strconv.ParseBool(string(v.Bytes()))
			if err != nil {
				return fmt.Errorf("failed to parse bool from bytes: %v", err)
			}
		} else {
			return fmt.Errorf("cannot convert slice of %v to bool", v.Type().Elem())
		}
	default:
		return fmt.Errorf("cannot convert %v to bool", v.Type())
	}

	reflect.ValueOf(s.value).Elem().SetBool(b)
	return nil
}

func (s *GenericScanner[T]) scanTime(src any) error {
	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	var t time.Time
	switch {
	case v.Type() == reflect.TypeOf(time.Time{}):
		t = v.Interface().(time.Time)
	case v.Kind() == reflect.String:
		var err error
		t, err = time.Parse(time.RFC3339Nano, v.String())
		if err != nil {
			return fmt.Errorf("failed to parse time from string: %v", err)
		}
	case v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8:
		var err error
		t, err = time.Parse(time.RFC3339Nano, string(v.Bytes()))
		if err != nil {
			return fmt.Errorf("failed to parse time from bytes: %v", err)
		}
	default:
		return fmt.Errorf("cannot convert %v to time.Time", v.Type())
	}

	reflect.ValueOf(s.value).Elem().Set(reflect.ValueOf(t))
	return nil
}
