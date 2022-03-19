package txrx

import (
	"context"
	"fmt"
	"os"

	"github.com/josephbudd/okp/backend/model/keyservice"
	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const checkCurrentKeyTestF = "checkCurrentKeyTestRX: %s"

func init() {
	addListener(message.CheckCurrentKeyTestID, checkCurrentKeyTestRX)
}

func checkCurrentKeyTestRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	checkCurrentKeyTestMsg := msg.(*message.CheckCurrentKeyTest)
	appState := state.NewBackendState()
	var stateMsg *state.Message
	var err, fatal error
	defer func() {
		if err == nil && fatal == nil {
			if checkCurrentKeyTestMsg.Testing && checkCurrentKeyTestMsg.Passed {
				stateMsg, fatal = appState.PassCurrentKeyTest()
			}
		}
		switch {
		case err != nil:
			checkCurrentKeyTestMsg.Error = true
			checkCurrentKeyTestMsg.ErrorMessage = fmt.Sprintf(checkCurrentKeyTestF, err.Error())
		case fatal != nil:
			checkCurrentKeyTestMsg.Fatal = true
			checkCurrentKeyTestMsg.ErrorMessage = fmt.Sprintf(checkCurrentKeyTestF, fatal.Error())
		default:
			if stateMsg != nil {
				// New state because the user passed the test.
				appState.Dispatch(stateMsg)
			}
		}
		message.BackEndToFrontEnd <- checkCurrentKeyTestMsg
	}()

	keyTest := appState.CurrentKeyTest()
	checkCurrentKeyTestMsg.DitDahs = keyTest.DitDahs
	checkCurrentKeyTestMsg.Text = keyTest.Text
	testing := os.Getenv("CWT_TESTING") == "first" || os.Getenv("CWT_TESTING") == "last"

	// If testing this app then just return the correct answer.
	if checkCurrentKeyTestMsg.Testing && testing {
		checkCurrentKeyTestMsg.Passed = true
		checkCurrentKeyTestMsg.Copy = keyTest.Text
		return
	}

	var keycodes []*record.KeyCode
	if keycodes, fatal = stores.KeyCode.GetAll(); fatal != nil {
		return
	}

	// Each guess contains 4 guesses.
	//  * compressed-characters guess
	//  * combinded-characters guess
	//  * word guess
	//  * sentence guess.
	guess := keyservice.CopyMilliSeconds(checkCurrentKeyTestMsg.Times, appState.Speed().KeyWPM, keycodes)
	switch {
	case guess.Character != nil:
		// Guess #1
		// Is this word a compressed key code character? Ex: "A"
		checkCurrentKeyTestMsg.Passed = guess.Character.Character == keyTest.Text
		checkCurrentKeyTestMsg.Copy = guess.Character.Character
		checkCurrentKeyTestMsg.DitDahs = guess.Character.DitDah
	case guess.Compressed != nil:
		// Guess #2
		// Is this word a compressed key code character? Ex: "DE" or "SOS"
		checkCurrentKeyTestMsg.Passed = guess.Compressed.Character == keyTest.Text
		checkCurrentKeyTestMsg.Copy = guess.Compressed.Character
		checkCurrentKeyTestMsg.DitDahs = guess.Compressed.DitDah
	case guess.Combined != nil:
		// Guess #3
		// Is this word a combined key code character? Ex: "C Q"
		checkCurrentKeyTestMsg.Passed = guess.Combined.Character == keyTest.Text
		checkCurrentKeyTestMsg.Copy = guess.Combined.Character
		checkCurrentKeyTestMsg.DitDahs = guess.Combined.DitDah
	case len(guess.CopiedWord) > 0:
		// Guess #4
		// Is this a word? Ex: "Hello"
		checkCurrentKeyTestMsg.Passed = guess.CopiedWord == keyTest.Text
		checkCurrentKeyTestMsg.Copy = guess.CopiedWord
		checkCurrentKeyTestMsg.DitDahs = guess.CopiedWordDD
	case len(guess.CopiedSentence) > 0:
		// Guess #5
		// Is this a sentence? Ex: "Hello world."
		checkCurrentKeyTestMsg.Passed = guess.CopiedSentence == keyTest.Text
		checkCurrentKeyTestMsg.Copy = guess.CopiedSentence
		checkCurrentKeyTestMsg.DitDahs = guess.CopiedSentenceDD
	default:
		// No guesses. Can this happen?
	}
}
