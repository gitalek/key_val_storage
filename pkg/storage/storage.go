package storage

import (
	"encoding/json"
	"os"
	"sync"
)

const (
	ErrKeyNotFound    storageError = "storage: there is no such key in the storage"
	ErrNilDBInitState storageError = "storage: provided nil database init state"
)

type storageError string

func (e storageError) Error() string {
	return string(e)
}

type dbState map[string]string

type Storage struct {
	storage dbState
	mu      sync.RWMutex
}

// NewStorage is a Storage constructor.
func NewStorage(initState dbState) (*Storage, error) {
	if initState == nil {
		return nil, ErrNilDBInitState
	}
	return &Storage{storage: initState}, nil
}

// NewStorageFromFile is an alternative Storage constructor.
func NewStorageFromFile(pathToFile string, emptyInitStateAllowed bool) (*Storage, error) {
	state, err := loadState(pathToFile)
	if err != nil {
		if emptyInitStateAllowed {
			state := make(dbState)
			return NewStorage(state)
		}
		return nil, err
	}
	return NewStorage(state)
}

func (s *Storage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.storage[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return val, nil
}

func (s *Storage) List() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.storage
}

// Delete method returns removed value and error
// if there was no given key in the storage
func (s *Storage) Delete(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.storage[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	delete(s.storage, key)
	return val, nil
}

func (s *Storage) Upsert(items map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range items {
		s.storage[k] = v
	}
}

// Backup method returns copy of storage
func (s *Storage) Backup() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bu := make(map[string]string, len(s.storage))
	for k, v := range s.storage {
		bu[k] = v
	}
	return bu
}

func loadState(pathToFile string) (dbState, error) {
	var state dbState
	f, err := os.Open(pathToFile)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(f).Decode(&(state))
	if err != nil {
		return nil, err
	}
	return state, nil
}
