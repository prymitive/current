package jstream

import (
	"encoding/json"
	"io"
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

func Object(keys ...NamedIterator) object {
	return object{keys: keys}
}

type object struct {
	pos  position
	keys []NamedIterator
}

func (o *object) Next(dec *json.Decoder) (err error) {
	switch o.pos {
	case posFirst:
		if err = requireToken(dec, mapStart); err != nil {
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
				break
			}
			for _, key := range o.keys {
				if key.Name() == tok {
					if err = stream(dec, key); err != nil {
						return err
					}
					goto NEXT
				}
			}
		NEXT:
		}
		o.pos = posLast
	case posLast:
		if err = requireToken(dec, mapEnd); err != nil {
			return err
		}
		o.pos = posEOF
	case posEOF:
		return io.EOF

	}
	return nil
}
