package start

import (
	"github.com/josephbudd/okp/shared/state"
	"github.com/josephbudd/okp/shared/store"
)

func initState(stores *store.Stores) (err error) {
	// Initialize the application state.
	err = state.Init(stores)
	return
}
