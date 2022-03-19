package copyservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/josephbudd/okp/backend/model/sound"
)

var keyRunning bool
var keyCtxCancel context.CancelFunc
var errKeyRunning = fmt.Errorf("already keying")
var lock sync.Mutex

// Keying returns if the package is keying.
func Keying() (keying bool) {
	keying = keyRunning
	return
}

// Key a valid ditdah sentence, word or character at wpm after pausing.
// After the playing is completely finished, it returns the error.
// Param ditdah is a valid ditdah sentence, word or character where
//  * a sentence is words joined by help.ditdahWordSeperator.
//  * a word is a characters joined by help.ditdahCharacterSeperator.
//  * a character is a combination of "." and "-".
// Param wpm is the words per minute.
// Param pauseSeconds is the amount of time in seconds to pause before playing the morse code.
// Returns the error.
func Key(ctx context.Context, ctxCancel context.CancelFunc, ditdah string, wpm, wordspacewpm, pauseSeconds uint64) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("copyservice.KeyWords: %w", err)
		}
	}()

	if keyRunning {
		err = errKeyRunning
		return
	}
	if ditdah, err = sound.ValidateDitDah(ditdah); err != nil {
		return
	}
	keyCtxCancel = ctxCancel
	lock.Lock()
	keyRunning = true
	lock.Unlock()
	timeout := time.After(time.Second * time.Duration(pauseSeconds))
	<-timeout
	err = sound.PlayCW(ctx, ditdah, wpm, wordspacewpm)
	lock.Lock()
	keyRunning = false
	lock.Unlock()
	return
}

// StopKeying stops the keying.
func StopKeying() {
	if keyRunning {
		lock.Lock()
		keyRunning = false
		lock.Unlock()
		keyCtxCancel()
	}
}
