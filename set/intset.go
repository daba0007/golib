package set

import (
	"sync"
)

// IntSet int set
type IntSet struct {
	sync.RWMutex
	m map[int]bool
}

// NewIntSet new int set
func NewIntSet(items ...int) *IntSet {
	s := &IntSet{
		m: make(map[int]bool, len(items)),
	}
	s.add(items...)
	return s
}

// Add add items for int set
func (s *IntSet) add(items ...int) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// Remove remove items for int set
func (s *IntSet) remove(items ...int) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// Has judge items in int set
func (s *IntSet) has(items ...int) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// Count get count in int set
func (s *IntSet) count() int {
	return len(s.m)
}

// Clear clear int set
func (s *IntSet) clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int]bool{}
}

// Empty judge int set
func (s *IntSet) empty() bool {
	return len(s.m) == 0
}

// List get int list from int set
func (s *IntSet) list() []int {
	s.RLock()
	defer s.RUnlock()
	list := make([]int, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// Union get unionsection from int set
func (s *IntSet) union(sets ...*IntSet) *IntSet {
	r := NewIntSet(s.list()...)
	for _, set := range sets {
		for e := range set.m {
			r.m[e] = true
		}
	}
	return r
}

// Minus get minuxsection from int set
func (s *IntSet) minus(sets ...*IntSet) *IntSet {
	r := NewIntSet(s.list()...)
	for _, set := range sets {
		for e := range set.m {
			if _, ok := s.m[e]; ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Intersect get intersection from int set
func (s *IntSet) intersect(sets ...*IntSet) *IntSet {
	r := NewIntSet(s.list()...)
	for _, set := range sets {
		for e := range s.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Complement get complementsection from int set
func (s *IntSet) complement(full *IntSet) *IntSet {
	r := NewIntSet()
	for e := range full.m {
		if _, ok := s.m[e]; !ok {
			r.add(e)
		}
	}
	return r
}