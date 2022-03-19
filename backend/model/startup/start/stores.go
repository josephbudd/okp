package start

import (
	"fmt"

	"github.com/josephbudd/okp/backend/model"
	"github.com/josephbudd/okp/backend/model/startup/keycodes"
	"github.com/josephbudd/okp/backend/model/startup/plans"
	"github.com/josephbudd/okp/backend/model/startup/plans/amateur"
	"github.com/josephbudd/okp/shared/options/wpm"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

func initStores() (stores *store.Stores, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("startup.initStores: %w", err)
		}
	}()

	// Stores.
	stores = store.New()
	if err = stores.Open(); err != nil {
		return
	}
	if err = stores.Close(); err != nil {
		return
	}

	// KeyCodes.
	var kcRR []*record.KeyCode
	if kcRR, err = keycodes.CreateKeyCodes(stores); err != nil {
		return
	}

	// Plans.
	// Create each lesson plan.
	if err = plans.Create(stores, kcRR); err != nil {
		return
	}

	// Courses.
	var courses []*record.Course
	if courses, err = stores.Course.GetAll(); err != nil {
		return
	}
	var course *record.Course
	if len(courses) == 0 {
		// There is no default course so create it.
		if course, err = createDefaultCourse(stores); err != nil {
			return
		}
	} else {
		course = courses[0]
	}

	// State.
	var state *record.State
	if state, err = stores.State.Get(1); err != nil {
		return
	}
	if state != nil {
		return
	}
	// There is no state so create it using the default course.
	state = record.NewState()
	state.CurrentCourseID = course.ID
	err = stores.State.Update(state)
	return
}

func createDefaultCourse(stores *store.Stores) (course *record.Course, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("createDefaultCourse: %w", err)
		}
	}()

	// Get the speed. 10 wpm spread out like 5 wpm.
	var speedID int
	if speedID, err = wpm.ByWPMSpread(7, 5); err != nil {
		return
	}

	var plans []*record.Plan
	if plans, err = stores.Plan.GetAll(); err != nil {
		return
	}
	var plan *record.Plan
	// Search for the amateur plan. If not found use the last plan.
	for _, plan = range plans {
		if plan.Name == amateur.PlanName {
			break
		}
	}
	if plan == nil {
		err = fmt.Errorf("plan named %q not found", amateur.PlanName)
		return
	}
	// Create the first course using the amateur plan.
	course, err = model.NewCourse(
		"Default", "The default course.",
		speedID, // Use the first speed record.
		plan,
		stores,
	)
	return
}
