// (Fortio) Sets.
//
// (c) 2023 Fortio Authors
// See LICENSE

// Sets and Set type and operations in go 1.18+ generics.
// (pending built in support in golang core)
package sets // import "fortio.org/sets"

import (
	"fmt"
	"sort"
	"strings"
)

type Set[T comparable] map[T]struct{}

// SetFromSlice constructs a Set from a slice.
func FromSlice[T comparable](items []T) Set[T] {
	// best pre-allocation if there are no duplicates
	res := make(map[T]struct{}, len(items))
	for _, item := range items {
		res[item] = struct{}{}
	}
	return res
}

func (s Set[T]) Clone() Set[T] {
	res := make(Set[T], len(s))
	for k := range s {
		res.Add(k)
	}
	return res
}

func (s Set[T]) Add(item ...T) {
	for _, i := range item {
		s[i] = struct{}{}
	}
}

func (s Set[T]) Has(item T) bool {
	_, found := s[item]
	return found
}

func (s Set[T]) Remove(item ...T) {
	for _, i := range item {
		delete(s, i)
	}
}

func New[T comparable](item ...T) Set[T] {
	res := make(Set[T], len(item))
	res.Add(item...)
	return res
}

func (s Set[T]) String() string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, fmt.Sprintf("%v", k))
	}
	sort.Strings(keys)
	return  strings.Join(keys, ",")
}

// RemoveCommon removes elements from both sets that are in both,
// leaving only the delta. Useful for Notifier on Set so that
// oldValue has what has been removed and newValue has what has been added.
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
