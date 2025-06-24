package enums

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

// 测试 anyToString 函数
func TestAnyToString(t *testing.T) {
	// 测试结构体
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	tests := []struct {
		name      string
		input     any
		expected  string
		shouldErr bool
	}{
		// 整数类型
		{"int8", int8(-123), "-123", false},
		{"int16", int16(12345), "12345", false},
		{"int32", int32(-2147483648), "-2147483648", false},
		{"int64", int64(9223372036854775807), "9223372036854775807", false},
		{"int", int(0), "0", false},
		{"uint8", uint8(255), "255", false},
		{"uint16", uint16(65535), "65535", false},
		{"uint32", uint32(4294967295), "4294967295", false},
		{"uint64", uint64(18446744073709551615), "18446744073709551615", false},
		{"uint", uint(123), "123", false},

		// 浮点数类型
		{"float32", float32(3.14), "3.14", false},
		{"float64", float64(3.141592653589793), "3.141592653589793", false},
		{"float32 zero", float32(0), "0", false},
		{"float64 negative", float64(-2.71828), "-2.71828", false},

		// 布尔类型
		{"bool true", true, "true", false},
		{"bool false", false, "false", false},

		// 字符串和字节类型
		{"string", "hello world", "hello world", false},
		{"string empty", "", "", false},
		{"[]byte", []byte("hello"), "hello", false},
		{"[]byte empty", []byte{}, "", false},

		// 复杂类型（JSON序列化）
		{"struct", TestStruct{Name: "test", Value: 42}, `{"name":"test","value":42}`, false},
		{"slice", []int{1, 2, 3}, "[1,2,3]", false},
		{"map", map[string]int{"a": 1, "b": 2}, "", false}, // JSON map顺序不确定，测试时特殊处理
		{"nil", nil, "null", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := anyToString(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("anyToString(%v) should return error but got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("anyToString(%v) returned error: %v", tt.input, err)
				return
			}

			// 特殊处理map类型，因为JSON序列化顺序不确定
			if tt.name == "map" {
				var mapResult map[string]int
				if json.Unmarshal([]byte(result), &mapResult) != nil {
					t.Errorf("anyToString(%v) result is not valid JSON: %s", tt.input, result)
				}
				return
			}

			if result != tt.expected {
				t.Errorf("anyToString(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// 测试 parseStringValue 函数
func TestParseStringValue(t *testing.T) {
	t.Run("整数类型", func(t *testing.T) {
		var intVal int
		err := parseStringValue("12345", &intVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if intVal != 12345 {
			t.Errorf("expected 12345, got %d", intVal)
		}

		var int8Val int8
		err = parseStringValue("-128", &int8Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if int8Val != -128 {
			t.Errorf("expected -128, got %d", int8Val)
		}

		// 测试所有整数类型
		var int16Val int16
		err = parseStringValue("32767", &int16Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if int16Val != 32767 {
			t.Errorf("expected 32767, got %d", int16Val)
		}

		var int32Val int32
		err = parseStringValue("-2147483648", &int32Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if int32Val != -2147483648 {
			t.Errorf("expected -2147483648, got %d", int32Val)
		}

		var int64Val int64
		err = parseStringValue("9223372036854775807", &int64Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if int64Val != 9223372036854775807 {
			t.Errorf("expected 9223372036854775807, got %d", int64Val)
		}

		// 测试无符号整数类型
		var uintVal uint
		err = parseStringValue("12345", &uintVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if uintVal != 12345 {
			t.Errorf("expected 12345, got %d", uintVal)
		}

		var uint8Val uint8
		err = parseStringValue("255", &uint8Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if uint8Val != 255 {
			t.Errorf("expected 255, got %d", uint8Val)
		}

		var uint16Val uint16
		err = parseStringValue("65535", &uint16Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if uint16Val != 65535 {
			t.Errorf("expected 65535, got %d", uint16Val)
		}

		var uint32Val uint32
		err = parseStringValue("4294967295", &uint32Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if uint32Val != 4294967295 {
			t.Errorf("expected 4294967295, got %d", uint32Val)
		}

		var uint64Val uint64
		err = parseStringValue("18446744073709551615", &uint64Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if uint64Val != 18446744073709551615 {
			t.Errorf("expected 18446744073709551615, got %d", uint64Val)
		}
	})

	t.Run("浮点数类型", func(t *testing.T) {
		var float32Val float32
		err := parseStringValue("3.14", &float32Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if float32Val != 3.14 {
			t.Errorf("expected 3.14, got %f", float32Val)
		}

		var float64Val float64
		err = parseStringValue("-2.718281828", &float64Val)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if float64Val != -2.718281828 {
			t.Errorf("expected -2.718281828, got %f", float64Val)
		}
	})

	t.Run("布尔类型", func(t *testing.T) {
		var boolVal bool
		err := parseStringValue("true", &boolVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if !boolVal {
			t.Errorf("expected true, got %t", boolVal)
		}

		err = parseStringValue("false", &boolVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if boolVal {
			t.Errorf("expected false, got %t", boolVal)
		}

		err = parseStringValue("1", &boolVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if !boolVal {
			t.Errorf("expected true for '1', got %t", boolVal)
		}
	})

	t.Run("字符串和字节类型", func(t *testing.T) {
		var strVal string
		err := parseStringValue("hello world", &strVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if strVal != "hello world" {
			t.Errorf("expected 'hello world', got %q", strVal)
		}

		var bytesVal []byte
		err = parseStringValue("hello", &bytesVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if string(bytesVal) != "hello" {
			t.Errorf("expected 'hello', got %q", string(bytesVal))
		}
	})

	t.Run("JSON类型", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		var structVal TestStruct
		err := parseStringValue(`{"name":"test","value":42}`, &structVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if structVal.Name != "test" || structVal.Value != 42 {
			t.Errorf("expected {Name: test, Value: 42}, got %+v", structVal)
		}

		// 测试 JSON 数组
		var sliceVal []int
		err = parseStringValue(`[1,2,3]`, &sliceVal)
		if err != nil {
			t.Errorf("parseStringValue failed: %v", err)
		}
		if len(sliceVal) != 3 || sliceVal[0] != 1 || sliceVal[1] != 2 || sliceVal[2] != 3 {
			t.Errorf("expected [1,2,3], got %v", sliceVal)
		}
	})

	t.Run("错误情况", func(t *testing.T) {
		var intVal int
		err := parseStringValue("not_a_number", &intVal)
		if err == nil {
			t.Errorf("parseStringValue should fail for invalid integer")
		}

		var floatVal float64
		err = parseStringValue("not_a_float", &floatVal)
		if err == nil {
			t.Errorf("parseStringValue should fail for invalid float")
		}

		var boolVal bool
		err = parseStringValue("not_a_bool", &boolVal)
		if err == nil {
			t.Errorf("parseStringValue should fail for invalid bool")
		}

		// 测试溢出情况
		var int8Val int8
		err = parseStringValue("1000", &int8Val) // 超过 int8 最大值 127
		if err == nil {
			t.Errorf("parseStringValue should fail for int8 overflow")
		}

		var uint8Val uint8
		err = parseStringValue("256", &uint8Val) // 超过 uint8 最大值 255
		if err == nil {
			t.Errorf("parseStringValue should fail for uint8 overflow")
		}

		// 测试无效 JSON
		var structVal struct {
			Name string `json:"name"`
		}
		err = parseStringValue(`{invalid json}`, &structVal)
		if err == nil {
			t.Errorf("parseStringValue should fail for invalid JSON")
		}
	})
}

// 测试 anyToBinary 函数
func TestAnyToBinary(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		expectedLen int
		shouldErr   bool
	}{
		// 整数类型
		{"int8", int8(-123), 1, false},
		{"int16", int16(1234), 2, false},
		{"int32", int32(123456), 4, false},
		{"int64", int64(123456789), 8, false},
		{"int", int(123), 8, false}, // 统一按64位处理
		{"uint8", uint8(255), 1, false},
		{"uint16", uint16(65535), 2, false},
		{"uint32", uint32(4294967295), 4, false},
		{"uint64", uint64(18446744073709551615), 8, false},
		{"uint", uint(123), 8, false}, // 统一按64位处理

		// 浮点数类型
		{"float32", float32(3.14), 4, false},
		{"float64", float64(3.141592653589793), 8, false},

		// 其他类型
		{"bool true", true, 1, false},
		{"bool false", false, 1, false},
		{"string", "hello", 5, false},
		{"[]byte", []byte("world"), 5, false},
		{"empty string", "", 0, false},
		{"empty []byte", []byte{}, 0, false},

		// 复杂类型（JSON序列化）
		{"slice", []int{1, 2, 3}, 7, false}, // "[1,2,3]"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := anyToBinary(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("anyToBinary(%v) should return error but got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("anyToBinary(%v) returned error: %v", tt.input, err)
				return
			}

			if len(result) != tt.expectedLen {
				t.Errorf("anyToBinary(%v) length = %d, want %d", tt.input, len(result), tt.expectedLen)
			}
		})
	}

	// 测试二进制数据的正确性
	t.Run("二进制数据验证", func(t *testing.T) {
		// 测试大端字节序
		data, err := anyToBinary(int16(0x1234))
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		expected := []byte{0x12, 0x34}
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("anyToBinary(0x1234) = %v, want %v", data, expected)
		}

		// 测试bool值
		trueData, _ := anyToBinary(true)
		falseData, _ := anyToBinary(false)
		if !reflect.DeepEqual(trueData, []byte{1}) {
			t.Errorf("anyToBinary(true) = %v, want [1]", trueData)
		}
		if !reflect.DeepEqual(falseData, []byte{0}) {
			t.Errorf("anyToBinary(false) = %v, want [0]", falseData)
		}

		// 测试[]byte副本
		original := []byte("test")
		data, _ = anyToBinary(original)
		original[0] = 'X' // 修改原始数据
		if data[0] == 'X' {
			t.Errorf("anyToBinary should return a copy, not share underlying array")
		}
	})
}

// 测试 parseBinaryValue 函数
func TestParseBinaryValue(t *testing.T) {
	t.Run("空数据错误", func(t *testing.T) {
		var intVal int
		err := parseBinaryValue([]byte{}, &intVal)
		if err == nil {
			t.Errorf("parseBinaryValue should fail for empty data")
		}
	})

	t.Run("数据不足错误", func(t *testing.T) {
		var int16Val int16
		err := parseBinaryValue([]byte{0x12}, &int16Val) // 只有1字节，需要2字节
		if err == nil {
			t.Errorf("parseBinaryValue should fail for insufficient data")
		}
	})

	t.Run("整数类型解析", func(t *testing.T) {
		// int16 大端字节序
		var int16Val int16
		err := parseBinaryValue([]byte{0x12, 0x34}, &int16Val)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if int16Val != 0x1234 {
			t.Errorf("expected 0x1234, got 0x%x", int16Val)
		}

		// 负数
		var int8Val int8
		err = parseBinaryValue([]byte{0x85}, &int8Val) // -123
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if int8Val != -123 {
			t.Errorf("expected -123, got %d", int8Val)
		}
	})

	t.Run("浮点数类型解析", func(t *testing.T) {
		// 先序列化再反序列化
		original := float32(3.14159)
		data, _ := anyToBinary(original)

		var result float32
		err := parseBinaryValue(data, &result)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if result != original {
			t.Errorf("expected %f, got %f", original, result)
		}
	})

	t.Run("布尔类型解析", func(t *testing.T) {
		var boolVal bool

		err := parseBinaryValue([]byte{1}, &boolVal)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if !boolVal {
			t.Errorf("expected true, got %t", boolVal)
		}

		err = parseBinaryValue([]byte{0}, &boolVal)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if boolVal {
			t.Errorf("expected false, got %t", boolVal)
		}

		// 非零值应该被解析为true
		err = parseBinaryValue([]byte{255}, &boolVal)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if !boolVal {
			t.Errorf("expected true for non-zero value, got %t", boolVal)
		}
	})

	t.Run("字符串和字节解析", func(t *testing.T) {
		var strVal string
		err := parseBinaryValue([]byte("hello"), &strVal)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if strVal != "hello" {
			t.Errorf("expected 'hello', got %q", strVal)
		}

		var bytesVal []byte
		original := []byte("world")
		err = parseBinaryValue(original, &bytesVal)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if string(bytesVal) != "world" {
			t.Errorf("expected 'world', got %q", string(bytesVal))
		}

		// 测试副本
		original[0] = 'X'
		if bytesVal[0] == 'X' {
			t.Errorf("parseBinaryValue should create a copy, not share underlying array")
		}
	})

	t.Run("JSON类型解析", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		data := []byte(`{"name":"test","value":42}`)
		var structVal TestStruct
		err := parseBinaryValue(data, &structVal)
		if err != nil {
			t.Errorf("parseBinaryValue failed: %v", err)
		}
		if structVal.Name != "test" || structVal.Value != 42 {
			t.Errorf("expected {Name: test, Value: 42}, got %+v", structVal)
		}
	})
}

// 测试完整的二进制序列化往返
func TestBinaryRoundTrip(t *testing.T) {
	// 为每种类型分别测试，避免泛型类型推断问题
	t.Run("int8", func(t *testing.T) {
		original := int8(-123)
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result int8
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("int16", func(t *testing.T) {
		original := int16(12345)
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result int16
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("int64", func(t *testing.T) {
		original := int64(9223372036854775807)
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result int64
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("uint64", func(t *testing.T) {
		original := uint64(18446744073709551615)
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result uint64
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("float32", func(t *testing.T) {
		original := float32(3.14159)
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result float32
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("float64", func(t *testing.T) {
		original := float64(2.718281828459045)
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result float64
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("bool", func(t *testing.T) {
		for _, original := range []bool{true, false} {
			data, err := anyToBinary(original)
			if err != nil {
				t.Fatalf("anyToBinary failed: %v", err)
			}
			var result bool
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Fatalf("parseBinaryValue failed: %v", err)
			}
			if original != result {
				t.Errorf("round-trip failed: original %v, got %v", original, result)
			}
		}
	})

	t.Run("string", func(t *testing.T) {
		original := "hello world"
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result string
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if original != result {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("[]byte", func(t *testing.T) {
		original := []byte("binary data")
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result []byte
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if !reflect.DeepEqual(original, result) {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})

	t.Run("complex_types", func(t *testing.T) {
		// JSON类型
		type TestStruct struct {
			Key    string `json:"key"`
			Number int    `json:"number"`
		}
		original := TestStruct{Key: "value", Number: 42}
		data, err := anyToBinary(original)
		if err != nil {
			t.Fatalf("anyToBinary failed: %v", err)
		}
		var result TestStruct
		err = parseBinaryValue(data, &result)
		if err != nil {
			t.Fatalf("parseBinaryValue failed: %v", err)
		}
		if !reflect.DeepEqual(original, result) {
			t.Errorf("round-trip failed: original %v, got %v", original, result)
		}
	})
}

// 边界值测试
func TestBoundaryValues(t *testing.T) {
	t.Run("int8边界值", func(t *testing.T) {
		for _, val := range []int8{-128, 127} {
			str, err := anyToString(val)
			if err != nil {
				t.Errorf("anyToString(%v) failed: %v", val, err)
			}

			data, err := anyToBinary(val)
			if err != nil {
				t.Errorf("anyToBinary(%v) failed: %v", val, err)
			}

			var result int8
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Errorf("parseBinaryValue for %v failed: %v", val, err)
			}

			if val != result {
				t.Errorf("boundary value round-trip failed: %v != %v", val, result)
			}

			t.Logf("边界值 %v -> 字符串: %s, 二进制长度: %d", val, str, len(data))
		}
	})

	t.Run("uint64边界值", func(t *testing.T) {
		for _, val := range []uint64{0, 18446744073709551615} {
			str, err := anyToString(val)
			if err != nil {
				t.Errorf("anyToString(%v) failed: %v", val, err)
			}

			data, err := anyToBinary(val)
			if err != nil {
				t.Errorf("anyToBinary(%v) failed: %v", val, err)
			}

			var result uint64
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Errorf("parseBinaryValue for %v failed: %v", val, err)
			}

			if val != result {
				t.Errorf("boundary value round-trip failed: %v != %v", val, result)
			}

			t.Logf("边界值 %v -> 字符串: %s, 二进制长度: %d", val, str, len(data))
		}
	})

	t.Run("int64边界值", func(t *testing.T) {
		for _, val := range []int64{-9223372036854775808, 9223372036854775807} {
			str, err := anyToString(val)
			if err != nil {
				t.Errorf("anyToString(%v) failed: %v", val, err)
			}

			data, err := anyToBinary(val)
			if err != nil {
				t.Errorf("anyToBinary(%v) failed: %v", val, err)
			}

			var result int64
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Errorf("parseBinaryValue for %v failed: %v", val, err)
			}

			if val != result {
				t.Errorf("boundary value round-trip failed: %v != %v", val, result)
			}

			t.Logf("边界值 %v -> 字符串: %s, 二进制长度: %d", val, str, len(data))
		}
	})

	t.Run("特殊浮点值", func(t *testing.T) {
		specialFloats := []float64{
			0.0, -0.0, math.Inf(1), math.Inf(-1), math.NaN(),
			math.SmallestNonzeroFloat64, math.MaxFloat64,
		}

		for _, val := range specialFloats {
			str, err := anyToString(val)
			if err != nil {
				t.Errorf("anyToString(%v) failed: %v", val, err)
			}

			data, err := anyToBinary(val)
			if err != nil {
				t.Errorf("anyToBinary(%v) failed: %v", val, err)
			}

			var result float64
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Errorf("parseBinaryValue for %v failed: %v", val, err)
			}

			// 特殊处理NaN
			if math.IsNaN(val) {
				if !math.IsNaN(result) {
					t.Errorf("NaN round-trip failed: expected NaN, got %v", result)
				}
			} else if val != result {
				t.Errorf("float round-trip failed: %v != %v", val, result)
			}

			t.Logf("特殊浮点值 %v -> 字符串: %s, 二进制长度: %d", val, str, len(data))
		}
	})
}
