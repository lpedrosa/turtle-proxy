package delay

import (
	"log"
	"sync"
)

type Rule struct {
	Path  string
	Delay int
}

//-----------------------
// Matcher
//-----------------------

type RuleMatcher interface {
	Add(method string, path string) string
	Match(method string, path string) bool
	Remove(id string)
	Clear()
}

type noopMatcher struct{}

func (n *noopMatcher) Add(method string, path string) string {
	return ""
}

func (n *noopMatcher) Match(method string, path string) bool {
	return false
}

func (n *noopMatcher) Remove(id string) {
}

func (n *noopMatcher) Clear() {
}

//-----------------------
// Storage
//-----------------------

type RuleStorage struct {
	storage map[string]*Rule
	matcher RuleMatcher
	sLock   sync.RWMutex
}

func DefaultStorage() *RuleStorage {
	// FIXME implement a matcher
	matcher := &noopMatcher{}
	return NewStorage(matcher)
}

func NewStorage(matcher RuleMatcher) *RuleStorage {
	ds := &RuleStorage{}
	ds.storage = make(map[string]*Rule)
	ds.matcher = matcher
	return ds
}

func (ds *RuleStorage) Store(path string, delay int) {
	if delay < 0 {
		panic("delay cannot be negative!")
	}

	// only one caller can write at a time
	ds.sLock.Lock()
	defer ds.sLock.Unlock()

	ds.matcher.Add("GET", path)
	ds.storage[path] = &Rule{Path: path, Delay: delay}
}

func (ds *RuleStorage) Get(path string) (rule *Rule, ok bool) {
	// only one caller can write at a time
	// still a write lock because we will delete the entry later
	ds.sLock.Lock()
	defer ds.sLock.Unlock()

	ok = ds.matcher.Match("GET", path)
	if !ok {
		return nil, ok
	}

	rule, ok = ds.storage[path]

	if !ok {
		log.Panicf("No rule found for match: %s\n", path)
	}

	return rule, ok
}

func (ds *RuleStorage) Remove(path string) {
}

func (ds *RuleStorage) Clear() {
}
