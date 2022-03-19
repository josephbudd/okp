package message

var ToneID = NextID()

type Tone struct {
	id      uint64 // both ways
	name    string // both ways
	GroupID uint64 // both ways

	TurnOn bool

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewTone constructs a new NewTone message.
// The records are sorted by record.Number
func NewTone(groupID uint64, turnOn bool) (msg *Tone) {
	msg = &Tone{
		id:      ToneID,
		name:    "Tone",
		GroupID: groupID,
		TurnOn:  turnOn,
	}
	return
}

// Tone implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *Tone) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *Tone) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *Tone) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *Tone) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
