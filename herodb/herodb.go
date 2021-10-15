package herodb

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotfound = errors.New("herodb: key not found")
)

type persister interface {
	Load(ctx context.Context) (map[string]string, error)
	Save(ctx context.Context, data map[string]string) error
}

// New instances a service
func New() *Store {
	return &Store{db: make(map[string]string)}
}

type Store struct {
	sync.RWMutex
	db         map[string]string
	persistent persister
}

// MakePersistent adds persistent storage provider and syncs given iterval
func (store *Store) MakePersistent(ctx context.Context, p persister, syncInterval time.Duration) error {
	if store.persistent != nil {
		return fmt.Errorf("already exists")
	}

	// try to restore from persistent storage
	db, err := p.Load(ctx)
	if err != nil {
		return fmt.Errorf("unable to restore from persistent storage: %w", err)
	}

	store.db = db
	ticker := time.NewTicker(syncInterval)
	go func() {
		for {
			select {
			// stop and quit
			case <-ctx.Done():
				ticker.Stop()
				store.Lock()
				if err := p.Save(ctx, store.db); err != nil {
					panic(err)
				}
				store.Unlock()
				return
			case <-ticker.C:
				store.Lock()
				if err := p.Save(ctx, store.db); err != nil {
					panic(err)
				}
				store.Unlock()
			}
		}
	}()
	return nil
}

// Get given key
func (s *Store) Get(ctx context.Context, key string) (string, error) {
	s.RLock()
	val, ok := s.db[key]
	s.RUnlock()
	if !ok {
		return "", ErrNotfound
	}

	return val, nil
}

// Set given key
func (s *Store) Set(ctx context.Context, key string, value string) error {
	s.Lock()
	s.db[key] = value
	s.Unlock()

	return nil
}

func (s *Store) Flush(ctx context.Context) error {
	return nil
}
