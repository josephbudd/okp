package txrx

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

const checkCurrentCopyTestF = "checkCurrentCopyTestRX: %s"

func init() {
	addListener(message.CheckCurrentCopyTestID, checkCurrentCopyTestRX)
}

func checkCurrentCopyTestRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	checkCurrentCopyTestMsg := msg.(*message.CheckCurrentCopyTest)
	appState := state.NewBackendState()
	var stateMsg *state.Message
	var err, fatal error
	defer func() {
		if err == nil && fatal == nil {
			if checkCurrentCopyTestMsg.Testing && checkCurrentCopyTestMsg.Passed {
				stateMsg, fatal = appState.PassCurrentCopyTest()
			}
		}
		switch {
		case err != nil:
			checkCurrentCopyTestMsg.Error = true
			checkCurrentCopyTestMsg.ErrorMessage = fmt.Sprintf(checkCurrentCopyTestF, err.Error())
		case fatal != nil:
			checkCurrentCopyTestMsg.Fatal = true
			checkCurrentCopyTestMsg.ErrorMessage = fmt.Sprintf(checkCurrentCopyTestF, fatal.Error())
		default:
			if stateMsg != nil {
				// New state because the user passed the test.
				appState.Dispatch(stateMsg)
			}
		}
		message.BackEndToFrontEnd <- checkCurrentCopyTestMsg
	}()

	copyTest := appState.CurrentKeyTest()
	checkCurrentCopyTestMsg.Text = copyTest.Text
	checkCurrentCopyTestMsg.DitDahs = copyTest.DitDahs
	testing := os.Getenv("CWT_TESTING") == "first" || os.Getenv("CWT_TESTING") == "last"

	// If testing this app then just return the correct answer.
	if checkCurrentCopyTestMsg.Testing && testing {
		checkCurrentCopyTestMsg.Passed = true
		checkCurrentCopyTestMsg.Copy = copyTest.Text
		return
	}

	userCopy := strings.ToUpper(strings.TrimSpace(checkCurrentCopyTestMsg.Copy))
	log.Printf("copyTest.Text:%q, userCopy:%q", copyTest.Text, userCopy)
	checkCurrentCopyTestMsg.Passed = copyTest.Text == userCopy
}
