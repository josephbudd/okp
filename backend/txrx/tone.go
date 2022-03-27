package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/backend/model/keyservice"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

const toneF = "toneRX: %s"

func init() {
	addListener(message.ToneID, toneRX)
}

func toneRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	toneMsg := msg.(*message.Tone)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			toneMsg.Error = true
			toneMsg.ErrorMessage = fmt.Sprintf(toneF, err.Error())
			message.BackEndToFrontEnd <- toneMsg
		case fatal != nil:
			toneMsg.Fatal = true
			toneMsg.ErrorMessage = fmt.Sprintf(toneF, fatal.Error())
			message.BackEndToFrontEnd <- toneMsg
		default:
			message.BackEndToFrontEnd <- toneMsg
		}
	}()

	if !toneMsg.TurnOn {
		// Turn off the tone.
		keyservice.StopTone()
		return
	}
	// Return the message with no error.
	message.BackEndToFrontEnd <- toneMsg

	// Get the current speed.
	appState := state.NewBackendState()
	speed := appState.Speed()
	// Start the tone.
	errorCh := make(chan error)
	keyservice.StartTone(ctx, speed.KeyWPM, errorCh)
	// Wait for sound.Tone(toneCtx, wpm, errCh) to return the error if any.
	fatal = <-errorCh
	// Return any tone errors.
}
