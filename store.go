// Store stores a string as a line in a file
// use Init(filename string) to set a store with the given filename
package stringstore

import (
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	sync.RWMutex
	path string
}

// New starts the store at the file location.
// It will read any values already stored in the file
func New(path string) (*Store, error) {
	var err error
	out := Store{path: path}
	err = os.MkdirAll(filepath.Dir(out.path), 0644)
	if err != nil {
		return &out, err
	}
	return &out, err
}

// Add adds the string to the store.
// It uses append to only write the new line to the file and not overwrite the whole file
func (this *Store) Add(filename string) error {
	this.Lock()
	defer this.Unlock()

	// Open the file in append mode, creating it if it doesn't exist
	file, err := os.OpenFile(this.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// write new line
	_, err = file.WriteString(filename + "\n")
	return err
}

// Pop removes the last line from the store and returns it, returns an empty string if store is empty.
// It uses truncate to efficiently delete only last line and not overwrite the whole file.
func (this *Store) Pop() (string, error) {
	var out string
	this.Lock()
	defer this.Unlock()

	// open the file for reading and writing
	file, err := os.OpenFile(this.path, os.O_RDWR, 0644)
	if err != nil {
		return out, err
	}
	defer file.Close()
	// If the file is empty, there's nothing to do
	stat, err := file.Stat()
	if err != nil {
		return out, err
	}
	if stat.Size() == 0 {
		return out, nil
	}

	// get the offset of the line to remove
	out, offset, err := getLastLine(file)
	if err != nil {
		return out, err
	}

	// truncate the file at the offset to remove the line
	err = file.Truncate(offset)
	if err != nil {
		return out, err
	}
	return out, nil
}

// getLastLine starts from the end of the file and returns the position of the last line
func getLastLine(file *os.File) (string, int64, error) {
	var out string
	// Get the file size.
	stat, err := file.Stat()
	if err != nil {
		return out, 0, err
	}
	// if file is empty return an empty string
	if stat.Size() <= 0 {
		return out, 0, nil
	}

	var text []byte
	var offset int64
	// Start from the end of the file.
	for offset = stat.Size() - 1; offset >= 0; offset-- {
		b := make([]byte, 1) // because of concurrency we need to make a new buffer each time
		_, err := file.ReadAt(b, offset)
		if err != nil {
			return out, 0, err
		}

		// Stop if we find a line break.
		if b[0] == '\n' {
			// skip if we find newline as the first character
			if len(text) == 0 {
				continue
			} else {
				break
			}
		}

		// Have to prepend the byte to the slice since we are reading the file backwards
		text = append(b, text...)
	}
	out = string(text)
	if offset < 0 {
		return out, 0, nil
	}
	return out, offset + 1, nil
}
