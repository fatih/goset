package set

import "sync"

// set defines a thread safe set data structure.
type set struct {
	s setNonTS
	l *sync.RWMutex // we name it because we don't want to expose it
}

// New creates and initialize a new set. It's accept a variable number of
// arguments to populate the initial set. If nothing passed a set with zero
// size is created.
func newTS() set {
	return set{
		s: newNonTS(),
		l: new(sync.RWMutex),
	}
}

// Add includes the specified items (one or more) to the set. The underlying
// set s is modified. If passed nothing it silently returns.
func (s set) Add(items ...interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s.Add(items...)
}

// Remove deletes the specified items from the set.  The underlying set s is
// modified. If passed nothing it silently returns.
func (s set) Remove(items ...interface{}) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s.Remove(items...)
}

// Pop  deletes and return an item from the set. The underlying set s is
// modified. If set is empty, nil is returned.
func (s set) Pop() interface{} {
	s.l.RLock()
	for item := range s.s {
		s.l.RUnlock()
		s.Remove(item)
		return item
	}
	s.l.RUnlock()
	return nil
}

// Has looks for the existence of items passed. It returns false if nothing is
// passed. For multiple items it returns true only if all of  the items exist.
func (s set) Has(items ...interface{}) bool {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.Has(items...)
}

// Size returns the number of items in a set.
func (s set) Size() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.Size()
}

// IsEqual test whether s and t are the same in size and have the same items.
func (s set) IsEqual(t Interface) bool {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.IsEqual(t)
}

// IsSubset tests whether t is a subset of s.
func (s set) IsSubset(t Interface) (subset bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.IsSubset(t)
}

// Each traverses the items in the set, calling the provided function for each
// set member. Traversal will continue until all items in the set have been
// visited, or if the closure returns false.
func (s set) Each(f func(item interface{}) bool) {
	s.l.RLock()
	defer s.l.RUnlock()

	s.s.Each(f)
}

// List returns a slice of all items. There is also StringSlice() and
// IntSlice() methods for returning slices of type string or int.
func (s set) List() []interface{} {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.List()
}

// Copy returns a new set with a copy of s.
func (s set) Copy() Interface {
	s.l.RLock()
	defer s.l.RUnlock()

	// could not directly return s.s.Copy()
	// because s.s.Copy is typeof setNonTS
	return set{
		s: s.s.Copy().(setNonTS),
		l: new(sync.RWMutex),
	}
}

// IsSuperset tests whether t is a superset of s.
func (s set) IsSuperset(t Interface) bool {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.IsSuperset(t)
}

func (s set) IsEmpty() bool {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.IsEmpty()
}

// Merge is like Union, however it modifies the current set it's applied on
// with the given t set.
func (s set) Merge(t Interface) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s.Merge(t)
}

// it's not the opposite of Merge.
// Separate removes the set items containing in t from set s. Please aware that
func (s set) Separate(t Interface) {
	s.l.Lock()
	defer s.l.Unlock()

	s.s.Separate(t)
}

// String returns a string representation of s
func (s set) String() string {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.s.String()
}
