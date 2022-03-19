package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const getCourseToRemoveF = "getCourseToRemoveRX: %s"

func init() {
	addListener(message.GetCourseToRemoveID, getCourseToRemoveRX)
}

func getCourseToRemoveRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	getCourseToRemoveMsg := msg.(*message.GetCourseToRemove)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			getCourseToRemoveMsg.Error = true
			getCourseToRemoveMsg.ErrorMessage = fmt.Sprintf(getCourseToRemoveF, err.Error())
		case fatal != nil:
			getCourseToRemoveMsg.Fatal = true
			getCourseToRemoveMsg.ErrorMessage = fmt.Sprintf(getCourseToRemoveF, fatal.Error())
		default:
		}
		message.BackEndToFrontEnd <- getCourseToRemoveMsg
	}()

	var plans []*record.Plan
	if plans, fatal = stores.Plan.GetAll(); fatal != nil {
		return
	}
	var course *record.Course
	if course, fatal = stores.Course.Get(getCourseToRemoveMsg.RecordID); fatal != nil {
		return
	}
	getCourseToRemoveMsg.Record = record.ToCourseRemove(course, plans)
}
