package alsa

import (
	"fmt"
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

// Format is the type used for specifying sample formats.
type Format C.snd_pcm_format_t

// The range of sample formats supported by ALSA.
const (
	FormatS8        = C.SND_PCM_FORMAT_S8
	FormatU8        = C.SND_PCM_FORMAT_U8
	FormatS16LE     = C.SND_PCM_FORMAT_S16_LE
	FormatS16BE     = C.SND_PCM_FORMAT_S16_BE
	FormatU16LE     = C.SND_PCM_FORMAT_U16_LE
	FormatU16BE     = C.SND_PCM_FORMAT_U16_BE
	FormatS24LE     = C.SND_PCM_FORMAT_S24_LE
	FormatS24BE     = C.SND_PCM_FORMAT_S24_BE
	FormatU24LE     = C.SND_PCM_FORMAT_U24_LE
	FormatU24BE     = C.SND_PCM_FORMAT_U24_BE
	FormatS32LE     = C.SND_PCM_FORMAT_S32_LE
	FormatS32BE     = C.SND_PCM_FORMAT_S32_BE
	FormatU32LE     = C.SND_PCM_FORMAT_U32_LE
	FormatU32BE     = C.SND_PCM_FORMAT_U32_BE
	FormatFloatLE   = C.SND_PCM_FORMAT_FLOAT_LE
	FormatFloatBE   = C.SND_PCM_FORMAT_FLOAT_BE
	FormatFloat64LE = C.SND_PCM_FORMAT_FLOAT64_LE
	FormatFloat64BE = C.SND_PCM_FORMAT_FLOAT64_BE
)

var (
	// ErrOverrun signals an overrun error
	ErrOverrun = fmt.Errorf("overrun")
	// ErrUnderrun signals an underrun error
	ErrUnderrun = fmt.Errorf("underrun")
)

// BufferParams specifies the buffer parameters of a device.
type BufferParams struct {
	BufferFrames int
	PeriodFrames int
	Periods      int
}

type device struct {
	h            *C.snd_pcm_t
	Channels     int
	Format       Format
	Rate         int
	BufferParams BufferParams
	frames       int
}
