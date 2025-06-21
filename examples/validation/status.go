package validation

type status int

const (
	Step1Initialized tokenRequestStatus = 1000 // Run2; 1000, Step1 process started (PENDING)
	Step1Canceled    tokenRequestStatus = 9010 // ;9010, User manually canceled, process ended (CANCELED)
	Step1MarkAllowed tokenRequestStatus = 1001 // ;1001, Marked as approved (PENDING)
	Step1MarkDenied  tokenRequestStatus = 1002 // ;1002, Marked as denied (PENDING)
	Step1Failed      tokenRequestStatus = 8010 // ;8010, Mark failed (others disagreed), waiting for user manual handling, re-enter process 1000, (FAILED)
	Step1Denied      tokenRequestStatus = 7011 // ;7011, Already denied, process ended (REJECTED)

	Step2WaitingPayment   tokenRequestStatus = 2000 // ;2000, Step1 passed, Step2 process started, user needs to start transferring Token to specified account (PENDING)
	Step2WaitingTxConfirm tokenRequestStatus = 2001 // ;2001, User has transferred, waiting for confirmation (at this time get user's TxHash) (PENDING)
	Step2Failed           tokenRequestStatus = 8020 // ;8020, Transfer inconsistent or transaction failed on chain, need to automatically return to process 2000 (automatically handled by program)

	Step3Initialized tokenRequestStatus = 3000 // ;3000, Step2 received, Step3 process started (fill in bank TransactionID) (PENDING)
	Step3MarkAllowed tokenRequestStatus = 3001 // ;3001, Marked as approved (PENDING)
	Step3Failed      tokenRequestStatus = 8030 // ;8030, Mark failed (others disagreed), waiting for user manual handling, re-enter process 3000 (FAILED)

	Step4Success tokenRequestStatus = 4000 // ;4000
)

// goenums: -json -text -binary -serde/string
type stringStatus int

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
