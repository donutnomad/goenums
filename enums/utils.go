package enums

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func anyToString(value any) (string, error) {
	if value == nil {
		return "null", nil
	}

	// 使用反射获取值的类型信息
	v := reflect.ValueOf(value)

	// 如果是指针，获取其元素
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "null", nil
		}
		v = v.Elem()
	}

	// 获取底层类型的Kind
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32:
		// 对于 float32，使用 32 位精度格式化，避免 float64 转换带来的精度问题
		f32 := float32(v.Float())
		return strconv.FormatFloat(float64(f32), 'g', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 { // []byte
			return string(v.Bytes()), nil
		}
		fallthrough
	default:
		marshal, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	}
}

// Helper functions for parsing values
func parseStringValue[T any](str string, value *T) error {
	switch v := any(value).(type) {
	case *int:
		parsed, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*v = int(parsed)
	case *int8:
		parsed, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return err
		}
		*v = int8(parsed)
	case *int16:
		parsed, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return err
		}
		*v = int16(parsed)
	case *int32:
		parsed, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return err
		}
		*v = int32(parsed)
	case *int64:
		parsed, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		*v = parsed
	case *uint:
		parsed, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*v = uint(parsed)
	case *uint8:
		parsed, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(parsed)
	case *uint16:
		parsed, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(parsed)
	case *uint32:
		parsed, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(parsed)
	case *uint64:
		parsed, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		*v = parsed
	case *float32:
		parsed, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}
		*v = float32(parsed)
	case *float64:
		parsed, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*v = parsed
	case *bool:
		parsed, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		*v = parsed
	case *string:
		*v = str
	case *[]byte:
		*v = []byte(str)
	default:
		if err := json.Unmarshal([]byte(str), value); err != nil {
			return err
		}
	}
	return nil
}

// anyToBinary 将任意类型转换为二进制格式
// 使用大端字节序（network byte order）作为标准
func anyToBinary(value any) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	v := reflect.ValueOf(value)

	// 如果是指针，获取其元素
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("nil pointer")
		}
		v = v.Elem()
	}

	// 获取底层类型的Kind
	switch v.Kind() {
	case reflect.Int8:
		return []byte{byte(v.Int())}, nil
	case reflect.Int16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(v.Int()))
		return buf, nil
	case reflect.Int32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v.Int()))
		return buf, nil
	case reflect.Int64, reflect.Int:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v.Int()))
		return buf, nil
	case reflect.Uint8:
		return []byte{byte(v.Uint())}, nil
	case reflect.Uint16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(v.Uint()))
		return buf, nil
	case reflect.Uint32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v.Uint()))
		return buf, nil
	case reflect.Uint64, reflect.Uint:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v.Uint())
		return buf, nil
	case reflect.Float32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, math.Float32bits(float32(v.Float())))
		return buf, nil
	case reflect.Float64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(v.Float()))
		return buf, nil
	case reflect.Bool:
		if v.Bool() {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	case reflect.String:
		return []byte(v.String()), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte 类型
			result := make([]byte, v.Len())
			copy(result, v.Bytes())
			return result, nil
		}
		fallthrough
	default:
		// 对于复杂类型，使用JSON序列化
		return json.Marshal(value)
	}
}

// parseBinaryValue 从二进制数据解析为指定类型
func parseBinaryValue[T any](data []byte, value *T) error {
	if len(data) == 0 {
		return fmt.Errorf("empty binary data")
	}

	// 使用反射获取目标类型的信息
	v := reflect.ValueOf(value).Elem()

	// 获取底层类型的Kind
	switch v.Kind() {
	case reflect.Int8:
		if len(data) < 1 {
			return fmt.Errorf("insufficient data for int8")
		}
		v.SetInt(int64(data[0]))
	case reflect.Int16:
		if len(data) < 2 {
			return fmt.Errorf("insufficient data for int16")
		}
		v.SetInt(int64(binary.BigEndian.Uint16(data)))
	case reflect.Int32:
		if len(data) < 4 {
			return fmt.Errorf("insufficient data for int32")
		}
		v.SetInt(int64(binary.BigEndian.Uint32(data)))
	case reflect.Int64, reflect.Int:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for int64")
		}
		v.SetInt(int64(binary.BigEndian.Uint64(data)))
	case reflect.Uint8:
		if len(data) < 1 {
			return fmt.Errorf("insufficient data for uint8")
		}
		v.SetUint(uint64(data[0]))
	case reflect.Uint16:
		if len(data) < 2 {
			return fmt.Errorf("insufficient data for uint16")
		}
		v.SetUint(uint64(binary.BigEndian.Uint16(data)))
	case reflect.Uint32:
		if len(data) < 4 {
			return fmt.Errorf("insufficient data for uint32")
		}
		v.SetUint(uint64(binary.BigEndian.Uint32(data)))
	case reflect.Uint64, reflect.Uint:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for uint64")
		}
		v.SetUint(binary.BigEndian.Uint64(data))
	case reflect.Float32:
		if len(data) < 4 {
			return fmt.Errorf("insufficient data for float32")
		}
		bits := binary.BigEndian.Uint32(data)
		v.SetFloat(float64(math.Float32frombits(bits)))
	case reflect.Float64:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for float64")
		}
		bits := binary.BigEndian.Uint64(data)
		v.SetFloat(math.Float64frombits(bits))
	case reflect.Bool:
		if len(data) < 1 {
			return fmt.Errorf("insufficient data for bool")
		}
		v.SetBool(data[0] != 0)
	case reflect.String:
		v.SetString(string(data))
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte 类型
			newSlice := reflect.MakeSlice(v.Type(), len(data), len(data))
			reflect.Copy(newSlice, reflect.ValueOf(data))
			v.Set(newSlice)
			return nil
		}
		fallthrough
	default:
		return json.Unmarshal(data, value)
	}
	return nil
}

// convertToTargetType converts an interface{} value to the target type
func convertToTargetType[R comparable](value any, target *R) error {
	if value == nil {
		return fmt.Errorf("cannot convert nil value")
	}

	// 使用反射获取目标类型的信息
	v := reflect.ValueOf(target).Elem()

	// 获取底层类型的Kind
	switch v.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			v.SetString(str)
		} else {
			str, err := anyToString(value)
			if err != nil {
				return err
			}
			v.SetString(str)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := v.Int()
		if v.OverflowInt(val) {
			return fmt.Errorf("value %v overflows %s", val, v.Type())
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := v.Uint()
		uval := uint64(val)
		if v.OverflowUint(uval) {
			return fmt.Errorf("value %v overflows %s", val, v.Type())
		}
		v.SetUint(uval)
	case reflect.Float32, reflect.Float64:
		val := v.Float()
		if v.OverflowFloat(val) {
			return fmt.Errorf("value %v overflows %s", val, v.Type())
		}
		v.SetFloat(val)
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			v.SetBool(b)
		} else if str, ok := value.(string); ok {
			b, err := strconv.ParseBool(str)
			if err != nil {
				return fmt.Errorf("cannot convert string %q to bool: %v", str, err)
			}
			v.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert %T to bool", value)
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte 类型
			if b, ok := value.([]byte); ok {
				newSlice := reflect.MakeSlice(v.Type(), len(b), len(b))
				reflect.Copy(newSlice, reflect.ValueOf(b))
				v.Set(newSlice)
			} else if str, ok := value.(string); ok {
				v.SetBytes([]byte(str))
			} else {
				return fmt.Errorf("cannot convert %T to []byte", value)
			}
		} else {
			return fmt.Errorf("cannot convert %T to %s", value, v.Type())
		}
	default:
		return fmt.Errorf("unsupported type %s", v.Type())
	}
	return nil
}
