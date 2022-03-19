package message

import "github.com/josephbudd/okp/shared/store/record"

var GetCourseToRemoveID = NextID()

type GetCourseToRemove struct {
	id      uint64
	name    string
	GroupID uint64 // both ways

	RecordID uint64               // from front
	Record   *record.CourseRemove // to front

	Error        bool   // to front
	Fatal        bool   // to front
	ErrorMessage string // to front
}

// NewGetRemoveCourse constructs a new NewGetRemoveCourse message.
func NewGetRemoveCourse(groupID uint64, recordID uint64) (msg *GetCourseToRemove) {
	msg = &GetCourseToRemove{
		id:       GetCourseToRemoveID,
		name:     "GetCourseToRemove",
		GroupID:  groupID,
		RecordID: recordID,
	}
	return
}

// GetCourseToRemove implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *GetCourseToRemove) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *GetCourseToRemove) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *GetCourseToRemove) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *GetCourseToRemove) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
