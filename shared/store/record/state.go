package record

type State struct {
	ID              uint64
	CurrentCourseID uint64
}

func NewState() (state *State) {
	state = &State{
		ID: 1,
	}
	return
}

type CurrentState struct {
	CurrentCourse   *CourseCurrent
	CurrentHomeWork *HomeWorkCurrent
}
