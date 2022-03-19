package message

var CheckCurrentKeyTestID = NextID()

type CheckCurrentKeyTest struct {
	id      uint64 // both ways
	name    string // both ways
	GroupID uint64 // both ways

	Times   []int64 // both ways
	Testing bool    // both ways

	Copy    string // to front
	Text    string // to front
	DitDahs string // to front
	Passed  bool   // to front

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewCheckCurrentKeyTest constructs a new New CheckCurrentKeyTest message.
func NewCheckCurrentKeyTest(groupID uint64, times []int64, testing bool) (msg *CheckCurrentKeyTest) {
	msg = &CheckCurrentKeyTest{
		id:      CheckCurrentKeyTestID,
		name:    "CheckCurrentKeyTest",
		GroupID: groupID,
		Times:   times,
		Testing: testing,
	}
	return
}

// CheckCurrentKeyTest implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *CheckCurrentKeyTest) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *CheckCurrentKeyTest) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *CheckCurrentKeyTest) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *CheckCurrentKeyTest) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
