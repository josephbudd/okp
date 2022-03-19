package keyservice

import (
	"github.com/josephbudd/okp/backend/model/sound"
	"github.com/josephbudd/okp/shared/store/record"
)

// Milliseconds converts the solution to keyed times.
// This allows the renderer process to fake the correct answer.
// ..- ...-  ==> ..-...
// .-- -..-  ==> .-- ...
// -..- -.-- ==> ... -.--
func Milliseconds(solution [][]*record.KeyCode, wpm uint64) (milliSeconds []int64) {
	// pauses
	ditMS, dahMS, ddPauseMS, charPauseMS, wordPauseMS := sound.MilliSecondsFromWPM(wpm)
	// build the milliseconds and the solution
	milliSeconds = make([]int64, 0, 1024)
	for _, word := range solution {
		milliSeconds = append(milliSeconds, wordPauseMS)
		for i, char := range word {
			if i > 0 {
				// The second and following characters in the word
				//  are separated from the previous character
				//  with a character pause.
				milliSeconds = append(milliSeconds, charPauseMS)
			}
			ditdah := char.DitDah
			lastK := len(ditdah) - 1
			for k := 0; k <= lastK; k++ {
				dd := ditdah[k]
				// The sound.
				switch dd {
				case '.':
					milliSeconds = append(milliSeconds, ditMS)
				case '-':
					milliSeconds = append(milliSeconds, dahMS)
				}
				// The pause.
				if k < lastK {
					switch ditdah[k+1] {
					case ' ':
						// Space marks the end of the character.
						milliSeconds = append(milliSeconds, charPauseMS)
						// Skip over the ' ' to the next '.' or '-'.
						k++
					default:
						// The next char is a '.' or a '-'.
						milliSeconds = append(milliSeconds, ddPauseMS)
					}
				}
			}
		}
	}
	return
}
