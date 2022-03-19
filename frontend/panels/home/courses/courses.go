package courses

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/josephbudd/okp/frontend/panels/home/courses/addcourse"
	"github.com/josephbudd/okp/frontend/panels/home/courses/currentcourse"
	"github.com/josephbudd/okp/frontend/panels/home/courses/editcourse"
	"github.com/josephbudd/okp/frontend/panels/home/courses/removecourse"
)

func Init(ctx context.Context, ctxCancel context.CancelFunc, app fyne.App, w fyne.Window) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("courses.NewPanel: %w", err)
		}
	}()

	if err = currentcourse.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	if err = addcourse.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	if err = editcourse.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	if err = removecourse.Init(ctx, ctxCancel, app, w); err != nil {
		return
	}
	return
}

func Content() (content *fyne.Container) {
	selectGroup := currentcourse.Content()
	addGroup := addcourse.Content()
	editGroup := editcourse.Content()
	removeGroup := removecourse.Content()
	tabs := container.NewAppTabs(
		container.NewTabItem("Current", selectGroup),
		container.NewTabItem("Add", addGroup),
		container.NewTabItem("Edit", editGroup),
		container.NewTabItem("Remove", removeGroup),
	)
	content = container.New(
		layout.NewMaxLayout(),
		tabs,
	)
	return
}
