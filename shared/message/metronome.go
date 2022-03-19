package message

var MetronomeID = NextID()

type Metronome struct {
	id      uint64 // both ways
	name    string // both ways
	GroupID uint64 // both ways

	TurnOn bool

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewMetronome constructs a new NewMetronome message.
// The records are sorted by record.Number
func NewMetronome(groupID uint64, turnOn bool) (msg *Metronome) {
	msg = &Metronome{
		id:      MetronomeID,
		name:    "Metronome",
		GroupID: groupID,
		TurnOn:  turnOn,
	}
	return
}

// Metronome implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *Metronome) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *Metronome) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *Metronome) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *Metronome) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
