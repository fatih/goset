package set

import (
	"fmt"
	"strings"
)

// Provides a common set baseline for both threadsafe and non-ts Sets.
type setNonTS map[interface{}]struct{} // struct{} doesn't take up space

// NewNonTS creates and initializes a new non-threadsafe Set.
func newNonTS() setNonTS {
	return setNonTS(make(map[interface{}]struct{}))
}

// Add includes the specified items (one or more) to the set. The underlying
// Set s is modified. If passed nothing it silently returns.
func (s setNonTS) Add(items ...interface{}) {
	for _, item := range items {
		s[item] = keyExists
	}
}

// Remove deletes the specified items from the set.  The underlying Set s is
// modified. If passed nothing it silently returns.
func (s setNonTS) Remove(items ...interface{}) {
	for _, item := range items {
		delete(s, item)
	}
}

// Pop  deletes and return an item from the set. The underlying Set s is
// modified. If set is empty, nil is returned.
func (s setNonTS) Pop() interface{} {
	for item := range s {
		delete(s, item)
		return item
	}
	return nil
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s setNonTS) Has(items ...interface{}) bool {
	// assume checked for empty item, which not exist
	for _, item := range items {
		if _, has := s[item]; !has {
			return false
		}
	}
	return len(items) != 0
}

// Size returns the number of items in a set.
func (s setNonTS) Size() int {
	return len(s)
}

// IsEmpty reports whether the Set is empty.
func (s setNonTS) IsEmpty() bool {
	return s.Size() == 0
}

// IsEqual test whether s and t are the same in size and have the same items.
func (s setNonTS) IsEqual(t Interface) bool {
	// Force locking only if given set is threadsafe.
	if conv, ok := t.(set); ok {
		conv.l.RLock()
		defer conv.l.RUnlock()
	}

	// return false if they are no the same size
	if len(s) != t.Size() {
		return false
	}

	equal := true
	t.Each(func(item interface{}) bool {
		_, equal = s[item]
		return equal // if false, Each() will end
	})

	return equal
}

// IsSubset tests whether t is a subset of s.
func (s setNonTS) IsSubset(t Interface) (subset bool) {
	subset = true

	t.Each(func(item interface{}) bool {
		_, subset = s[item]
		return subset
	})

	return
}

// IsSuperset tests whether t is a superset of s.
func (s setNonTS) IsSuperset(t Interface) bool {
	return t.IsSubset(s)
}

// Each traverses the items in the Set, calling the provided function for each
// set member. Traversal will continue until all items in the Set have been
// visited, or if the closure returns false.
func (s setNonTS) Each(f func(item interface{}) bool) {
	for item := range s {
		if !f(item) {
			break
		}
	}
}

// Copy returns a new Set with a copy of s.
func (s setNonTS) Copy() Interface {
	u := newNonTS()
	for item := range s {
		u.Add(item)
	}
	return u
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s setNonTS) List() []interface{} {
	list := make([]interface{}, 0, len(s))

	for item := range s {
		list = append(list, item)
	}

	return list
}

// String returns a string representation of s
func (s setNonTS) String() string {
	t := make([]string, len(s.List()))
	for i, item := range s.List() {
		t[i] = fmt.Sprintf("%v", item)
	}

	return fmt.Sprintf("[%s]", strings.Join(t, ", "))
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s setNonTS) Merge(t Interface) {
	t.Each(func(item interface{}) bool {
		s[item] = keyExists
		return true
	})
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s setNonTS) Separate(t Interface) {
	s.Remove(t.List()...)
}
