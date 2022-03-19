package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

const setCurrentCourseF = "setCurrentCourseRX: %s"

func init() {
	addListener(message.SetCurrentCourseID, setCurrentCourseRX)
}

func setCurrentCourseRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	setCurrentCourseMsg := msg.(*message.SetCurrentCourse)
	var err, fatal error
	appState := state.NewBackendState()
	var stateMsg *state.Message
	defer func() {
		switch {
		case err != nil:
			setCurrentCourseMsg.Error = true
			setCurrentCourseMsg.ErrorMessage = fmt.Sprintf(setCurrentCourseF, err.Error())
			// Send the message back.
			message.BackEndToFrontEnd <- setCurrentCourseMsg
		case fatal != nil:
			setCurrentCourseMsg.Fatal = true
			setCurrentCourseMsg.ErrorMessage = fmt.Sprintf(setCurrentCourseF, fatal.Error())
			// Send the message back.
			message.BackEndToFrontEnd <- setCurrentCourseMsg
		default:
			// Send the message back.
			message.BackEndToFrontEnd <- setCurrentCourseMsg
			appState.Dispatch(stateMsg)
			message.BackEndToFrontEnd <- message.NewCourseListReboot(stores)
		}
	}()

	stateMsg, fatal = appState.SetCurrentCourseID(setCurrentCourseMsg.RecordID)
}
