package sound

import (
	"math"
	"time"

	alsa "github.com/josephbudd/okp/backend/model/goalsa"
)

const (
	elementsInParis = uint64(50)
	// P = di da da di = 1 {1} 3 {1} 3 {1} 1 (3) = 14 elements
	// A = di da = 1 {1} 3 (3) = 8 elements
	// R = di da di = 1 {1} 3 {1} 1 (3) = 10 elements
	// I = di di = 1 {1} 1 (3) = 6 elements
	// S = di di di = 1 {1} 1 {1} 1 [7] = 12 elements
	// Total = 50 elements
	// {1} = interDitDah pause
	// (3) = intercharacter pause
	// [7] = interword pause
	// Source Credit: http://www.kent-engineers.com/codespeed.htm
)

// MilliSecondsFromWPM returns the important cw time spanws as an int64 number of milliseconds.
func MilliSecondsFromWPM(wpm uint64) (dit, dah, ddPause, charPause, wordPause int64) {
	ditD, dahD, ddPauseD, charPauseD, wordPauseD := DurationsFromWPM(wpm)
	dit = ditD.Milliseconds()
	dah = dahD.Milliseconds()
	ddPause = ddPauseD.Milliseconds()
	charPause = charPauseD.Milliseconds()
	wordPause = wordPauseD.Milliseconds()
	return
}

// Float64MilliSecondsFromWPM returns the important cw time spans as float64s.
func Float64MilliSecondsFromWPM(wpm uint64) (dit, dah, ddPause, charPause, wordPause float64) {
	ditI64, dahI64, ddPauseI64, charPauseI64, wordPauseI64 := MilliSecondsFromWPM(wpm)
	dit = float64(ditI64)
	dah = float64(dahI64)
	ddPause = float64(ddPauseI64)
	charPause = float64(charPauseI64)
	wordPause = float64(wordPauseI64)
	return
}

func SamplesPerElement(wpm uint64, device *alsa.PlaybackDevice) (samplesPerElement int) {
	elementsPerSecond := elementsPerMinute(wpm) / 60
	samplesPerElement = device.Rate / int(elementsPerSecond)
	return
}

func Float64ElementsPerMinute(wpm uint64) (float64ElementsPerMinute float64) {
	float64ElementsPerMinute = float64(elementsPerMinute(wpm))
	return
}

func SecondsPerElement(wpm uint64) (seconds float64) {
	dit, _, _, _, _ := DurationsFromWPM(wpm)
	seconds = dit.Seconds()
	return
}

// DurationsFromWPM returns the important cw time spanws as time.Durations.
func DurationsFromWPM(wpm uint64) (dit, dah, ddPause, charPause, wordPause time.Duration) {
	d := int64(time.Minute) / elementsPerMinute(wpm)
	dit = time.Duration(d)
	dah = 3 * dit
	ddPause = dit
	charPause = dit * 3
	wordPause = dit * 7
	return
}

func elementsPerMinute(wpm uint64) (elementsPerMinute int64) {
	elementsPerMinute = int64(wpm * elementsInParis)
	return
}

// buildSoundLE16 makes a sinus wave of 16 bit little endian
func buildSoundLE16(device *alsa.PlaybackDevice, nSeconds float64) (raw []byte) {
	pi2 := math.Pi * 2.0
	ffreq := float64(SampleRate)
	nSamplesF64 := nSeconds * ffreq
	nSamples := int(math.Floor(nSamplesF64))
	maxUint16F64 := float64(math.MaxUint16)

	raw = make([]byte, nSamples*2)
	i := 0
	for t := 0; t < nSamples; t++ {
		amplF64 := maxUint16F64 * 0.5 * (1.0 + math.Sin(pi2*float64(t)*1000.0/ffreq))
		amplF64 = math.Floor(amplF64)
		ampl := int16(math.Floor(amplF64))
		// LittleEndian
		lend := byte(ampl >> 8)
		raw[i] = lend
		i++
		bend := byte(ampl & 0xFF)
		raw[i] = bend
		i++
	}
	return

	/*
		#include <math.h>
		#include <stdio.h>
		#include <stdlib.h>
		#include <stdint.h>

		int main(void) {
		    FILE *f;
		    const double PI2 = 2 * acos(-1.0);
		    const double SAMPLE_FREQ = 44100;
		    const unsigned int NSAMPLES = 4 * SAMPLE_FREQ;
		    uint16_t ampl;
		    uint8_t bytes[2];
		    unsigned int t;

		    f = fopen("out.raw", "wb");
		    for (t = 0; t < NSAMPLES; ++t) {
		        ampl = UINT16_MAX * 0.5 * (1.0 + sin(PI2 * t * 1000.0 / SAMPLE_FREQ));
		        bytes[0] = ampl >> 8;
		        bytes[1] = ampl & 0xFF;
		        fwrite(bytes, 2, sizeof(uint8_t), f);
		    }
		    fclose(f);
		    return EXIT_SUCCESS;
		}	*/
}
