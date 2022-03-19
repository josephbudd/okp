package files

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed images/courses_lettered_piano_keys.png
var courses_lettered_piano_keys []byte

//go:embed images/training_lettered_piano_keys.png
var training_lettered_piano_keys []byte

func TrainingLetteredPianoKeys() (res fyne.Resource) {
	res = fyne.NewStaticResource("images/training_lettered_piano_keys.png", training_lettered_piano_keys)
	return
}

func CoursesLetteredPianoKeys() (res fyne.Resource) {
	res = fyne.NewStaticResource("images/courses_lettered_piano_keys.png", courses_lettered_piano_keys)
	return
}
