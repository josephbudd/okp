package txrx

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/shared/message"
	"github.com/josephbudd/okp/shared/store"
)

const removeCourseF = "removeCourseRX: %s"

func init() {
	addListener(message.RemoveCourseID, removeCourseRX)
}

func removeCourseRX(ctx context.Context, ctxCancel context.CancelFunc, stores *store.Stores, msg interface{}) {

	removeCourseMsg := msg.(*message.RemoveCourse)
	var err, fatal error
	defer func() {
		switch {
		case err != nil:
			removeCourseMsg.Error = true
			removeCourseMsg.ErrorMessage = fmt.Sprintf(removeCourseF, err.Error())
			message.BackEndToFrontEnd <- removeCourseMsg
		case fatal != nil:
			removeCourseMsg.Fatal = true
			removeCourseMsg.ErrorMessage = fmt.Sprintf(removeCourseF, fatal.Error())
			message.BackEndToFrontEnd <- removeCourseMsg
		default:
			message.BackEndToFrontEnd <- removeCourseMsg
			// No errors so tell the course lists to reboot.
			message.BackEndToFrontEnd <- message.NewCourseListReboot(stores)
		}
	}()

	fatal = stores.Course.Remove(removeCourseMsg.RecordID)
}
