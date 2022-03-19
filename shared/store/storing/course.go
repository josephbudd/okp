package storing

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/yaml.v2"

	"github.com/josephbudd/okp/shared/paths"
	"github.com/josephbudd/okp/shared/store/record"
)

type byCourseName []*record.Course

func (bcn byCourseName) Len() int           { return len(bcn) }
func (bcn byCourseName) Swap(i, j int)      { bcn[i], bcn[j] = bcn[j], bcn[i] }
func (bcn byCourseName) Less(i, j int) bool { return bcn[i].Name < bcn[j].Name }

type CourseData struct {
	LastID  uint64
	Records []*record.Course
}

// CourseStore is the API of the Course store.
// It is the implementation of the interface in /domain/store/storer/Course.go.
type CourseStore struct {
	uri  fyne.URI
	lock sync.Mutex
	data CourseData
}

// NewCourseStore constructs a new CourseStore.
// Param db is an open bolt data-store.
// Returns a pointer to the new CourseStore.
func NewCourseStore() (store *CourseStore) {
	store = &CourseStore{
		uri: paths.StoreURI("course.yaml"),
	}
	return
}

// Open opens the bolt data-store.
// Keeps the file in memory.
// Returns the error.
func (store *CourseStore) Open() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.Open: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	if err = store.readAll(); err != nil {
		return
	}
	sort.Sort(byCourseName(store.data.Records))

	return
}

// Close closes the data-store.
// Returns the error.
func (store *CourseStore) Close() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.Close: %w", err)
		}
	}()

	// The store is always closed.
	return
}

// Get retrieves one *record.Course from the data-store.
// Param id is the record ID.
// Returns a *record.Course and error.
// When no record is found, the returned *record.Course is nil and the returned error is nil.
func (store *CourseStore) Get(id uint64) (r *record.Course, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.Get: %w", err)
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

// GetAll retrieves all of the *record.Course records from the data-store.
// Returns a slice of *record.Course and error.
// When no records are found, the returned slice length is 0 and the returned error is nil.
func (store *CourseStore) GetAll() (rr []*record.Course, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.GetAll: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	rr = store.data.Records
	return
}

// Update updates the *record.Course in the data-store.
// Param newR is the *record.Course to be updated.
// If newR is a new record then r.ID is updated as well.
// Returns the error.
func (store *CourseStore) Update(newR *record.Course) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.Update: %w", err)
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
	sort.Sort(byCourseName(store.data.Records))
	// Write the file backend.
	err = store.writeAll()
	return
}

// UpdateAll updates a slice of *record.Course in the data-store.
// Param newRR is the slice of *record.Course to be updated.
// If any record in newRR is new then it's ID is updated as well.
// Returns the error.
func (store *CourseStore) UpdateAll(newRR []*record.Course) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.UpdateAll: %w", err)
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
	sort.Sort(byCourseName(store.data.Records))
	// Write the file backend.
	err = store.writeAll()
	return
}

// Remove removes the record.Course from the data-store.
// Param id is the record ID of the record.Course to be removed.
// If the record is not found returns a nil error.
// Returns the error.
func (store *CourseStore) Remove(id uint64) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.Remove: %w", err)
		}
	}()

	store.lock.Lock()
	defer store.lock.Unlock()

	// Find the record.
	var found bool
	var records []*record.Course
	for i, r := range store.data.Records {
		if r.ID == id {
			l := len(store.data.Records)
			records = make([]*record.Course, l-1)
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

func (store *CourseStore) readAll() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.readAll: %w", err)
		}
	}()

	// If the file doesn't exists then setup the data.
	var exists bool
	if exists, err = storage.Exists(store.uri); err != nil {
		return
	}
	if !exists {
		store.data.Records = make([]*record.Course, 0, 1024)
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

func (store *CourseStore) writeAll() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("CourseStore.writeAll: %w", err)
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
