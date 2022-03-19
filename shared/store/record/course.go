package record

import (
	"fmt"

	"github.com/josephbudd/okp/shared/options/wpm"
)

const (
	DefaultCoursePassCopyCount uint64 = 3
	DefaultCoursePassKeyCount  uint64 = 3
	DefaultCourseDelaySeconds  uint64 = 2
	DefaultCourseCopyWPM       uint64 = 13
	DefaultCourseCopySpaceWPM  uint64 = 20
	DefaultCourseKeyWPM        uint64 = 5
)

// Course is the record used in storage.
type Course struct {
	ID                  uint64
	Name                string
	Description         string
	SpeedID             int
	CurrentLessonNumber uint64
	DelaySeconds        uint64
	Completed           bool
	PlanID              uint64
	HomeWorks           []HomeWork
}

// CourseAdd is for the form when the user creates a course.
type CourseAdd struct {
	Name         string
	Description  string
	DelaySeconds uint64
	SpeedID      int
	PlanID       uint64
}

// CourseEdit is for the form when the user edits a course.
type CourseEdit struct {
	ID                       uint64
	Name                     string
	Description              string
	SpeedDescription         string
	PlanDescription          string
	CurrentLessonDescription string
	DelaySeconds             uint64
	Completed                bool
}

// CourseRemove is for display when the user confirms removal of a course.
type CourseRemove struct {
	ID                       uint64
	Name                     string
	Description              string
	SpeedDescription         string
	PlanDescription          string
	CurrentLessonDescription string
	DelaySeconds             uint64
	Completed                bool
}

// CourseCurrent is for displaying the current course.
type CourseCurrent struct {
	ID                  uint64
	SpeedID             int
	PlanID              uint64
	Name                string
	Description         string
	PlanDescription     string
	SpeedDescription    string
	CurrentLessonNumber uint64
	DelaySeconds        uint64
	Completed           bool
	HomeWorks           []HomeWork
}

// CourseOption is for when the user selects a course from a list.
type CourseOption struct {
	SortedIndex              int
	ID                       uint64
	Name                     string
	Description              string
	SpeedDescription         string
	PlanDescription          string
	CurrentLessonDescription string
	DelaySeconds             uint64
	Completed                bool
}

// Constructors.

// NewCourse constructs a new Course.
func NewCourse() (course *Course) {
	course = &Course{
		DelaySeconds: DefaultCourseDelaySeconds,
	}
	return
}

// NewCourseAdd constructs a new CourseAdd.
func NewCourseAdd() (course *CourseAdd) {
	course = &CourseAdd{
		DelaySeconds: DefaultCourseDelaySeconds,
	}
	return
}

// Conversion funcs.

// ToCourseEdit converts a *Course to a *CourseEdit.
func ToCourseEdit(course *Course, plans []*Plan) (e *CourseEdit) {
	var homework HomeWork
	var currentLessonDescription string
	if course.Completed {
		currentLessonDescription = "Completed"
	} else {
		for _, homework = range course.HomeWorks {
			if homework.LessonNumber == course.CurrentLessonNumber {
				currentLessonDescription = homework.LessonDescription
				break
			}
		}
	}
	var plan *Plan
	for _, plan = range plans {
		if plan.ID == course.PlanID {
			break
		}
	}

	var speed wpm.Option
	course.SpeedID, speed = wpm.ByID(course.SpeedID)
	e = &CourseEdit{
		ID:                       course.ID,
		Name:                     course.Name,
		Description:              course.Description,
		PlanDescription:          plan.Description,
		SpeedDescription:         speed.String(),
		CurrentLessonDescription: currentLessonDescription,
		Completed:                course.Completed,
		DelaySeconds:             course.DelaySeconds,
	}
	return
}

// ToCourseRemove converts a *Course to a *CourseRemove.
func ToCourseRemove(course *Course, plans []*Plan) (r *CourseRemove) {
	var homework HomeWork
	var currentLessonDescription string
	if course.Completed {
		currentLessonDescription = "Completed"
	} else {
		for _, homework = range course.HomeWorks {
			if homework.LessonNumber == course.CurrentLessonNumber {
				currentLessonDescription = homework.LessonDescription
				break
			}
		}
	}
	var plan *Plan
	for _, plan = range plans {
		if plan.ID == course.PlanID {
			break
		}
	}
	var speed wpm.Option
	course.SpeedID, speed = wpm.ByID(course.SpeedID)
	r = &CourseRemove{
		ID:                       course.ID,
		Name:                     course.Name,
		Description:              course.Description,
		PlanDescription:          plan.String(),
		SpeedDescription:         speed.String(),
		CurrentLessonDescription: currentLessonDescription,
		Completed:                course.Completed,
		DelaySeconds:             course.DelaySeconds,
	}
	return
}

// ToCourseCurrent converts a *Course to a *CourseCurrent which is used for state.
func ToCourseCurrent(course *Course, plans []*Plan) (r *CourseCurrent) {
	var plan *Plan
	for _, plan = range plans {
		if plan.ID == course.PlanID {
			break
		}
	}
	var speed wpm.Option
	course.SpeedID, speed = wpm.ByID(course.SpeedID)
	r = &CourseCurrent{
		ID:                  course.ID,
		Name:                course.Name,
		Description:         course.Description,
		PlanDescription:     plan.String(),
		SpeedDescription:    speed.String(),
		SpeedID:             course.SpeedID,
		PlanID:              course.PlanID,
		CurrentLessonNumber: course.CurrentLessonNumber,
		Completed:           course.Completed,
		DelaySeconds:        course.DelaySeconds,
	}
	return
}

// ToCourseOption converts a *Course to a *CourseOption.
func ToCourseOption(course *Course, sortedIndex int, plans []*Plan) (o *CourseOption) {
	var homework HomeWork
	var currentLessonDescription string
	if course.Completed {
		currentLessonDescription = "Completed"
	} else {
		for _, homework = range course.HomeWorks {
			if homework.LessonNumber == course.CurrentLessonNumber {
				currentLessonDescription = homework.LessonDescription
				break
			}
		}
	}
	var plan *Plan
	for _, plan = range plans {
		if plan.ID == course.PlanID {
			break
		}
	}
	var speed wpm.Option
	course.SpeedID, speed = wpm.ByID(course.SpeedID)
	o = &CourseOption{
		SortedIndex:              sortedIndex,
		ID:                       course.ID,
		Name:                     course.Name,
		Description:              course.Description,
		PlanDescription:          plan.Description,
		SpeedDescription:         speed.String(),
		CurrentLessonDescription: currentLessonDescription,
		Completed:                course.Completed,
		DelaySeconds:             course.DelaySeconds,
	}
	return
}

// FromCourseAdd convertsa *CourseAdd to a *Course.
// Returns a new course record that must be added to the store with Update(r).
func FromCourseAdd(add *CourseAdd) (course *Course) {
	course = &Course{
		Name:         add.Name,
		Description:  add.Description,
		DelaySeconds: add.DelaySeconds,
		PlanID:       add.PlanID,
		SpeedID:      add.SpeedID,
	}
	return
}

// FromCourseEdit converts a *CourseEdit to *Course.
// Returns an updated or new course record that must be added to the store with Update(r).
func FromCourseEdit(edit *CourseEdit, all []*Course) (course *Course) {
	for _, course = range all {
		if edit.ID == course.ID {
			course.Name = edit.Name
			course.Description = edit.Description
			course.DelaySeconds = edit.DelaySeconds
			return
		}
	}
	// Not found. This course will be added.
	course = NewCourse()
	course.Name = edit.Name
	course.Description = edit.Description
	// course.DelaySeconds = edit.DelaySeconds
	return
}

const (
	NameFieldName        = "Name"
	DescriptionFieldName = "Description"
	KWPMFieldName        = "Key Speed @WPM"
	CWPMFieldName        = "Copy Speed / word spacing @WPM"
	DelaySecondsName     = "Key delay seconds"
	PlanFieldName        = "Lesson Plan"
	PassKeyCountName     = "Must pass each key test"  // 1 time, 2 times, 3 times.
	PassCopyCountName    = "Must pass each copy test" // 1 time, 2 times, 3 times.
)

// Check for errors

func (add *CourseAdd) Errors() (err error) {
	if len(add.Name) == 0 {
		err = fmt.Errorf("%q is a required", NameFieldName)
		return
	}
	if len(add.Description) == 0 {
		err = fmt.Errorf("%q is a required", DescriptionFieldName)
		return
	}
	return
}
