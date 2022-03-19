package storer

import (
	"github.com/josephbudd/okp/shared/store/record"
)

// PlanStorer defines the behavior (API) of a store of /shared/store/record.Plan records.
type PlanStorer interface {

	// Open opens the bolt data-store.
	// Keeps the file in memory.
	// Returns the error.
	Open() (err error)

	// Close closes the data-store.
	// Returns the error.
	Close() (err error)

	// Get retrieves one *record.Plan from the data-store.
	// Param id is the record ID.
	// Returns a *record.Plan and error.
	// When no record is found, the returned *record.Plan is nil and the returned error is nil.
	Get(id uint64) (r *record.Plan, err error)

	// GetAll retrieves all of the *record.Plan records from the data-store.
	// Returns a slice of *record.Plan and error.
	// When no records are found, the returned slice length is 0 and the returned error is nil.
	GetAll() (rr []*record.Plan, err error)

	// Update updates the *record.Plan in the data-store.
	// Param newR is the *record.Plan to be updated.
	// If newR is a new record then r.ID is updated as well.
	// Returns the error.
	Update(newR *record.Plan) (err error)

	// UpdateAll updates a slice of *record.Plan in the data-store.
	// Param newRR is the slice of *record.Plan to be updated.
	// If any record in newRR is new then it's ID is updated as well.
	// Returns the error.
	UpdateAll(newRR []*record.Plan) (err error)

	// Remove removes the record.Plan from the data-store.
	// Param id is the record ID of the record.Plan to be removed.
	// If the record is not found returns a nil error.
	// Returns the error.
	Remove(id uint64) (err error)
}
