package delay

import (
	"testing"
)

func TestGorillaMatcher(t *testing.T) {
	matcher := newGorillaMatcher()

	pathExpression := "/path/{id}"
	method := "GET"

	matcher.Add(method, pathExpression)

	t.Run("Should return pathExpression for explicit path", func(t *testing.T) {
		concretePath := "/path/12093"

		expression, ok := matcher.Match(method, concretePath)

		if !ok {
			t.Fatalf("Expected to find a match for: %s. Got none.", concretePath)
		}

		if expression != pathExpression {
			t.Fatalf("Expected expression to be: %s. Got %s.", pathExpression, expression)
		}
	})

	t.Run("Should not math unrelated paths", func(t *testing.T) {
		concretePath := "/pat"

		_, ok := matcher.Match(method, concretePath)

		if ok {
			t.Fatalf("Expected to not find a match for: %s. Got one.", concretePath)
		}
	})
	t.Run("Should not match subpaths ", func(t *testing.T) {
		concretePath := "/path/12093/lol/201983"

		_, ok := matcher.Match(method, concretePath)

		if ok {
			t.Fatalf("Expected to not find a match for: %s. Got one.", concretePath)
		}
	})
}

type TestRule struct {
	Method string
	Path   string
}

func TestGorillaMatcherMultipleRules(t *testing.T) {
	matcher := newGorillaMatcher()

	rule1 := TestRule{Method: "GET", Path: "/path"}
	rule2 := TestRule{Method: "GET", Path: "/path/{id}"}
	rule3 := TestRule{Method: "GET", Path: "/path/{id}/sub"}
	rule4 := TestRule{Method: "GET", Path: "/path/{id}/sub/{id}"}

	rules := []TestRule{rule1, rule2, rule3, rule4}

	for _, r := range rules {
		matcher.Add(r.Method, r.Path)
	}

	t.Run("Explicit path should match root", func(t *testing.T) {
		concretePath := "/path"
		method := "GET"

		expression, ok := matcher.Match(method, concretePath)

		if !ok {
			t.Fatalf("Expected to find a match for: %s. Got none.", concretePath)
		}

		if expression != rule1.Path {
			t.Fatalf("Expected expression to be: %s. Got %s.", rule1.Path, expression)
		}
	})
	t.Run("1 level child leaf example match", func(t *testing.T) {
		concretePath := "/path/12"
		method := "GET"

		expression, ok := matcher.Match(method, concretePath)

		if !ok {
			t.Fatalf("Expected to find a match for: %s. Got none.", concretePath)
		}

		if expression != rule2.Path {
			t.Fatalf("Expected expression to be: %s. Got %s.", rule2.Path, expression)
		}
	})
	t.Run("1 level child with children example match", func(t *testing.T) {
		concretePath := "/path/12/sub"
		method := "GET"

		expression, ok := matcher.Match(method, concretePath)

		if !ok {
			t.Fatalf("Expected to find a match for: %s. Got none.", concretePath)
		}

		if expression != rule3.Path {
			t.Fatalf("Expected expression to be: %s. Got %s.", rule3.Path, expression)
		}
	})
	t.Run("2 level child leaf example match", func(t *testing.T) {
		concretePath := "/path/12/sub/something"
		method := "GET"

		expression, ok := matcher.Match(method, concretePath)

		if !ok {
			t.Fatalf("Expected to find a match for: %s. Got none.", concretePath)
		}

		if expression != rule4.Path {
			t.Fatalf("Expected expression to be: %s. Got %s.", rule4.Path, expression)
		}
	})

}
