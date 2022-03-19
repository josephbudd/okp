package message

var RemoveCourseID = NextID()

type RemoveCourse struct {
	id      uint64
	name    string
	GroupID uint64 // both ways

	RecordID uint64 // to back

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewRemoveCourse constructs a new NewRemoveCourse message.
func NewRemoveCourse(groupID uint64, id uint64) (msg *RemoveCourse) {
	msg = &RemoveCourse{
		id:       RemoveCourseID,
		name:     "RemoveCourse",
		GroupID:  groupID,
		RecordID: id,
	}
	return
}

// RemoveCourse implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *RemoveCourse) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *RemoveCourse) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *RemoveCourse) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *RemoveCourse) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
