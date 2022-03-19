package storing

import (
	"bytes"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/yaml.v2"

	"github.com/josephbudd/okp/shared/paths"
	"github.com/josephbudd/okp/shared/store/record"
)

type PlanData struct {
	LastID  uint64
	Records []*record.Plan
}

// PlanStore is the API of the Plan store.
// It is the implementation of the interface in /domain/store/storer/Plan.go.
type PlanStore struct {
	uri  fyne.URI
	lock sync.Mutex
	data PlanData
}

// NewPlanStore constructs a new PlanStore.
// Param db is an open bolt data-store.
// Returns a pointer to the new PlanStore.
func NewPlanStore() (store *PlanStore) {
	store = &PlanStore{
		uri: paths.StoreURI("plan.yaml"),
	}
	return
}

// Open opens the bolt data-store.
// Keeps the file in memory.
// Returns the error.
func (store *PlanStore) Open() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.Open: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	err = store.readAll()
	return
}

// Close closes the data-store.
// Returns the error.
func (store *PlanStore) Close() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.Close: %w", err)
		}
	}()

	// The store is always closed.
	return
}

// Get retrieves one *record.Plan from the data-store.
// Param id is the record ID.
// Returns a *record.Plan and error.
// When no record is found, the returned *record.Plan is nil and the returned error is nil.
func (store *PlanStore) Get(id uint64) (r *record.Plan, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.Get: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	for _, r = range store.data.Records {
		if r.ID == id {
			return
		}
	}
	// Not found. No error.
	r = nil
	return
}

// GetAll retrieves all of the *record.Plan records from the data-store.
// Returns a slice of *record.Plan and error.
// When no records are found, the returned slice length is 0 and the returned error is nil.
func (store *PlanStore) GetAll() (rr []*record.Plan, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.GetAll: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	rr = store.data.Records
	return
}

// Update updates the *record.Plan in the data-store.
// Param newR is the *record.Plan to be updated.
// If newR is a new record then r.ID is updated as well.
// Returns the error.
func (store *PlanStore) Update(newR *record.Plan) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.Update: %w", err)
		}
	}()

	// This record must have no more than record.MaxPlanLessonsCount lessons.
	l := len(newR.Lessons)
	if l > int(record.MaxPlanLessonsCount) {
		err = fmt.Errorf("the plan has %d lessons but only %d are allowed", l, record.MaxPlanLessonsCount)
		return
	}

	store.lock.Lock()
	defer store.lock.Unlock()

	// Add or replace the record.
	if newR.ID == 0 {
		// Adding a new record so append it.
		store.data.LastID++
		newR.ID = store.data.LastID
		store.data.Records = append(store.data.Records, newR)
	} else {
		// Updating an existing record so replace it.
		var found bool
		for i, r := range store.data.Records {
			if r.ID == newR.ID {
				found = true
				store.data.Records[i] = newR
				break
			}
		}
		if !found {
			store.data.Records = append(store.data.Records, newR)
		}
	}
	// Write the file backend.
	err = store.writeAll()
	return
}

// UpdateAll updates a slice of *record.Plan in the data-store.
// Param newRR is the slice of *record.Plan to be updated.
// If any record in newRR is new then it's ID is updated as well.
// Returns the error.
func (store *PlanStore) UpdateAll(newRR []*record.Plan) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.UpdateAll: %w", err)
		}
	}()

	err = fmt.Errorf("deprecated use Update")
	return
}

// Remove removes the record.Plan from the data-store.
// Param id is the record ID of the record.Plan to be removed.
// If the record is not found returns a nil error.
// Returns the error.
func (store *PlanStore) Remove(id uint64) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.Remove: %w", err)
		}
	}()

	err = fmt.Errorf("can not remove the plan record")
	return
}

func (store *PlanStore) readAll() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.readAll: %w", err)
		}
	}()

	// If the file doesn't exists then setup the data.
	var exists bool
	if exists, err = storage.Exists(store.uri); err != nil {
		return
	}
	if !exists {
		store.data.Records = make([]*record.Plan, 0, 1024)
		return
	}

	// Open.
	var rc fyne.URIReadCloser
	if rc, err = storage.Reader(store.uri); err != nil {
		return
	}
	defer func() {
		closeErr := rc.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Read.
	buffer := bytes.Buffer{}
	if _, err = buffer.ReadFrom(rc); err != nil {
		return
	}
	err = yaml.Unmarshal(buffer.Bytes(), &store.data)
	return
}

func (store *PlanStore) writeAll() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("PlanStore.writeAll: %w", err)
		}
	}()

	// Open.
	var wc fyne.URIWriteCloser
	if wc, err = storage.Writer(store.uri); err != nil {
		return
	}
	defer func() {
		closeErr := wc.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Convert.
	var bb []byte
	if bb, err = yaml.Marshal(&store.data); err != nil {
		return
	}

	// Write.
	_, err = wc.Write(bb)
	return
}
