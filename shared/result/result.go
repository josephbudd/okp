package result

import "github.com/josephbudd/okp/shared/store/record"

// TestResult represents a user input and the correct answer.
type TestResult struct {
	Input   *record.KeyCode
	Control *record.KeyCode
}
