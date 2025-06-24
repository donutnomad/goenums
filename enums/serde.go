package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func MarshalJSON[R comparable, T comparable, E Enum[R, T]](e E, b any) ([]byte, error) {
	if e.SerdeFormat() == FormatName {
		return json.Marshal(e.Name())
	}
	bs, err := anyToString(b)
	if err != nil {
		return nil, err
	}
	return []byte(bs), nil
}

func UnmarshalJSON[R comparable, T comparable, E Enum[R, T]](e E, bs []byte) (*E, error) {
	if e.SerdeFormat() == FormatName {
		var name string
		if err := json.Unmarshal(bs, &name); err != nil {
			return nil, err
		}
		return findNameOrValue(e, name, true, string(bs))
	}
	var rawValue R
	if err := json.Unmarshal(bs, &rawValue); err != nil {
		return nil, err
	}
	return findNameOrValue(e, rawValue, false, string(bs))
}

func SQLValue[R comparable, T comparable, E Enum[R, T]](e E) (driver.Value, error) {
	if e.SerdeFormat() == FormatName {
		return e.Name(), nil
	}
	val := any(e.Val())
	if v, ok := toInt64(val); ok {
		return v, nil
	} else if v, ok := toFloat64(val); ok {
		return v, nil
	} else if v, ok := val.([]byte); ok {
		return v, nil
	} else if v, ok := val.(bool); ok {
		return v, nil
	} else if v, ok := val.(string); ok {
		return v, nil
	} else {
		marshal, err := json.Marshal(val)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	}
}

func SQLScan[R comparable, T comparable, E Enum[R, T]](e E, src any) (*E, error) {
	if e.SerdeFormat() == FormatName {
		var name string
		err := NewScanner[string](&name).Scan(src)
		if err != nil {
			return nil, err
		}
		return findNameOrValue(e, name, true, src)
	}

	var rawValue R
	err := NewScanner[R](&rawValue).Scan(src)
	if err != nil {
		return nil, err
	}
	return findNameOrValue(e, rawValue, false, src)
}

func MarshalText[R comparable, T comparable, E Enum[R, T]](e E, b any) ([]byte, error) {
	if e.SerdeFormat() == FormatName {
		return []byte(e.Name()), nil
	}
	bs, err := anyToString(b)
	if err != nil {
		return nil, err
	}
	return []byte(bs), nil
}

func UnmarshalText[R comparable, T comparable, E Enum[R, T]](e E, bs []byte) (*E, error) {
	str := string(bs)
	if e.SerdeFormat() == FormatName {
		return findNameOrValue(e, str, true, str)
	}

	var rawValue R
	err := parseStringValue(str, &rawValue)
	if err != nil {
		return nil, err
	}
	return findNameOrValue(e, rawValue, false, string(bs))
}

func MarshalBinary[R comparable, T comparable, E Enum[R, T]](e E, b any) ([]byte, error) {
	if e.SerdeFormat() == FormatName {
		return []byte(e.Name()), nil
	}
	return anyToBinary(b)
}

func UnmarshalBinary[R comparable, T comparable, E Enum[R, T]](e E, bs []byte) (*E, error) {
	if e.SerdeFormat() == FormatName {
		name := string(bs)
		return findNameOrValue(e, name, true, string(bs))
	}

	var rawValue R
	err := parseBinaryValue(bs, &rawValue)
	if err != nil {
		return nil, err
	}
	return findNameOrValue(e, rawValue, false, string(bs))
}

func findNameOrValue[R comparable, T comparable, E Enum[R, T], V any](e E, value V, isName bool, src any) (*E, error) {
	if isName {
		ret, ok := e.FromName(any(value).(string))
		if ok {
			if en, ok := any(ret).(E); ok {
				return &en, nil
			}
		}
		return nil, fmt.Errorf("unknown constants %v", src)
	}
	ret, ok := e.FromValue(any(value).(R))
	if ok {
		if en, ok := any(ret).(E); ok {
			return &en, nil
		}
	}
	return nil, fmt.Errorf("unknown constants %v", src)
}

// YAMLNode represents a YAML node interface to avoid importing yaml package directly
// This allows the package to work without requiring yaml dependency
type YAMLNode interface {
	// Decode decodes the node into the provided value
	Decode(interface{}) error
}

// MarshalYAML implements YAML marshaling for enums
// Returns the value that should be marshaled to YAML
func MarshalYAML[R comparable, T comparable, E Enum[R, T]](e E, b any) (interface{}, error) {
	if e.SerdeFormat() == FormatName {
		return e.Name(), nil
	}

	// For value format, we need to return the actual value
	val := any(e.Val())
	if v, ok := toInt64(val); ok {
		return v, nil
	} else if v, ok := toFloat64(val); ok {
		return v, nil
	} else if v, ok := val.(bool); ok {
		return v, nil
	} else if v, ok := val.(string); ok {
		return v, nil
	} else if v, ok := val.([]byte); ok {
		return string(v), nil
	} else {
		// For complex types, convert to string representation
		str, err := anyToString(val)
		if err != nil {
			return nil, err
		}
		return str, nil
	}
}

// UnmarshalYAML implements YAML unmarshaling for enums using the new Node interface
func UnmarshalYAML[R comparable, T comparable, E Enum[R, T]](e E, node YAMLNode) (*E, error) {
	if e.SerdeFormat() == FormatName {
		var name string
		if err := node.Decode(&name); err != nil {
			return nil, fmt.Errorf("failed to decode YAML node as string: %w", err)
		}
		return findNameOrValue(e, name, true, name)
	}

	// For value format, try to decode as the raw value type
	var rawValue R
	if err := node.Decode(&rawValue); err != nil {
		// If direct decoding fails, try to decode as interface{} and convert
		var value interface{}
		if err2 := node.Decode(&value); err2 != nil {
			return nil, fmt.Errorf("failed to decode YAML node: %w", err)
		}

		// Convert the decoded value to the target type
		if err := convertToTargetType(value, &rawValue); err != nil {
			return nil, fmt.Errorf("failed to convert YAML value to target type: %w", err)
		}
	}

	return findNameOrValue(e, rawValue, false, rawValue)
}

// convertToTargetType converts an interface{} value to the target type
func convertToTargetType[R comparable](value interface{}, target *R) error {
	switch t := any(target).(type) {
	case *string:
		if str, ok := value.(string); ok {
			*t = str
		} else {
			str, err := anyToString(value)
			if err != nil {
				return err
			}
			*t = str
		}
	case *int:
		if v, ok := toInt64(value); ok {
			*t = int(v)
		} else {
			return fmt.Errorf("cannot convert %T to int", value)
		}
	case *int8:
		if v, ok := toInt64(value); ok {
			*t = int8(v)
		} else {
			return fmt.Errorf("cannot convert %T to int8", value)
		}
	case *int16:
		if v, ok := toInt64(value); ok {
			*t = int16(v)
		} else {
			return fmt.Errorf("cannot convert %T to int16", value)
		}
	case *int32:
		if v, ok := toInt64(value); ok {
			*t = int32(v)
		} else {
			return fmt.Errorf("cannot convert %T to int32", value)
		}
	case *int64:
		if v, ok := toInt64(value); ok {
			*t = v
		} else {
			return fmt.Errorf("cannot convert %T to int64", value)
		}
	case *uint:
		if v, ok := toInt64(value); ok && v >= 0 {
			*t = uint(v)
		} else {
			return fmt.Errorf("cannot convert %T to uint", value)
		}
	case *uint8:
		if v, ok := toInt64(value); ok && v >= 0 && v <= 255 {
			*t = uint8(v)
		} else {
			return fmt.Errorf("cannot convert %T to uint8", value)
		}
	case *uint16:
		if v, ok := toInt64(value); ok && v >= 0 && v <= 65535 {
			*t = uint16(v)
		} else {
			return fmt.Errorf("cannot convert %T to uint16", value)
		}
	case *uint32:
		if v, ok := toInt64(value); ok && v >= 0 && v <= 4294967295 {
			*t = uint32(v)
		} else {
			return fmt.Errorf("cannot convert %T to uint32", value)
		}
	case *uint64:
		if v, ok := toInt64(value); ok && v >= 0 {
			*t = uint64(v)
		} else {
			return fmt.Errorf("cannot convert %T to uint64", value)
		}
	case *float32:
		if v, ok := toFloat64(value); ok {
			*t = float32(v)
		} else {
			return fmt.Errorf("cannot convert %T to float32", value)
		}
	case *float64:
		if v, ok := toFloat64(value); ok {
			*t = v
		} else {
			return fmt.Errorf("cannot convert %T to float64", value)
		}
	case *bool:
		if v, ok := value.(bool); ok {
			*t = v
		} else {
			return fmt.Errorf("cannot convert %T to bool", value)
		}
	case *[]byte:
		if v, ok := value.(string); ok {
			*t = []byte(v)
		} else if v, ok := value.([]byte); ok {
			*t = v
		} else {
			return fmt.Errorf("cannot convert %T to []byte", value)
		}
	default:
		// For complex types, try JSON conversion
		jsonData, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for conversion: %w", err)
		}
		if err := json.Unmarshal(jsonData, target); err != nil {
			return fmt.Errorf("failed to unmarshal value: %w", err)
		}
	}
	return nil
}
