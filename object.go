package jstream

import (
	"encoding/json"
	"io"
)

type Entry[T Iterator] struct {
	name string
	iter T
}

func (e Entry[T]) Name() string {
	return e.name
}

func (e *Entry[T]) Next(dec *json.Decoder) (err error) {
	return e.iter.Next(dec)
}

func Key[T Iterator](name string, iter T) *Entry[T] {
	return &Entry[T]{name: name, iter: iter}
}

type Object struct {
	pos  position
	Keys []NamedIterator
}

func (o *Object) Next(dec *json.Decoder) (err error) {
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
			for _, key := range o.Keys {
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
