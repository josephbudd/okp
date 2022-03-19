package record

import "fmt"

const (
	MaxPlanLessonsCount uint64 = 1024
)

type Plan struct {
	ID          uint64
	Name        string
	Description string
	Lessons     []Lesson
}

func (plan Plan) String() (s string) {
	s = fmt.Sprintf("%s.\n%s", plan.Name, plan.Description)
	return
}

func NewPlan(name, desc string) (plan *Plan) {
	plan = &Plan{
		Name:        name,
		Description: desc,
	}
	return
}

type CurrentPlan struct {
	Name        string
	Description string
	Lessons     []*Lesson
}

type PlanOption struct {
	ID          uint64
	Name        string
	Description string
}

// ToPlanOption converts a *Plan to a *PlanOption.
func ToPlanOption(plan *Plan) (o PlanOption) {
	o = PlanOption{
		ID:          plan.ID,
		Name:        plan.Name,
		Description: plan.Description,
	}
	return
}

func (p PlanOption) String() (s string) {
	s = fmt.Sprintf("%s: %s", p.Name, p.Description)
	return
}
