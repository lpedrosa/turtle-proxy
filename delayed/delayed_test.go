package delayed

import (
	"testing"
)

func TestEmptyDelayedStore(t *testing.T) {
	t.Run("Retrieving from an empty storage should not retrive anything", func(t *testing.T) {
		storage := NewStorage()

		_, ok := storage.GetAndRemove("non-existing")

		if ok {
			t.Fatal("Expected Storage to be empty!")
		}
	})
}

func TestStoreOperation(t *testing.T) {
	t.Run("Store returns an error when storing the same element twice", func(t *testing.T) {
		storage := NewStorage()

		delayed := DelayedDownload{}

		if err := storage.Store(delayed); err != nil {
			t.Fatal("Expected first Store to succeed")
		}

		if err := storage.Store(delayed); err == nil {
			t.Fatal("Expected Store operation to return an error")
		}
	})
	t.Run("Store returns an error when storing an element with the same id twice", func(t *testing.T) {
		storage := NewStorage()

		delayed := DelayedDownload{Id: "1", Delay: 2}
		delayed2 := DelayedDownload{Id: "1", Delay: 5}

		if err := storage.Store(delayed); err != nil {
			t.Fatal("Expected first Store to succeed")
		}

		if err := storage.Store(delayed2); err == nil {
			t.Fatal("Expected Store operation to return an error")
		}
	})
}

func TestGetAndRemoveOperation(t *testing.T) {
	t.Run("GetAndRemove works as expected", func(t *testing.T) {
		storage := NewStorage()

		delayed := DelayedDownload{Id: "1"}

		storage.Store(delayed)
		delayed2, _ := storage.GetAndRemove(delayed.Id)

		if delayed.Id != delayed2.Id {
			t.Fatal("Expected GetAndRemove to return the original element")
		}
	})
	t.Run("GetAndRemove removes element from storage", func(t *testing.T) {
		storage := NewStorage()

		delayed := DelayedDownload{Id: "1"}

		storage.Store(delayed)
		_, ok := storage.GetAndRemove(delayed.Id)
		_, ok2 := storage.GetAndRemove(delayed.Id)

		if !ok {
			t.Fatal("Expected GetAndRemove to succeed the first time")
		}

		if ok2 {
			t.Fatal("Expected GetAndRemove to fail the second time")
		}
	})
}
