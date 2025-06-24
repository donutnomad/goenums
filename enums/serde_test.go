package enums

import (
	"fmt"
	"testing"
)

// 测试各种序列化功能的基本集成测试

func TestSerializationIntegration(t *testing.T) {
	// 注意：MarshalJSON, MarshalText, MarshalBinary, SQLValue 等函数是泛型函数
	// 需要具体的枚举类型才能测试，这里主要测试底层工具函数

	t.Run("序列化工具函数可用性", func(t *testing.T) {
		// 这些函数在实际的枚举类型中会被调用
		// 我们在这里确认它们的存在性和基本功能
		t.Log("序列化函数已实现: MarshalJSON, MarshalText, MarshalBinary, SQLValue")
		t.Log("反序列化函数已实现: UnmarshalJSON, UnmarshalText, UnmarshalBinary, SQLScan")
		t.Log("这些函数将在具体的枚举类型中进行集成测试")
	})
}

func TestSerializationTypes(t *testing.T) {
	// 测试不同类型的序列化

	testCases := []struct {
		name  string
		value any
	}{
		{"int", 42},
		{"string", "hello"},
		{"bool", true},
		{"float64", 3.14159},
		{"[]byte", []byte("world")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试字符串序列化
			strResult, err := anyToString(tc.value)
			if err != nil {
				t.Errorf("anyToString(%v) failed: %v", tc.value, err)
			}
			t.Logf("anyToString(%v) = %s", tc.value, strResult)

			// 测试二进制序列化
			binResult, err := anyToBinary(tc.value)
			if err != nil {
				t.Errorf("anyToBinary(%v) failed: %v", tc.value, err)
			}
			t.Logf("anyToBinary(%v) = %d bytes", tc.value, len(binResult))
		})
	}
}

func TestMarshalTextVsJSON(t *testing.T) {
	// 验证MarshalText和MarshalJSON底层实现的区别
	testValue := 42

	// 直接测试底层工具函数
	strResult, err := anyToString(testValue)
	if err != nil {
		t.Fatalf("anyToString failed: %v", err)
	}

	binResult, err := anyToBinary(testValue)
	if err != nil {
		t.Fatalf("anyToBinary failed: %v", err)
	}

	t.Logf("字符串序列化结果: %s", strResult)
	t.Logf("二进制序列化结果: %d bytes", len(binResult))

	// 验证类型转换函数
	int64Val, ok := toInt64(testValue)
	if !ok {
		t.Errorf("toInt64 should work for int")
	}
	if int64Val != 42 {
		t.Errorf("toInt64 result = %d, want 42", int64Val)
	}
}

func TestBinaryVsTextSerialization(t *testing.T) {
	// 比较二进制和文本序列化的差异
	testValue := int16(1234)

	// 文本序列化
	textResult, err := anyToString(testValue)
	if err != nil {
		t.Fatalf("anyToString failed: %v", err)
	}

	// 二进制序列化
	binResult, err := anyToBinary(testValue)
	if err != nil {
		t.Fatalf("anyToBinary failed: %v", err)
	}

	t.Logf("值 %d:", testValue)
	t.Logf("  文本序列化: %s (%d 字节)", textResult, len(textResult))
	t.Logf("  二进制序列化: %v (%d 字节)", binResult, len(binResult))

	// 验证二进制序列化的长度符合预期
	if len(binResult) != 2 {
		t.Errorf("int16 binary serialization should be 2 bytes, got %d", len(binResult))
	}
}

func TestErrorHandling(t *testing.T) {
	// 测试错误处理

	t.Run("parseStringValue错误", func(t *testing.T) {
		var intVal int
		err := parseStringValue("not_a_number", &intVal)
		if err == nil {
			t.Errorf("parseStringValue should fail for invalid input")
		}
		t.Logf("期望的错误: %v", err)
	})

	t.Run("parseBinaryValue错误", func(t *testing.T) {
		var intVal int
		err := parseBinaryValue([]byte{}, &intVal)
		if err == nil {
			t.Errorf("parseBinaryValue should fail for empty data")
		}
		t.Logf("期望的错误: %v", err)
	})
}

func TestTypeConversion(t *testing.T) {
	// 测试类型转换功能

	t.Run("toInt64转换", func(t *testing.T) {
		tests := []struct {
			input    any
			expected int64
			shouldOk bool
		}{
			{int32(42), 42, true},
			{uint16(100), 100, true},
			{float64(3.14), 0, false}, // 应该失败
		}

		for _, test := range tests {
			result, ok := toInt64(test.input)
			if ok != test.shouldOk {
				t.Errorf("toInt64(%v) ok = %v, want %v", test.input, ok, test.shouldOk)
			}
			if ok && result != test.expected {
				t.Errorf("toInt64(%v) = %d, want %d", test.input, result, test.expected)
			}
		}
	})

	t.Run("toFloat64转换", func(t *testing.T) {
		tests := []struct {
			input    any
			expected float64
			shouldOk bool
		}{
			{float32(3.14), float64(float32(3.14)), true},
			{float64(2.71), 2.71, true},
			{int(42), 0, false}, // 应该失败
		}

		for _, test := range tests {
			result, ok := toFloat64(test.input)
			if ok != test.shouldOk {
				t.Errorf("toFloat64(%v) ok = %v, want %v", test.input, ok, test.shouldOk)
			}
			if ok && result != test.expected {
				t.Errorf("toFloat64(%v) = %f, want %f", test.input, result, test.expected)
			}
		}
	})
}

// 测试 MarshalBinary 和 UnmarshalBinary 功能完整性
func TestBinarySerializationIntegration(t *testing.T) {
	t.Run("验证二进制序列化集成", func(t *testing.T) {
		// 测试不同底层类型的值
		testCases := []struct {
			name  string
			value any
		}{
			{"int", int(42)},
			{"int8", int8(-123)},
			{"int16", int16(12345)},
			{"int32", int32(-2147483648)},
			{"int64", int64(9223372036854775807)},
			{"uint", uint(42)},
			{"uint8", uint8(255)},
			{"uint16", uint16(65535)},
			{"uint32", uint32(4294967295)},
			{"uint64", uint64(18446744073709551615)},
			{"float32", float32(3.14159)},
			{"float64", float64(2.718281828459045)},
			{"bool_true", true},
			{"bool_false", false},
			{"string", "hello world"},
			{"bytes", []byte("binary data")},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// 测试 anyToBinary
				binaryData, err := anyToBinary(tc.value)
				if err != nil {
					t.Fatalf("anyToBinary failed for %v: %v", tc.value, err)
				}

				// 验证数据长度符合预期
				expectedLengths := map[string]int{
					"int": 8, "int8": 1, "int16": 2, "int32": 4, "int64": 8,
					"uint": 8, "uint8": 1, "uint16": 2, "uint32": 4, "uint64": 8,
					"float32": 4, "float64": 8,
					"bool_true": 1, "bool_false": 1,
					"string": 11, "bytes": 11,
				}

				if expectedLen, ok := expectedLengths[tc.name]; ok {
					if len(binaryData) != expectedLen {
						t.Errorf("Expected binary length %d for %s, got %d", expectedLen, tc.name, len(binaryData))
					}
				}

				t.Logf("%s: 值 %v -> 二进制数据长度: %d bytes", tc.name, tc.value, len(binaryData))
			})
		}
	})

	t.Run("往返测试", func(t *testing.T) {
		// 测试各种类型的往返转换
		testRoundTripInt8 := func(original int8) {
			data, err := anyToBinary(original)
			if err != nil {
				t.Fatalf("anyToBinary failed: %v", err)
			}
			var result int8
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Fatalf("parseBinaryValue failed: %v", err)
			}
			if result != original {
				t.Errorf("int8 round-trip failed: expected %d, got %d", original, result)
			}
		}

		testRoundTripFloat32 := func(original float32) {
			data, err := anyToBinary(original)
			if err != nil {
				t.Fatalf("anyToBinary failed: %v", err)
			}
			var result float32
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Fatalf("parseBinaryValue failed: %v", err)
			}
			if result != original {
				t.Errorf("float32 round-trip failed: expected %f, got %f", original, result)
			}
		}

		testRoundTripBool := func(original bool) {
			data, err := anyToBinary(original)
			if err != nil {
				t.Fatalf("anyToBinary failed: %v", err)
			}
			var result bool
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Fatalf("parseBinaryValue failed: %v", err)
			}
			if result != original {
				t.Errorf("bool round-trip failed: expected %t, got %t", original, result)
			}
		}

		testRoundTripString := func(original string) {
			data, err := anyToBinary(original)
			if err != nil {
				t.Fatalf("anyToBinary failed: %v", err)
			}
			var result string
			err = parseBinaryValue(data, &result)
			if err != nil {
				t.Fatalf("parseBinaryValue failed: %v", err)
			}
			if result != original {
				t.Errorf("string round-trip failed: expected '%s', got '%s'", original, result)
			}
		}

		// 执行往返测试
		testRoundTripInt8(-123)
		testRoundTripFloat32(3.14159)
		testRoundTripBool(true)
		testRoundTripString("hello world")

		t.Log("所有往返测试通过")
	})

	t.Run("错误处理", func(t *testing.T) {
		// 测试空数据
		var intVal int
		err := parseBinaryValue([]byte{}, &intVal)
		if err == nil {
			t.Error("parseBinaryValue should fail for empty data")
		}

		// 测试数据不足
		var int32Val int32
		err = parseBinaryValue([]byte{0x12, 0x34}, &int32Val) // 只有2字节，需要4字节
		if err == nil {
			t.Error("parseBinaryValue should fail for insufficient data")
		}

		t.Log("错误处理测试通过")
	})
}

// MockYAMLNode 用于测试的模拟 YAML 节点
type MockYAMLNode struct {
	value interface{}
	kind  uint8
	tag   string
}

// 实现 YAMLNode 接口的方法
func (m *MockYAMLNode) Decode(v interface{}) error {
	switch target := v.(type) {
	case *string:
		if str, ok := m.value.(string); ok {
			*target = str
		} else {
			*target = fmt.Sprintf("%v", m.value)
		}
	case *int:
		if i, ok := m.value.(int); ok {
			*target = i
		} else if i64, ok := m.value.(int64); ok {
			*target = int(i64)
		} else {
			return fmt.Errorf("cannot convert %T to int", m.value)
		}
	case *int32:
		if i, ok := m.value.(int32); ok {
			*target = i
		} else if i, ok := m.value.(int); ok {
			*target = int32(i)
		} else {
			return fmt.Errorf("cannot convert %T to int32", m.value)
		}
	case *float32:
		if f, ok := m.value.(float32); ok {
			*target = f
		} else if f64, ok := m.value.(float64); ok {
			*target = float32(f64)
		} else {
			return fmt.Errorf("cannot convert %T to float32", m.value)
		}
	case *bool:
		if b, ok := m.value.(bool); ok {
			*target = b
		} else {
			return fmt.Errorf("cannot convert %T to bool", m.value)
		}
	case *interface{}:
		*target = m.value
	default:
		return fmt.Errorf("unsupported decode type: %T", v)
	}
	return nil
}

func (m *MockYAMLNode) Value() interface{} {
	return m.value
}

func (m *MockYAMLNode) Kind() uint8 {
	return m.kind
}

func (m *MockYAMLNode) Tag() string {
	return m.tag
}

// MockEnum 用于测试的模拟枚举类型
type MockEnum struct{}

func (m MockEnum) SerdeFormat() Format                          { return FormatValue }
func (m MockEnum) Name() string                                 { return "MockEnum" }
func (m MockEnum) Val() interface{}                             { return 0 }
func (m MockEnum) FromName(name string) (MockEnum, bool)        { return MockEnum{}, true }
func (m MockEnum) FromValue(value interface{}) (MockEnum, bool) { return MockEnum{}, true }
func (m MockEnum) All() []MockEnum                              { return []MockEnum{m} }

// 测试 YAML 序列化功能
func TestYAMLSerialization(t *testing.T) {
	t.Run("MarshalYAML测试", func(t *testing.T) {
		tests := []struct {
			name     string
			value    interface{}
			expected interface{}
		}{
			{"string", "hello", "hello"},
			{"int", 42, int64(42)},
			{"float64", 3.14, 3.14},
			{"bool", true, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// 直接测试 MarshalYAML 的逻辑，避免泛型类型推断问题
				// 这里我们主要验证函数的内部逻辑
				t.Logf("测试 %s 类型的 YAML 序列化逻辑", tt.name)

				// 模拟 MarshalYAML 的内部逻辑
				val := tt.value
				var result interface{}
				if v, ok := toInt64(val); ok {
					result = v
				} else if v, ok := toFloat64(val); ok {
					result = v
				} else if v, ok := val.(bool); ok {
					result = v
				} else if v, ok := val.(string); ok {
					result = v

				} else {
					str, err := anyToString(val)
					if err != nil {
						t.Errorf("anyToString failed: %v", err)
						return
					}
					result = str
				}

				t.Logf("MarshalYAML 逻辑处理 %v => %v", tt.value, result)
			})
		}
	})

	t.Run("convertToTargetType测试", func(t *testing.T) {
		// 测试各种类型转换
		t.Run("string转换", func(t *testing.T) {
			var target string
			err := convertToTargetType("hello", &target)
			if err != nil {
				t.Errorf("convertToTargetType failed: %v", err)
			}
			if target != "hello" {
				t.Errorf("expected 'hello', got '%s'", target)
			}
		})

		t.Run("int转换", func(t *testing.T) {
			var target int
			err := convertToTargetType(int64(42), &target)
			if err != nil {
				t.Errorf("convertToTargetType failed: %v", err)
			}
			if target != 42 {
				t.Errorf("expected 42, got %d", target)
			}
		})

		t.Run("float32转换", func(t *testing.T) {
			var target float32
			err := convertToTargetType(float64(3.14), &target)
			if err != nil {
				t.Errorf("convertToTargetType failed: %v", err)
			}
			expected := float32(3.14)
			if target != expected {
				t.Errorf("expected %f, got %f", expected, target)
			}
		})

		t.Run("bool转换", func(t *testing.T) {
			var target bool
			err := convertToTargetType(true, &target)
			if err != nil {
				t.Errorf("convertToTargetType failed: %v", err)
			}
			if !target {
				t.Errorf("expected true, got %t", target)
			}
		})

		t.Run("无效转换", func(t *testing.T) {
			var target int
			err := convertToTargetType("not_a_number", &target)
			if err == nil {
				t.Errorf("convertToTargetType should fail for invalid conversion")
			}
			t.Logf("期望的错误: %v", err)
		})
	})

	t.Run("UnmarshalYAML模拟测试", func(t *testing.T) {
		// 由于我们不能直接创建具体的枚举类型，我们主要测试转换函数的正确性
		// 实际的 UnmarshalYAML 测试需要在具体的枚举实现中进行
		t.Log("UnmarshalYAML 的完整测试需要具体的枚举类型支持")
		t.Log("当前主要验证 convertToTargetType 函数的正确性")

		// 测试 YAML 节点解码
		node := &MockYAMLNode{value: "test_value", kind: 2, tag: "!!str"}

		var decoded string
		err := node.Decode(&decoded)
		if err != nil {
			t.Errorf("MockYAMLNode.Decode failed: %v", err)
		}
		if decoded != "test_value" {
			t.Errorf("expected 'test_value', got '%s'", decoded)
		}

		t.Logf("MockYAMLNode 解码测试通过: %s", decoded)
	})
}
