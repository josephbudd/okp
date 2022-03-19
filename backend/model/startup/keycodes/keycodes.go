package keycodes

import (
	"fmt"

	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

/*
	morse code PROSIGNS are chars run together.
	morse code ABBREVIATIONS are normal character words.
*/

type initDatum struct {
	Name          string
	Char          string
	DitDah        string
	IsWord        bool
	IsCompression bool
	IsNotReal     bool
}

var (
	cwLetters = []initDatum{
		{Name: "Alpha", Char: "A", DitDah: ".-"},
		{Name: "Bravo", Char: "B", DitDah: "-..."},
		{Name: "Charlie", Char: "C", DitDah: "-.-."},
		{Name: "Delta", Char: "D", DitDah: "-.."},
		{Name: "Echo", Char: "E", DitDah: "."},
		{Name: "Foxtrot", Char: "F", DitDah: "..-."},
		{Name: "Golf", Char: "G", DitDah: "--."},
		{Name: "Hotel", Char: "H", DitDah: "...."},
		{Name: "India", Char: "I", DitDah: ".."},
		{Name: "Juliet", Char: "J", DitDah: ".---"},
		{Name: "Kilo", Char: "K", DitDah: "-.-"},
		{Name: "Lima", Char: "L", DitDah: ".-.."},
		{Name: "Mike", Char: "M", DitDah: "--"},
		{Name: "November", Char: "N", DitDah: "-."},
		{Name: "Oscar", Char: "O", DitDah: "---"},
		{Name: "Papa", Char: "P", DitDah: ".--."},
		{Name: "Quebec", Char: "Q", DitDah: "--.-"},
		{Name: "Romeo", Char: "R", DitDah: ".-."},
		{Name: "Sierra", Char: "S", DitDah: "..."},
		{Name: "Tango", Char: "T", DitDah: "-"},
		{Name: "Uniform", Char: "U", DitDah: "..-"},
		{Name: "Victor", Char: "V", DitDah: "...-"},
		{Name: "Wiskey", Char: "W", DitDah: ".--"},
		{Name: "X-Ray", Char: "X", DitDah: "-..-"},
		{Name: "Yankee", Char: "Y", DitDah: "-.--"},
		{Name: "Zulu", Char: "Z", DitDah: "--.."},
	}

	cwNumbers = []initDatum{
		{Name: "1", Char: "1", DitDah: ".----"},
		{Name: "2", Char: "2", DitDah: "..---"},
		{Name: "3", Char: "3", DitDah: "...--"},
		{Name: "4", Char: "4", DitDah: "....-"},
		{Name: "5", Char: "5", DitDah: "....."},
		{Name: "6", Char: "6", DitDah: "-...."},
		{Name: "7", Char: "7", DitDah: "--..."},
		{Name: "8", Char: "8", DitDah: "---.."},
		{Name: "9", Char: "9", DitDah: "----."},
		{Name: "0", Char: "0", DitDah: "-----"},
	}

	cwPunctuation = []initDatum{
		{Name: "Period, Decimal Point", Char: ".", DitDah: ".-.-.-"},
		{Name: "Comma", Char: ",", DitDah: "--..--"},
		{Name: "Slash", Char: "/", DitDah: "-..-."},
		{Name: "Plus", Char: "+", DitDah: ".-.-."},
		{Name: "Equals or New Paragraph", Char: "=", DitDah: "-...-", IsWord: true},
		{Name: "Question Mark.", Char: "?", DitDah: "..--.."},
		{Name: "Please Repeat.", Char: "?", DitDah: "..--..", IsWord: true},
		{Name: "Open Paren", Char: "(", DitDah: "-.--."},
		{Name: "Close Paren", Char: ")", DitDah: "-.--.-"},
		{Name: "Dash", Char: "-", DitDah: "-....-"},
		{Name: "Double Quote", Char: "\"", DitDah: ".-..-."},
		{Name: "Underline", Char: "_", DitDah: "..--.-"},
		{Name: "Single Quote", Char: "'", DitDah: ".----."},
		{Name: "Colon", Char: ":", DitDah: "---..."},
		{Name: "Semicolon", Char: ";", DitDah: "-.-.-."},
		{Name: "Dollar Sign", Char: "$", DitDah: "...-..-"},
		{Name: "At Sign", Char: "@", DitDah: ".--.-."},
	}

	// Combinations are combinations.
	// Abbreviations.
	cwCombination = []initDatum{
		{Name: "Calling Anyone", Char: "CQ", DitDah: "-.-. --.-", IsWord: true},
		{Name: "This Is", Char: "DE", DitDah: "-.. .", IsWord: true},
		{Name: "Back To You", Char: "BTU", DitDah: "-... - ..-", IsWord: true},
		{Name: "Break, Pause.", Char: "BK", DitDah: "-... -.-", IsWord: true},
		{Name: "Closing Down", Char: "CL", DitDah: "-.-. .-..", IsWord: true},
		{Name: "Yes, Correct", Char: "C", DitDah: "-.-.", IsWord: true},
		{Name: "No", Char: "N", DitDah: "-.", IsWord: true},
		{Name: "Roger. Received as transmitted", Char: "R", DitDah: ".-.", IsWord: true},
		{Name: "Weather", Char: "WX", DitDah: ".-- -..-", IsWord: true},                       // New
		{Name: "Readablity, strength, tone.", Char: "RST", DitDah: ".-. ... -", IsWord: true}, // New
	}

	// Compressions are combinations of chars compressed into a single character.
	cwCompression = []initDatum{
		{Name: "Over (to all)", Char: "K", DitDah: "-.-", IsWord: true},
		{Name: "Over (only to the NAMED caller)", Char: "KN", DitDah: "-.--.", IsWord: true, IsCompression: true}, // "KN" == "("
		{Name: "NewLine", Char: "AA", DitDah: ".-.-", IsWord: true, IsCompression: true},
		{Name: "No Reply Expected, NewMessage", Char: "AR", DitDah: ".-.-.", IsWord: true, IsCompression: true},
		{Name: "Wait", Char: "AS", DitDah: ".-...", IsWord: true, IsCompression: true},
		{Name: "NewParagraph", Char: "BT", DitDah: "-...-", IsWord: true, IsCompression: true},
		{Name: "StartOfTransmission", Char: "CT", DitDah: "-.-.-", IsWord: true, IsCompression: true},
		{Name: "Error", Char: "HH", DitDah: "........", IsWord: true, IsCompression: true},
		{Name: "End Of Contact", Char: "SK", DitDah: "...-.-", IsWord: true, IsCompression: true},
		{Name: "Understood, Verified", Char: "SN", DitDah: "...-.", IsWord: true, IsCompression: true},
		{Name: "Save Our Ship", Char: "SOS", DitDah: "...---...", IsWord: true, IsCompression: true},
		{Name: "Warning", Char: "EW", DitDah: ".-..-", IsWord: true, IsCompression: true},
	}
)

// CreateKeyCodes creates the default keycode records.
func CreateKeyCodes(stores *store.Stores) (keycodes []*record.KeyCode, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("model.CreateKeyCodes: %w", err)
		}
	}()

	// Only initialize if the store is empty.
	if keycodes, err = stores.KeyCode.GetAll(); err != nil {
		return
	}
	if len(keycodes) > 0 {
		// Already initialized.
		return
	}
	// Need to initialize.
	keycodes = makeInitialKeyCodes()
	err = stores.KeyCode.UpdateAll(keycodes)
	return
}

func makeInitialKeyCodes() (keycodes []*record.KeyCode) {
	var d initDatum
	var r *record.KeyCode
	var countKeyCodes = len(cwLetters) + len(cwNumbers) + len(cwPunctuation) + len(cwCombination) + len(cwCompression)
	keycodes = make([]*record.KeyCode, 0, countKeyCodes)
	for _, d = range cwLetters {
		r = &record.KeyCode{
			Name:      d.Name,
			Character: d.Char,
			DitDah:    d.DitDah,
		}
		keycodes = append(keycodes, r)
	}
	for _, d = range cwNumbers {
		r = &record.KeyCode{
			Name:      d.Name,
			Character: d.Char,
			DitDah:    d.DitDah,
		}
		keycodes = append(keycodes, r)
	}
	for _, d = range cwPunctuation {
		r = &record.KeyCode{
			Name:      d.Name,
			Character: d.Char,
			DitDah:    d.DitDah,
			IsWord:    d.IsWord,
		}
		keycodes = append(keycodes, r)
	}
	for _, d = range cwCombination {
		r = &record.KeyCode{
			Name:      d.Name,
			Character: d.Char,
			DitDah:    d.DitDah,
			IsWord:    d.IsWord,
		}
		keycodes = append(keycodes, r)
	}
	for _, d = range cwCompression {
		r = &record.KeyCode{
			Name:          d.Name,
			Character:     d.Char,
			DitDah:        d.DitDah,
			IsWord:        d.IsWord,
			IsCompression: d.IsCompression,
		}
		keycodes = append(keycodes, r)
	}
	return
}
