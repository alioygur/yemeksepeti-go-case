package herodb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// NewFilepersist instances file persist service
func NewFilepersist(file string) (*FilePersist, error) {
	if !fileExists(file) {
		f, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		f.WriteString(`{}`)
		f.Close()
	}
	return &FilePersist{filename: file}, nil
}

type FilePersist struct {
	filename string
}

// Load reads and decodes db file content then returns it
func (f *FilePersist) Load(ctx context.Context) (map[string]string, error) {
	fr, err := os.Open(f.filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer fr.Close()

	db := make(map[string]string)
	if err := json.NewDecoder(fr).Decode(&db); err != nil {
		return nil, fmt.Errorf("unable to decode file: %w", err)
	}
	return db, nil
}

// Save persists given data to file
func (f *FilePersist) Save(ctx context.Context, data map[string]string) error {
	fr, err := os.OpenFile(f.filename, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file for persist: %w", err)
	}
	defer fr.Close()
	if err := json.NewEncoder(fr).Encode(data); err != nil {
		return fmt.Errorf("unable to encode: %w", err)
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
