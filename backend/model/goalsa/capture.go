package alsa

import (
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

// CaptureDevice is an ALSA device configured to record audio.
type CaptureDevice struct {
	device
}

// NewCaptureDevice creates a new CaptureDevice object.
func NewCaptureDevice(deviceName string, channels int, format Format, rate int, bufferParams BufferParams) (c *CaptureDevice, err error) {
	c = new(CaptureDevice)
	if err = c.createDevice(deviceName, channels, format, rate, false, bufferParams); err != nil {
		err = fmt.Errorf("alsa.NewCaptureDevice: %w", err)
	}
	return
}

// Read reads samples into a buffer and returns the amount read.
func (c *CaptureDevice) Read(buffer interface{}) (nSamplesRead int, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("alsa.Read: %w", err)
		}
	}()

	if c.checkBufferType(buffer); err != nil {
		return
	}

	frames, bufPtr := c.FrameBuffer(buffer)
	ret := C.snd_pcm_readi(c.h, bufPtr, frames)

	if ret == -C.EPIPE {
		C.snd_pcm_prepare(c.h)
		err = ErrOverrun
		return
	} else if ret < 0 {
		err = createError("read error", C.int(ret))
		return
	}
	nSamplesRead = int(ret) * c.Channels
	return
}

func (c *CaptureDevice) checkBufferType(buffer interface{}) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CaptureDevice.checkBufferType: %w", err)
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
		if c.FormatSampleSize() != 1 {
			err = sizeError
			return
		}
	case reflect.Int16:
		if c.FormatSampleSize() != 2 {
			err = sizeError
			return
		}
	case reflect.Int32, reflect.Float32:
		if c.FormatSampleSize() != 4 {
			err = sizeError
			return
		}
	case reflect.Float64:
		if c.FormatSampleSize() != 8 {
			err = sizeError
			return
		}
	default:
		err = fmt.Errorf("buffer format not supported")
		return
	}
	return
}

func (c *CaptureDevice) FrameBuffer(buffer interface{}) (frames C.snd_pcm_uframes_t, bufPtr unsafe.Pointer) {
	val := reflect.ValueOf(buffer)
	length := val.Len()
	sliceData := val.Slice(0, length)

	frames = C.snd_pcm_uframes_t(length / c.Channels)
	bufPtr = unsafe.Pointer(sliceData.Index(0).Addr().Pointer())
	return
}

func (c *CaptureDevice) NChannels() (nChannels int) {
	nChannels = c.Channels
	return
}
