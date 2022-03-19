package removecourse

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"

	"github.com/josephbudd/okp/frontend/txrx"
	"github.com/josephbudd/okp/shared/message"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "removecourse"
	return
}

func (m *messenger) listen() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("messenger.listen: %w", err)
		}
	}()

	if err = txrx.AddListener(message.GetCourseToRemoveID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.RemoveCourseID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.CourseListRebootID, m); err != nil {
		return
	}
	return
}

func (m *messenger) Listen(msg interface{}) {
	// A message sent from the main process to the view.
	switch msg := msg.(type) {
	case *message.GetCourseToRemove:
		m.getRemoveCourseRX(msg)
	case *message.RemoveCourse:
		m.courseRemoveRX(msg)
	case *message.CourseListReboot:
		m.courseListRebootRX(msg)
	default:
	}
}

// CourseListReboot message. Sends the new course list options.

func (m *messenger) courseListRebootRX(msg *message.CourseListReboot) {
	if len(msg.Records) == 0 {
		showNotReadyPanel()
		return
	}
	sPanel.list.Reboot(msg.Records)
	showSelectPanel()
}

// GetCourseToRemove message. Get the record to remove.

func (m *messenger) getRemoveCourseTX(recordID uint64) {
	msg := message.NewGetRemoveCourse(groupID, recordID)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) getRemoveCourseRX(msg *message.GetCourseToRemove) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
	fPanel.fillForm(msg.Record)
	showFormPanel()
}

// RemoveCourse message. Send the removed record.

func (m *messenger) courseRemoveTX(id uint64) {
	msg := message.NewRemoveCourse(groupID, id)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) courseRemoveRX(msg *message.RemoveCourse) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
	dialog.ShowInformation("Success", "Course removed.", window)
}
