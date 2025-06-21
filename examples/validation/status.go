package validation

type status int

//go:generate ../../goenums -f -c status.go

const (
	// invalid UNKNOWN
	// Customer11111
	unknown status = iota
	// FAILED
	// Hello
	failed
	passed    // PASSED; I am a single line commentI am a single line commentI am a single line commentI am a single line commentI am a single line commentI am a single line commentI am a single line comment
	skipped   // SKIPPED
	scheduled // SCHEDULED
	running   // RUNNING
	booked    // BOOKED
)
