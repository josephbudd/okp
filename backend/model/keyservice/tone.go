package keyservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/josephbudd/okp/backend/model/sound"
)

var (
	toneRunning       bool
	toneCtx           context.Context
	toneCancel        context.CancelFunc
	toneLock          sync.Mutex
	errAlreadyRunning = fmt.Errorf("a tone is already playing")
)

// StartTone starts the tone.
// If the tone is already running
//   then it just sends a nil error through the errCh
//   to stop the error handler
//   and returns.
func StartTone(ctx context.Context, wpm uint64, errCh chan error) {
	if toneRunning {
		// Currently already running the tone.
		// Send the nil error to the caller so it goes away.
		err := fmt.Errorf("keyservice.StartTone: %w", errAlreadyRunning)
		errCh <- err
		return
	}
	toneLock.Lock()
	toneCtx, toneCancel = context.WithCancel(ctx)
	// The tone is not running so start it.
	toneRunning = true
	toneLock.Unlock()
	go sound.Tone(toneCtx, wpm, errCh)
}

// StopTone stops the tone.
func StopTone() {
	if !toneRunning {
		// Nothing to stop.
		// Return to the caller so it goes away.
		return
	}
	toneLock.Lock()
	toneRunning = false
	toneLock.Unlock()
	toneCancel()
}
