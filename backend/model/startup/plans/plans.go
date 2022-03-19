package plans

import (
	"fmt"

	"github.com/josephbudd/okp/backend/model/startup/plans/amateur"
	"github.com/josephbudd/okp/backend/model/startup/plans/cq"
	"github.com/josephbudd/okp/backend/model/startup/plans/numbers"

	// "github.com/josephbudd/okp/backend/model/startup/plans/original"
	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

// Create creates each plan.
// Make sure to add a call to your new plan's func Create here.
func Create(stores *store.Stores, keycodes []*record.KeyCode) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("plans.Create: %w", err)
		}
	}()

	keyCodeCharMap := make(map[string]*record.KeyCode, len(keycodes))
	keyCodeWordMap := make(map[string]*record.KeyCode, len(keycodes))
	keyCodeIDMap := make(map[uint64]*record.KeyCode, len(keycodes))
	for _, kc := range keycodes {
		keyCodeIDMap[kc.ID] = kc
		switch {
		case kc.IsNotReal:
			// Ignore if not real.
		case kc.IsCompression, kc.IsWord:
			keyCodeWordMap[kc.Character] = kc
		default:
			keyCodeCharMap[kc.Character] = kc
		}
	}

	// Create each plan if needed.
	// if err = original.Create(stores, keyCodeCharMap, keyCodeNameMap, keyCodeIDMap); err != nil {
	// 	return
	// }
	if err = amateur.Create(stores, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	if err = cq.Create(stores, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	if err = numbers.Create(stores, keyCodeCharMap); err != nil {
		return
	}
	return
}
