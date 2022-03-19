package message

var SetCurrentCourseID = NextID()

type SetCurrentCourse struct {
	id      uint64
	name    string
	GroupID uint64 // both ways

	RecordID uint64 // to back

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewSetCurrentCourse constructs a new NewSetCurrentCourse message.
func NewSetCurrentCourse(groupID uint64, recordID uint64) (msg *SetCurrentCourse) {
	msg = &SetCurrentCourse{
		id:       SetCurrentCourseID,
		name:     "SetCurrentCourse",
		GroupID:  groupID,
		RecordID: recordID,
	}
	return
}

// SetCurrentCourse implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *SetCurrentCourse) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *SetCurrentCourse) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *SetCurrentCourse) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *SetCurrentCourse) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
