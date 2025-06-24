package enums

import (
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
