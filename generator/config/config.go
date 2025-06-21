// Package config defines configuration options for the enum generation process.
//
// The Configuration struct centralizes all options that influence the behavior
// of enum parsing and generation:
//
//   - Failfast: Strict validation mode for enum values
//   - Legacy: Compatibility mode for older Go versions
//   - Insensitive: Case flexibility in string parsing
//   - Verbose: Extended logging for debugging
//
// This package allows configuration to be passed consistently through the
// system, ensuring all components respect the same settings.
package config

// SerializationType defines the type of serialization/deserialization to use
type SerializationType int

const (
	// SerdeString uses string representation (default)
	SerdeString SerializationType = iota
	// SerdeBytes uses byte representation
	SerdeBytes
	// SerdePrimitive uses the underlying primitive type (int, float, etc.)
	SerdePrimitive
)

// EnumTypeConfig holds configuration for a specific enum type
type EnumTypeConfig struct {
	// TypeName is the name of the enum type
	TypeName string

	// UppercaseFields controls whether container struct field names should be uppercase.
	// When true, field names like STEP1INITIALIZED are generated.
	// When false (default), field names like Step1Initialized are generated in camelCase.
	UppercaseFields bool

	// GenerateNameConstants controls whether to generate enum name constants.
	// When true, generates a string type (e.g., TokenRequestStatusName) with const values
	// for each enum name, and uses these constants in the NamesMap instead of string slicing.
	GenerateNameConstants bool

	// Handlers defines which interfaces to implement for this enum type
	Handlers Handlers

	// SerializationType defines how this enum should be serialized/deserialized
	SerializationType SerializationType
}

// Configuration holds all the settings that control enum generation behavior.
// It is passed to both parsers and generators to ensure consistent behavior
// throughout the generation process.
type Configuration struct {
	// Failfast enables strict validation of enum values during parsing and generation.
	// When true, the system will return errors for invalid enum values rather than
	// silently handling them.
	Failfast bool

	// Insensitive enables case-insensitive matching when parsing enum string values.
	// When true, enum values can be matched regardless of case (e.g., "RED" == "red").
	Insensitive bool

	// Legacy enables compatibility with Go versions before 1.23.
	// When true, the generated code will not use features like range-over-func
	// that are only available in Go 1.21+.
	Legacy bool

	// Verbose enables detailed logging throughout the enum generation process.
	// When true, additional information about parsing and generation steps will
	// be logged, which is useful for debugging.
	Verbose bool

	// OutputFormat is the format of the output file.
	OutputFormat string

	// Filenames is the list of paths provided to the reader
	Filenames []string

	// Constraints is the flag to generate the constraints or not
	Constraints bool

	// Handlers defines the behavior of the enum generation process.
	// DEPRECATED: Use EnumTypeConfigs instead for per-type configuration
	Handlers Handlers

	// EnumTypeConfigs holds configuration for individual enum types
	// This allows different enum types in the same file to have different configurations
	EnumTypeConfigs map[string]EnumTypeConfig
}

// GetEnumTypeConfig returns the configuration for a specific enum type
// Falls back to global configuration if no specific config is found
func (c *Configuration) GetEnumTypeConfig(typeName string) EnumTypeConfig {
	if config, exists := c.EnumTypeConfigs[typeName]; exists {
		return config
	}

	// Fallback to global configuration for backward compatibility
	return EnumTypeConfig{
		TypeName:          typeName,
		SerializationType: SerdeString, // Default to string serialization
	}
}

type Handlers struct {
	JSON   bool
	Text   bool
	YAML   bool
	SQL    bool
	Binary bool
}
