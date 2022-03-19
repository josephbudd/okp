package alsa

import (
	"fmt"
	"runtime"
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

func (d *device) createDevice(deviceName string, channels int, format Format, rate int, playback bool, bufferParams BufferParams) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("createDevice: %w", err)
		}
	}()

	deviceCString := C.CString(deviceName)
	defer C.free(unsafe.Pointer(deviceCString))
	var ret C.int
	if playback {
		ret = C.snd_pcm_open(&d.h, deviceCString, C.SND_PCM_STREAM_PLAYBACK, 0)
	} else {
		ret = C.snd_pcm_open(&d.h, deviceCString, C.SND_PCM_STREAM_CAPTURE, 0)
	}
	if ret < 0 {
		err = fmt.Errorf("could not open ALSA device %s", deviceName)
		return
	}
	runtime.SetFinalizer(d, (*device).Close)
	var hwParams *C.snd_pcm_hw_params_t
	ret = C.snd_pcm_hw_params_malloc(&hwParams)
	if ret < 0 {
		err = createError("could not alloc hw params", ret)
		return
	}
	defer C.snd_pcm_hw_params_free(hwParams)
	ret = C.snd_pcm_hw_params_any(d.h, hwParams)
	if ret < 0 {
		err = createError("could not set default hw params", ret)
		return
	}
	ret = C.snd_pcm_hw_params_set_access(d.h, hwParams, C.SND_PCM_ACCESS_RW_INTERLEAVED)
	if ret < 0 {
		err = createError("could not set access params", ret)
		return
	}
	ret = C.snd_pcm_hw_params_set_format(d.h, hwParams, C.snd_pcm_format_t(format))
	if ret < 0 {
		err = createError("could not set format params", ret)
		return
	}
	ret = C.snd_pcm_hw_params_set_channels(d.h, hwParams, C.uint(channels))
	if ret < 0 {
		err = createError("could not set channels params", ret)
		return
	}
	ret = C.snd_pcm_hw_params_set_rate(d.h, hwParams, C.uint(rate), 0)
	if ret < 0 {
		err = createError("could not set rate params", ret)
		return
	}
	var bufferSize = C.snd_pcm_uframes_t(bufferParams.BufferFrames)
	if bufferParams.BufferFrames == 0 {
		// Default buffer size: max buffer size
		ret = C.snd_pcm_hw_params_get_buffer_size_max(hwParams, &bufferSize)
		if ret < 0 {
			err = createError("could not get buffer size", ret)
			return
		}
	}
	ret = C.snd_pcm_hw_params_set_buffer_size_near(d.h, hwParams, &bufferSize)
	if ret < 0 {
		err = createError("could not set buffer size", ret)
		return
	}
	// Default period size: 1/8 of a second
	var periodFrames = C.snd_pcm_uframes_t(rate / 8)
	if bufferParams.PeriodFrames > 0 {
		periodFrames = C.snd_pcm_uframes_t(bufferParams.PeriodFrames)
	} else if bufferParams.Periods > 0 {
		periodFrames = C.snd_pcm_uframes_t(int(bufferSize) / bufferParams.Periods)
	}
	ret = C.snd_pcm_hw_params_set_period_size_near(d.h, hwParams, &periodFrames, nil)
	if ret < 0 {
		err = createError("could not set period size", ret)
		return
	}
	var periods = C.uint(0)
	ret = C.snd_pcm_hw_params_get_periods(hwParams, &periods, nil)
	if ret < 0 {
		err = createError("could not get periods", ret)
		return
	}
	ret = C.snd_pcm_hw_params(d.h, hwParams)
	if ret < 0 {
		err = createError("could not set hw params", ret)
		return
	}
	d.frames = int(periodFrames)
	d.Channels = channels
	d.Format = format
	d.Rate = rate
	d.BufferParams.BufferFrames = int(bufferSize)
	d.BufferParams.PeriodFrames = int(periodFrames)
	d.BufferParams.Periods = int(periods)
	return
}

// Close closes a device and frees the resources associated with it.
func (d *device) Close() {
	if d.h != nil {
		C.snd_pcm_drain(d.h)
		C.snd_pcm_close(d.h)
		d.h = nil
	}
	runtime.SetFinalizer(d, nil)
}

func (d *device) NChannels() (nChannels int) {
	switch d.Format {
	case FormatS8:
		nChannels = 2
	case FormatU8:
		nChannels = 1

	case FormatS16LE:
		nChannels = 2
	case FormatS16BE:
		nChannels = 2
	case FormatU16LE:
		nChannels = 1
	case FormatU16BE:
		nChannels = 1

	case FormatS24LE:
		nChannels = 2
	case FormatS24BE:
		nChannels = 2
	case FormatU24LE:
		nChannels = 1
	case FormatU24BE:
		nChannels = 1

	case FormatS32LE:
		nChannels = 2
	case FormatS32BE:
		nChannels = 2
	case FormatU32LE:
		nChannels = 1
	case FormatU32BE:
		nChannels = 1

	case FormatFloatLE:
		nChannels = d.Channels
	case FormatFloatBE:
		nChannels = d.Channels

	case FormatFloat64LE:
		nChannels = d.Channels
	case FormatFloat64BE:
		nChannels = d.Channels
	default:
		// nChannels = 0
	}
	return
}

func (d *device) IsLittleEndian() (isLittleEndian bool) {
	switch d.Format {
	case FormatS8:
	case FormatU8:

	case FormatS16LE:
		isLittleEndian = true
	case FormatS16BE:
	case FormatU16LE:
		isLittleEndian = true
	case FormatU16BE:

	case FormatS24LE:
		isLittleEndian = true
	case FormatS24BE:
	case FormatU24LE:
		isLittleEndian = true
	case FormatU24BE:

	case FormatS32LE:
		isLittleEndian = true
	case FormatS32BE:
	case FormatU32LE:
		isLittleEndian = true
	case FormatU32BE:

	case FormatFloatLE:
		isLittleEndian = true
	case FormatFloatBE:

	case FormatFloat64LE:
		isLittleEndian = true
	case FormatFloat64BE:

	default:
		// isLittleEndian = false
	}
	return
}

func (d *device) IsFloat() (isFloat bool) {
	switch d.Format {
	case FormatS8:
	case FormatU8:

	case FormatS16LE:
	case FormatS16BE:
	case FormatU16LE:
	case FormatU16BE:

	case FormatS24LE:
	case FormatS24BE:
	case FormatU24LE:
	case FormatU24BE:

	case FormatS32LE:
	case FormatS32BE:
	case FormatU32LE:
	case FormatU32BE:

	case FormatFloatLE:
		isFloat = true
	case FormatFloatBE:
		isFloat = true

	case FormatFloat64LE:
		isFloat = true
	case FormatFloat64BE:
		isFloat = true

	default:
		// isFloat = false
	}
	return
}

func (d device) FormatSampleSize() (nbytes int) {
	switch d.Format {
	case FormatS8, FormatU8:
		nbytes = 1
	case FormatS16LE, FormatS16BE, FormatU16LE, FormatU16BE:
		nbytes = 2
	case FormatS24LE, FormatS24BE, FormatU24LE, FormatU24BE:
		nbytes = 4
	case FormatS32LE, FormatS32BE, FormatU32LE, FormatU32BE, FormatFloatLE, FormatFloatBE:
		nbytes = 4
	case FormatFloat64LE, FormatFloat64BE:
		nbytes = 8
	default:
	}
	// if nbytes == 0 {
	// 	panic("unsupported format")
	// }
	return
}
