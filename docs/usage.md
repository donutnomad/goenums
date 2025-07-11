---
layout: default
title: Usage
---

Using goenums involves a simple workflow:

1. Define your enum constants as you normally would with Go's `iota`
2. Add a `go:generate` directive to invoke goenums
3. Run `go generate` to create the enum implementation
4. Use the generated code in your project

```
$ goenums -h
   ____ _____  ___  ____  __  ______ ___  _____
  / __ '/ __ \/ _ \/ __ \/ / / / __ '__ \/ ___/
 / /_/ / /_/ /  __/ / / / /_/ / / / / / (__  ) 
 \__, /\____/\___/_/ /_/\__,_/_/ /_/ /_/____/  
/____/
Usage: goenums [options] file.go[,file2.go,...]
Options:
  -c
  -constraints
    	Specify whether to generate the float and integer constraints or import 'golang.org/x/exp/constraints' (default: false - imports)
  -f
  -failfast
    	Enable failfast mode - fail on generation of invalid enum while parsing (default: false)
  -h
  -help
    	Print help information
  -i
  -insensitive
    	Generate case insensitive string parsing (default: false)
  -l
  -legacy
    	Generate legacy code without Go 1.23+ iterator support (default: false)
  -o string
  -output string
    	Specify the output format (default: go)
  -v
  -version
    	Print version information
  -vv
  -verbose
    	Enable verbose mode - prints out the generated code (default: false)
```

## Adding a `go:generate` Directive

To use goenums, add a `go:generate` directive to your Go source file:
```go
package validation

type status int

//go:generate goenums status.go
const (
    unknown   status = iota // invalid Unknown
    pending                 // Pending
    approved                // Approved
    rejected                // Rejected
    completed               // Completed
)
```
## Running Code Generation

To generate the enum implementations, run the following command:

```bash
$ go generate ./...
```
This will create the enum implementations in the same directory as the source file.

# What Gets Generated

goenums will create a new file alongside your enum definition, named after the type. For example, a `status` type will generate a `statuses_enums.go` file in the same directory.

This file will contain:

  - A type-safe wrapper struct around your enum
  - A singleton container with all valid enum values
  - String conversion methods
  - Parsing functions for various input types (including numeric types)
  - JSON marshaling/unmarshaling
  - YAML marshaling/unmarshaling
  - Database scanning/valuing
  - Binary marshaling/unmarshaling
  - Text marshaling/unmarshaling
  - Validation functions
  - Iteration helpers
  - Compile-time validation

## Using the Generated Code

After generating the enum implementations, you can use the generated code in your Go project.

```go
// Import the package containing your enum
import "yourpackage/validation"

// Access enum values via the container
status := validation.Statuses.PASSED

// Convert to string
statusName := status.String() // "PASSED"

// Parse from string
parsed := validation.ParseStatus("FAILED")

// Check validity
if !parsed.IsValid() {
    // Handle invalid status
}

// JSON marshaling/unmarshaling
type Task struct {
    ID     int              `json:"id"`
    Status validation.Status `json:"status"`
}

// Use in exhaustive function (great for tests)
validation.ExhaustiveStatuses(func(status validation.Status) {
    // Process each status
    switch status {
    case validation.Statuses.PASSED:
        // Handle passed
    case validation.Statuses.FAILED:
        // Handle failed
    // ...handle other cases
    }
})

// Iterate using modern Go 1.23+ range-over-func
for status := range validation.Statuses.All() {
    fmt.Printf("Status: %s\n", status)
}
// Legacy iteration
for _, status := range validation.Statuses.All() {
    fmt.Printf("Status: %s\n", status)
}
```

Here is some more [Examples]({{ '/examples' | relative_url }}) →



Now you can use the generated `status` enums type in your code:

```go
/// Access enum constants safely
myStatus := validation.Statuses.PASSED

// Convert to string
fmt.Println(myStatus.String()) // "PASSED"

// Parse from various sources
parsed, _ := validation.ParseStatus("SKIPPED")

// Validate enum values
if !parsed.IsValid() {
    fmt.Println("Invalid status")
}
```

## Basic Command Syntax

```bash
goenums [options] file.go[,file2.go,...]
```

Where <filename> is the Go source file containing your enum definitions.

# Adding a `go:generate` Directive

To use goenums, you need to add a `go:generate` directive to your Go source file. 

```go
package validation

type status int

//go:generate goenums status.go
const (
	unknown status = iota // invalid
	failed
	passed
	skipped
	scheduled
	running
	booked
)
```

## Running Code Generation

To generate the enum implementations, run the following command:

```bash
go generate ./...
```

This will create the enum implementations in the same directory as the source file.

# What Gets Generated

goenums will create a new file alongside your enum definition, named after the type. For example, a `status` type will generate a `statuses_enums.go` file in the same directory.

This file will contain:

  - A type-safe wrapper struct around your enum
  - A singleton container with all valid enum values
  - String conversion methods
  - Parsing functions for various input types (including numeric types)
  - JSON marshaling/unmarshaling
  - YAML marshaling/unmarshaling
  - Database scanning/valuing
  - Binary marshaling/unmarshaling
  - Text marshaling/unmarshaling
  - Validation functions
  - Iteration helpers
  - Compile-time validation

## Using the Generated Code

After generating the enum implementations, you can use the generated code in your Go project.

```go
// Import the package containing your enum
import "yourpackage/validation"

// Access enum values via the container
status := validation.Statuses.PASSED

// Convert to string
statusName := status.String() // "PASSED"

// Parse from string
parsed, err := validation.ParseStatus("FAILED")
if err != nil {
    // Handle error
}

// Check validity
if !parsed.IsValid() {
    // Handle invalid status
}

// JSON marshaling/unmarshaling
type Task struct {
    ID     int              `json:"id"`
    Status validation.Status `json:"status"`
}

// Use in exhaustive function (great for tests)
validation.ExhaustiveStatuses(func(status validation.Status) {
    // Process each status
    switch status {
    case validation.Statuses.PASSED:
        // Handle passed
    case validation.Statuses.FAILED:
        // Handle failed
    // ...handle other cases
    }
})

// Iterate using modern Go 1.23+ range-over-func
for status := range validation.Statuses.All() {
    fmt.Printf("Status: %s\n", status)
}
```

Here is some more [Examples]({{ '/examples' | relative_url }}) → 