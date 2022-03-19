package sound

import (
	"context"
	"fmt"
	"math"
	"time"

	alsa "github.com/josephbudd/okp/backend/model/goalsa"
	"github.com/josephbudd/okp/shared/help"
)

var (
	tau = math.Pi * 2
)

// PlayCWF plays ditdah at wpm words per minute.
//  * a sentence is words joined by help.ditdahWordSeperator.
//  * a word is a characters joined by help.ditdahCharacterSeperator.
//  * a character is a combination of "." and "-".
// Param ctx is the context for this play.
// Param ditdah is the ditdah sentence, word or character.
// Param wpm is the words per minute.
func PlayCWF(ctx context.Context, ditdah string, wpm, wordSpaceWPM uint64) (err error) {

	// Open the player and play the dit-dahs.
	var device *alsa.PlaybackDevice
	// sampling rate must be 44100 for cd sound range and quality.
	// sampling rate is not frequency.
	if device, err = alsa.NewPlaybackDevice("default", 2, alsa.FormatS32LE, SampleRate, alsa.BufferParams{}); err != nil {
		return
	}
	_, charSepR := help.DitdahCharacterSeperator()
	_, wordSepR := help.DitdahWordSeperator()
	err = cwF(ctx, device, wpm, wordSpaceWPM, ditdah, charSepR, wordSepR)
	device.Close()
	return
}

// cw converts ".- -."  to sound
func cwF(ctx context.Context, device *alsa.PlaybackDevice, wpm, spacewpm uint64, ditdah string, charSep, wordSep rune) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("cw: %w", err)
		}
	}()

	// word == "paris" == 50 elements.
	nElementsPerMinute := Float64ElementsPerMinute(wpm)
	nElementsPerSecond := math.Floor(nElementsPerMinute / 60.0)
	nElementsPerHalfMinute := math.Floor(nElementsPerMinute / 2.0)
	secondsPerElement := SecondsPerElement(wpm)

	// build the cw data to be played
	ditSound := buildCWSoundF32(device, wpm, 1)
	ditFrames, ditBufPtr := device.FrameBuffer(ditSound)
	dahSound := buildCWSoundF32(device, wpm, 3)
	dahFrames, dahBufPtr := device.FrameBuffer(dahSound)
	oneSecondSound := buildCWSoundF32(device, wpm, nElementsPerSecond)
	oneSecondFrames, oneSecondBufPtr := device.FrameBuffer(oneSecondSound)
	thirySecondSound := buildCWSoundF32(device, wpm, nElementsPerHalfMinute)
	thirySecondFrames, thirySecondBufPtr := device.FrameBuffer(thirySecondSound)

	// wordSilenceCount adjust the word separation using the spacing wpm.
	// a word separator is 7 elements of silence.
	// so the selence for the loop below is 7(adjusted) - 1
	wordSilenceCount := int64(math.Floor(float64(7)*(float64(spacewpm)/float64(wpm)))) - 1
	// loop through the runes in the ditdah string.
	var soundCount int64
	var silenceCount int64
	for _, r := range ditdah {
		switch r {
		case charSep:
			// char separator : 3 elements of silence.
			soundCount = 0
			silenceCount = 2 // 3 - the 1 silence that followed the last dit or dah.
		case wordSep:
			// word separator : 7 elements of silence.
			soundCount = 0
			// 7 - the 1 silence that followed the last dit or dah.
			// silenceCount = 6
			silenceCount = wordSilenceCount
		case '.':
			// period is used as a dit.
			// dit : 1 element of sound followed by 1 element of silence
			// play the dit.
			if _, err = device.WriteBuffer(ditFrames, ditBufPtr); err != nil {
				return
			}
			soundCount = 1
			silenceCount = 1
		case '-':
			// dash is used as a dah.
			// dah : 3 elements of sound followed by 1 element of silence.
			// play the dah.
			if _, err = device.WriteBuffer(dahFrames, dahBufPtr); err != nil {
				return
			}
			soundCount = 3
			silenceCount = 1
		case 's':
			// 1 second dit for testing timing.
			if _, err = device.WriteBuffer(oneSecondFrames, oneSecondBufPtr); err != nil {
				return
			}
			soundCount = int64(nElementsPerSecond)
			silenceCount = 0
		case 'm':
			// 30 second dah for testing timing.
			if _, err = device.WriteBuffer(thirySecondFrames, thirySecondBufPtr); err != nil {
				return
			}
			soundCount = int64(nElementsPerHalfMinute)
			silenceCount = 0
		}
		pauseCount := secondsPerElement * float64(silenceCount+soundCount)
		pause := time.Duration(pauseCount * float64(time.Second))
		if pause > 0 {
			timeout := time.After(pause)
			<-timeout
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	return
}

func buildCWSoundF32(device *alsa.PlaybackDevice, wpm uint64, nElements float64) (data []float32) {
	samplesPerElement := SamplesPerElement(wpm, device)
	nSamples := samplesPerElement * int(nElements)
	dataSize := nSamples * device.Channels
	data = make([]float32, dataSize)
	var angle float64 = tau / float64(nSamples)
	// each sample is a frame of device.Channels int16.
	for i := 0; i < nSamples; i++ {
		for j := 1; j <= device.Channels; j++ {
			index := (device.Channels * i) + j - 1
			sample := math.Sin(angle * float64(i))
			// data[index] = int16((i%(j*128))*100 - 1000)
			data[index] = float32(sample)
		}
	}
	return
}
