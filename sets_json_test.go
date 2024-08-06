// Copyright (c) Fortio Authors, All Rights Reserved
// See LICENSE for licensing terms. (Apache-2.0)

package sets_test

import (
	"encoding/json"
	"testing"

	"fortio.org/assert"
	"fortio.org/sets"
)

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

func TestBadJson(t *testing.T) {
	jsonStr := `[
		"a,b",
		"c,d"
	]`
	s := sets.New[int]()
	err := json.Unmarshal([]byte(jsonStr), &s)
	assert.Error(t, err)
}
