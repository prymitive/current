package current

import (
	"encoding/json"
	"fmt"
	"strings"
)

type entry[T Streamer] struct {
	name string
	str  T
}

func (e entry[T]) Name() string {
	return e.name
}

func (e *entry[T]) Stream(dec *json.Decoder) (err error) {
	return e.str.Stream(dec)
}

// revive:disable:unexported-return
func Key[T Streamer](name string, str T) *entry[T] {
	return &entry[T]{name: name, str: str}
}

// Object decodes a json object, mostly useful when different keys
// must be decoded in a different way.
// It takes a list of current.Key instances for each of the keys you
// want to decode.
// Example:
//
// Let's say we want to decode:
//
// {"name": "bob", "age": 4}
//
//	current.Object(
//		current.Key("name", current.Value(func(s string, isNull bool) {
//			fmt.Printf("name is %q", s)
//		}),
//		current.Key("age", current.Value(func(i int, isNull bool) {
//			fmt.Printf("age is %d", i)
//		}),
//	)
//
// revive:disable:unexported-return
func Object(keys ...NamedStreamer) *object {
	return &object{keys: keys}
}

type object struct {
	keys []NamedStreamer
}

func (o object) String() string {
	keys := make([]string, 0, len(o.keys))
	for _, key := range o.keys {
		keys = append(keys, key.Name())
	}
	return fmt.Sprintf("Object{%s}", strings.Join(keys, ","))
}

func (o *object) Stream(dec *json.Decoder) (err error) {
	if err = requireToken(dec, mapStart, o); err != nil {
		return err
	}

	var tok json.Token
	for {
		if tok, err = dec.Token(); err != nil {
			return err
		}
		if tok == mapEnd {
			return nil
		}
		for _, key := range o.keys {
			if key.Name() == tok {
				if err = key.Stream(dec); err != nil {
					return err
				}
				break
			}
		}
	}
}
