package enums

import (
	"encoding/json"
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
