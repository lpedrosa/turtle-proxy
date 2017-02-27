package delay

import (
	"log"
	"sync"
)

//-----------------------
// Rule
//-----------------------

type Rule struct {
	Method        string
	Path          string
	RequestDelay  int
	ResponseDelay int
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

//-----------------------
// Storage
//-----------------------

type storageKey string

func newKey(method string, path string) storageKey {
	return storageKey(method + ":" + path)
}

type ruleWithMatch struct {
	matchId string
	*Rule
}

type RuleStorage struct {
	storage map[storageKey]*ruleWithMatch
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
	rs.storage = make(map[storageKey]*ruleWithMatch)
	rs.matcher = matcher
	return rs
}

func (rs *RuleStorage) Store(rule Rule) {
	// TODO validate the rest of the rule?
	if rule.RequestDelay < 0 || rule.ResponseDelay < 0 {
		panic("delays cannot be negative!")
	}

	if rule.RequestDelay == 0 && rule.ResponseDelay == 0 {
		log.Printf("rule: discarding rule with no delay %#v", rule)
		return
	}

	// only one caller can write at a time
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	method := rule.Method
	path := rule.Path

	// add rule to matcher
	matchRuleId := rs.matcher.Add(method, path)

	// add rule to storage
	key := newKey(method, path)
	rs.storage[key] = &ruleWithMatch{matchId: matchRuleId, Rule: &rule}
}

func (rs *RuleStorage) Get(method string, path string) (rule *Rule, ok bool) {
	// fine to use a read lock here
	rs.sLock.RLock()
	defer rs.sLock.RUnlock()

	ok = rs.matcher.Match(method, path)
	if !ok {
		return nil, ok
	}

	key := newKey(method, path)
	ruleWithMatch, ok := rs.storage[key]

	if !ok {
		log.Panicf("rule: No rule found for match: %s\n", path)
	}

	return ruleWithMatch.Rule, ok
}

func (rs *RuleStorage) Remove(method string, path string) {
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	key := newKey(method, path)
	ruleWithMatch, ok := rs.storage[key]
	if !ok {
		// nothing to remove
		return
	}

	matchId := ruleWithMatch.matchId
	rs.matcher.Remove(matchId)

	delete(rs.storage, key)
}

func (rs *RuleStorage) Clear() {
	rs.sLock.Lock()
	defer rs.sLock.Unlock()

	// clear storage
	rs.storage = make(map[storageKey]*ruleWithMatch)

	// clear matcher
	rs.matcher.Clear()
}
