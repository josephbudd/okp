package record

/*

	TODO:

	You need to complete this record definition.

*/

// KeyCode is the KeyCode record.
type KeyCode struct {
	ID            uint64
	Name          string
	Character     string
	DitDah        string
	IsWord        bool
	IsCompression bool
	IsNotReal     bool
}
