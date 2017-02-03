package main

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
)

type DelayedDownload struct {
	Slug  string
	URL   *url.URL
	Delay uint
}

type DelayedStorage struct {
	storage map[string]DelayedDownload
	sLock   sync.RWMutex
}

func newDelayedStorage() *DelayedStorage {
	ds := &DelayedStorage{}
	ds.storage = make(map[string]DelayedDownload)

	return ds
}

func (ds *DelayedStorage) Store(delayed DelayedDownload) error {
	// only one caller can write at a time
	ds.sLock.Lock()
	defer ds.sLock.Unlock()

	storageKey := delayed.Slug

	_, ok := ds.storage[storageKey]

	if !ok {
		errMsg := fmt.Sprintf("Cannot store existing key %s", storageKey)
		return errors.New(errMsg)
	}

	ds.storage[storageKey] = delayed
	return nil
}

func (ds *DelayedStorage) GetAndRemove(slug string) (delayed *DelayedDownload, ok bool) {
	// only one caller can write at a time
	// still a write lock because we will delete the entry later
	ds.sLock.Lock()
	defer ds.sLock.Unlock()

	*delayed, ok = ds.storage[slug]

	// we do not need the element again
	delete(ds.storage, slug)

	return delayed, ok
}
