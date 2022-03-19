package message

import "github.com/josephbudd/okp/shared/store/record"

var AddCourseID = NextID()

type AddCourse struct {
	id      uint64
	name    string
	GroupID uint64 // both ways
	Record  *record.CourseAdd

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewAddCourse constructs a new NewAddCourse message.
func NewAddCourse(groupID uint64, r *record.CourseAdd) (msg *AddCourse) {
	msg = &AddCourse{
		id:      AddCourseID,
		name:    "AddCourse",
		GroupID: groupID,
		Record:  r,
	}
	return
}

// AddCourse implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *AddCourse) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *AddCourse) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *AddCourse) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *AddCourse) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
