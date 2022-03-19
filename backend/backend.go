package backend

import (
	"context"
	"fmt"

	"github.com/josephbudd/okp/backend/model/startup"
)

// Start starts the backend.
func Start(ctx context.Context, ctxCancel context.CancelFunc, errCh chan error) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("backend.Start: %w", err)
		}
	}()

	err = startup.Start(ctx, ctxCancel)
	return
}
