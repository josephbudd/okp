package home

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/josephbudd/okp/frontend/panels/home/courses"
	"github.com/josephbudd/okp/frontend/panels/home/training"
	"github.com/josephbudd/okp/frontend/widget/safebutton"
	"github.com/josephbudd/okp/frontend/widget/subpanel"
)

var (
	window fyne.Window
)

// Init initializes the panel groups.
func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("home.Init: %w", err)
		}
	}()

	window = w

	// Courses button.
	if err = courses.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	// Training button.
	if err = training.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	return
}

func Content() (content *fyne.Container, err error) {

	back := func() {
		window.SetContent(content)
	}

	// Courses button.
	coursesGroupContent := courses.Content()
	coursesContent := subpanel.Content("Courses", back, coursesGroupContent)

	// Training button.
	var trainingGroupContent *fyne.Container
	if trainingGroupContent, err = training.Content(); err != nil {
		return
	}
	trainingContent := subpanel.Content("Training", back, trainingGroupContent)

	// Main menu.
	// mainMenu := makeMainMenu(coursesContent, trainingContent)
	// window.SetMainMenu(mainMenu)

	coursesButton := safebutton.New("Courses", func() { window.SetContent(coursesContent) })
	trainingButton := safebutton.New("Training", func() { window.SetContent(trainingContent) })
	// coursesButton := imgbutton.New("Courses", func() { window.SetContent(coursesContent) }, files.CoursesLetteredPianoKeys())
	// trainingButton := imgbutton.New("Training", func() { window.SetContent(trainingContent) }, files.TrainingLetteredPianoKeys())
	content = container.NewCenter(container.NewHBox(coursesButton, layout.NewSpacer(), trainingButton))

	return
}

// func makeMainMenu(coursesContent, trainingContent *fyne.Container) (mm *fyne.MainMenu) {
// 	courses := fyne.NewMenuItem("Courses", func() { window.SetContent(coursesContent) })
// 	training := fyne.NewMenuItem("Training", func() { window.SetContent(trainingContent) })
// 	var menu *fyne.Menu
// 	if fyne.CurrentDevice().IsMobile() {
// 		// Mobile.
// 		menu = fyne.NewMenu("okp", courses, training)
// 	} else {
// 		// Not mobile.
// 		quit := fyne.NewMenuItem(
// 			"Quit",
// 			func() {
// 				ctxCancelFunc()
// 				window.Close()
// 			},
// 		)
// 		menu = fyne.NewMenu("okp", courses, training, fyne.NewMenuItemSeparator(), quit)
// 	}
// 	mm = fyne.NewMainMenu(menu)
// 	return
// }
