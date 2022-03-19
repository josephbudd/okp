package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

const initF = "initRX: %s"

func init() {
	addListener(message.InitID, initRX)
}

func initRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	initMsg := msg.(*message.Init)
	var stateMsg *state.Message
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			initMsg.Error = true
			initMsg.ErrorMessage = fmt.Sprintf(initF, err.Error())
			message.BackEndToFrontEnd <- initMsg
		case fatal != nil:
			initMsg.Fatal = true
			initMsg.ErrorMessage = fmt.Sprintf(initF, fatal.Error())
			message.BackEndToFrontEnd <- initMsg
		default:
			// No errors.
		}
	}()

	appState := state.NewBackendState()
	stateMsg = appState.ToInitialStateMessage()
	appState.Dispatch(stateMsg)
	message.BackEndToFrontEnd <- message.NewCourseListReboot(stores)
	message.BackEndToFrontEnd <- message.NewPlanListReboot(stores)
}
