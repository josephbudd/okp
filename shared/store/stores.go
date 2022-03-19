package store

import (
	"fmt"
	"strings"

	"github.com/josephbudd/okp/shared/store/storing"
)

// Stores is each of the application's storers.
type Stores struct {
	// Local yaml stores.
	State   *storing.StateStore
	Course  *storing.CourseStore
	KeyCode *storing.KeyCodeStore
	Plan    *storing.PlanStore
}

// New constructs a new Stores.
func New() (stores *Stores) {
	stores = &Stores{
		State:   storing.NewStateStore(),
		Course:  storing.NewCourseStore(),
		KeyCode: storing.NewKeyCodeStore(),
		Plan:    storing.NewPlanStore(),
	}
	return
}

// Open opens every store.
// It returns all of the errors as one single error.
func (stores *Stores) Open() (err error) {

	errList := make([]string, 0, 10)
	defer func() {
		if len(errList) > 0 {
			msg := strings.Join(errList, "\n")
			err = fmt.Errorf("stores.Open: %s", msg)
		}
	}()

	// Local yaml stores.
	if err = stores.State.Open(); err != nil {
		errList = append(errList, err.Error())
	}
	if err = stores.Course.Open(); err != nil {
		errList = append(errList, err.Error())
	}
	if err = stores.KeyCode.Open(); err != nil {
		errList = append(errList, err.Error())
	}
	if err = stores.Plan.Open(); err != nil {
		errList = append(errList, err.Error())
	}

	return
}

// Close closes every store.
// It returns all of the errors as one single error.
func (stores *Stores) Close() (err error) {

	errList := make([]string, 0, 4)
	defer func() {
		if len(errList) > 0 {
			msg := strings.Join(errList, "\n")
			err = fmt.Errorf("stores.Close: %s", msg)
		}
	}()

	// Local yaml stores.
	if err = stores.State.Close(); err != nil {
		errList = append(errList, err.Error())
	}
	if err = stores.Course.Close(); err != nil {
		errList = append(errList, err.Error())
	}
	if err = stores.KeyCode.Close(); err != nil {
		errList = append(errList, err.Error())
	}
	if err = stores.Plan.Close(); err != nil {
		errList = append(errList, err.Error())
	}

	return
}
