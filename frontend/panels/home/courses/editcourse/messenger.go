package editcourse

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"github.com/josephbudd/okp/frontend/txrx"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/store/record"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "editcourse"
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
	if err = txrx.AddListener(message.CourseToEditID, m); err != nil {
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
	case *message.GetCourseToEdit:
		m.getCourseToEditRX(msg)
	case *message.CourseToEdit:
		m.courseToEditRX(msg)
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

// GetCourseToEdit message. Get the record to edit.

func (m *messenger) getCourseToEditTX(recordID uint64) {
	msg := message.NewGetCourseToEdit(groupID, recordID)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) getCourseToEditRX(msg *message.GetCourseToEdit) {
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

// CourseEdit message. Send the edited record.

func (m *messenger) courseEditTX(r *record.CourseEdit) {
	msg := message.NewCourseToEdit(groupID, r)
	message.FrontEndToBackEnd <- msg
}

func (m *messenger) courseToEditRX(msg *message.CourseToEdit) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
	showSelectPanel()
	dialog.ShowInformation("Success", "Course saved.", window)
	// No error so let the form panel handle this.
}
