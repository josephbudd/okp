package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const courseToEditF = "courseToEditRX: %s"

func init() {
	addListener(message.CourseToEditID, courseToEditRX)
}

func courseToEditRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	var stateMsg *state.Message
	appState := state.NewBackendState()
	courseToEditMsg := msg.(*message.CourseToEdit)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			courseToEditMsg.Error = true
			courseToEditMsg.ErrorMessage = fmt.Sprintf(courseToEditF, err.Error())
			message.BackEndToFrontEnd <- courseToEditMsg
		case fatal != nil:
			courseToEditMsg.Fatal = true
			courseToEditMsg.ErrorMessage = fmt.Sprintf(courseToEditF, fatal.Error())
			message.BackEndToFrontEnd <- courseToEditMsg
		default:
			message.BackEndToFrontEnd <- courseToEditMsg
			// Send the new state if there is one.
			if stateMsg != nil {
				appState.Dispatch(stateMsg)
			}
			// Tell the course lists to reboot.
			message.BackEndToFrontEnd <- message.NewCourseListReboot(stores)
		}
	}()

	var rr []*record.Course
	if rr, fatal = stores.Course.GetAll(); fatal != nil {
		return
	}
	r := record.FromCourseEdit(courseToEditMsg.Record, rr)
	if fatal = stores.Course.Update(r); fatal != nil {
		return
	}
	if appState.CurrentCourseID() == r.ID {
		// Just edited the current course so sync the application state.
		stateMsg, fatal = appState.Sync()
	}
}
