// Copyright (c) Fortio Authors, All Rights Reserved
// See LICENSE for licensing terms. (Apache-2.0)

package sets_test

import (
	"testing"

	"fortio.org/assert"
	"fortio.org/sets"
)

func TestSetToString(t *testing.T) {
	s := sets.Set[string]{"z": {}, "a": {}, "c": {}, "b": {}}
	assert.Equal(t, "a,b,c,z", s.String())
}

func TestArrayToSet(t *testing.T) {
	a := []string{"z", "a", "c", "b"}
	s := sets.FromSlice(a)
	assert.Equal(t, "a,b,c,z", s.String())
	assert.Equal(t, s.Sorted(), []string{"a", "b", "c", "z"})
}

func TestRemoveCommon(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setAA := setA.Clone()
	setBB := setB.Clone()
	sets.RemoveCommon(setAA, setBB)
	assert.Equal(t, "a,c", setAA.String())   // removed
	assert.Equal(t, "e,f,g", setBB.String()) // added
	// Swap order to exercise the optimization on length of iteration
	// also check clone is not modifying the original etc
	setAA = setB.Clone() // putting B in AA on purpose and vice versa
	setBB = setA.Clone()
	sets.RemoveCommon(setAA, setBB)
	assert.Equal(t, "a,c", setBB.String())
	assert.Equal(t, "e,f,g", setAA.String())
	assert.True(t, setBB.Has("c"))
	setBB.Remove("c")
	assert.False(t, setBB.Has("c"))
}
