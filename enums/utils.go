package enums

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

func toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case int:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case uint:
		return int64(v), true
	default:
		return 0, false
	}
}

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func anyToString(value any) (string, error) {
	var str string
	switch v := value.(type) {
	case int8:
		str = strconv.FormatInt(int64(v), 10)
	case int16:
		str = strconv.FormatInt(int64(v), 10)
	case int32:
		str = strconv.FormatInt(int64(v), 10)
	case int64:
		str = strconv.FormatInt(v, 10)
	case int:
		str = strconv.FormatInt(int64(v), 10)
	case uint8:
		str = strconv.FormatUint(uint64(v), 10)
	case uint16:
		str = strconv.FormatUint(uint64(v), 10)
	case uint32:
		str = strconv.FormatUint(uint64(v), 10)
	case uint64:
		str = strconv.FormatUint(v, 10)
	case uint:
		str = strconv.FormatUint(uint64(v), 10)
	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		str = strconv.FormatBool(v)
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		marshal, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		str = string(marshal)
	}
	return str, nil
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

func parseFloat64Value[T any](num float64, value *T) error {
	switch v := any(value).(type) {
	case *int:
		*v = int(num)
	case *int8:
		*v = int8(num)
	case *int16:
		*v = int16(num)
	case *int32:
		*v = int32(num)
	case *int64:
		*v = int64(num)
	case *uint:
		*v = uint(num)
	case *uint8:
		*v = uint8(num)
	case *uint16:
		*v = uint16(num)
	case *uint32:
		*v = uint32(num)
	case *uint64:
		*v = uint64(num)
	case *float32:
		*v = float32(num)
	case *float64:
		*v = num
	default:
		return json.Unmarshal([]byte(strconv.FormatFloat(num, 'f', -1, 64)), value)
	}
	return nil
}

// anyToBinary 将任意类型转换为二进制格式
// 使用大端字节序（network byte order）作为标准
func anyToBinary(value any) ([]byte, error) {
	switch v := value.(type) {
	case int8:
		return []byte{byte(v)}, nil
	case int16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(v))
		return buf, nil
	case int32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		return buf, nil
	case int64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		return buf, nil
	case int:
		// int 在不同平台可能是32位或64位，统一用64位存储
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		return buf, nil
	case uint8:
		return []byte{v}, nil
	case uint16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, v)
		return buf, nil
	case uint32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, v)
		return buf, nil
	case uint64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v)
		return buf, nil
	case uint:
		// uint 统一用64位存储
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v))
		return buf, nil
	case float32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, math.Float32bits(v))
		return buf, nil
	case float64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(v))
		return buf, nil
	case bool:
		if v {
			return []byte{1}, nil
		}
		return []byte{0}, nil
	case string:
		return []byte(v), nil
	case []byte:
		// 直接返回副本，避免修改原始数据
		result := make([]byte, len(v))
		copy(result, v)
		return result, nil
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

	switch v := any(value).(type) {
	case *int8:
		if len(data) < 1 {
			return fmt.Errorf("insufficient data for int8")
		}
		*v = int8(data[0])
	case *int16:
		if len(data) < 2 {
			return fmt.Errorf("insufficient data for int16")
		}
		*v = int16(binary.BigEndian.Uint16(data))
	case *int32:
		if len(data) < 4 {
			return fmt.Errorf("insufficient data for int32")
		}
		*v = int32(binary.BigEndian.Uint32(data))
	case *int64:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for int64")
		}
		*v = int64(binary.BigEndian.Uint64(data))
	case *int:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for int")
		}
		*v = int(binary.BigEndian.Uint64(data))
	case *uint8:
		if len(data) < 1 {
			return fmt.Errorf("insufficient data for uint8")
		}
		*v = data[0]
	case *uint16:
		if len(data) < 2 {
			return fmt.Errorf("insufficient data for uint16")
		}
		*v = binary.BigEndian.Uint16(data)
	case *uint32:
		if len(data) < 4 {
			return fmt.Errorf("insufficient data for uint32")
		}
		*v = binary.BigEndian.Uint32(data)
	case *uint64:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for uint64")
		}
		*v = binary.BigEndian.Uint64(data)
	case *uint:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for uint")
		}
		*v = uint(binary.BigEndian.Uint64(data))
	case *float32:
		if len(data) < 4 {
			return fmt.Errorf("insufficient data for float32")
		}
		bits := binary.BigEndian.Uint32(data)
		*v = math.Float32frombits(bits)
	case *float64:
		if len(data) < 8 {
			return fmt.Errorf("insufficient data for float64")
		}
		bits := binary.BigEndian.Uint64(data)
		*v = math.Float64frombits(bits)
	case *bool:
		if len(data) < 1 {
			return fmt.Errorf("insufficient data for bool")
		}
		*v = data[0] != 0
	case *string:
		*v = string(data)
	case *[]byte:
		// 创建副本，避免共享底层数组
		*v = make([]byte, len(data))
		copy(*v, data)
	default:
		// 对于复杂类型，尝试JSON反序列化
		return json.Unmarshal(data, value)
	}
	return nil
}
