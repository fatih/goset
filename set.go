// Goset is a thread safe SET data structure implementation
package goset

import (
	"fmt"
	"strings"
	"sync"
)

type Set struct {
	m map[interface{}]struct{}
	l sync.RWMutex // we name it because we don't want to expose it
}

// New creates and initialize a new Set. It's accept a variable number of
// arguments to populate the initial set. If nothing passed a Set with zero
// size is created.
func New(items ...interface{}) *Set {
	s := &Set{
		m: make(map[interface{}]struct{}), // struct{} doesn't take up space
	}
	s.Add(items...)
	return s
}

// Add includes the specified items (one or more) to the set. If passed nothing
// it silently returns.
func (s *Set) Add(items ...interface{}) {
	if s == nil {
		// we can get rid of this when Add returns its result...
		panic("can not add to nil")
	}
	if len(items) == 0 {
		return
	}
	s.l.Lock()
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	s.l.Unlock()
}

// Remove deletes the specified items from the set. If passed nothing it
// silently returns.
func (s *Set) Remove(items ...interface{}) {
	if s == nil {
		return
	}
	if len(items) == 0 {
		return
	}
	s.l.Lock()
	for _, item := range items {
		delete(s.m, item)
	}
	s.l.Unlock()
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
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

// Size returns the number of items in a set.
func (s *Set) Size() (size int) {
	if s == nil {
		return 0
	}
	s.l.RLock()
	size = len(s.m)
	s.l.RUnlock()
	return size
}

// Clear removes all items from the set.
func (s *Set) Clear() {
	if s == nil {
		return
	}
	s.l.Lock()
	s.m = make(map[interface{}]struct{})
	s.l.Unlock()
}

// IsEmpty checks for emptiness of the set.
func (s *Set) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in size and have the same items.
func (s *Set) IsEqual(t *Set) (equal bool) {
	if s == nil || t == nil || s == t {
		return s == t
	}
	s.l.RLock()
	t.l.RLock()
	ss := len(s.m)
	st := len(t.m)
	equal = ss == st && ss == len(s.Union(t).m)
	t.l.Unlock()
	s.l.Unlock()
	return equal
}

// IsSubset tests t is a subset of s.
func (s *Set) IsSubset(t *Set) (subset bool) {
	if s == nil || t == nil || s == t {
		return s != nil || s == t
	}
	s.l.RLock()
	t.l.RLock()
	for _, item := range t.m {
		if _, subset = s.m[item]; !subset {
			break
		}
	}
	t.l.Unlock()
	s.l.Unlock()
	return subset
}

// IsSuperset tests if t is a superset of s.
func (s *Set) IsSuperset(t *Set) bool {
	return t.IsSubset(s)
}

// String representation of s
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
		return []interface{}{}
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
	return New(s.List()...)
}

// Union is the merger of two sets. It returns a new set with the element in s
// and t combined.
func (s *Set) Union(t *Set) *Set {
	if t == nil {
		return s.Copy()
		if s == nil && t == nil {
			return nil
		}
	}
	if s == nil {
		return t.Copy()
	}
	u := New(t.List()...)
	u.Add(s.List()...)
	return u
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s *Set) Merge(t *Set) {
	if t == nil {
		return
	}
	if s == nil {
		panic("can not merge with nil") // same as with Add
	}
	s.Lock()
	t.Rlock()
	for _, item := range t.m {
		s.m[item] = struct{}
	}
	t.Unlock()
	s.Unlock()
}

// Separate removes the set items containing in t from set s. Please be aware that
// it's not the opposite of Merge.
func (s *Set) Separate(t *Set) {
	if s == nil || t == nil {
		return
	}
	return s.Remove(t.List()...)
}

// Intersection returns a new set which contains items which is in both s and t.
func (s *Set) Intersection(t *Set) *Set {
	u := New()
	for _, item := range s.List() {
		if t.Has(item) {
			u.Add(item)
		}
	}
	for _, item := range t.List() {
		if s.Has(item) {
			u.Add(item)
		}
	}
	return u
}

// Intersection returns a new set which contains items which are both s but not in t.
func (s *Set) Difference(t *Set) *Set {
	u := New()
	for _, item := range s.List() {
		if !t.Has(item) {
			u.Add(item)
		}
	}
	return u
}

// Symmetric returns a new set which s is the difference of items  which are in
// one of either, but not in both.
func (s *Set) SymmetricDifference(t *Set) *Set {
	u := s.Difference(t)
	v := t.Difference(s)
	return u.Union(v)
}

// StringSlice is a helper function that returns a slice of strings of s. If
// the set contains mixed types of items only items of type string are returned.
func (s *Set) StringSlice() []string {
	slice := make([]string, 0)
	for _, item := range s.List() {
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
	slice := make([]int, 0)
	for _, item := range s.List() {
		v, ok := item.(int)
		if !ok {
			continue
		}

		slice = append(slice, v)
	}
	return slice
}
