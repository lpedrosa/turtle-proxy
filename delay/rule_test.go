package delay

import (
	"testing"
)

func TestRuleStorage(t *testing.T) {
	t.Run("Get from an empty storage should not retrive anything", func(t *testing.T) {
		storage := DefaultStorage()

		_, ok := storage.Get("POST", "/non-existing/route")

		if ok {
			t.Fatal("Expected Storage to be empty!")
		}
	})
	t.Run("Get works as for existing delay", func(t *testing.T) {
		storage := DefaultStorage()

		method := "GET"
		path := "/route"
		delay := 10

		r := Rule{Method: method, Path: path, RequestDelay: delay}

		storage.Store(r)
		rule, ok := storage.Get(method, path)

		if !ok {
			t.Fatalf("Expected storage to contain rule with path: %s", path)
		}

		if rule.Method != method {
			t.Fatalf("Expected rule method to be: %s. Got: %s", method, rule.Method)
		}

		if rule.Path != path {
			t.Fatalf("Expected rule path to be: %s. Got: %s", path, rule.Path)
		}

		if rule.RequestDelay != delay {
			t.Fatalf("Expected rule request delay to be: %d. Got: %d", delay, rule.RequestDelay)
		}

		if rule.ResponseDelay != 0 {
			t.Fatalf("Expected rule request delay to be: %d. Got: %d", delay, rule.ResponseDelay)
		}
	})
	t.Run("Remove removes element from storage", func(t *testing.T) {
		storage := DefaultStorage()

		method := "GET"
		path := "/route"

		rule := Rule{Method: method, Path: path, RequestDelay: 1}

		storage.Store(rule)
		storage.Remove(method, path)
		_, ok := storage.Get(method, path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}
	})
	t.Run("Clear removes all elements from storage", func(t *testing.T) {
		storage := DefaultStorage()

		method := "GET"
		path := "/route"

		rule := Rule{Method: method, Path: path, RequestDelay: 1}

		storage.Store(rule)
		storage.Clear()
		_, ok := storage.Get(method, path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}
	})
	t.Run("Remove one rule should keep other rules intact", func(t *testing.T) {
		storage := DefaultStorage()

		method := "GET"

		path := "/route"
		rule := Rule{Method: method, Path: path, RequestDelay: 1}
		storage.Store(rule)

		path2 := "/route/{id}"
		rule2 := Rule{Method: method, Path: path2, RequestDelay: 1}
		storage.Store(rule2)

		storage.Remove(method, path)
		_, ok := storage.Get(method, path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}

		_, ok = storage.Get(method, path2)
		if !ok {
			t.Fatalf("Expected storage to contain rule with path: %s", path2)
		}
	})
	t.Run("Rule with no delay are discarded", func(t *testing.T) {
		storage := DefaultStorage()

		method := "GET"
		path := "/route"

		rule := Rule{Method: method, Path: path}

		storage.Store(rule)
		_, ok := storage.Get(method, path)

		if ok {
			t.Fatalf("Expected storage to not contain rule with path: %s", path)
		}
	})
}

func TestRuleStorageComplexPaths(t *testing.T) {
	t.Run("Get works for sub-path", func(t *testing.T) {
		storage := DefaultStorage()

		method := "GET"
		path := "/route/{everything:.*}"
		delay := 10

		r := Rule{Method: method, Path: path, RequestDelay: delay}

		storage.Store(r)

		subPath := "/route/109823/something"
		rule, ok := storage.Get(method, subPath)

		if !ok {
			t.Fatalf("Expected storage to contain rule with path: %s", path)
		}

		if rule.Method != method {
			t.Fatalf("Expected rule method to be: %s. Got: %s", method, rule.Method)
		}

		if rule.Path != path {
			t.Fatalf("Expected rule path to be: %s. Got: %s", path, rule.Path)
		}

		if rule.RequestDelay != delay {
			t.Fatalf("Expected rule request delay to be: %d. Got: %d", delay, rule.RequestDelay)
		}

		if rule.ResponseDelay != 0 {
			t.Fatalf("Expected rule request delay to be: %d. Got: %d", delay, rule.ResponseDelay)
		}
	})
}
