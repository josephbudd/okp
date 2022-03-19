package currentcourse

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"

	"github.com/josephbudd/okp/frontend/txrx"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "selectcourse"
	return
}

func (m *messenger) listen() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("messenger.listen: %w", err)
		}
	}()

	if err = txrx.AddListener(message.GetCourseToEditID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.SetCurrentCourseID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.CourseListRebootID, m); err != nil {
		return
	}
	// Listen for state changes.
	appState.AddListener(m)
	return
}

func (m *messenger) Listen(msg interface{}) {
	// A message sent from the main process to the view.
	switch msg := msg.(type) {
	case *message.SetCurrentCourse:
		m.setCurrentCourseRX(msg)
	case *message.CourseListReboot:
		m.courseListRebootRX(msg)
	default:
	}
}

// CourseListReboot message. Sends the new course list options.

func (m *messenger) courseListRebootRX(msg *message.CourseListReboot) {
	countCourses = len(msg.Records)
	if countCourses > 0 {
		sPanel.list.Reboot(msg.Records)
	}
}

// State

// StateRX gets the message from the state.
func (m *messenger) StateRX(msg state.Message) {
	fPanel.fillForm()
	showFormPanel()
}

// Set Current course using it's sorted index.

func (m *messenger) setCurrentCourseTX(id uint64) {
	msg := message.NewSetCurrentCourse(groupID, id)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) setCurrentCourseRX(msg *message.SetCurrentCourse) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
}
