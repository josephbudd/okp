package sound

import (
	"context"
	"fmt"
	"log"

	alsa "github.com/josephbudd/okp/backend/model/goalsa"
)

const (
	SampleRate = 44100
	Frequency  = 440
	Format     = alsa.FormatS16LE
)

var (
	playBackDevice *alsa.PlaybackDevice
	deviceErr      error
)

// PlaybackDevice returns the play back device and any evxisting device errors.
func PlaybackDevice() (device *alsa.PlaybackDevice) {
	device = playBackDevice
	return
}

// Open the play back device.
func Open(ctx context.Context) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("sound.Open: %w", err)
			deviceErr = err
		}
	}()

	// if playBackDevice, deviceErr = alsa.NewPlaybackDevice("default", 2, alsa.FormatS32LE, SampleRate, alsa.BufferParams{}); err != nil {
	// 	return
	// }

	if playBackDevice, deviceErr = alsa.NewPlaybackDevice(
		"default",
		1, alsa.FormatS16LE,
		SampleRate,
		alsa.BufferParams{},
	); err != nil {
		return
	}

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			Close()
		}
	}(ctx)
	return
}

// Close the play back device.
func Close() {
	if deviceErr != nil {
		log.Println("alsa play back device previously closed")
		return
	}
	playBackDevice.Close()
	log.Println("alsa play back device closed")
}
