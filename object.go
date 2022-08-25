package current

import (
	"encoding/json"
	"fmt"
	"io"
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

func Key[T Streamer](name string, str T) *entry[T] {
	return &entry[T]{name: name, str: str}
}

func Object(commit func(), keys ...NamedStreamer) *object {
	return &object{keys: keys, commit: commit}
}

type object struct {
	keys   []NamedStreamer
	commit func()
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
			o.commit()
			return io.EOF
		}
		for _, key := range o.keys {
			if key.Name() == tok {
				if err = Stream(dec, key); err != nil {
					return err
				}
				break
			}
		}
	}
}
