package message

var KeyID = NextID()

type Key struct {
	id      uint64 // both ways
	name    string // both ways
	GroupID uint64 // both ways

	Run     bool // both ways
	Testing bool // both ways

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewKey constructs a new New Key message.
func NewKey(groupID uint64, run bool, testing bool) (msg *Key) {
	msg = &Key{
		id:      KeyID,
		name:    "Key",
		GroupID: groupID,
		Run:     run,
		Testing: testing,
	}
	return
}

// Key implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *Key) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *Key) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *Key) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *Key) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
