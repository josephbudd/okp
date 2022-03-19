package keyservice

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/josephbudd/okp/backend/model/startup/start"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

var (
	stores  *store.Stores
	initErr error
)

func init() {
	// os.Setenv("USETESTPATH", "1")
	// stores, initErr = start.Init()
	// log.Printf("initErr is %q", initErr.Error())
	// log.Printf("stores is %#v", stores)
}

func TestCopyMilliSeconds(t *testing.T) {
	os.Setenv("USETESTPATH", "1")
	if stores, initErr = start.Init(); initErr != nil {
		t.Fatalf("initErr is %q", initErr.Error())
	}
	log.Printf("initErr is %q", initErr.Error())
	log.Printf("stores is %#v", stores)

	var err error
	if err = stores.Open(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = stores.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	var solution [][]*record.KeyCode
	var keyCodeRecords []*record.KeyCode
	if solution, keyCodeRecords, err = makeSolutionKeyCodeRecords(); err != nil {
		t.Fatal(err)
	}

	wpm := uint64(20)
	type args struct {
		milliSeconds []int64
		wpm          uint64
		keyCodes     []*record.KeyCode
		source       [][]*record.KeyCode
	}
	type data struct {
		name            string
		args            args
		wantDitdahWords [][]string
	}
	ls := len(solution)
	tests := make([]data, ls)
	for i := range solution {
		i1 := i + 1
		tests[i] = data{
			name: toname(solution[i]),
			args: args{
				milliSeconds: Milliseconds(solution[i:i1], wpm),
				source:       solution[i:i1],
				wpm:          wpm,
				keyCodes:     keyCodeRecords,
			},
			wantDitdahWords: [][]string{
				toditdahstring(solution[i]),
			},
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guess := CopyMilliSeconds(tt.args.milliSeconds, tt.args.wpm, tt.args.keyCodes)
			if !guessedCorrectly(t, guess, tt.args.source) {
				t.Fatalf("unable to guess %s", tt.name)
			}
		})
	}
}

// Param guess is the CopyGuess.
// Param solutionWords is each word's keycodes.
//
func guessedCorrectly(t *testing.T, guess CopyGuess, solutionWords [][]*record.KeyCode) (matchedWord bool) {
	// Each item in solutionWords is a separate word from the user's keying.
	var nSolutionWords int
	var lFirstWord int
	if nSolutionWords = len(solutionWords); nSolutionWords > 0 {
		lFirstWord = len(solutionWords[0])
	}
	switch {
	case nSolutionWords == 1 && lFirstWord == 1 && guess.Compressed.ID == solutionWords[0][0].ID:
		matchedWord = true
	case nSolutionWords == 1 && lFirstWord == 1 && guess.Combined.ID == solutionWords[0][0].ID:
		matchedWord = true
	case nSolutionWords == 1 && len(guess.Word) == lFirstWord:
		solutionWord := solutionWords[0]
		matchedWord = true
		for i, kc := range guess.Word {
			if kc.ID != solutionWord[i].ID {
				matchedWord = false
				break
			}
		}
	case nSolutionWords > 1 && len(guess.Sentence) == nSolutionWords:
		// Each item in guess.Sentence is a word that should correspond to the item in solution.
		matchedWord = true
		for i, sentenceWord := range guess.Sentence {
			solutionWord := solutionWords[i]
			if len(solutionWord) == len(sentenceWord) {
				for j, kc := range sentenceWord {
					if kc.ID != solutionWord[j].ID {
						matchedWord = false
						break
					}
				}
			} else {
				matchedWord = false
			}
			if !matchedWord {
				break
			}
		}
	}
	return
}

func toditdahstring(rr []*record.KeyCode) (ditdahs []string) {
	ditdahs = make([]string, len(rr))
	for i, r := range rr {
		ditdahs[i] = r.DitDah
	}
	return
}

func toname(rr []*record.KeyCode) (name string) {
	ss := make([]string, len(rr))
	for i, r := range rr {
		ss[i] = r.Character
	}
	name = strings.Join(ss, ",")
	return
}

func rr2Ditdahs(rr []*record.KeyCode) (ss []string) {
	ss = make([]string, len(rr))
	for i, r := range rr {
		ss[i] = r.DitDah
	}
	return
}

func Test_SameAs(t *testing.T) {
	type args struct {
		src CopyGuess
		dst CopyGuess
	}
	type datum struct {
		name string
		args args
		want bool
	}
	var data = []datum{
		{
			name: "first",
			args: args{
				src: CopyGuess{},
				dst: CopyGuess{},
			},
			want: true,
		},
	}

	for i, datum := range data {
		got := datum.args.src.SameAs(datum.args.dst)
		if got != datum.want {
			t.Fatalf("no match at %d", i)
		}
	}
}

func Test_ditdahWordGuesses(t *testing.T) {
	solution, keyCodeRecords, err := makeSolutionKeyCodeRecords()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ditdahChars [][]string
		records     []*record.KeyCode
	}
	tests := []struct {
		name  string
		args  args
		guess CopyGuess
		want  bool
	}{
		// TODO: Add test cases.
		{
			name: "first",
			args: args{
				ditdahChars: [][]string{rr2Ditdahs(solution[0])},
				records:     keyCodeRecords,
			},
			guess: DitdahWordGuesses([][]string{rr2Ditdahs(solution[0])}, keyCodeRecords),
			want:  true,
		},
	}
	for i, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				copyGuess := DitdahWordGuesses(tt.args.ditdahChars, tt.args.records)
				if copyGuess.SameAs(tt.guess) != tt.want {
					t.Fatalf("no match at %d", i)
				}
			},
		)
	}
}

func makeSolutionKeyCodeRecords() (solution [][]*record.KeyCode, keyCodeRecords []*record.KeyCode, err error) {

	chars, words := splitKeyCodes(keyCodeRecords)
	solution = [][]*record.KeyCode{
		chars[0:5],
		words[0:1],

		chars[5:10],
		words[1:2],

		chars[10:15],
		words[2:3],

		chars[15:20],
		words[3:4],

		chars[20:25],
		words[4:5], // "BTU"
	}
	return
}

func splitKeyCodes(rr []*record.KeyCode) (chars, words []*record.KeyCode) {
	l := len(rr)
	chars = make([]*record.KeyCode, 0, l)
	words = make([]*record.KeyCode, 0, l)
	for _, r := range rr {
		if r.IsWord && !r.IsNotReal {
			words = append(words, r)
		} else if !r.IsWord && !r.IsNotReal {
			chars = append(chars, r)
		}
	}
	return
}
