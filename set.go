// Goset is a thread safe SET data structure implementation.
// The thread safety encompasses all operations on one set.
// Operations on multiple sets are consistent in that the elements
// of each set used was valid at exactly one point in time between the
// start and the end of the operation.

package goset

import (
	"fmt"
	"strings"
	"sync"
)

type Set struct {
	m map[interface{}]struct{} // struct{} doesn't take up space
	l sync.RWMutex             // we name it because we don't want to expose it
}

var keyexists = struct{}{}

// New creates and initialize a new Set. It accepts a variable number of
// arguments to populate the initial set.
func New(items ...interface{}) *Set {
	return (*Set)(nil).Add(items...)
}

// Add includes the specified items to the set.
// The Set may be modified or a new one may be created.
// The resulting Set is returned.
func (s *Set) Add(items ...interface{}) *Set {
	if len(items) == 0 {
		return s
	}
	if s == nil {
		s = &Set{
			m: make(map[interface{}]struct{}),
		}
	}
	s.l.Lock()
	for _, item := range items {
		s.m[item] = keyexists
	}
	s.l.Unlock()
	return s
}

// Remove deletes the specified items from the set.
// The Set may be modified or a new one may be created.
// The resulting Set is returned.
func (s *Set) Remove(items ...interface{}) *Set {
	if s == nil || len(items) == 0 {
		return s
	}
	empty := false
	s.l.Lock()
	for _, item := range items {
		delete(s.m, item)
	}
	empty = len(s.m) == 0
	s.l.Unlock()
	if empty {
		return nil
	}
	return s
}

// Has reports whether all arguments are contained in the Set
func (s *Set) Has(items ...interface{}) (has bool) {
	if s == nil {
		return false
	}
	// assume checked for empty item, which not exist
	if len(items) == 0 {
		return false
	}
	s.l.RLock()
	for _, item := range items {
		if _, has = s.m[item]; !has {
			break
		}
	}
	s.l.RUnlock()
	return has
}

// Size returns the number of items in the Set
func (s *Set) Size() (size int) {
	if s == nil {
		return 0
	}
	s.l.RLock()
	size = len(s.m)
	s.l.RUnlock()
	return size
}

// Clear returns an empty Set
func (s *Set) Clear() *Set {
	return nil
}

// IsEmpty reports whether the Set is empty
func (s *Set) IsEmpty() bool {
	return s == nil
}

// IsEqual reports whether s is contained in t and t is contained in s
func (s *Set) IsEqual(t *Set) (equal bool) {
	if s == nil || t == nil || s == t {
		return s == t
	}
	s.l.RLock()
	t.l.RLock()
	if equal = len(s.m) == len(t.m); equal {
		for item := range t.m {
			if _, equal = s.m[item]; !equal {
				break
			}
		}
	}
	t.l.RUnlock()
	s.l.RUnlock()
	return equal
}

// IsSubset reports whether t is a subset of s
func (s *Set) IsSubset(t *Set) (subset bool) {
	if s == nil || t == nil || s == t {
		return s != nil || s == t
	}
	s.l.RLock()
	t.l.RLock()
	subset = true
	for item := range t.m {
		if _, subset = s.m[item]; !subset {
			break
		}
	}
	t.l.RUnlock()
	s.l.RUnlock()
	return subset
}

// IsSuperset reports whether t is a superset of s
func (s *Set) IsSuperset(t *Set) bool {
	return t.IsSubset(s)
}

// String returns a string representation of s
func (s *Set) String() string {
	if s == nil {
		return "[]"
	}
	items := s.List()
	t := make([]string, 0, len(items))
	for _, item := range items {
		t = append(t, fmt.Sprintf("%v", item))
	}
	return fmt.Sprintf("[%s]", strings.Join(t, ", "))
}

// List returns a slice of all items
func (s *Set) List() (list []interface{}) {
	if s == nil {
		return list
	}
	s.l.RLock()
	list = make([]interface{}, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	s.l.RUnlock()
	return list
}

// Copy returns a new Set with a copy of s.
func (s *Set) Copy() *Set {
	if s == nil {
		return nil
	}
	return New(s.List()...)
}

// Union is the merger of two sets. It returns a new set with the element in s
// and t combined.
func (s *Set) Union(t *Set) *Set {
	if t == nil {
		if s == nil && t == nil {
			return nil
		}
		return s.Copy()
	}
	if s == nil {
		return t.Copy()
	}
	return s.Copy().Merge(t)
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *Set) Merge(t *Set) *Set {
	if s == nil {
		return t.Copy()
	}
	if t == nil {
		return s
	}
	s.l.Lock()
	t.l.RLock()
	for item := range t.m {
		s.m[item] = keyexists
	}
	t.l.RUnlock()
	s.l.Unlock()
	return s
}

// Exclude removes the set items contained in t from set s. Please be aware that
// it's not the opposite of Merge.
func (s *Set) Exclude(t *Set) *Set {
	if s == nil {
		return nil
	}
	if t == nil {
		return s
	}
	return s.Remove(t.List()...)
}

// Intersection returns a new set which contains items which is in both s and t.
func (s *Set) Intersection(t *Set) *Set {
	if s == nil || t == nil {
		return nil
	}
	if s == t {
		return s.Copy()
	}
	c := s.Copy()
	return c.Exclude(c.Difference(t))
}

// Intersection returns a new set which contains items which are both s but not in t.
func (s *Set) Difference(t *Set) *Set {
	return s.Copy().Exclude(t)
}

// Symmetric returns a new set which s is the difference of items  which are in
// one of either, but not in both.
func (s *Set) SymmetricDifference(t *Set) *Set {
	cs := s.Copy()
	ct := t.Copy()
	return cs.Union(ct).Exclude(cs.Intersection(ct))
}

// StringSlice is a helper function that returns a slice of strings of s. If
// the set contains mixed types of items only items of type string are returned.
func (s *Set) StringSlice() []string {
	list := s.List()
	slice := make([]string, 0, len(list))
	for _, item := range list {
		v, ok := item.(string)
		if !ok {
			continue
		}
		slice = append(slice, v)
	}
	return slice
}

// IntSlice is a helper function that returns a slice of ints of s. If
// the set contains mixed types of items only items of type int are returned.
func (s *Set) IntSlice() []int {
	list := s.List()
	slice := make([]int, 0, len(list))
	for _, item := range list {
		v, ok := item.(int)
		if !ok {
			continue
		}
		slice = append(slice, v)
	}
	return slice
}
