// (Fortio) Sets.
//
// (c) 2023 Fortio Authors
// See LICENSE

// Sets and Set[T] type and operations of any comparable type (go 1.18+ generics)
// [Intersection], [Union], [Subset], difference aka [Minus], [XOR],
// JSON serialization and deserialization and more.
package sets // import "fortio.org/sets"

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

// Set defines a low memory footprint set of any comparable type. Based on `map[T]struct{}`.
type Set[T comparable] map[T]struct{}

// New returns a new set containing the given elements.
func New[T comparable](item ...T) Set[T] {
	// best pre-allocation if there are no duplicates
	res := make(Set[T], len(item))
	res.Add(item...)
	return res
}

// FromSlice constructs a Set from a slice.
// [Elements] is the inverse function, getting back a slice from the Set.
// This is a short cut/alias for New[T](items...).
func FromSlice[T comparable](items []T) Set[T] {
	return New(items...)
}

// Clone returns a copy of the set.
func (s Set[T]) Clone() Set[T] {
	res := make(Set[T], len(s))
	for k := range s {
		res.Add(k)
	}
	return res
}

// Add items to the set.
func (s Set[T]) Add(item ...T) {
	for _, i := range item {
		s[i] = struct{}{}
	}
}

// Has returns true if the item is present in the set.
func (s Set[T]) Has(item T) bool {
	_, found := s[item]
	return found
}

// Remove items from the set.
func (s Set[T]) Remove(item ...T) {
	for _, i := range item {
		delete(s, i)
	}
}

// Union returns a new set that has all the elements of all the sets.
// Note that Union(s1) == s1.Clone() and Union[T]() == New[T]().
func Union[T comparable](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return New[T]()
	}
	res := sets[0].Clone()
	for _, s := range sets[1:] {
		for k := range s {
			res.Add(k)
		}
	}
	return res
}

// Intersection returns a new set that has the elements common to all the input sets.
func Intersection[T comparable](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return New[T]()
	}
	res := sets[0].Clone()
	for _, s := range sets[1:] {
		if len(res) == 0 { // no point in continuing if already empty
			return res
		}
		for k := range res {
			if !s.Has(k) {
				res.Remove(k)
			}
		}
	}
	return res
}

// Elements returns a slice of the elements in the set.
func (s Set[T]) Elements() []T {
	res := make([]T, 0, len(s))
	for k := range s {
		res = append(res, k)
	}
	return res
}

// Subset returns true if all elements of s are in the passed in set.
func (s Set[T]) Subset(bigger Set[T]) bool {
	for k := range s {
		if !bigger.Has(k) {
			return false
		}
	}
	return true
}

// Minus mutates the receiver to remove all the elements of the passed in set.
// If you want a copy use s.Clone().Minus(other). Returns the receiver for chaining.
func (s Set[T]) Minus(other Set[T]) Set[T] {
	for k := range other {
		s.Remove(k)
	}
	return s
}

// Plus is similar to Union but mutates the receiver. Added for symmetry with Minus.
// Returns the receiver for chaining.
func (s Set[T]) Plus(others ...Set[T]) Set[T] {
	for _, o := range others {
		s.Add(o.Elements()...)
	}
	return s
}

// Equals returns true if the two sets have the same elements.
func (s Set[T]) Equals(other Set[T]) bool {
	return len(s) == len(other) && s.Subset(other)
}

// Len returns the number of elements in the set (same as len(s) but as a method).
func (s Set[T]) Len() int {
	return len(s)
}

// Clear removes all elements from the set.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// String() returns a coma separated list of the elements in the set.
// This is mostly for troubleshooting/debug output unless the [T] serializes
// to a string that doesn't contain commas.
func (s Set[T]) String() string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, fmt.Sprint(k))
	}
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

// RemoveCommon removes elements from both sets that are in both,
// leaving only the delta. Useful when a is an old value and b is new
// and you want to apply some operation on all removed and added elements.
func RemoveCommon[T comparable](a, b Set[T]) {
	if len(a) > len(b) {
		a, b = b, a
	}
	for e := range a {
		if _, found := b[e]; found {
			delete(a, e)
			delete(b, e)
		}
	}
}

// XOR is an alias for [RemoveCommon], efficiently removes from each set the common
// elements.
func XOR[T comparable](a, b Set[T]) {
	RemoveCommon(a, b)
}

// -- Serialization

// MarshalJSON implements the json.Marshaler interface and only gets the elements as an array.
func (s Set[T]) MarshalJSON() ([]byte, error) {
	// How to handle all ordered at once??
	switch v := any(s).(type) {
	case Set[string]:
		return json.Marshal(Sort(v))
	case Set[int]:
		return json.Marshal(Sort(v))
	case Set[int8]:
		return json.Marshal(Sort(v))
	case Set[int64]:
		return json.Marshal(Sort(v))
	case Set[float64]:
		return json.Marshal(Sort(v))
	default:
		return json.Marshal(s.Elements())
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface turns the slice back to a Set.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*s = New[T](items...)
	return nil
}

// -- Additional operations on sets of ordered types

// Sort returns a sorted slice of the elements in the set.
// Only applicable for when the type is sortable.
func Sort[Q constraints.Ordered](s Set[Q]) []Q {
	keys := s.Elements()
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
