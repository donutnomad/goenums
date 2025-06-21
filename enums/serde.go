package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func MarshalJSON[R comparable, T comparable, E Enum[R, T]](e E, b any) ([]byte, error) {
	if e.Format() == FormatName {
		return json.Marshal(e.Name())
	}
	bs, err := anyToString(b)
	if err != nil {
		return nil, err
	}
	return []byte(bs), nil
}

func UnmarshalJSON[R comparable, T comparable, E Enum[R, T]](e E, bs []byte) (*E, error) {
	if e.Format() == FormatName {
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
	if e.Format() == FormatName {
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
	if e.Format() == FormatName {
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
	return MarshalJSON(e, b)
}

func UnmarshalText[R comparable, T comparable, E Enum[R, T]](e E, bs []byte) (*E, error) {
	return UnmarshalJSON(e, bs)
}

func MarshalBinary[R comparable, T comparable, E Enum[R, T]](e E, b any) ([]byte, error) {
	return MarshalJSON(e, b)
}

func UnmarshalBinary[R comparable, T comparable, E Enum[R, T]](e E, bs []byte) (*E, error) {
	return UnmarshalJSON(e, bs)
}

func findNameOrValue[R comparable, T comparable, E Enum[R, T], V any](e E, value V, isName bool, src any) (*E, error) {
	if isName {
		ret, ok := e.FindByName(any(value).(string))
		if ok {
			if en, ok := any(ret).(E); ok {
				return &en, nil
			}
		}
		return nil, fmt.Errorf("unknown constants %v", src)
	}
	ret, ok := e.FindByValue(any(value).(R))
	if ok {
		if en, ok := any(ret).(E); ok {
			return &en, nil
		}
	}
	return nil, fmt.Errorf("unknown constants %v", src)
}

//func MarshalYAML[R comparable, T comparable, E Enum[R, T]](e E, b any) (interface{}, error) {
//	if e.Format() == FormatName {
//		return e.Name(), nil
//	}
//	str, err := anyToString(b)
//	if err != nil {
//		return nil, err
//	}
//	return str, nil
//}
//
//func UnmarshalYAML[R comparable, T comparable, E Enum[R, T]](e E, node *yaml.Node) (*E, error) {
//	//var value interface{}
//	//if err := yaml.Unmarshal(data, &value); err != nil {
//	//	return nil, err
//	//}
//
//	// var bytes []byte
//	// switch v := value.(type) {
//	// case string:
//	// 	bytes = []byte(v)
//	// case []byte:
//	// 	bytes = v
//	// default:
//	// 	bytes = []byte(fmt.Sprintf("%v", v))
//	// }
//
//	// return UnmarshalText(e, bytes)
//	panic("don't know")
//}
