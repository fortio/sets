// (Fortio) Sets.
//
// (c) 2023 Fortio Authors
// See LICENSE

//go:build !no_json
// +build !no_json

package sets // import "fortio.org/sets"

import (
	"encoding/json"
)

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
