package txrx

import (
	"context"
	"fmt"
	"os"

	"github.com/josephbudd/okp/backend/model/copyservice"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

const KeyF = "KeyRX: %s"

func init() {
	addListener(message.KeyID, keyRX)
}

func keyRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	keyMsg := msg.(*message.Key)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			keyMsg.Error = true
			keyMsg.ErrorMessage = fmt.Sprintf(KeyF, err.Error())
			message.BackEndToFrontEnd <- keyMsg
		case fatal != nil:
			keyMsg.Fatal = true
			keyMsg.ErrorMessage = fmt.Sprintf(KeyF, fatal.Error())
			message.BackEndToFrontEnd <- keyMsg
		default:
			message.BackEndToFrontEnd <- keyMsg
		}
	}()

	if !keyMsg.Run {
		// Run is false which means stop keying.
		copyservice.StopKeying()
		return
	}

	// Start keying.
	// Get the current ditdahs from the current copy test and play them.
	appState := state.NewBackendState()
	copyTest := appState.CurrentCopyTest()
	// Start keying. Return when the keying is finished.
	newCtx, newCtxCancel := context.WithCancel(ctx)
	// Start the go func which will key the morse code.
	// Then return with no error indicating to the front end that the morse code is being keyed.
	copyWPM := appState.Speed().CopyWPM
	copySpaceWPM := appState.Speed().CopySpaceWPM
	testing := os.Getenv("CWT_TESTING") == "first" || os.Getenv("CWT_TESTING") == "last"
	if keyMsg.Testing && testing {
		// Only testing the front end.
		copyWPM = 30
		copySpaceWPM = copyWPM
	}
	go func(ctx context.Context, ctxCancel context.CancelFunc, ditdah string, wpm, spacewpm, pauseSeconds uint64) {
		err = copyservice.Key(ctx, ctxCancel, ditdah, wpm, spacewpm, pauseSeconds)
		// Send any errors back to the front end.
		if err != nil {
			keyMsg.Error = true
			keyMsg.ErrorMessage = fmt.Sprintf(KeyF, err.Error())
		}
		keyMsg.Run = false // Done keying.
		message.BackEndToFrontEnd <- keyMsg
	}(newCtx, newCtxCancel, copyTest.DitDahs, copyWPM, copySpaceWPM, appState.Delay())
}
