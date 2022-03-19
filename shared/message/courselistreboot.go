package message

import (
	"fmt"

	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

var CourseListRebootID = NextID()

type CourseListReboot struct {
	id      uint64                 // to front
	name    string                 // to front
	Records []*record.CourseOption // to front

	Error        bool
	Fatal        bool
	ErrorMessage string
}

// NewCourseListReboot constructs a new CourseListReboot message.
func NewCourseListReboot(stores *store.Stores) (msg *CourseListReboot) {

	msg = &CourseListReboot{
		id:   CourseListRebootID,
		name: "CourseListReboot",
	}

	var fatal error
	defer func() {
		if fatal != nil {
			fatal = fmt.Errorf("NewCourseListReboot : %w", fatal)
			msg.Fatal = true
			msg.ErrorMessage = fatal.Error()
		}
	}()

	beState := state.NewBackendState()
	currentCourse := beState.CurrentCourse()

	// Plans.
	var plans []*record.Plan
	if plans, fatal = stores.Plan.GetAll(); fatal != nil {
		return
	}
	// Get the sorted list of courses.
	// The current course is not included.
	var rr []*record.Course
	if rr, fatal = stores.Course.GetAll(); fatal != nil {
		return
	}
	options := make([]*record.CourseOption, 0, len(rr))
	for i, r := range rr {
		if r.ID != currentCourse.ID {
			option := record.ToCourseOption(r, i, plans)
			options = append(options, option)
		}
	}
	msg.Records = options
	return
}

// CourseListReboot implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *CourseListReboot) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *CourseListReboot) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *CourseListReboot) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *CourseListReboot) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
