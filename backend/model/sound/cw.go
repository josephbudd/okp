package sound

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	alsa "github.com/josephbudd/okp/backend/model/goalsa"
	"github.com/josephbudd/okp/shared/help"
)

// ValidateDitDah checks that the ditdah is only made if '.', '-', separators.
func ValidateDitDah(ditdah string) (valid string, err error) {

	valids := make([]string, 0, len(ditdah))
	errs := make([]string, 0, len(ditdah))
	defer func() {
		valid = strings.Join(valids, "")
		if len(errs) > 0 {
			err = fmt.Errorf(strings.Join(errs, "\n"))
			err = fmt.Errorf("ValidateDitDah: %w", err)
		}
	}()

	// Check for bad chars in the dit-day string.
	charSepS, charSepR := help.DitdahCharacterSeperator()
	wordSepS, wordSepR := help.DitdahWordSeperator()
	for i, r := range ditdah {
		switch r {
		case '.', '-', charSepR, wordSepR:
			valids = append(valids, string(r))
		default:
			msg := fmt.Sprintf("%q is not a valid dit dah character at position %d. It should be a \".\", a \"-\", a %q or a %q", r, i, charSepS, wordSepS)
			errs = append(errs, msg)
		}
	}
	return
}

// PlayCW plays ditdah at wpm words per minute.
//  * a sentence is words joined by help.ditdahWordSeperator.
//  * a word is a characters joined by help.ditdahCharacterSeperator.
//  * a character is a combination of "." and "-".
// Param ctx is the context for this play.
// Param ditdah is the ditdah sentence, word or character.
// Param wpm is the words per minute.
func PlayCW(ctx context.Context, ditdah string, wpm, wordSpaceWPM uint64) (err error) {
	// cw converts ".- -."  to sound

	var playBackError error
	defer func() {
		if playBackError != nil {
			err = fmt.Errorf("sound.PlayCW: %w", playBackError)
			if err != deviceErr {
				deviceErr = err
			}
		}
	}()

	if playBackError = deviceErr; playBackError != nil {
		return
	}
	// Get the play back device.
	device := PlaybackDevice()

	// The separator characters.
	_, charSepR := help.DitdahCharacterSeperator()
	_, wordSepR := help.DitdahWordSeperator()

	// word == "paris" == 50 elements.
	nElementsPerMinute := Float64ElementsPerMinute(wpm)
	nElementsPerSecond := math.Floor(nElementsPerMinute / 60.0)
	nElementsPerHalfMinute := math.Floor(nElementsPerMinute / 2.0)
	secondsPerElement := SecondsPerElement(wpm)

	// Sounds.
	ditSound := buildCWSound(device, wpm, 1)
	ditFrames, ditBufPtr := device.FrameBuffer(ditSound)
	dahSound := buildCWSound(device, wpm, 3)
	dahFrames, dahBufPtr := device.FrameBuffer(dahSound)
	oneSecondSound := buildCWSound(device, wpm, uint64(nElementsPerSecond))
	oneSecondFrames, oneSecondBufPtr := device.FrameBuffer(oneSecondSound)
	thirySecondSound := buildCWSound(device, wpm, uint64(nElementsPerHalfMinute))
	thirySecondFrames, thirySecondBufPtr := device.FrameBuffer(thirySecondSound)
	// Waits for a sound to end.
	ditSoundWait := time.Duration(secondsPerElement * float64(time.Second))
	dahSoundWait := time.Duration((3.0 * secondsPerElement) * float64(time.Second))
	oneSecondSoundWait := time.Second
	thirtySecondSoundWait := time.Duration(time.Second / 2)
	// The silent pauses.
	preCharSilentWait := time.Duration(secondsPerElement * float64(time.Second))
	charSeparatorSilentWait := time.Duration((3.0 * secondsPerElement) * float64(time.Second))
	wordSeparatorSilentWait := time.Duration((7.0 * secondsPerElement) * float64(time.Second))
	// Timers.
	var preCharTimeout <-chan time.Time
	var soundTimeout <-chan time.Time
	var pauseTimeout <-chan time.Time

	for _, r := range ditdah {
		switch r {
		case charSepR:
			// char separator : 3 elements of silence.
			pauseTimeout = time.After(charSeparatorSilentWait)
			<-pauseTimeout
		case wordSepR:
			// word separator : 7 elements of silence.
			pauseTimeout = time.After(wordSeparatorSilentWait)
			<-pauseTimeout
		case '.':
			// Period is used as a dit.
			// Precharacter pause
			preCharTimeout = time.After(preCharSilentWait)
			<-preCharTimeout
			// Dit sound.
			if _, playBackError = device.WriteBuffer(ditFrames, ditBufPtr); playBackError != nil {
				return
			}
			// Wait while the dit sound plays.
			soundTimeout = time.After(ditSoundWait)
			<-soundTimeout
		case '-':
			// Dash is used as a dah.
			// Precharacter pause
			preCharTimeout = time.After(preCharSilentWait)
			<-preCharTimeout
			// Dah sound.
			if _, playBackError = device.WriteBuffer(dahFrames, dahBufPtr); playBackError != nil {
				return
			}
			// Wait while the dah sound plays.
			soundTimeout = time.After(dahSoundWait)
			<-soundTimeout
		case 's':
			// 1 second sound.
			if _, playBackError = device.WriteBuffer(oneSecondFrames, oneSecondBufPtr); playBackError != nil {
				return
			}
			// Wait while the 1 second sound plays.
			soundTimeout := time.After(oneSecondSoundWait)
			<-soundTimeout
		case 'm':
			// 30 second sound.
			if _, playBackError = device.WriteBuffer(thirySecondFrames, thirySecondBufPtr); playBackError != nil {
				return
			}
			// Wait while the 30 second sound plays.
			soundTimeout := time.After(thirtySecondSoundWait)
			<-soundTimeout
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	return
}

func buildCWSound(device *alsa.PlaybackDevice, wpm uint64, nElements uint64) (data []int16) {
	samplesPerElement := SamplesPerElement(wpm, device)
	nSamples := samplesPerElement * int(nElements)
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
