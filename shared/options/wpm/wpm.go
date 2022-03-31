package wpm

import (
	"fmt"
)

var (
	options    []Option
	optionsLen int
)

func init() {
	options = []Option{
		// 7 wpm.
		{
			CopyWPM:      7,
			CopySpaceWPM: 5,
			KeyWPM:       7,
		},
		{
			CopyWPM:      7,
			CopySpaceWPM: 7,
			KeyWPM:       7,
		},
		// 10 wpm.
		{
			CopyWPM:      10,
			CopySpaceWPM: 5,
			KeyWPM:       7,
		},
		{
			CopyWPM:      10,
			CopySpaceWPM: 7,
			KeyWPM:       7,
		},
		{
			CopyWPM:      10,
			CopySpaceWPM: 5,
			KeyWPM:       10,
		},
		{
			CopyWPM:      10,
			CopySpaceWPM: 10,
			KeyWPM:       10,
		},
		// 13 wpm.
		{
			CopyWPM:      13,
			CopySpaceWPM: 5,
			KeyWPM:       7,
		},
		{
			CopyWPM:      13,
			CopySpaceWPM: 7,
			KeyWPM:       7,
		},
		{
			CopyWPM:      13,
			CopySpaceWPM: 5,
			KeyWPM:       13,
		},
		{
			CopyWPM:      13,
			CopySpaceWPM: 10,
			KeyWPM:       13,
		},
		{
			CopyWPM:      13,
			CopySpaceWPM: 13,
			KeyWPM:       13,
		},
		// 20 wpm.
		{
			CopyWPM:      20,
			CopySpaceWPM: 5,
			KeyWPM:       20,
		},
		{
			CopyWPM:      20,
			CopySpaceWPM: 10,
			KeyWPM:       20,
		},
		{
			CopyWPM:      20,
			CopySpaceWPM: 20,
			KeyWPM:       20,
		},
		// 30 wpm.
		{
			CopyWPM:      30,
			CopySpaceWPM: 5,
			KeyWPM:       30,
		},
		{
			CopyWPM:      30,
			CopySpaceWPM: 10,
			KeyWPM:       30,
		},
		{
			CopyWPM:      30,
			CopySpaceWPM: 20,
			KeyWPM:       30,
		},
		{
			CopyWPM:      30,
			CopySpaceWPM: 30,
			KeyWPM:       30,
		},
	}
	optionsLen = len(options)
}

// Option represents a wpm option.
type Option struct {
	KeyWPM       uint64 // The user will key at this rate.
	CopyWPM      uint64 // The app will key at this rate.
	CopySpaceWPM uint64 // The app will key at this rate.
}

// String stringifies the wpm option for a selection list.
func (o Option) String() (s string) {
	var sc string
	if o.CopyWPM == o.CopySpaceWPM {
		sc = fmt.Sprintf("Copy words which will be keyed at %d WPM.", o.CopyWPM)
	} else {
		sc = fmt.Sprintf("Copy words which will be keyed at %d WPM but spaced as if at %d WPM.", o.CopyWPM, o.CopySpaceWPM)
	}
	s = fmt.Sprintf("The user must:\n * %s\n * Key words at %d WPM.", sc, o.KeyWPM)
	return
}

// ByID returns the option at the index.
func ByID(indexID int) (id int, option Option) {
	if id = indexID; id < 0 || id >= optionsLen {
		id = 0
	}
	option = options[id]
	return
}

// ByText returns the option where text == option.String().
func ByText(text string) (id int, option Option) {
	for id, option = range options {
		if option.String() == text {
			return
		}
	}
	id = 0
	option = options[0]
	return
}

// ByWPMSpread returns the option where the wpm and spread match.
func ByWPMSpread(copyWPM, spread, keyWPM uint64) (id int, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("wpm.ByWPMSpread: %w", err)
		}
	}()

	var option Option
	for id, option = range options {
		if option.CopyWPM == copyWPM && option.CopySpaceWPM == spread && option.KeyWPM == keyWPM {
			return
		}
	}
	err = fmt.Errorf("copyWPM %d, spread %d, keyWPM %d not found in wpm.options", copyWPM, spread, keyWPM)
	return
}

// TextValue represents a text value pair.
type TextValue struct {
	Text  string
	Value int
}

// TextValuePairs returns the text value pairs for each option.
func TextValuePairs() (pairs []TextValue) {
	pairs = make([]TextValue, optionsLen)
	for i, option := range options {
		pairs[i] = TextValue{
			Text:  option.String(),
			Value: i,
		}
	}
	return
}

// Texts returns the text each option.
func Texts() (texts []string) {
	texts = make([]string, optionsLen)
	for i, option := range options {
		texts[i] = option.String()
	}
	return
}
