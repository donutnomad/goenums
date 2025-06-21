package enums

import (
	"testing"
	"time"
)

func TestGenericScanner_String(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		expected string
		wantErr  bool
	}{
		{
			name:     "string to string",
			src:      "hello",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "bytes to string",
			src:      []byte("world"),
			expected: "world",
			wantErr:  false,
		},
		{
			name:     "nil to string",
			src:      nil,
			expected: "",
			wantErr:  false,
		},
		{
			name:     "invalid type to string",
			src:      123,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			scanner := NewScanner(&result)
			err := scanner.Scan(tt.src)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenericScanner_Int(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		expected int64
		wantErr  bool
	}{
		{
			name:     "int64 to int64",
			src:      int64(123),
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "int to int64",
			src:      int(456),
			expected: 456,
			wantErr:  false,
		},
		{
			name:     "string to int64",
			src:      "789",
			expected: 789,
			wantErr:  false,
		},
		{
			name:     "bytes to int64",
			src:      []byte("321"),
			expected: 321,
			wantErr:  false,
		},
		{
			name:     "float64 to int64",
			src:      float64(12.5),
			expected: 12,
			wantErr:  false,
		},
		{
			name:     "bool true to int64",
			src:      true,
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "bool false to int64",
			src:      false,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "uint64 to int64",
			src:      uint64(999),
			expected: 999,
			wantErr:  false,
		},
		{
			name:     "invalid string to int64",
			src:      "not_a_number",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "nil to int64",
			src:      nil,
			expected: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int64
			scanner := NewScanner(&result)
			err := scanner.Scan(tt.src)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenericScanner_IntTypes(t *testing.T) {
	// Test different integer types
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "int32",
			test: func(t *testing.T) {
				var result int32
				scanner := NewScanner(&result)
				err := scanner.Scan(int64(123))
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != 123 {
					t.Errorf("expected 123, got %v", result)
				}
			},
		},
		{
			name: "uint16",
			test: func(t *testing.T) {
				var result uint16
				scanner := NewScanner(&result)
				err := scanner.Scan(int64(456))
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != 456 {
					t.Errorf("expected 456, got %v", result)
				}
			},
		},
		{
			name: "int8",
			test: func(t *testing.T) {
				var result int8
				scanner := NewScanner(&result)
				err := scanner.Scan(int64(127))
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != 127 {
					t.Errorf("expected 127, got %v", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestGenericScanner_Float(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		expected float64
		wantErr  bool
	}{
		{
			name:     "float64 to float64",
			src:      float64(123.456),
			expected: 123.456,
			wantErr:  false,
		},
		{
			name:     "float32 to float64",
			src:      float32(78.9),
			expected: float64(float32(78.9)), // Consider precision loss
			wantErr:  false,
		},
		{
			name:     "int to float64",
			src:      int(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "string to float64",
			src:      "3.14159",
			expected: 3.14159,
			wantErr:  false,
		},
		{
			name:     "bytes to float64",
			src:      []byte("2.718"),
			expected: 2.718,
			wantErr:  false,
		},
		{
			name:     "invalid string to float64",
			src:      "not_a_float",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "nil to float64",
			src:      nil,
			expected: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result float64
			scanner := NewScanner(&result)
			err := scanner.Scan(tt.src)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenericScanner_Float32(t *testing.T) {
	var result float32
	scanner := NewScanner(&result)
	err := scanner.Scan(float64(123.456))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := float32(123.456)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGenericScanner_Bool(t *testing.T) {
	tests := []struct {
		name     string
		src      any
		expected bool
		wantErr  bool
	}{
		{
			name:     "bool true to bool",
			src:      true,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "bool false to bool",
			src:      false,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "int 1 to bool",
			src:      int(1),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "int 0 to bool",
			src:      int(0),
			expected: false,
			wantErr:  false,
		},
		{
			name:     "int64 non-zero to bool",
			src:      int64(42),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "float64 non-zero to bool",
			src:      float64(3.14),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "float64 zero to bool",
			src:      float64(0.0),
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string true to bool",
			src:      "true",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string false to bool",
			src:      "false",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "bytes true to bool",
			src:      []byte("true"),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "invalid string to bool",
			src:      "maybe",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "nil to bool",
			src:      nil,
			expected: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			scanner := NewScanner(&result)
			err := scanner.Scan(tt.src)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenericScanner_Time(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 10, 30, 45, 123456789, time.UTC)
	timeString := testTime.Format(time.RFC3339Nano)

	tests := []struct {
		name     string
		src      any
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "time.Time to time.Time",
			src:      testTime,
			expected: testTime,
			wantErr:  false,
		},
		{
			name:     "string to time.Time",
			src:      timeString,
			expected: testTime,
			wantErr:  false,
		},
		{
			name:     "bytes to time.Time",
			src:      []byte(timeString),
			expected: testTime,
			wantErr:  false,
		},
		{
			name:     "invalid string to time.Time",
			src:      "not_a_time",
			expected: time.Time{},
			wantErr:  true,
		},
		{
			name:     "nil to time.Time",
			src:      nil,
			expected: time.Time{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result time.Time
			scanner := NewScanner(&result)
			err := scanner.Scan(tt.src)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !result.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenericScanner_UnsupportedType(t *testing.T) {
	type CustomType struct {
		Value string
	}

	var result CustomType
	scanner := NewScanner(&result)
	err := scanner.Scan("test")

	if err == nil {
		t.Errorf("expected error for unsupported type but got none")
	}

	expectedErrMsg := "unsupported target type:"
	if err != nil && len(err.Error()) < len(expectedErrMsg) {
		t.Errorf("error message too short: %v", err.Error())
	}
}

func TestGenericScanner_NilHandling(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "nil to string",
			test: func(t *testing.T) {
				var result string
				scanner := NewScanner(&result)
				err := scanner.Scan(nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// nil should keep original value unchanged
			},
		},
		{
			name: "nil to int",
			test: func(t *testing.T) {
				var result int
				scanner := NewScanner(&result)
				err := scanner.Scan(nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "nil to bool",
			test: func(t *testing.T) {
				var result bool
				scanner := NewScanner(&result)
				err := scanner.Scan(nil)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestGenericScanner_EdgeCases(t *testing.T) {
	t.Run("empty string to int", func(t *testing.T) {
		var result int
		scanner := NewScanner(&result)
		err := scanner.Scan("")
		if err == nil {
			t.Errorf("expected error for empty string to int")
		}
	})

	t.Run("empty bytes to float", func(t *testing.T) {
		var result float64
		scanner := NewScanner(&result)
		err := scanner.Scan([]byte(""))
		if err == nil {
			t.Errorf("expected error for empty bytes to float")
		}
	})

	t.Run("very large int64", func(t *testing.T) {
		var result int64
		scanner := NewScanner(&result)
		largeInt := int64(9223372036854775807) // max int64
		err := scanner.Scan(largeInt)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != largeInt {
			t.Errorf("expected %v, got %v", largeInt, result)
		}
	})
}
