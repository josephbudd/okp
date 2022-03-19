package addcourse

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"github.com/josephbudd/okp/frontend/txrx"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/store/record"
)

type messenger struct{}

func (m *messenger) GroupName() (groupName string) {
	groupName = "addcourse"
	return
}

func (m *messenger) listen() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("listen: %w", err)
		}
	}()

	if err = txrx.AddListener(message.PlanListRebootID, m); err != nil {
		return
	}
	if err = txrx.AddListener(message.AddCourseID, m); err != nil {
		return
	}
	return
}

func (m *messenger) Listen(msg interface{}) {
	switch msg := msg.(type) {
	case *message.AddCourse:
		m.courseAddRX(msg)
	case *message.PlanListReboot:
		m.planListRebootRX(msg)
	}
}

// PlanList Reboot

func (m *messenger) planListRebootRX(msg *message.PlanListReboot) {
	fPanel.form.RebootPlan(msg.Records)
}

// CourseAdd message.

func (m *messenger) courseAddTX(r *record.CourseAdd) {
	message.FrontEndToBackEnd <- message.NewAddCourse(groupID, r)
}

func (m *messenger) courseAddRX(msg *message.AddCourse) {
	if msg.GroupID != groupID {
		return
	}
	if msg.Error {
		dialog.ShowInformation("Error", msg.ErrorMessage, window)
		return
	}
	dialog.ShowInformation("Success", "Course added.", window)
	fPanel.form.Clear()
}
