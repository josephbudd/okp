package message

import "github.com/josephbudd/okp/shared/store/record"

var GetCourseToEditID = NextID()

type GetCourseToEdit struct {
	id      uint64
	name    string
	GroupID uint64 // both ways

	RecordID uint64
	Record   *record.CourseEdit // to front

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewGetCourseToEdit constructs a new NewGetCourseToEdit message.
func NewGetCourseToEdit(groupID uint64, recordID uint64) (msg *GetCourseToEdit) {
	msg = &GetCourseToEdit{
		id:       GetCourseToEditID,
		name:     "GetCourseToEdit",
		GroupID:  groupID,
		RecordID: recordID,
	}
	return
}

// GetCourseToEdit implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *GetCourseToEdit) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *GetCourseToEdit) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *GetCourseToEdit) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *GetCourseToEdit) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
