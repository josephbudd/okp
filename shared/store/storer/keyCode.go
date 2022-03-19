package storer

import (
	"github.com/josephbudd/okp/shared/store/record"
)

// KeyCodeStorer defines the behavior (API) of a store of /shared/store/record.KeyCode records.
type KeyCodeStorer interface {

	// Open opens the bolt data-store.
	// Keeps the file in memory.
	// Returns the error.
	Open() (err error)

	// Close closes the data-store.
	// Returns the error.
	Close() (err error)

	// Get retrieves one *record.KeyCode from the data-store.
	// Param id is the record ID.
	// Returns a *record.KeyCode and error.
	// When no record is found, the returned *record.KeyCode is nil and the returned error is nil.
	Get(id uint64) (r *record.KeyCode, err error)

	// GetAll retrieves all of the *record.KeyCode records from the data-store.
	// Returns a slice of *record.KeyCode and error.
	// When no records are found, the returned slice length is 0 and the returned error is nil.
	GetAll() (rr []*record.KeyCode, err error)

	// Update updates the *record.KeyCode in the data-store.
	// Param newR is the *record.KeyCode to be updated.
	// If newR is a new record then r.ID is updated as well.
	// Returns the error.
	Update(newR *record.KeyCode) (err error)

	// UpdateAll updates a slice of *record.KeyCode in the data-store.
	// Param newRR is the slice of *record.KeyCode to be updated.
	// If any record in newRR is new then it's ID is updated as well.
	// Returns the error.
	UpdateAll(newRR []*record.KeyCode) (err error)

	// Remove removes the record.KeyCode from the data-store.
	// Param id is the record ID of the record.KeyCode to be removed.
	// If the record is not found returns a nil error.
	// Returns the error.
	Remove(id uint64) (err error)
}
