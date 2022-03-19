package sound

import (
	"context"
	"fmt"
	"time"

	alsa "github.com/josephbudd/okp/backend/model/goalsa"
)

// Metronome clicks an element beat.
func Metronome(ctx context.Context, wpm uint64, errCh chan error) {

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
	ditDuration, _, _, _, _ := DurationsFromWPM(wpm)
	buffer := buildMetronomeSound(device, wpm)
	frames, bufPtr := device.FrameBuffer(buffer)
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

func buildMetronomeSound(device *alsa.PlaybackDevice, wpm uint64) (data []int16) {
	samplesPerElement := SamplesPerElement(wpm, device)
	// this metronome tick is 1/10 of an element.
	nSamples := samplesPerElement / 10
	dataSize := nSamples * device.Channels
	data = make([]int16, dataSize)
	// each sample is a frame of device.Channels int16.
	for i := 0; i < nSamples; i++ {
		for j := 1; j <= device.Channels; j++ {
			index := (device.Channels * i) + j - 1
			data[index] = int16((i%(j*128))*100 - 1000)
		}
	}
	return
}
