package sound

import (
	"context"
	"fmt"
	"time"
)

// Tone plays a sound.
func Tone(ctx context.Context, wpm uint64, errCh chan error) {

	var playBackError error
	defer func() {
		var err error
		if playBackError != nil {
			err = fmt.Errorf("sound.Metronome: %w", playBackError)
			if err != deviceErr {
				deviceErr = err
			}
		}
		// Signal end of func with or without an error.
		errCh <- err
	}()

	device := PlaybackDevice()
	buffer := buildCWSound(device, wpm, 1)
	frames, bufPtr := device.FrameBuffer(buffer)
	ditDuration, _, _, _, _ := DurationsFromWPM(wpm)
	timer := time.NewTimer(ditDuration)
	for {
		if _, playBackError = device.WriteBuffer(frames, bufPtr); playBackError != nil {
			return
		}
		select {
		case <-timer.C:
			timer.Reset(ditDuration)
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}
