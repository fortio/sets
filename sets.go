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

	"golang.org/x/exp/constraints"
)

type Set[T constraints.Ordered] map[T]struct{}

// SetFromSlice constructs a Set from a slice.
func FromSlice[T constraints.Ordered](items []T) Set[T] {
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

func New[T constraints.Ordered](item ...T) Set[T] {
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
	return strings.Join(keys, ",")
}

// RemoveCommon removes elements from both sets that are in both,
// leaving only the delta. Useful for Notifier on Set so that
// oldValue has what has been removed and newValue has what has been added.
func RemoveCommon[T constraints.Ordered](a, b Set[T]) {
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

func (s Set[T]) Sorted() []T {
	keys := make([]T, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
