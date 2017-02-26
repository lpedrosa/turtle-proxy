package delay

import (
	"testing"
)

func TestRuleStorage(t *testing.T) {
	t.Run("Get from an empty storage should not retrive anything", func(t *testing.T) {
		storage := DefaultStorage()

		_, ok := storage.Get("/non-existing/route")

		if ok {
			t.Fatal("Expected Storage to be empty!")
		}
	})
	t.Run("Get works as for existing delay", func(t *testing.T) {
		storage := DefaultStorage()

		path := "/route"
		delay := 10

		storage.Store(path, delay)
		rule, ok := storage.Get(path)

		if !ok {
			t.Fatalf("Expected storage to contain rule with path: %s", path)
		}

		if rule.Path != path {
			t.Fatalf("Expected rule path to be: %s. Got: %s", path, rule.Path)
		}

		if rule.Path != path {
			t.Fatalf("Expected rule delay to be: %d. Got: %d", delay, rule.Delay)
		}
	})
	t.Run("Remove removes element from storage", func(t *testing.T) {
		storage := DefaultStorage()

		path := "/route"

		storage.Store(path, 10)
		storage.Remove(path)
		_, ok := storage.Get(path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}
	})
	t.Run("Clear removes all elements from storage", func(t *testing.T) {
		storage := DefaultStorage()

		path := "/route"

		storage.Store(path, 10)
		storage.Clear()
		_, ok := storage.Get(path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}
	})
	t.Run("Remove one rule should keep other rules intact", func(t *testing.T) {
		storage := DefaultStorage()

		path := "/route"
		storage.Store(path, 10)

		path2 := "/route/{id}"
		storage.Store(path2, 10)

		storage.Remove(path)
		_, ok := storage.Get(path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}

		_, ok = storage.Get(path2)
		if !ok {
			t.Fatalf("Expected storage to contain rule with path: %s", path2)
		}
	})
}
