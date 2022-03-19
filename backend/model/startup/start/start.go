package start

import (
	"fmt"

	"github.com/josephbudd/okp/shared/store"
)

// Init is where the back end really starts.
func Init() (stores *store.Stores, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("start.Start: %w", err)
		}
	}()

	// File paths.
	if err = initPaths(); err != nil {
		return
	}

	// Stores.
	if stores, err = initStores(); err != nil {
		return
	}

	// State.
	err = initState(stores)
	return
}
