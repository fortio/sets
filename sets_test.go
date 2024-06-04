// Copyright (c) Fortio Authors, All Rights Reserved
// See LICENSE for licensing terms. (Apache-2.0)

package sets_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"fortio.org/assert"
	"fortio.org/sets"
)

func TestSetToString(t *testing.T) {
	s := sets.Set[string]{"z": {}, "a": {}, "c": {}, "b": {}}
	assert.Equal(t, "a,b,c,z", s.String())
	assert.Equal(t, s.Len(), 4)
	s.Clear()
	assert.Equal(t, "", s.String())
	assert.Equal(t, s.Len(), 0)
}

func TestArrayToSet(t *testing.T) {
	a := []string{"z", "a", "c", "b"}
	s := sets.FromSlice(a)
	assert.Equal(t, "a,b,c,z", s.String())
	assert.Equal(t, sets.Sort(s), []string{"a", "b", "c", "z"})
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
	assert.True(t, setAA.Equals(setB))
	assert.True(t, setB.Equals(setAA))
	assert.False(t, setAA.Equals(setA))
	assert.False(t, setB.Equals(setBB))
	sets.XOR(setAA, setBB)
	assert.Equal(t, "a,c", setBB.String())
	assert.Equal(t, "e,f,g", setAA.String())
	assert.True(t, setBB.Has("c"))
	setBB.Remove("c")
	assert.False(t, setBB.Has("c"))
}

func TestMinus(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setAB := setA.Clone().Minus(setB)
	setBA := setB.Clone().Minus(setA)
	assert.Equal(t, "a,c", setAB.String())
	assert.Equal(t, "e,f,g", setBA.String())
}

func TestPlus(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setAB := setA.Clone().Plus(setB)
	setBA := setB.Clone().Plus(setA)
	assert.Equal(t, "a,b,c,d,e,f,g", setAB.String())
	assert.Equal(t, "a,b,c,d,e,f,g", setBA.String())
}

func TestUnion(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setC := sets.Union(sets.Union[string](), setA, setB)
	assert.Equal(t, "a,b,c,d,e,f,g", setC.String())
}

func TestIntersection1(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setC := sets.Intersection(setA, setB)
	assert.Equal(t, "b,d", setC.String())
}

func TestIntersection2(t *testing.T) {
	assert.Equal(t, len(sets.Intersection[string]()), 0)
	setA := sets.New("a", "b", "c")
	setB := sets.New("d", "e", "f")
	// cover stop early when empty intersection is reached, ie 3rd set won't be looked at
	setC := sets.Intersection(setA, setB, setA)
	assert.Equal(t, "", setC.String())
}

func TestSubset(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setC := sets.New("b", "d")
	assert.True(t, setC.Subset(setA))
	assert.True(t, setA.Subset(setA))
	assert.False(t, setA.Subset(setC))
	assert.False(t, setA.Subset(setB))
	assert.False(t, setB.Subset(setA))
}

func TestJSON(t *testing.T) {
	setA := sets.New("c,d", "a b", "y\000z", "mno")
	b, err := json.Marshal(setA)
	assert.NoError(t, err)
	assert.Equal(t, `["a b","c,d","mno","y\u0000z"]`, string(b))
	jsonStr := `[
		"a,b",
		"c,d"
	]`
	setB := sets.New[string]()
	err = json.Unmarshal([]byte(jsonStr), &setB)
	assert.NoError(t, err)
	assert.Equal(t, setB.Len(), 2)
	assert.True(t, setB.Has("a,b"))
	assert.True(t, setB.Has("c,d"))
	setI := sets.New(3, 42, 7, 10)
	b, err = json.Marshal(setI)
	assert.NoError(t, err)
	assert.Equal(t, `[3,7,10,42]`, string(b))
	smallIntSet := sets.New[int8](66, 65, 67) // if using byte, aka uint8, one gets base64("ABC")
	b, err = json.Marshal(smallIntSet)
	assert.NoError(t, err)
	t.Logf("smallIntSet: %q", string(b))
	assert.Equal(t, `[65,66,67]`, string(b))
	floatSet := sets.New[float64](2.3, 1.1, -7.6, 42)
	b, err = json.Marshal(floatSet)
	assert.NoError(t, err)
	t.Logf("floatSet: %q", string(b))
	assert.Equal(t, `[-7.6,1.1,2.3,42]`, string(b))
	i64Set := sets.New[int64](2, 1, -7, 42)
	b, err = json.Marshal(i64Set)
	assert.NoError(t, err)
	t.Logf("i64Set: %q", string(b))
	assert.Equal(t, `[-7,1,2,42]`, string(b))
}

type foo struct {
	X int
}

func TestNonOrderedJSON(t *testing.T) {
	s := sets.New(
		foo{3},
		foo{1},
		foo{2},
		foo{4},
	)
	b, err := json.Marshal(s)
	t.Logf("b: %s", string(b))
	assert.NoError(t, err)
	// though I guess given it could be in any order it could be accidentally sorted too
	assert.NotEqual(t, `[{"X":1},{"X":2},{"X":3},{"X":4}]`, string(b))
	u := sets.New[foo]()
	err = json.Unmarshal(b, &u)
	assert.NoError(t, err)
	assert.Equal(t, 4, u.Len())
	assert.True(t, s.Equals(u))
}

func TestGenerate(t *testing.T) {
	setA := sets.New("a", "b", "c")
	res := sets.Tuplets(setA, 1)
	assert.Equal(t, res, [][]string{{"a"}, {"b"}, {"c"}}, "should match single/identical")
	res = sets.Tuplets(setA, 2)
	assert.Equal(t, res, [][]string{{"a", "b"}, {"a", "c"}, {"b", "a"}, {"b", "c"}, {"c", "a"}, {"c", "b"}}, "should match pairs")
	res = sets.Tuplets(setA, 3)
	assert.Equal(t, res, [][]string{{"a", "b", "c"}, {"a", "c", "b"}, {"b", "a", "c"}, {"b", "c", "a"}, {"c", "a", "b"}, {"c", "b", "a"}}, "should match triplets")
}

func TestBadJson(t *testing.T) {
	jsonStr := `[
		"a,b",
		"c,d"
	]`
	s := sets.New[int]()
	err := json.Unmarshal([]byte(jsonStr), &s)
	assert.Error(t, err)
}

func setup(b *testing.B, n int) sets.Set[int64] {
	s := sets.Set[int64]{}
	max := 8 * int64(n)
	i := 0
	for ; len(s) != n; i++ {
		// Add random elements to the set.
		s.Add(rand.Int63n(max)) // set is somewhat sparse
	}
	b.Logf("Took %d iterations to fill set", i)
	return s
}

var s1000 sets.Set[int64]

func BenchmarkSetSort1000(b *testing.B) {
	if s1000 == nil {
		s1000 = setup(b, 1000)
		b.ResetTimer()
	}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s1 := s1000.Clone()
		b.StartTimer()
		r := sets.Sort(s1)
		if len(r) != s1.Len() {
			b.Fatalf("unexpected length change: %d", len(r))
		}
	}
}
