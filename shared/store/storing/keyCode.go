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

type KeyCodeData struct {
	LastID  uint64
	Records []*record.KeyCode
}

// KeyCodeStore is the API of the KeyCode store.
// It is the implementation of the interface in /domain/store/storer/KeyCode.go.
type KeyCodeStore struct {
	uri  fyne.URI
	lock sync.Mutex
	data KeyCodeData
}

// NewKeyCodeStore constructs a new KeyCodeStore.
// Param db is an open bolt data-store.
// Returns a pointer to the new KeyCodeStore.
func NewKeyCodeStore() (store *KeyCodeStore) {
	store = &KeyCodeStore{
		uri: paths.StoreURI("keycode.yaml"),
	}
	return
}

// Open opens the bolt data-store.
// Keeps the file in memory.
// Returns the error.
func (store *KeyCodeStore) Open() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.Open: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	err = store.readAll()
	return
}

// Close closes the data-store.
// Returns the error.
func (store *KeyCodeStore) Close() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.Close: %w", err)
		}
	}()

	// The store is always closed.
	return
}

// Get retrieves one *record.KeyCode from the data-store.
// Param id is the record ID.
// Returns a *record.KeyCode and error.
// When no record is found, the returned *record.KeyCode is nil and the returned error is nil.
func (store *KeyCodeStore) Get(id uint64) (r *record.KeyCode, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.Get: %w", err)
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

// GetAll retrieves all of the *record.KeyCode records from the data-store.
// Returns a slice of *record.KeyCode and error.
// When no records are found, the returned slice length is 0 and the returned error is nil.
func (store *KeyCodeStore) GetAll() (rr []*record.KeyCode, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.GetAll: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	rr = store.data.Records
	return
}

// Update updates the *record.KeyCode in the data-store.
// Param newR is the *record.KeyCode to be updated.
// If newR is a new record then r.ID is updated as well.
// Returns the error.
func (store *KeyCodeStore) Update(newR *record.KeyCode) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.Update: %w", err)
		}
	}()

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

// UpdateAll updates a slice of *record.KeyCode in the data-store.
// Param newRR is the slice of *record.KeyCode to be updated.
// If any record in newRR is new then it's ID is updated as well.
// Returns the error.
func (store *KeyCodeStore) UpdateAll(newRR []*record.KeyCode) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.UpdateAll: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	for _, newR := range newRR {
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
	}
	// Write the file backend.
	err = store.writeAll()
	return
}

// Remove removes the record.KeyCode from the data-store.
// Param id is the record ID of the record.KeyCode to be removed.
// If the record is not found returns a nil error.
// Returns the error.
func (store *KeyCodeStore) Remove(id uint64) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.Remove: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	// Find the record.
	var found bool
	var records []*record.KeyCode
	for i, r := range store.data.Records {
		if r.ID == id {
			l := len(store.data.Records)
			records = make([]*record.KeyCode, l-1)
			if i > 0 {
				copy(records, store.data.Records[:i])
			}
			j := i
			// Skip over the unwanted record.
			for i++; i < l; i++ {
				records[j] = store.data.Records[i]
				j++
			}
			found = true
			break
		}
	}
	if !found {
		// No error if not found.
		return
	}
	store.data.Records = records
	err = store.writeAll()
	return
}

func (store *KeyCodeStore) readAll() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.readAll: %w", err)
		}
	}()

	// If the file doesn't exists then setup the data.
	var exists bool
	if exists, err = storage.Exists(store.uri); err != nil {
		return
	}
	if !exists {
		store.data.Records = make([]*record.KeyCode, 0, 1024)
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

func (store *KeyCodeStore) writeAll() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeStore.writeAll: %w", err)
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
