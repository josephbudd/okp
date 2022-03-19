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

func createError(errorMsg string, errorCode C.int) (err error) {
	strError := C.GoString(C.snd_strerror(errorCode))
	err = fmt.Errorf("%s: %s", errorMsg, strError)
	return
}
