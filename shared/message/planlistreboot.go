package message

import (
	"fmt"

	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

var PlanListRebootID = NextID()

type PlanListReboot struct {
	id      uint64              // to front
	name    string              // to front
	Records []record.PlanOption // to front

	Error        bool
	Fatal        bool
	ErrorMessage string
}

// NewPlanListReboot constructs a new PlanListReboot message.
func NewPlanListReboot(stores *store.Stores) (msg *PlanListReboot) {

	msg = &PlanListReboot{
		id:   PlanListRebootID,
		name: "PlanListReboot",
	}

	var fatal error
	defer func() {
		if fatal != nil {
			fatal = fmt.Errorf("NewPlanListReboot : %w", fatal)
			msg.Fatal = true
			msg.ErrorMessage = fatal.Error()
		}
	}()

	// Plans.
	var plans []*record.Plan
	if plans, fatal = stores.Plan.GetAll(); fatal != nil {
		return
	}
	options := make([]record.PlanOption, len(plans))
	for i, plan := range plans {
		options[i] = record.ToPlanOption(plan)
	}
	msg.Records = options
	return
}

// PlanListReboot implements the MSGer interface with ID and MSG.

// ID returns the message's id
func (msg *PlanListReboot) ID() (id uint64) {
	id = msg.id
	return
}

// Name returns the message's id
func (msg *PlanListReboot) Name() (name string) {
	name = msg.name
	return
}

// Message returns the message's id
func (msg *PlanListReboot) MSG() (m interface{}) {
	m = msg
	return
}

// IsFatal return if there was a fatal error and it's message.
func (msg *PlanListReboot) FatalError() (fatal bool, message string) {
	fatal = msg.Fatal
	message = msg.ErrorMessage
	return
}
