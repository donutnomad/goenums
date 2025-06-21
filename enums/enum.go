package enums

import (
	"iter"
)

type Format int

const (
	FormatName  Format = iota // Serialize as enum name (e.g. "Red")
	FormatValue               // Serialize as value (e.g. 0)
)

// Enum interface definition
type Enum[R comparable, Self comparable] interface {
	Val() R
	Values() iter.Seq[Self]
	IsValid() bool
	FindByName(name string) (Self, bool) // Return complete enum instance
	FindByValue(value R) (Self, bool)    // Return complete enum instance
	Format() Format
	Name() string // Enum name, required value
	String() string
}
