package main

import (
	"encoding/json"
	"os"
)

// A simple key-value store that is written to a single file.
// It does not play well with other processes that also try to write to the database file.
// It is however reasonably resistant to tampering and corruption.
type FileStore struct {
	store    map[string][]string
	filename string
}

// Create a new FileStore.
func NewFileStore(filename string) (*FileStore, error) {
	fileStore := &FileStore{
		store:    make(map[string][]string),
		filename: filename,
	}
	err := fileStore.Load()
	return fileStore, err
}

// Load the store from the database file.
func (fileStore *FileStore) Load() error {
	file, err := os.Open(fileStore.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&fileStore.store)
}

// Save the current store to the database file.
func (fileStore *FileStore) Save() error {
	file, err := os.Create(fileStore.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(fileStore.store)
}

func (fileStore *FileStore) Get(key string) ([]string, bool) {
	value, ok := fileStore.store[key]
	return value, ok
}

func (fileStore *FileStore) Set(key string, value []string) error {
	fileStore.store[key] = value
	return fileStore.Save()
}

func (fileStore *FileStore) Delete(key string) error {
	delete(fileStore.store, key)
	return fileStore.Save()
}
