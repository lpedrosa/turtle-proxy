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

type ruleWithMatch struct {
	matchId string
	*Rule
}

type RuleStorage struct {
	storage map[string]*ruleWithMatch
	matcher RuleMatcher
	sLock   sync.RWMutex
}

func DefaultStorage() *RuleStorage {
	// use gorilla mux as a matcher
	matcher := newGorillaMatcher()
	return NewStorage(matcher)
}

func NewStorage(matcher RuleMatcher) *RuleStorage {
	rs := &RuleStorage{}
	rs.storage = make(map[string]*ruleWithMatch)
	rs.matcher = matcher
	return rs
}

func (rs *RuleStorage) Store(path string, delay int) {
	if delay < 0 {
		panic("delay cannot be negative!")
	}

	// only one caller can write at a time
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	matchRuleId := rs.matcher.Add("GET", path)
	rule := &Rule{Path: path, Delay: delay}
	rs.storage[path] = &ruleWithMatch{matchId: matchRuleId, Rule: rule}
}

func (rs *RuleStorage) Get(path string) (rule *Rule, ok bool) {
	// fine to use a read lock here
	rs.sLock.RLock()
	defer rs.sLock.RUnlock()

	ok = rs.matcher.Match("GET", path)
	if !ok {
		return nil, ok
	}

	ruleWithMatch, ok := rs.storage[path]

	if !ok {
		log.Panicf("No rule found for match: %s\n", path)
	}

	return ruleWithMatch.Rule, ok
}

func (rs *RuleStorage) Remove(path string) {
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	ruleWithMatch, ok := rs.storage[path]
	if !ok {
		// nothing to remove
		return
	}

	matchId := ruleWithMatch.matchId
	rs.matcher.Remove(matchId)

	delete(rs.storage, path)
}

func (rs *RuleStorage) Clear() {
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	// clear storage
	rs.storage = make(map[string]*ruleWithMatch)

	// clear matcher
	rs.matcher.Clear()
}
