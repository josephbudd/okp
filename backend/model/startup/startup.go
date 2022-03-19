package startup

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/backend/model/sound"
	"github.com/josephbudd/okp/backend/model/startup/start"
	"github.com/josephbudd/okp/backend/txrx"
	"github.com/josephbudd/okp/shared/store"
)

// Start is where the back end really starts.
func Start(ctx context.Context, ctxCancel context.CancelFunc) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("startup.Start: %w", err)
		}
	}()

	var stores *store.Stores
	if stores, err = start.Init(); err != nil {
		return
	}

	if err = sound.Open(ctx); err != nil {
		return
	}

	// Messages.
	txrx.Listen(ctx, ctxCancel, stores)
	return
}
