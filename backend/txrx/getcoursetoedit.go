package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const getCourseToEditF = "getCourseToEditRX: %s"

func init() {
	addListener(message.GetCourseToEditID, getCourseToEditRX)
}

func getCourseToEditRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	getCourseToEditMsg := msg.(*message.GetCourseToEdit)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			getCourseToEditMsg.Error = true
			getCourseToEditMsg.ErrorMessage = fmt.Sprintf(getCourseToEditF, err.Error())
		case fatal != nil:
			getCourseToEditMsg.Fatal = true
			getCourseToEditMsg.ErrorMessage = fmt.Sprintf(getCourseToEditF, fatal.Error())
		default:
		}
		message.BackEndToFrontEnd <- getCourseToEditMsg
	}()

	var course *record.Course
	if course, fatal = stores.Course.Get(getCourseToEditMsg.RecordID); fatal != nil {
		return
	}
	var plans []*record.Plan
	if plans, fatal = stores.Plan.GetAll(); fatal != nil {
		return
	}
	getCourseToEditMsg.Record = record.ToCourseEdit(course, plans)
}
