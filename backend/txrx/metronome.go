package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/backend/model/keyservice"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

const metronomeF = "metronomeRX: %s"

func init() {
	addListener(message.MetronomeID, metronomeRX)
}

func metronomeRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	metronomeMsg := msg.(*message.Metronome)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			metronomeMsg.Error = true
			metronomeMsg.ErrorMessage = fmt.Sprintf(metronomeF, err.Error())
			message.BackEndToFrontEnd <- metronomeMsg
		case fatal != nil:
			metronomeMsg.Fatal = true
			metronomeMsg.ErrorMessage = fmt.Sprintf(metronomeF, fatal.Error())
			message.BackEndToFrontEnd <- metronomeMsg
		default:
			message.BackEndToFrontEnd <- metronomeMsg
		}
	}()

	if !metronomeMsg.TurnOn {
		// Turn off the metronome.
		keyservice.StopMetronome()
		return
	}
	// Return the message with no error.
	message.BackEndToFrontEnd <- metronomeMsg

	// Get the current speeds.
	appState := state.NewBackendState()
	speed := appState.Speed()
	// Start the metronome.
	errorCh := make(chan error)
	keyservice.StartMetronome(ctx, speed.KeyWPM, errorCh)
	// keyservice.StartMetronome calls go sound.Metronome(metronomeCtx, wpm, errCh)
	// Wait for sound.Metronome(metronomeCtx, wpm, errCh) to return the error if any.
	fatal = <-errorCh
	// Return any metronome errors.
}
