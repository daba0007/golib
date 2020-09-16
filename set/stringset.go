package set

import (
	"sync"
)

// StringSet string set
type StringSet struct {
	sync.RWMutex
	m map[string]bool
}

// NewStringSet new string set
func NewStringSet(items ...string) *StringSet {
	s := &StringSet{
		m: make(map[string]bool, len(items)),
	}
	s.add(items...)
	return s
}

// Add add items for string set
func (s *StringSet) add(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// Remove remove items for string set
func (s *StringSet) remove(items ...string) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// Has judge items in string set
func (s *StringSet) has(items ...string) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// Count get count in string set
func (s *StringSet) count() int {
	return len(s.m)
}

// Clear clear string set
func (s *StringSet) clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[string]bool{}
}

// Empty judge string set
func (s *StringSet) empty() bool {
	return len(s.m) == 0
}

// List get string list from string set
func (s *StringSet) list() []string {
	s.RLock()
	defer s.RUnlock()
	list := make([]string, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// Union get unionsection from string set
func (s *StringSet) union(sets ...*StringSet) *StringSet {
	r := NewStringSet(s.list()...)
	for _, set := range sets {
		for e := range set.m {
			r.m[e] = true
		}
	}
	return r
}

// Minus get minuxsection from string set
func (s *StringSet) minus(sets ...*StringSet) *StringSet {
	r := NewStringSet(s.list()...)
	for _, set := range sets {
		for e := range set.m {
			if _, ok := s.m[e]; ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Intersect get intersection from string set
func (s *StringSet) intersect(sets ...*StringSet) *StringSet {
	r := NewStringSet(s.list()...)
	for _, set := range sets {
		for e := range s.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Complement get complementsection from string set
func (s *StringSet) complement(full *StringSet) *StringSet {
	r := NewStringSet()
	for e := range full.m {
		if _, ok := s.m[e]; !ok {
			r.add(e)
		}
	}
	return r
}
