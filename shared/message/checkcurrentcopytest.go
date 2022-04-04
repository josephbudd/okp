package message

import "github.com/josephbudd/okp/shared/state"

var CheckCurrentCopyTestID = NextID()

type CheckCurrentCopyTest struct {
	id      uint64 // both ways
	name    string // both ways
	GroupID uint64 // both ways

	Copy string // to back

	Text    string // to front
	DitDahs string // to front
	Passed  bool   // to front

	State state.Message

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewCheckCurrentCopyTest constructs a new New CheckCurrentCopyTest message.
func NewCheckCurrentCopyTest(groupID uint64, userCopy string) (msg *CheckCurrentCopyTest) {
	msg = &CheckCurrentCopyTest{
		id:      CheckCurrentCopyTestID,
		name:    "CheckCurrentCopyTest",
		GroupID: groupID,
		Copy:    userCopy,
	}
	return
}

// CheckCurrentCopyTest implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *CheckCurrentCopyTest) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *CheckCurrentCopyTest) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *CheckCurrentCopyTest) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *CheckCurrentCopyTest) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
