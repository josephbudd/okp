package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/backend/model"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const addCourseF = "addCourseRX: %s"

func init() {
	addListener(message.AddCourseID, addCourseRX)
}

func addCourseRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	addCourseMsg := msg.(*message.AddCourse)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			addCourseMsg.Error = true
			addCourseMsg.ErrorMessage = fmt.Sprintf(addCourseF, err.Error())
			message.BackEndToFrontEnd <- addCourseMsg
		case fatal != nil:
			addCourseMsg.Fatal = true
			addCourseMsg.ErrorMessage = fmt.Sprintf(addCourseF, fatal.Error())
			message.BackEndToFrontEnd <- addCourseMsg
		default:
			message.BackEndToFrontEnd <- addCourseMsg
			// No errors so tell the course lists to reboot.
			message.BackEndToFrontEnd <- message.NewCourseListReboot(stores)
		}
	}()

	if len(addCourseMsg.Record.Name) == 0 {
		err = fmt.Errorf("Name is a required field.")
		return
	}
	if len(addCourseMsg.Record.Description) == 0 {
		err = fmt.Errorf("Description is a required field.")
		return
	}

	var plan *record.Plan
	if plan, fatal = stores.Plan.Get(addCourseMsg.Record.PlanID); fatal != nil {
		return
	}
	_, fatal = model.NewCourse(
		addCourseMsg.Record.Name, addCourseMsg.Record.Description,
		addCourseMsg.Record.SpeedID,
		plan,
		stores,
	)
}
