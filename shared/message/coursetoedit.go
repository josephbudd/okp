package message

import "github.com/josephbudd/okp/shared/store/record"

var CourseToEditID = NextID()

type CourseToEdit struct {
	id      uint64
	name    string
	GroupID uint64 // both ways
	Record  *record.CourseEdit

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewCourseToEdit constructs a new NewCourseToEdit message.
func NewCourseToEdit(groupID uint64, r *record.CourseEdit) (msg *CourseToEdit) {
	msg = &CourseToEdit{
		id:      CourseToEditID,
		name:    "CourseToEdit",
		GroupID: groupID,
		Record:  r,
	}
	return
}

// CourseToEdit implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *CourseToEdit) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *CourseToEdit) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *CourseToEdit) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *CourseToEdit) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
