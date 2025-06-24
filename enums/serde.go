package enums

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	if val == nil {
		return nil, errors.New("nil value")
	}

	v := reflect.ValueOf(val)

	// 如果是指针，获取其元素
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, errors.New("nil value")
		}
		v = v.Elem()
	}

	// 获取底层类型的Kind
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Slice:
		return v.Bytes(), nil
	default:
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
		err := NewScanner(&name).Scan(src)
		if err != nil {
			return nil, err
		}
		return findNameOrValue(e, name, true, src)
	}

	var rawValue R
	err := NewScanner(&rawValue).Scan(src)
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
	if val == nil {
		return nil, errors.New("nil value")
	}

	v := reflect.ValueOf(val)

	// 如果是指针，获取其元素
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, errors.New("nil value")
		}
		v = v.Elem()
	}

	// 获取底层类型的Kind
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Slice:
		return v.Bytes(), nil
	default:
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
