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
		// Copy @ 7 wpm. Key @ 7 wpm.
		{
			Name:         "7-5",
			KeyWPM:       7,
			CopyWPM:      7,
			CopySpaceWPM: 5,
		},
		{
			Name:         "7-7",
			KeyWPM:       7,
			CopyWPM:      7,
			CopySpaceWPM: 7,
		},
		// Copy @ 10 wpm. Key @ 10 wpm.
		{
			Name:         "10-5",
			KeyWPM:       10,
			CopyWPM:      10,
			CopySpaceWPM: 5,
		},
		{
			Name:         "10-10",
			KeyWPM:       10,
			CopyWPM:      10,
			CopySpaceWPM: 10,
		},
		// Copy @ 13 wpm. Key @ 13 wpm.
		{
			Name:         "13-5",
			KeyWPM:       13,
			CopyWPM:      13,
			CopySpaceWPM: 5,
		},
		{
			Name:         "13-10",
			KeyWPM:       13,
			CopyWPM:      13,
			CopySpaceWPM: 10,
		},
		{
			Name:         "13-13",
			KeyWPM:       13,
			CopyWPM:      13,
			CopySpaceWPM: 13,
		},
		// Copy @ 20 wpm. Key @ 20 wpm.
		{
			Name:         "20-5",
			KeyWPM:       20,
			CopyWPM:      20,
			CopySpaceWPM: 5,
		},
		{
			Name:         "20-10",
			KeyWPM:       20,
			CopyWPM:      20,
			CopySpaceWPM: 10,
		},
		{
			Name:         "20-20",
			KeyWPM:       20,
			CopyWPM:      20,
			CopySpaceWPM: 20,
		},
		// Copy @ 30 wpm. Key @ 30 wpm.
		{
			Name:         "30-5",
			KeyWPM:       30,
			CopyWPM:      30,
			CopySpaceWPM: 5,
		},
		{
			Name:         "30-10",
			KeyWPM:       30,
			CopyWPM:      30,
			CopySpaceWPM: 10,
		},
		{
			Name:         "30-20",
			KeyWPM:       30,
			CopyWPM:      30,
			CopySpaceWPM: 20,
		},
		{
			Name:         "30-30",
			KeyWPM:       30,
			CopyWPM:      30,
			CopySpaceWPM: 30,
		},
	}
	optionsLen = len(options)
}

// Option represents a wpm option.
type Option struct {
	Name         string // A unique descriptive name.
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

// ByWPMSpread returns the option where name == option.Name.
func ByWPMSpread(keyCopyWPM, spread uint64) (id int, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("wpm.ByWPMSpread: %w", err)
		}
	}()

	var option Option
	for id, option = range options {
		if option.KeyWPM == keyCopyWPM {
			return
		}
	}
	err = fmt.Errorf("speed %d, spread %d, not found in wpm.options", keyCopyWPM, spread)
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

// ByName returns the option where name == option.Name.
func ByName(name string) (id int, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("wpm.ByName: %w", err)
		}
	}()

	var option Option
	for id, option = range options {
		if option.Name == name {
			return
		}
	}
	id = -1
	err = fmt.Errorf("name %q not found in wpm.options", name)
	return
}
