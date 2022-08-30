package current

import (
	"encoding/json"
	"fmt"
)

// Map decodes a json object into a simple map[string]T
// It takes a list of current.Key instances for each of the keys you
// want to decode.
// Example:
//
// Let's say we want to decode:
//
// {"tags" : {"name": "bob", "alias": "b"}}
//
//	current.Object(
//		current.Key("tags", current.Map(func(k,v string) {
//			fmt.Printf("tag %s=%s", k,v)
//		}),
//	)
//
// revive:disable:unexported-return
func Map[T any](commit func(k string, v T)) *jmap[T] {
	return &jmap[T]{commit: commit}
}

type jmap[T any] struct {
	commit func(k string, v T)
	zero   T
}

func (m jmap[T]) String() string {
	return fmt.Sprintf("Map[%T]", m.zero)
}

func (m *jmap[T]) Stream(dec *json.Decoder) (err error) {
	if err = requireToken(dec, mapStart, m); err != nil {
		return err
	}

	var tok json.Token
	var key string
	var val T
	for dec.More() {
		if tok, err = dec.Token(); err != nil {
			return err
		}
		key = tok.(string)

		if err = dec.Decode(&val); err != nil {
			return err
		}

		m.commit(key, val)
	}

	if err = requireToken(dec, mapEnd, m); err != nil {
		return err
	}

	return nil
}
