package jstream

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type entry[T Iterator] struct {
	name string
	iter T
}

func (e entry[T]) Name() string {
	return e.name
}

func (e *entry[T]) Next(dec *json.Decoder) (err error) {
	return e.iter.Next(dec)
}

func Key[T Iterator](name string, iter T) *entry[T] {
	return &entry[T]{name: name, iter: iter}
}

func Object(commit func(), keys ...NamedIterator) *object {
	return &object{keys: keys, commit: commit}
}

type object struct {
	pos    position
	keys   []NamedIterator
	commit func()
}

func (o object) String() string {
	keys := make([]string, 0, len(o.keys))
	for _, key := range o.keys {
		keys = append(keys, key.Name())
	}
	return fmt.Sprintf("Object{%s}", strings.Join(keys, ","))
}

func (o *object) Next(dec *json.Decoder) (err error) {
	switch o.pos {
	case posFirst:
		if err = requireToken(dec, mapStart, o); err != nil {
			return err
		}
		o.pos = posDecoding
	case posDecoding:
		var tok json.Token
		for {
			if tok, err = dec.Token(); err != nil {
				return err
			}
			if tok == mapEnd {
				o.pos = posEOF
				return nil
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
	case posEOF:
		o.commit()
		return io.EOF
	}
	return nil
}
