package alsa

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"
)

/*
// Use the newer ALSA API
#define ALSA_PCM_NEW_HW_PARAMS_API
#cgo LDFLAGS: -lasound
#cgo CFLAGS: -Iinclude
#include <alsa/asoundlib.h>
#include <stdint.h>
*/
import "C"

// PlaybackDevice is an ALSA device configured to playback audio.
type PlaybackDevice struct {
	device
}

// NewDefaultPlaybackDevice creates the default play back device.
// It implements 16 bit stereo.
func NewDefaultPlaybackDevice() (p *PlaybackDevice, err error) {
	p = new(PlaybackDevice)
	if err = p.createDevice("default", 2, FormatS16BE, 44100, true, BufferParams{}); err != nil {
		err = fmt.Errorf("alsa.NewDefaultPlaybackDevice: %w", err)
	}
	return
}

// NewPlaybackDevice creates a new PlaybackDevice object.
func NewPlaybackDevice(deviceName string, channels int, format Format, rate int, bufferParams BufferParams) (p *PlaybackDevice, err error) {
	p = new(PlaybackDevice)
	if err = p.createDevice(deviceName, channels, format, rate, true, bufferParams); err != nil {
		err = fmt.Errorf("alsa.NewPlaybackDevice: %w", err)
	}
	return
}

// Write writes a buffer of data to a playback device.
// Returns the number of sample written and the error after writing is completed.
// func (p *PlaybackDevice) Write(buffer interface{}) (nSamplesWritten int, err error) {

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("PlaybackDevice.Write: %w", err)
// 		}
// 	}()

// 	if err = p.checkBufferType(buffer); err != nil {
// 		return
// 	}
// 	frames, bufPtr := p.FrameBuffer(buffer)
// 	nSamplesWritten, err = p.WriteBuffer(frames, bufPtr)
// 	return
// }

// WriteWhile writes a buffer of data to a playback device.
// Returns the number of sample written and the error after writing is completed.
func (p *PlaybackDevice) WriteWhile(ctx context.Context, errChan chan error, buffer interface{}) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlaybackDevice.WriteWhile: %w", err)
			errChan <- err
		}
	}()

	if err = p.checkBufferType(buffer); err != nil {
		return
	}
	frames, bufPtr := p.FrameBuffer(buffer)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err = p.WriteBuffer(frames, bufPtr); err != nil {
				return
			}
		}
	}
}

func (p *PlaybackDevice) checkBufferType(buffer interface{}) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlaybackDevice.checkBufferType: %w", err)
		}
	}()

	bufferType := reflect.TypeOf(buffer)
	if !(bufferType.Kind() == reflect.Array || bufferType.Kind() == reflect.Slice) {
		err = fmt.Errorf("buffer is not an array type")
		return
	}

	sizeError := fmt.Errorf("buffer does not match sample size")
	switch bufferType.Elem().Kind() {
	case reflect.Int8:
		if p.FormatSampleSize() != 1 {
			err = sizeError
			return
		}
	case reflect.Int16:
		if p.FormatSampleSize() != 2 {
			err = sizeError
			return
		}
	case reflect.Int32, reflect.Float32:
		if p.FormatSampleSize() != 4 {
			err = sizeError
			return
		}
	case reflect.Float64:
		if p.FormatSampleSize() != 8 {
			err = sizeError
			return
		}
	default:
		err = fmt.Errorf("buffer format not supported")
		return
	}
	return
}

func (p *PlaybackDevice) FrameBuffer(buffer interface{}) (frames C.snd_pcm_uframes_t, bufPtr unsafe.Pointer) {
	val := reflect.ValueOf(buffer)
	length := val.Len()
	sliceData := val.Slice(0, length)

	frames = C.snd_pcm_uframes_t(length / p.Channels)
	bufPtr = unsafe.Pointer(sliceData.Index(0).Addr().Pointer())
	return
}

func (p *PlaybackDevice) WriteBuffer(frames C.snd_pcm_uframes_t, bufPtr unsafe.Pointer) (nSamplesWritten int, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlaybackDevice.WriteBuffer: %w", err)
		}
	}()

	var ret C.long
	if ret = C.snd_pcm_writei(p.h, bufPtr, frames); ret == -C.EPIPE {
		C.snd_pcm_prepare(p.h)
		ret = C.snd_pcm_writei(p.h, bufPtr, frames)
	}
	if ret < 0 {
		err = createError("write error", C.int(ret))
		return
	}
	// No errors.
	nSamplesWritten = int(ret) * p.Channels
	return
}
